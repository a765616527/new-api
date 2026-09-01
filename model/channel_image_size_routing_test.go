package model

import (
	"testing"

	taskdto "github.com/QuantumNous/new-api/dto"
	relaydto "github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/stretchr/testify/assert"
)

func TestFilterChannelsByImageSizeTier(t *testing.T) {
	channelSyncLock.Lock()
	previousChannels := channelsIDM
	channelsIDM = map[int]*Channel{
		1: {Id: 1, OtherSettings: `{"gpt_image_2_size_models":{"1k":"vendor-image-1k"}}`},
		2: {Id: 2, OtherSettings: `{"gpt_image_2_size_models":{"2k":"vendor-image-2k"}}`},
		3: {Id: 3, OtherSettings: `{}`},
		4: {Id: 4, OtherSettings: `{"gpt_image_2_size_models":{"1k":" ","2k":"","4k":"  "}}`},
		5: {Id: 5, OtherSettings: `{"gpt_image_2_size_models":{"4k":"vendor-image-4k"}}`},
	}
	channelSyncLock.Unlock()
	t.Cleanup(func() {
		channelSyncLock.Lock()
		channelsIDM = previousChannels
		channelSyncLock.Unlock()
	})

	channelSyncLock.RLock()
	filtered, emptiedBy := filterCandidateIDs(
		[]int{1, 2, 3, 4, 5},
		relaydto.GPTImage2Model,
		[]taskdto.ChannelFilter{{
			Kind:          taskdto.FilterImageSizeTier,
			ImageSizeTier: relaydto.ImageSizeTier2K,
		}},
	)
	channelSyncLock.RUnlock()

	assert.Equal(t, []int{2, 3, 4}, filtered)
	assert.Empty(t, emptiedBy)
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
			tier:     relaydto.ImageSizeTier4K,
			expected: relaydto.GPTImage2Model,
		},
		{
			name:     "all blank settings pass through",
			settings: `{"gpt_image_2_size_models":{"1k":" ","2k":"","4k":"  "}}`,
			tier:     relaydto.ImageSizeTier2K,
			expected: relaydto.GPTImage2Model,
		},
		{
			name:     "configured tier uses upstream model",
			settings: `{"gpt_image_2_size_models":{"2k":" vendor-image-2k "}}`,
			tier:     relaydto.ImageSizeTier2K,
			expected: "vendor-image-2k",
		},
		{
			name:     "partial settings do not fall back for blank tier",
			settings: `{"gpt_image_2_size_models":{"4k":"vendor-image-4k"}}`,
			tier:     relaydto.ImageSizeTier2K,
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
