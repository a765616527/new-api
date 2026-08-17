package service

import (
	"testing"

	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
	hosttypes "github.com/QuantumNous/new-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppendGPTImage2RequestInfoRecordsStructuredFields(t *testing.T) {
	savedModelPrices := ratio_setting.ModelPrice2JSONString()
	t.Cleanup(func() {
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(savedModelPrices))
	})
	require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(`{"gpt-image-2@4k":0.08}`))

	count := uint(3)
	info := &relaycommon.RelayInfo{
		OriginModelName: dto.GPTImage2Model,
		Request: &dto.ImageRequest{
			Model:   dto.GPTImage2Model,
			Size:    "auto",
			Quality: "high",
			N:       &count,
		},
		PriceData: hosttypes.PriceData{
			UsePrice:   true,
			ModelPrice: 0.08,
		},
	}
	other := map[string]interface{}{}

	appendGPTImage2RequestInfo(other, info)

	assert.Equal(t, "auto", other["image_size"])
	assert.Equal(t, "high", other["image_quality"])
	assert.Equal(t, "4K", other["image_size_tier"])
	assert.Equal(t, count, other["image_count"])
	assert.Equal(t, 0.08, other["image_unit_price"])
}

func TestAppendGPTImage2RequestInfoUsesAutoDefaults(t *testing.T) {
	info := &relaycommon.RelayInfo{
		OriginModelName: dto.GPTImage2Model,
		Request:         &dto.ImageRequest{Model: dto.GPTImage2Model},
	}
	other := map[string]interface{}{}

	appendGPTImage2RequestInfo(other, info)

	assert.Equal(t, "auto", other["image_size"])
	assert.Equal(t, "auto", other["image_quality"])
	assert.Equal(t, "4K", other["image_size_tier"])
	assert.Equal(t, uint(1), other["image_count"])
	assert.NotContains(t, other, "image_unit_price")
}
