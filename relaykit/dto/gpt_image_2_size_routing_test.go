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
		{name: "portrait under 1.5k uses 1k", size: "640x1040", expected: ImageSizeTier1K},
		{name: "below 1.5k uses 1k", size: "1500x1024", expected: ImageSizeTier1K},
		{name: "classic 1.5k boundary uses 1k", size: "1536x1024", expected: ImageSizeTier1K},
		{name: "above 1.5k uses 2k", size: "1537x1024", expected: ImageSizeTier2K},
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

func TestIsGPTImageSizeRoutedModel(t *testing.T) {
	assert.True(t, IsGPTImageSizeRoutedModel(GPTImage2Model))
	assert.True(t, IsGPTImageSizeRoutedModel(GPTImage25FlareModel))
	assert.True(t, IsGPTImageSizeRoutedModel(GPTImage25SunburstModel))
	assert.False(t, IsGPTImageSizeRoutedModel("gpt-image-2.5"))
	assert.False(t, IsGPTImageSizeRoutedModel("dall-e-3"))
}
