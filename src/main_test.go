package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test1(t *testing.T)  {
	result := test()
	require.True(t, result)
}

func Test2(t *testing.T)  {
	result := test2()
	require.False(t, result)
}
