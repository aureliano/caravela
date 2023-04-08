package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReleaseCompareTo(t *testing.T) {
	r1 := &Release{Name: "v0.1.0"}
	r2 := &Release{Name: "v0.1.1"}

	assert.Equal(t, r1.CompareTo(r2), -1)
}
