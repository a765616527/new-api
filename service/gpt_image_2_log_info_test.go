package service

import (
	"testing"

	"github.com/QuantumNous/new-api/model"
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
	other := model.NewLogOther()

	appendGPTImage2RequestInfo(other, info)
	snapshot := other.Snapshot()

	assert.Equal(t, "auto", snapshot["image_size"])
	assert.Equal(t, "high", snapshot["image_quality"])
	assert.Equal(t, "4K", snapshot["image_size_tier"])
	assert.Equal(t, count, snapshot["image_request_count"])
	assert.Equal(t, count, snapshot["image_count"])
	assert.Equal(t, 0.08, snapshot["image_unit_price"])
}

func TestAppendGPTImage2RequestInfoUsesAutoDefaults(t *testing.T) {
	info := &relaycommon.RelayInfo{
		OriginModelName: dto.GPTImage2Model,
		Request:         &dto.ImageRequest{Model: dto.GPTImage2Model},
	}
	other := model.NewLogOther()

	appendGPTImage2RequestInfo(other, info)
	snapshot := other.Snapshot()

	assert.Equal(t, "auto", snapshot["image_size"])
	assert.Equal(t, "auto", snapshot["image_quality"])
	assert.Equal(t, "4K", snapshot["image_size_tier"])
	assert.Equal(t, uint(1), snapshot["image_request_count"])
	assert.NotContains(t, snapshot, "image_count")
	assert.NotContains(t, snapshot, "image_unit_price")
}

func TestAppendGPTImage2RequestInfoRecordsFlareFields(t *testing.T) {
	info := &relaycommon.RelayInfo{
		OriginModelName: dto.GPTImage25FlareModel,
		Request:         &dto.ImageRequest{Model: dto.GPTImage25FlareModel, Size: "3840x2160"},
		PriceData: hosttypes.PriceData{
			UsePrice:   true,
			ModelPrice: 0.11,
		},
	}
	other := model.NewLogOther()

	appendGPTImage2RequestInfo(other, info)
	snapshot := other.Snapshot()

	assert.Equal(t, "3840x2160", snapshot["image_size"])
	assert.Equal(t, "4K", snapshot["image_size_tier"])
	assert.Equal(t, uint(1), snapshot["image_request_count"])
	assert.Equal(t, uint(1), snapshot["image_count"])
	assert.Equal(t, 0.11, snapshot["image_unit_price"])
}

func TestAppendGPTImage2RequestInfoDoesNotOverwriteSettledImageCount(t *testing.T) {
	count := uint(2)
	info := &relaycommon.RelayInfo{
		OriginModelName: dto.GPTImage2Model,
		Request:         &dto.ImageRequest{Model: dto.GPTImage2Model, N: &count},
		PriceData:       hosttypes.PriceData{UsePrice: true, ModelPrice: 0.04},
	}
	other := model.NewLogOther()
	other.SetPublic("image_count", 9)

	appendGPTImage2RequestInfo(other, info)
	snapshot := other.Snapshot()

	assert.Equal(t, count, snapshot["image_request_count"])
	assert.Equal(t, 9, snapshot["image_count"])
}
