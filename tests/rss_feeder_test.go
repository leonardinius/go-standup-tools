package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/leonardinius/go-standup-tools/pkg/activityfeed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseExamples(t *testing.T) {
	file, err := os.Open("testData/activity.xml")
	require.Nil(t, err, "Unexpected error")
	feed, err := activityfeed.ParseFromReader(file)
	require.Nil(t, err, "Unexpected error")
	assert.Equal(t, 2, feed.Len())

	itemStrings := make([]string, 0, feed.Len())
	for _, i := range feed {
		itemStrings = append(itemStrings, i.String())
	}

	// nolint: gocritic
	expectedItemJSON := strings.Join([]string{
		`{"Updated":"2022-11-15T15:59:36.288Z","Category":"","IssueID":"TP-141","Summary":"S1","Link":"https://jira.demo.net/browse/TP-141"}`,
		`{"Updated":"2022-11-15T15:59:09Z","Category":"created","IssueID":"TP-142","Summary":"T1","Link":"https://jira.demo.net/browse/TP-142"}`,
	}, ",\n")
	assert.Equal(t, expectedItemJSON, strings.Join(itemStrings, ",\n"))
}
