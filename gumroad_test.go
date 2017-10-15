package gumroad_test

import (
	"testing"

	"github.com/segmentio/go-env"
	"github.com/tj/assert"
	"github.com/tj/gumroad"
)

func TestLicenses_Verify(t *testing.T) {
	c := gumroad.New()

	t.Run("valid", func(t *testing.T) {
		l, err := c.Licenses.Verify("apex-up", env.MustGet("LICENSE"))
		assert.NoError(t, err, "verify")

		assert.False(t, l.Cancelled(), "cancelled")
		assert.False(t, l.Failed(), "failed")
	})

	t.Run("invalid", func(t *testing.T) {
		_, err := c.Licenses.Verify("apex-up", "40C376D3-3F814C92-A7E31F47-11111111")
		assert.EqualError(t, err, `That license does not exist for the provided product.`)

		e := err.(gumroad.Error)
		assert.Equal(t, 404, e.Status)
		assert.Equal(t, false, e.Success)
		assert.Equal(t, `That license does not exist for the provided product.`, e.Message)
	})
}
