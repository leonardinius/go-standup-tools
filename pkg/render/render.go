package render

import (
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/leonardinius/go-standup-tools/pkg/activityfeed"
)

const dateStringFormat string = "2006-01-02"

func RenderStandupMarkdown(feed activityfeed.ActivityFeedReport) string {
	itemStrings := make([]string, 0, feed.Len())
	var dateString string
	for _, item := range feed {
		itemDate := item.Updated.Format(dateStringFormat)
		if dateString != itemDate {
			dateString = itemDate
			itemStrings = append(itemStrings, fmt.Sprintf("- **%s**", dateString))
		}
		itemStrings = append(itemStrings,
			"- "+t(fmt.Sprintf("%s [%s](%s) - %s", t(strings.ToLower(item.Category)), t(item.IssueID), t(item.Link), t(item.Summary))))
	}

	return strings.Join(itemStrings, "\n")
}

func RenderStandupHTML(feed activityfeed.ActivityFeedReport) string {
	md := RenderStandupMarkdown(feed)
	return string(markdown.ToHTML([]byte(md), nil, nil))
}

func RenderStandupTXT(feed activityfeed.ActivityFeedReport) string {
	itemStrings := make([]string, 0, feed.Len())
	var dateString string
	for _, item := range feed {
		itemDate := item.Updated.Format(dateStringFormat)
		if dateString != itemDate {
			dateString = itemDate
			itemStrings = append(itemStrings, fmt.Sprintf("- **%s**", dateString))
		}
		itemStrings = append(itemStrings,
			"- "+t(fmt.Sprintf("%s %s - %s", t(strings.ToLower(item.Category)), t(item.IssueID), t(item.Summary))))
	}

	return strings.Join(itemStrings, "\n")
}

func t(s string) string {
	return strings.TrimSpace(s)
}
