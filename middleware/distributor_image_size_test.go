package middleware

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

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
