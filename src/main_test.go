package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test1(t *testing.T)  {
	result := test()
	require.True(t, result)
}

func Test2(t *testing.T)  {
	result := test2()
	require.False(t, result)
	_, err := os.Open("testdata/foo.json")
	require.Nil(t, err)

}
