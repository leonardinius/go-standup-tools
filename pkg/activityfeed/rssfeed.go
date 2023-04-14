package activityfeed

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
	"golang.org/x/exp/maps"
)

type Config struct {
	Host      string
	Username  string
	Password  string
	AccountID string
	SinceDate time.Time
	TillDate  time.Time
	Verbose   bool
}

type ReportItem struct {
	Updated  *time.Time
	Category string
	IssueID  string
	Summary  string
	Link     string
}

type ActivityFeedReport []*ReportItem

func (feed ActivityFeedReport) Len() int {
	return len(feed)
}

func (feed ActivityFeedReport) FilterItems(matches func(*ReportItem) bool) ActivityFeedReport {
	feedCopy := make(ActivityFeedReport, 0, feed.Len())
	for _, item := range feed {
		if matches(item) {
			feedCopy = append(feedCopy, item)
		}
	}

	return feedCopy
}

func (feed ActivityFeedReport) SortItems(compare func(item1, item2 *ReportItem) int) ActivityFeedReport {
	sort.Slice(feed, func(i, j int) bool {
		item1 := feed[i]
		item2 := feed[j]
		return compare(item1, item2) < 0
	})

	return feed
}

var _ fmt.Stringer = (*ReportItem)(nil)

func (i *ReportItem) String() string {
	bytes, _ := json.Marshal(i)
	return string(bytes)
}

func RssURL(domain, accountID string, maxResults int, updatedAfter, updatedBefore int64) string {
	urlVariables := map[string]string{
		"domain":        domain,
		"accountId":     accountID,
		"maxResults":    strconv.FormatInt(int64(maxResults), 10),
		"updatedAfter":  strconv.FormatInt(updatedAfter, 10),
		"updatedBefore": strconv.FormatInt(updatedBefore, 10),
	}
	return os.Expand(
		"${domain}/plugins/servlet/streams?"+
			strings.Join(
				[]string{
					"maxResults=${maxResults}",
					"streams=update-date+AFTER+${updatedAfter}",
					"streams=update-date+BEFORE+${updatedBefore}",
					"streams=account-id+IS+${accountId}",
					"os_authType=basic"},
				"&"),
		func(s string) string { return urlVariables[s] },
	)
}

func ParseFromURL(cfg *Config, ctx context.Context) (report ActivityFeedReport, err error) {
	fp := gofeed.NewParser()
	fp.AuthConfig = &gofeed.Auth{
		Username: cfg.Username,
		Password: cfg.Password,
	}

	cancelCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	updatedAfterFilter := cfg.SinceDate.UnixMilli()
	updatedBeforeFilter := cfg.TillDate.UnixMilli()
	var feed *gofeed.Feed

	const maxItemsPerRequest = 1000
	for {
		url := RssURL(cfg.Host, cfg.AccountID, maxItemsPerRequest, updatedAfterFilter, updatedBeforeFilter)
		if cfg.Verbose {
			log.Printf("[INFO ] Fetching %s", url)
		}
		if feed, err = fp.ParseURLWithContext(url, cancelCtx); err != nil {
			return nil, err
		}

		if _report, _err := parseFeedItems(feed); _err != nil {
			return nil, _err
		} else {
			report = append(report, _report...)
			if cfg.Verbose {
				log.Printf("[INFO ] page of %d items, total %d", _report.Len(), report.Len())
			}
			if _report.Len() < maxItemsPerRequest*0.2 {
				break
			}
			// may skip same millisecond events
			updatedBeforeFilter = int64(math.Min(float64(_report[0].Updated.UnixMilli()), float64(_report[_report.Len()-1].Updated.UnixMilli())))
			if updatedBeforeFilter < updatedAfterFilter {
				break
			}
		}
	}
	return report, err
}

func ParseFromReader(reader io.Reader) (report ActivityFeedReport, err error) {
	var feed *gofeed.Feed

	if feed, err = gofeed.NewParser().Parse(reader); err != nil {
		return nil, err
	}

	return parseFeedItems(feed)
}

func parseFeedItems(feed *gofeed.Feed) (report ActivityFeedReport, err error) {
	reportItems := make(ActivityFeedReport, 0, feed.Len())
	for index, item := range feed.Items {
		date := item.UpdatedParsed
		link := item.Link

		category := strings.Join(item.Categories, "")
		var activityExt *ext.Extension
		if activityArray, ok := item.Extensions["activity"]["target"]; !ok {
			if activityArray, ok = item.Extensions["activity"]["object"]; !ok {
				err = multierror.Append(err, fmt.Errorf("unrecognized extension item[%d], extension[activity]: %#v", index, maps.Keys(item.Extensions["activity"])))
				continue
			}
			activityExt = &activityArray[0]
		} else {
			activityExt = &activityArray[0]
		}

		title := activityExt.Children["title"][0].Value
		summary := activityExt.Children["summary"][0].Value

		reportItem := &ReportItem{
			Updated:  date,
			Category: category,
			IssueID:  title,
			Summary:  summary,
			Link:     link,
		}

		reportItems = append(reportItems, reportItem)
	}

	return reportItems, err
}
