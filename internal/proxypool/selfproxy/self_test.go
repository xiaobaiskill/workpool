package selfproxy

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestNewSelf(t *testing.T) {
	s := NewSelf("https://hooks.slack.com/services/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")
	client,_ := s.Pop()
	time.Sleep(time.Second)

	req, err := http.NewRequest("GET", "https://jimqaweb.mlytics.ai/cache.txt", nil)
	if err != nil {
		return
	}

	rep,err := client.Do(req)
	assert.NoError(t,err)
	assert.Equal(t,true,rep.StatusCode == http.StatusOK)
}
