package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNils(t *testing.T) {
	var a any
	var i int
	require.True(t, Nils(a))
	require.False(t, Nils(&i))
}
