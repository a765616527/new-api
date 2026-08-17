package helper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/pkg/billingexpr"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/QuantumNous/new-api/relaykit/types"
	"github.com/QuantumNous/new-api/setting/billing_setting"
	"github.com/QuantumNous/new-api/setting/config"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestModelPriceHelperTieredUsesPreloadedRequestInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	saved := map[string]string{}
	require.NoError(t, config.GlobalConfig.SaveToDB(func(key, value string) error {
		saved[key] = value
		return nil
	}))
	t.Cleanup(func() {
		require.NoError(t, config.GlobalConfig.LoadFromDB(saved))
	})

	require.NoError(t, config.GlobalConfig.LoadFromDB(map[string]string{
		"billing_setting.billing_mode": `{"tiered-test-model":"tiered_expr"}`,
		"billing_setting.billing_expr": `{"tiered-test-model":"param(\"stream\") == true ? tier(\"stream\", p * 3) : tier(\"base\", p * 2)"}`,
	}))

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodPost, "/api/channel/test/1", nil)
	req.Body = nil
	req.ContentLength = 0
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	ctx.Set("group", "default")

	info := &relaycommon.RelayInfo{
		OriginModelName: "tiered-test-model",
		UserGroup:       "default",
		UsingGroup:      "default",
		RequestHeaders:  map[string]string{"Content-Type": "application/json"},
		BillingRequestInput: &billingexpr.RequestInput{
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    []byte(`{"stream":true}`),
		},
	}

	priceData, err := ModelPriceHelper(ctx, info, 1000, &types.TokenCountMeta{
		BillingRatios: map[string]float64{"n": 3},
	})
	require.NoError(t, err)
	require.Equal(t, 1500, priceData.QuotaToPreConsume)
	require.NotNil(t, info.TieredBillingSnapshot)
	require.Equal(t, "stream", info.TieredBillingSnapshot.EstimatedTier)
	require.Equal(t, billing_setting.BillingModeTieredExpr, info.TieredBillingSnapshot.BillingMode)
	require.Equal(t, common.QuotaPerUnit, info.TieredBillingSnapshot.QuotaPerUnit)
}

func TestModelPriceHelperTieredPreConsumeMaxTokensFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	saved := map[string]string{}
	require.NoError(t, config.GlobalConfig.SaveToDB(func(key, value string) error {
		saved[key] = value
		return nil
	}))
	t.Cleanup(func() {
		require.NoError(t, config.GlobalConfig.LoadFromDB(saved))
	})

	require.NoError(t, config.GlobalConfig.LoadFromDB(map[string]string{
		"billing_setting.billing_mode":    `{"tiered-fallback-model":"tiered_expr"}`,
		"billing_setting.billing_expr":    `{"tiered-fallback-model":"tier(\"base\", p * 3 + c * 15)"}`,
		"group_ratio_setting.group_ratio": `{"default":1,"free":0}`,
	}))

	const promptTokens = 1000

	cases := []struct {
		name      string
		group     string
		maxTokens int
		expected  int
	}{
		{
			// max_tokens omitted in a paid group -> fall back to 8192 completion tokens.
			// p*3 + c*15 = 1000*3 + 8192*15 = 125880 -> /1e6 * 500000 = 62940
			name:      "non-free group falls back to 8192 completion tokens",
			group:     "default",
			maxTokens: 0,
			expected:  62940,
		},
		{
			// explicit max_tokens is used verbatim, no fallback.
			// 1000*3 + 100*15 = 4500 -> /1e6 * 500000 = 2250
			name:      "explicit max_tokens is used verbatim",
			group:     "default",
			maxTokens: 100,
			expected:  2250,
		},
		{
			// free group (ratio 0) stays zero; fallback is gated on non-zero group ratio.
			name:      "free group stays zero without fallback",
			group:     "free",
			maxTokens: 0,
			expected:  0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req
			ctx.Set("group", tc.group)

			info := &relaycommon.RelayInfo{
				OriginModelName: "tiered-fallback-model",
				UserGroup:       tc.group,
				UsingGroup:      tc.group,
				RequestHeaders:  map[string]string{"Content-Type": "application/json"},
				BillingRequestInput: &billingexpr.RequestInput{
					Headers: map[string]string{"Content-Type": "application/json"},
					Body:    []byte(`{}`),
				},
			}

			priceData, err := ModelPriceHelper(ctx, info, promptTokens, &types.TokenCountMeta{MaxTokens: tc.maxTokens})
			require.NoError(t, err)
			require.Equal(t, tc.expected, priceData.QuotaToPreConsume)
		})
	}
}

