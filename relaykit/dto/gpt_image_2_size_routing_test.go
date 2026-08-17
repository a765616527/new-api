package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGPTImage2SizeTier(t *testing.T) {
	tests := []struct {
		name     string
		size     string
		expected string
	}{
		{name: "omitted uses 4k", size: "", expected: ImageSizeTier4K},
		{name: "auto uses 4k", size: "auto", expected: ImageSizeTier4K},
		{name: "square 1k", size: "1024x1024", expected: ImageSizeTier1K},
		{name: "landscape 1k", size: "1024x768", expected: ImageSizeTier1K},
		{name: "classic 1.5k uses 2k", size: "1536x1024", expected: ImageSizeTier2K},
		{name: "2k boundary", size: "2048x1152", expected: ImageSizeTier2K},
		{name: "under 4k still uses 4k", size: "2560x1440", expected: ImageSizeTier4K},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tier, err := GPTImage2SizeTier(tt.size)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, tier)
		})
	}
}

func TestGPTImage2SizeTierRejectsInvalidSize(t *testing.T) {
	for _, size := range []string{"1024", "0x1024", "1024xnope", "1024×1024"} {
		t.Run(size, func(t *testing.T) {
			_, err := GPTImage2SizeTier(size)
			require.Error(t, err)
		})
	}
}
