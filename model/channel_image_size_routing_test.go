package model

import (
	"testing"

	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/stretchr/testify/assert"
)

func TestFilterChannelsByImageSizeTier(t *testing.T) {
	channelSyncLock.Lock()
	previousChannels := channelsIDM
	previousAdvancedConfigs := channel2advancedCustomConfig
	channelsIDM = map[int]*Channel{
		1: {Id: 1, OtherSettings: `{"gpt_image_2_size_models":{"1k":"vendor-image-1k"}}`},
		2: {Id: 2, OtherSettings: `{"gpt_image_2_size_models":{"2k":"vendor-image-2k"}}`},
		3: {Id: 3, OtherSettings: `{}`},
		4: {Id: 4, OtherSettings: `{"gpt_image_2_size_models":{"1k":" ","2k":"","4k":"  "}}`},
		5: {Id: 5, OtherSettings: `{"gpt_image_2_size_models":{"4k":"vendor-image-4k"}}`},
	}
	channel2advancedCustomConfig = map[int]*dto.AdvancedCustomConfig{}
	channelSyncLock.Unlock()
	t.Cleanup(func() {
		channelSyncLock.Lock()
		channelsIDM = previousChannels
		channel2advancedCustomConfig = previousAdvancedConfigs
		channelSyncLock.Unlock()
	})

	channelSyncLock.RLock()
	filtered := filterChannelsByRequestPathAndModel(
		[]int{1, 2, 3, 4, 5},
		"/v1/images/generations",
		dto.GPTImage2Model,
		dto.ImageSizeTier2K,
	)
	channelSyncLock.RUnlock()

	assert.Equal(t, []int{2, 3, 4}, filtered)
}

func TestGPTImage2UpstreamModelRoutingMode(t *testing.T) {
	tests := []struct {
		name     string
		settings string
		tier     string
		expected string
	}{
		{
			name:     "missing settings pass through",
			settings: `{}`,
			tier:     dto.ImageSizeTier4K,
			expected: dto.GPTImage2Model,
		},
		{
			name:     "all blank settings pass through",
			settings: `{"gpt_image_2_size_models":{"1k":" ","2k":"","4k":"  "}}`,
			tier:     dto.ImageSizeTier2K,
			expected: dto.GPTImage2Model,
		},
		{
			name:     "configured tier uses upstream model",
			settings: `{"gpt_image_2_size_models":{"2k":" vendor-image-2k "}}`,
			tier:     dto.ImageSizeTier2K,
			expected: "vendor-image-2k",
		},
		{
			name:     "partial settings do not fall back for blank tier",
			settings: `{"gpt_image_2_size_models":{"4k":"vendor-image-4k"}}`,
			tier:     dto.ImageSizeTier2K,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channel := &Channel{OtherSettings: tt.settings}

			assert.Equal(t, tt.expected, channel.GetGPTImage2UpstreamModel(tt.tier))
		})
	}
}
