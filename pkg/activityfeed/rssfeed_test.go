package activityfeed_test

import (
	"testing"

	"github.com/leonardinius/go-standup-tools/pkg/activityfeed"
	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	assert.Equal(
		t,
		"https://test-domain.com/plugins/servlet/streams?maxResults=10&streams=update-date+AFTER+0&streams=update-date+BEFORE+1&streams=account-id+IS+userId&os_authType=basic",
		activityfeed.RssURL("https://test-domain.com", "userId", 10, 0, 1))
}
