package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSSID(t *testing.T) {
	ssid, err := GetSSID()
	require.NoError(t, err)
	assert.NotEmpty(t, ssid)
}