func TestModelPriceHelperTieredRejectsPreConsumeOverflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	saved := map[string]string{}
	require.NoError(t, config.GlobalConfig.SaveToDB(func(key, value string) error {
		saved[key] = value
		return nil
	}))
	t.Cleanup(func() {
		require.NoError(t, config.GlobalConfig.LoadFromDB(saved))
	})

	require.NoError(t, config.GlobalConfig.LoadFromDB(map[string]string{
		"billing_setting.billing_mode":    `{"tiered-overflow-model":"tiered_expr"}`,
		"billing_setting.billing_expr":    `{"tiered-overflow-model":"tier(\"overflow\", p * 1000000000000000)"}`,
		"group_ratio_setting.group_ratio": `{"default":1}`,
	}))

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	ctx.Set("group", "default")
	info := &relaycommon.RelayInfo{
		OriginModelName: "tiered-overflow-model",
		UserGroup:       "default",
		UsingGroup:      "default",
		BillingRequestInput: &billingexpr.RequestInput{
			Body: []byte(`{}`),
		},
	}

	_, err := ModelPriceHelper(ctx, info, 1000, &types.TokenCountMeta{})

	var clamp *common.QuotaClamp
	require.ErrorAs(t, err, &clamp)
	require.Equal(t, "QuotaRound", clamp.Op)
	require.Equal(t, common.QuotaClampOverflow, clamp.Kind)
}

func TestModelPriceHelperRequestBillingRatiosOnlyApplyToFixedPrice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	savedModelPrices := ratio_setting.ModelPrice2JSONString()
	savedModelRatios := ratio_setting.ModelRatio2JSONString()
	t.Cleanup(func() {
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(savedModelPrices))
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(savedModelRatios))
	})

	modelPrices, err := common.Marshal(map[string]float64{
		"fixed-image-price":      0.04,
		"fractional-image-price": 0.0000012,
		"overflow-image-price":   float64(common.MaxQuota) / common.QuotaPerUnit / 2,
	})
	require.NoError(t, err)
	require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(string(modelPrices)))
	modelRatios, err := common.Marshal(map[string]float64{"ratio-image-price": 15})
	require.NoError(t, err)
	require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(string(modelRatios)))

	tests := []struct {
		name           string
		model          string
		wantQuota      int
		wantUsePrice   bool
		wantImageCount bool
	}{
		{
			name:           "fixed price applies image count",
			model:          "fixed-image-price",
			wantQuota:      180000,
			wantUsePrice:   true,
			wantImageCount: true,
		},
		{
			name:         "ratio price ignores request billing ratios",
			model:        "ratio-image-price",
			wantQuota:    15000,
			wantUsePrice: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Set("group", "default")
			info := &relaycommon.RelayInfo{
				OriginModelName: tt.model,
				UserGroup:       "default",
				UsingGroup:      "default",
			}
			meta := &types.TokenCountMeta{
				ImagePriceRatio: 3,
				BillingRatios:   map[string]float64{"n": 3},
			}

			priceData, err := ModelPriceHelper(ctx, info, 1000, meta)

			require.NoError(t, err)
			require.Equal(t, tt.wantQuota, priceData.QuotaToPreConsume)
			require.Equal(t, tt.wantUsePrice, priceData.UsePrice)
			require.Equal(t, tt.wantImageCount, priceData.HasOtherRatio("n"))
			require.Equal(t, priceData.OtherRatios(), info.PriceData.OtherRatios())
		})
	}

	newInfo := func(model string) (*gin.Context, *relaycommon.RelayInfo) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Set("group", "default")
		return ctx, &relaycommon.RelayInfo{
			OriginModelName: model,
			UserGroup:       "default",
			UsingGroup:      "default",
		}
	}
	meta := &types.TokenCountMeta{BillingRatios: map[string]float64{"n": 3}}

	ctx, info := newInfo("fractional-image-price")
	priceData, err := ModelPriceHelper(ctx, info, 0, meta)
	require.NoError(t, err)
	// 0.0000012 * 500000 * 3 = 1.8, then truncate once to 1.
	require.Equal(t, 1, priceData.QuotaToPreConsume)

	ctx, info = newInfo("overflow-image-price")
	_, err = ModelPriceHelper(ctx, info, 0, meta)
	var clamp *common.QuotaClamp
	require.ErrorAs(t, err, &clamp)
	require.Equal(t, "QuotaFromFloat", clamp.Op)
	require.Equal(t, common.QuotaClampOverflow, clamp.Kind)
	require.Nil(t, info.Billing)
}

