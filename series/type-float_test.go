package series

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatEq(t *testing.T) {
	f1 := (&floatElement{e: 1.42857142857, precision: 11})
	f2 := f1.Copy()
	assert.True(t, f1.Eq(f2), "Expected f1 to be equal to f2")
	assert.False(t, f1.WithPrecision(6).Eq(f2), "Expected f1 to be equal to f2")
	assert.True(t, f1.WithPrecision(2).Eq(f1.WithPrecision(2)), "Expected f1 to be equal to f2")

	f4 := f1.Copy().(*floatElement)
	f4.nan = true
	f5 := f4.Copy()
	assert.False(t, f4.Eq(f5), "NaN should not be equal to NaN")

	f6 := &floatElement{}
	f6.Set("1.428572")
	assert.False(t, f1.Eq(f6), "Expected f1 to be equal to f6")
	assert.True(t, f1.WithPrecision(5).Eq(f6.WithPrecision(5)), "Expected f1 to be equal to f6")

	f7 := &floatElement{}
	f7.Set(true)
	assert.False(t, f1.Eq(f7), "Expected f1 to be equal to f6")
	assert.True(t, f1.WithPrecision(0).Eq(f7), "Expected f1 to be equal to f6")
}
