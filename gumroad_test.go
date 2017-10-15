package gumroad_test

import (
	"testing"

	"github.com/segmentio/go-env"
	"github.com/tj/assert"
	"github.com/tj/gumroad"
)

func TestLicenses_Verify(t *testing.T) {
	c := gumroad.New()

	l, err := c.Licenses.Verify("apex-up", env.MustGet("LICENSE"))
	assert.NoError(t, err, "verify")

	assert.False(t, l.Cancelled(), "cancelled")
	assert.False(t, l.Failed(), "failed")
}
