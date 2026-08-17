package helper

import (
	"testing"

	rootcommon "github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelMappedHelperUsesImageSizeModelBeforeGeneralMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	rootcommon.SetContextKey(c, constant.ContextKeyChannelSizeMappedModel, "vendor-image-2k")
	rootcommon.SetContextKey(c, constant.ContextKeyChannelModelMapping, `{"gpt-image-2":"generic-image"}`)
	info := &relaycommon.RelayInfo{
		OriginModelName: dto.GPTImage2Model,
		ChannelMeta: &relaycommon.ChannelMeta{
			UpstreamModelName: dto.GPTImage2Model,
		},
	}
	request := &dto.ImageRequest{Model: dto.GPTImage2Model}

	err := ModelMappedHelper(c, info, request)

	require.NoError(t, err)
	assert.Equal(t, "vendor-image-2k", info.UpstreamModelName)
	assert.Equal(t, "vendor-image-2k", request.Model)
	assert.True(t, info.IsModelMapped)
}
