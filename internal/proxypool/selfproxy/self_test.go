package selfproxy

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewSelf(t *testing.T) {
	s := NewSelf("https://hooks.slack.com/services/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")
	client,_ := s.Pop()

	req, err := http.NewRequest("GET", "https://jimqaweb.mlytics.ai/cache.txt", nil)
	if err != nil {
		return
	}

	rep,err := client.Do(req)
	assert.NoError(t,err)
	assert.Equal(t,true,rep.StatusCode == http.StatusOK)
}