func TestModelPriceHelperUsesGPTImage2TierPrice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	savedModelPrices := ratio_setting.ModelPrice2JSONString()
	t.Cleanup(func() {
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(savedModelPrices))
	})

	prices, err := common.Marshal(map[string]float64{
		dto.GPTImage2Model: 0.03,
		ratio_setting.GPTImage2TierPriceKey(dto.ImageSizeTier1K): 0.01,
		ratio_setting.GPTImage2TierPriceKey(dto.ImageSizeTier4K): 0.04,
	})
	require.NoError(t, err)
	require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(string(prices)))

	tests := []struct {
		name          string
		size          string
		expectedPrice float64
		expectedQuota int
	}{
		{name: "1k tier", size: "1024x768", expectedPrice: 0.01, expectedQuota: 10000},
		{name: "missing 2k tier falls back to base", size: "1536x1024", expectedPrice: 0.03, expectedQuota: 30000},
		{name: "auto uses 4k tier", size: "auto", expectedPrice: 0.04, expectedQuota: 40000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Set("group", "default")
			count := uint(2)
			request := &dto.ImageRequest{Model: dto.GPTImage2Model, Size: tt.size, Quality: "high", N: &count}
			info := &relaycommon.RelayInfo{
				OriginModelName: dto.GPTImage2Model,
				UserGroup:       "default",
				UsingGroup:      "default",
				Request:         request,
			}

			priceData, err := ModelPriceHelper(ctx, info, 0, request.GetTokenCountMeta())

			require.NoError(t, err)
			require.True(t, priceData.UsePrice)
			require.Equal(t, tt.expectedPrice, priceData.ModelPrice)
			require.Equal(t, tt.expectedQuota, priceData.QuotaToPreConsume)
		})
	}
}

func TestModelPriceHelperRejectsIncompleteGPTImage2TierPricesWithoutFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	savedModelPrices := ratio_setting.ModelPrice2JSONString()
	t.Cleanup(func() {
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(savedModelPrices))
	})
	require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(`{"gpt-image-2@1k":0.01}`))

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set("group", "default")
	request := &dto.ImageRequest{Model: dto.GPTImage2Model, Size: "1536x1024"}
	info := &relaycommon.RelayInfo{
		OriginModelName: dto.GPTImage2Model,
		UserGroup:       "default",
		UsingGroup:      "default",
		Request:         request,
	}

	_, err := ModelPriceHelper(ctx, info, 0, request.GetTokenCountMeta())

	require.EqualError(t, err, "gpt-image-2 2K price is not configured and no fixed fallback price is available")
}
