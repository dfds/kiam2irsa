package k8s

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsPodUsingKiam(t *testing.T) {
	a := assert.New(t)
	a.False(isPodUsingKiam(false, false, false))
	a.False(isPodUsingKiam(false, false, true))
	a.False(isPodUsingKiam(false, true, true))

	a.False(isPodUsingKiam(true, true, true))
	a.True(isPodUsingKiam(true, true, false))
	a.True(isPodUsingKiam(true, false, false))

	a.True(isPodUsingKiam(true, false, true))
	a.False(isPodUsingKiam(false, true, false))
}

func TestIsPodUsingIrsa(t *testing.T) {
	a := assert.New(t)
	a.False(isPodUsingIrsa(false, false, false))
	a.False(isPodUsingIrsa(false, false, true))
	a.True(isPodUsingIrsa(false, true, true))

	a.False(isPodUsingIrsa(true, true, true))
	a.False(isPodUsingIrsa(true, true, false))
	a.False(isPodUsingIrsa(true, false, false))

	a.False(isPodUsingIrsa(true, false, true))
	a.False(isPodUsingIrsa(false, true, false))
}

func TestIsPodUsingBoth(t *testing.T) {
	a := assert.New(t)
	a.False(isPodUsingBoth(false, false, false))
	a.False(isPodUsingBoth(false, false, true))
	a.False(isPodUsingBoth(false, true, true))

	a.True(isPodUsingBoth(true, true, true))
	a.False(isPodUsingBoth(true, true, false))
	a.False(isPodUsingBoth(true, false, false))

	a.False(isPodUsingBoth(true, false, true))
	a.False(isPodUsingBoth(false, true, false))
}
