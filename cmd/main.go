package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/leonardinius/go-standup-tools/pkg/activityfeed"
	"github.com/leonardinius/go-standup-tools/pkg/clipboard"
	"github.com/leonardinius/go-standup-tools/pkg/render"
	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/lv"
	"github.com/tj/go-naturaldate"
)

const dateStringFormat string = "2006-01-02"

var (
	// COMMIT git commit
	COMMIT = "gitsha1"
	// BRANCH git branch
	BRANCH = "dirty"
)

// Opts with all cli commands and flags
type cliOpts struct {
	Host      string `long:"host" env:"JIRA_HOST" required:"true" description:"your jira server host"`
	Username  string `long:"username" env:"JIRA_USER" required:"true" description:"your jira user"`
	Password  string `long:"password" env:"JIRA_PASSWORD" required:"true" default:"testflo" description:"your jira api token or password"`
	AccountID string `long:"account-id" env:"JIRA_ACCOUNT_ID" required:"false" default:"" description:"your account ID"`
	Since     string `long:"since" required:"false" default:"" description:"human readable date, e.g. 'yesterday'"`
	Verbose   bool   `long:"verbose" short:"v" required:"false" description:"verbose output"`
}

func main() {
	var opts cliOpts
	p := flags.NewParser(&opts, flags.Default)

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			fmt.Printf("standup revision %s-%s\n", BRANCH, COMMIT)
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if opts.AccountID == "" {
		opts.AccountID = opts.Username
	}

	sinceDateTime, dateNow := standupReportDateRange()
	if opts.Since != "" {
		optionPast := naturaldate.WithDirection(naturaldate.Past)
		if sinceParsed, err := naturaldate.Parse(opts.Since, dateNow, optionPast); err == nil && sinceParsed != dateNow {
			sinceDateTime = sinceParsed
		} else {
			log.Fatalf("[ERROR] Err: unexpected input %s (%+v)", opts.Since, err)
		}
	}

	config := &activityfeed.Config{
		Host:      opts.Host,
		Username:  opts.Username,
		Password:  opts.Password,
		AccountID: opts.AccountID,
	}
	if opts.Verbose {
		log.Printf("[INFO ] %#v\n%s-%s", config, sinceDateTime.Format(dateStringFormat), dateNow.Format(dateStringFormat))
	}

	ctx := context.Background()
	feed, err := activityfeed.ParseFromURL(config, ctx)
	if err != nil {
		log.Fatalf("[ERROR] Err: %+v", err)
	}

	feed = feed.SortItems(sortReverse(sortItemsByUpdatedFn))
	feed = feed.FilterItems(func(item *activityfeed.ReportItem) bool {
		if item.Updated.After(sinceDateTime) || item.Updated.Equal(sinceDateTime) {
			return item.Updated.Before(dateNow)
		}
		return false
	})
	feed = feed.FilterItems(filterUniqIssueIDs())
	feed = feed.SortItems(sortItemsByUpdatedFn)

	txt := render.RenderStandupTXT(feed)
	log.Printf("[INFO ] Text:\n%s\n", txt)
	html := render.RenderStandupHTML(feed)
	if err := clipboard.CopyHTMLToClipboardAsRTF(html); err != nil {
		log.Fatalf("[ERROR] Err: %+v", err)
	}
	log.Printf("[INFO ] ✔️ copied to clipboard!")
}

func standupReportDateRange() (prevWorkDate, dateNow time.Time) {
	c := cal.NewBusinessCalendar()
	// add holidays that the business observes
	c.AddHoliday(lv.Holidays...)

	timeNow := time.Now()
	dateNow = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, time.UTC)

	prevWorkTime := c.WorkdaysFrom(dateNow, -1)
	prevWorkDate = time.Date(prevWorkTime.Year(), prevWorkTime.Month(), prevWorkTime.Day(), 0, 0, 0, 0, time.UTC)
	return prevWorkDate, dateNow
}

func sortItemsByUpdatedFn(item1, item2 *activityfeed.ReportItem) int {
	if item1.Updated.Before(*item2.Updated) {
		return -1
	}
	if item1.Updated.After(*item2.Updated) {
		return 1
	}
	return 0
}

func sortReverse(fn func(item1, item2 *activityfeed.ReportItem) int) func(item1, item2 *activityfeed.ReportItem) int {
	return func(item1, item2 *activityfeed.ReportItem) int {
		return fn(item2, item1)
	}
}

func filterUniqIssueIDs() func(item *activityfeed.ReportItem) bool {
	issueIDMap := make(map[string]string)
	return func(item *activityfeed.ReportItem) bool {
		key := fmt.Sprintf("%s-%s", item.Updated.Format(dateStringFormat), item.IssueID)
		if _, ok := issueIDMap[key]; !ok {
			issueIDMap[key] = key
			return true
		}
		return false
	}
}
