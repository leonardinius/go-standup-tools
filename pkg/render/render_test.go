package render_test

import (
	"strings"
	"testing"
	"time"

	"github.com/leonardinius/go-standup-tools/pkg/activityfeed"
	"github.com/leonardinius/go-standup-tools/pkg/render"
	"github.com/stretchr/testify/assert"
)

func TestRenderFunctions(t *testing.T) {
	yesterday := time.Date(2022, time.November, 25, 0, 0, 0, 0, time.UTC)
	today := time.Date(2022, time.November, 26, 0, 0, 0, 0, time.UTC)
	feed := activityfeed.ActivityFeedReport{
		&activityfeed.ReportItem{
			Updated:  &yesterday,
			Category: "updated",
			IssueID:  "TP-141",
			Summary:  "S1",
			Link:     "https://jira.demo.net/browse/TP-141",
		},
		&activityfeed.ReportItem{
			Updated:  &yesterday,
			Category: "created",
			IssueID:  "TP-142",
			Summary:  "S2",
			Link:     "https://jira.demo.net/browse/TP-142",
		},
		&activityfeed.ReportItem{
			Updated:  &today,
			Category: "created",
			IssueID:  "TP-143",
			Summary:  "S3",
			Link:     "https://jira.demo.net/browse/TP-143",
		},
	}

	t.Run("RenderStandupMarkdown", func(t *testing.T) {
		expectedHTML := joinLines([]string{
			`- **2022-11-25**`,
			`- updated [TP-141](https://jira.demo.net/browse/TP-141) - S1`,
			`- created [TP-142](https://jira.demo.net/browse/TP-142) - S2`,
			`- **2022-11-26**`,
			`- created [TP-143](https://jira.demo.net/browse/TP-143) - S3`,
		})
		assert.Equal(t, expectedHTML, strings.TrimSpace(render.RenderStandupMarkdown(feed)))
	})

	t.Run("RenderStandupHTML", func(t *testing.T) {
		expectedHTML := joinLines([]string{
			`<ul>`,
			`<li><strong>2022-11-25</strong></li>`,
			`<li>updated <a href="https://jira.demo.net/browse/TP-141">TP-141</a> - S1</li>`,
			`<li>created <a href="https://jira.demo.net/browse/TP-142">TP-142</a> - S2</li>`,
			`<li><strong>2022-11-26</strong></li>`,
			`<li>created <a href="https://jira.demo.net/browse/TP-143">TP-143</a> - S3</li>`,
			`</ul>`,
		})
		assert.Equal(t, expectedHTML, strings.TrimSpace(render.RenderStandupHTML(feed)))
	})

	t.Run("RenderStandupTXT", func(t *testing.T) {
		expectedTxt := joinLines([]string{
			`- **2022-11-25**`,
			`- updated TP-141 - S1`,
			`- created TP-142 - S2`,
			`- **2022-11-26**`,
			`- created TP-143 - S3`,
		})
		assert.Equal(t, expectedTxt, strings.TrimSpace(render.RenderStandupTXT(feed)))
	})
}

func joinLines(lines []string) string {
	return strings.Join(lines, "\n")
}
