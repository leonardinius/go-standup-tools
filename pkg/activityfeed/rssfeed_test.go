package activityfeed_test

import (
	"testing"

	"github.com/leonardinius/go-standup-tools/pkg/activityfeed"
	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	assert.Equal(
		t,
		"https://test-domain.com/activity?maxResults=10&streams=account-id+IS+userId&os_authType=basic",
		activityfeed.RssURL("https://test-domain.com", "userId", 10))
}
