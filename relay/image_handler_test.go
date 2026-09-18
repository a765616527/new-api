package relay

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMappedImagePassThroughJSONBodyRewritesOnlyModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(
		http.MethodPost,
		"/v1/images/generations",
		bytes.NewBufferString(`{"model":"gpt-image-2","size":"1536x1024","custom":{"keep":true}}`),
	)
	c.Request.Header.Set("Content-Type", "application/json")

	data, err := mappedImagePassThroughJSONBody(c, "vendor-image-2k")
	require.NoError(t, err)

	var decoded map[string]interface{}
	require.NoError(t, common.Unmarshal(data, &decoded))
	assert.Equal(t, "vendor-image-2k", decoded["model"])
	assert.Equal(t, "1536x1024", decoded["size"])
	assert.Equal(t, map[string]interface{}{"keep": true}, decoded["custom"])
}
