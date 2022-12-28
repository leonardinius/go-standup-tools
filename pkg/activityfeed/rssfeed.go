package activityfeed

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func RssURL(domain, accountID string, maxResults int) string {
	urlVariables := map[string]string{
		"domain":     domain,
		"accountId":  accountID,
		"maxResults": strconv.FormatInt(int64(maxResults), 10),
	}
	return os.Expand(
		"${domain}/activity?maxResults=${maxResults}&streams=account-id+IS+${accountId}&os_authType=basic",
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

	var feed *gofeed.Feed
	if feed, err = fp.ParseURLWithContext(RssURL(cfg.Host, cfg.AccountID, 99999), cancelCtx); err != nil {
		return nil, err
	}

	return parseFeedItems(feed)
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
