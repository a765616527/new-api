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
		[]int{1, 2, 3},
		"/v1/images/generations",
		dto.GPTImage2Model,
		dto.ImageSizeTier2K,
	)
	channelSyncLock.RUnlock()

	assert.Equal(t, []int{2}, filtered)
}
