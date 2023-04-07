package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	c := HttpClientDecorator{client: http.Client{}}
	res, err := c.Do(&http.Request{})
	assert.NotNil(t, err)
	assert.Nil(t, res)
}
