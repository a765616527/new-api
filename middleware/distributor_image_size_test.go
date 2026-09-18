package middleware

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetModelRequestClassifiesGPTImage2GenerationSize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", bytes.NewBufferString(`{"model":"gpt-image-2","size":"1536x1024"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	request, shouldSelect, err := getModelRequest(c)

	require.NoError(t, err)
	assert.True(t, shouldSelect)
	assert.Equal(t, dto.GPTImage2Model, request.Model)
	assert.Equal(t, dto.ImageSizeTier1K, request.Size)
}

func TestGetModelRequestRoutesAutoImageEditTo4K(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	require.NoError(t, writer.WriteField("model", dto.GPTImage2Model))
	require.NoError(t, writer.WriteField("size", "auto"))
	require.NoError(t, writer.Close())

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/edits", &body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	request, shouldSelect, err := getModelRequest(c)

	require.NoError(t, err)
	assert.True(t, shouldSelect)
	assert.Equal(t, dto.ImageSizeTier4K, request.Size)
}

func TestGetModelRequestClassifiesGPTImage25FlareSize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", bytes.NewBufferString(`{"model":"gpt-image-2.5-flare","size":"2048x1152"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	request, shouldSelect, err := getModelRequest(c)

	require.NoError(t, err)
	assert.True(t, shouldSelect)
	assert.Equal(t, dto.GPTImage25FlareModel, request.Model)
	assert.Equal(t, dto.ImageSizeTier2K, request.Size)
}

func TestSetupContextOmitsSizeMappedModelWhenRoutingUnset(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	channel := &model.Channel{Key: "sk-test", OtherSettings: `{}`}

	_ = SetupContextForSelectedChannel(c, channel, dto.GPTImage25FlareModel)

	assert.Empty(t, common.GetContextKeyString(c, constant.ContextKeyChannelSizeMappedModel))
}

func TestSetupContextSetsSizeMappedModelWhenRoutingConfigured(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	common.SetContextKey(c, constant.ContextKeyImageSizeTier, dto.ImageSizeTier1K)
	channel := &model.Channel{
		Key:           "sk-test",
		OtherSettings: `{"gpt_image_2_5_flare_size_models":{"1k":"flare-1k"}}`,
	}

	_ = SetupContextForSelectedChannel(c, channel, dto.GPTImage25FlareModel)

	assert.Equal(t, "flare-1k", common.GetContextKeyString(c, constant.ContextKeyChannelSizeMappedModel))
}
