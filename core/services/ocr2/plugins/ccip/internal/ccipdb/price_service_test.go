package db

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	cciporm "github.com/smartcontractkit/chainlink/v2/core/services/ccip"
	ccipmocks "github.com/smartcontractkit/chainlink/v2/core/services/ccip/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

func TestPriceService_priceCleanup(t *testing.T) {
	lggr := logger.TestLogger(t)
	jobId := int32(1)
	destChainSelector := uint64(12345)
	sourceChainSelector := uint64(67890)

	testCases := []struct {
		name            string
		gasPriceError   bool
		tokenPriceError bool
		expectedErr     bool
	}{
		{
			name:            "ORM called successfully",
			gasPriceError:   false,
			tokenPriceError: false,
			expectedErr:     false,
		},
		{
			name:            "gasPrice clear failed",
			gasPriceError:   true,
			tokenPriceError: false,
			expectedErr:     true,
		},
		{
			name:            "tokenPrice clear failed",
			gasPriceError:   false,
			tokenPriceError: true,
			expectedErr:     true,
		},
		{
			name:            "both ORM calls failed",
			gasPriceError:   true,
			tokenPriceError: true,
			expectedErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tests.Context(t)

			var gasPricesError error
			var tokenPricesError error
			if tc.gasPriceError {
				gasPricesError = fmt.Errorf("gas prices error")
			}
			if tc.tokenPriceError {
				tokenPricesError = fmt.Errorf("token prices error")
			}

			mockOrm := ccipmocks.NewORM(t)
			mockOrm.On("ClearGasPricesByDestChain", ctx, destChainSelector, int(priceExpireThreshold.Seconds())).Return(gasPricesError).Once()
			mockOrm.On("ClearTokenPricesByDestChain", ctx, destChainSelector, int(priceExpireThreshold.Seconds())).Return(tokenPricesError).Once()

			priceService := NewPriceService(
				lggr,
				mockOrm,
				jobId,
				destChainSelector,
				sourceChainSelector,
				"",
				nil,
				nil,
			).(*priceService)
			err := priceService.runCleanup(ctx)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPriceService_writeGasPrices(t *testing.T) {
	lggr := logger.TestLogger(t)
	jobId := int32(1)
	destChainSelector := uint64(12345)
	sourceChainSelector := uint64(67890)

	gasPrice := big.NewInt(1e18)

	expectedGasPriceUpdate := []cciporm.GasPriceUpdate{
		{
			SourceChainSelector: sourceChainSelector,
			GasPrice:            assets.NewWei(gasPrice),
		},
	}

	testCases := []struct {
		name          string
		gasPriceError bool
		expectedErr   bool
	}{
		{
			name:          "ORM called successfully",
			gasPriceError: false,
			expectedErr:   false,
		},
		{
			name:          "gasPrice clear failed",
			gasPriceError: true,
			expectedErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tests.Context(t)

			var gasPricesError error
			if tc.gasPriceError {
				gasPricesError = fmt.Errorf("gas prices error")
			}

			mockOrm := ccipmocks.NewORM(t)
			mockOrm.On("InsertGasPricesForDestChain", ctx, destChainSelector, jobId, expectedGasPriceUpdate).Return(gasPricesError).Once()

			priceService := NewPriceService(
				lggr,
				mockOrm,
				jobId,
				destChainSelector,
				sourceChainSelector,
				"",
				nil,
				nil,
			).(*priceService)
			err := priceService.writeGasPricesToDB(ctx, gasPrice)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPriceService_writeTokenPrices(t *testing.T) {
	lggr := logger.TestLogger(t)
	jobId := int32(1)
	destChainSelector := uint64(12345)
	sourceChainSelector := uint64(67890)

	tokenPrices := map[cciptypes.Address]*big.Int{
		"0x123": big.NewInt(2e18),
		"0x234": big.NewInt(3e18),
	}

	expectedTokenPriceUpdate := []cciporm.TokenPriceUpdate{
		{
			TokenAddr:  "0x123",
			TokenPrice: assets.NewWei(big.NewInt(2e18)),
		},
		{
			TokenAddr:  "0x234",
			TokenPrice: assets.NewWei(big.NewInt(3e18)),
		},
	}

	testCases := []struct {
		name            string
		tokenPriceError bool
		expectedErr     bool
	}{
		{
			name:            "ORM called successfully",
			tokenPriceError: false,
			expectedErr:     false,
		},
		{
			name:            "tokenPrice clear failed",
			tokenPriceError: true,
			expectedErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tests.Context(t)

			var tokenPricesError error
			if tc.tokenPriceError {
				tokenPricesError = fmt.Errorf("token prices error")
			}

			mockOrm := ccipmocks.NewORM(t)
			mockOrm.On("InsertTokenPricesForDestChain", ctx, destChainSelector, jobId, expectedTokenPriceUpdate).Return(tokenPricesError).Once()

			priceService := NewPriceService(
				lggr,
				mockOrm,
				jobId,
				destChainSelector,
				sourceChainSelector,
				"",
				nil,
				nil,
			).(*priceService)
			err := priceService.writeTokenPricesToDB(ctx, tokenPrices)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPriceService_observeGasPriceUpdates(t *testing.T) {
	lggr := logger.TestLogger(t)
	jobId := int32(1)
	destChainSelector := uint64(12345)
	sourceChainSelector := uint64(67890)
	sourceNativeToken := cciptypes.Address(utils.RandomAddress().String())

	testCases := []struct {
		name                 string
		sourceNativeToken    cciptypes.Address
		priceGetterRespData  map[cciptypes.Address]*big.Int
		priceGetterRespErr   error
		feeEstimatorRespFee  *big.Int
		feeEstimatorRespErr  error
		maxGasPrice          uint64
		expSourceGasPriceUSD *big.Int
		expErr               bool
	}{
		{
			name:              "base",
			sourceNativeToken: sourceNativeToken,
			priceGetterRespData: map[cciptypes.Address]*big.Int{
				sourceNativeToken: val1e18(100),
			},
			priceGetterRespErr:   nil,
			feeEstimatorRespFee:  big.NewInt(10),
			feeEstimatorRespErr:  nil,
			maxGasPrice:          1e18,
			expSourceGasPriceUSD: big.NewInt(1000),
			expErr:               false,
		},
		{
			name:                "price getter returned an error",
			sourceNativeToken:   sourceNativeToken,
			priceGetterRespData: nil,
			priceGetterRespErr:  fmt.Errorf("some random network error"),
			expErr:              true,
		},
		{
			name:              "price getter did not return source native gas price",
			sourceNativeToken: sourceNativeToken,
			priceGetterRespData: map[cciptypes.Address]*big.Int{
				"0x1": val1e18(100),
			},
			priceGetterRespErr: nil,
			expErr:             true,
		},
		{
			name:              "dynamic fee cap overrides legacy",
			sourceNativeToken: sourceNativeToken,
			priceGetterRespData: map[cciptypes.Address]*big.Int{
				sourceNativeToken: val1e18(100),
			},
			priceGetterRespErr:   nil,
			feeEstimatorRespFee:  big.NewInt(20),
			feeEstimatorRespErr:  nil,
			maxGasPrice:          1e18,
			expSourceGasPriceUSD: big.NewInt(2000),
			expErr:               false,
		},
		{
			name:              "nil gas price",
			sourceNativeToken: sourceNativeToken,
			priceGetterRespData: map[cciptypes.Address]*big.Int{
				sourceNativeToken: val1e18(100),
			},
			feeEstimatorRespFee: nil,
			maxGasPrice:         1e18,
			expErr:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			priceGetter := pricegetter.NewMockPriceGetter(t)
			defer priceGetter.AssertExpectations(t)

			gasPriceEstimator := prices.NewMockGasPriceEstimatorCommit(t)
			defer gasPriceEstimator.AssertExpectations(t)

			priceGetter.On("TokenPricesUSD", mock.Anything, []cciptypes.Address{tc.sourceNativeToken}).Return(tc.priceGetterRespData, tc.priceGetterRespErr)

			if tc.maxGasPrice > 0 {
				gasPriceEstimator.On("GetGasPrice", mock.Anything).Return(tc.feeEstimatorRespFee, tc.feeEstimatorRespErr)
				if tc.feeEstimatorRespFee != nil {
					pUSD := ccipcalc.CalculateUsdPerUnitGas(tc.feeEstimatorRespFee, tc.priceGetterRespData[tc.sourceNativeToken])
					gasPriceEstimator.On("DenoteInUSD", mock.Anything, mock.Anything).Return(pUSD, nil)
				}
			}

			priceService := NewPriceService(
				lggr,
				nil,
				jobId,
				destChainSelector,
				sourceChainSelector,
				tc.sourceNativeToken,
				priceGetter,
				nil,
			).(*priceService)
			priceService.gasPriceEstimator = gasPriceEstimator

			sourceGasPriceUSD, err := priceService.observeGasPriceUpdates(context.Background(), lggr)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.True(t, tc.expSourceGasPriceUSD.Cmp(sourceGasPriceUSD) == 0)
		})
	}
}

func TestPriceService_observeTokenPriceUpdates(t *testing.T) {
	lggr := logger.TestLogger(t)
	jobId := int32(1)
	destChainSelector := uint64(12345)
	sourceChainSelector := uint64(67890)

	const nTokens = 10
	tokens := make([]cciptypes.Address, nTokens)
	for i := range tokens {
		tokens[i] = cciptypes.Address(utils.RandomAddress().String())
	}
	sort.Slice(tokens, func(i, j int) bool { return tokens[i] < tokens[j] })

	testCases := []struct {
		name                string
		tokenDecimals       map[cciptypes.Address]uint8
		filterOutTokens     []cciptypes.Address
		priceGetterRespData map[cciptypes.Address]*big.Int
		priceGetterRespErr  error
		expTokenPricesUSD   map[cciptypes.Address]*big.Int
		expErr              bool
	}{
		{
			name: "base",
			tokenDecimals: map[cciptypes.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 12,
			},
			filterOutTokens: []cciptypes.Address{tokens[2]},
			priceGetterRespData: map[cciptypes.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
				tokens[2]: val1e18(300), // price getter returned a price for this token even though we didn't request it (should be skipped)
			},
			priceGetterRespErr: nil,
			expTokenPricesUSD: map[cciptypes.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200 * 1e6),
			},
			expErr: false,
		},
		{
			name: "price getter returned an error",
			tokenDecimals: map[cciptypes.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			priceGetterRespData: nil,
			priceGetterRespErr:  fmt.Errorf("some random network error"),
			expErr:              true,
		},
		{
			name: "price getter skipped a requested price",
			tokenDecimals: map[cciptypes.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			priceGetterRespData: map[cciptypes.Address]*big.Int{
				tokens[0]: val1e18(100),
			},
			priceGetterRespErr: nil,
			expErr:             true,
		},
		{
			name: "nil token price",
			tokenDecimals: map[cciptypes.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			filterOutTokens: []cciptypes.Address{tokens[2]},
			priceGetterRespData: map[cciptypes.Address]*big.Int{
				tokens[0]: nil,
				tokens[1]: val1e18(200),
				tokens[2]: val1e18(300),
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			priceGetter := pricegetter.NewMockPriceGetter(t)
			defer priceGetter.AssertExpectations(t)

			var destTokens []cciptypes.Address
			for tk := range tc.tokenDecimals {
				destTokens = append(destTokens, tk)
			}
			sort.Slice(destTokens, func(i, j int) bool {
				return destTokens[i] < destTokens[j]
			})
			var destDecimals []uint8
			for _, token := range destTokens {
				destDecimals = append(destDecimals, tc.tokenDecimals[token])
			}

			queryTokens := ccipcommon.FlattenUniqueSlice(destTokens)

			if len(queryTokens) > 0 {
				priceGetter.On("TokenPricesUSD", mock.Anything, queryTokens).Return(tc.priceGetterRespData, tc.priceGetterRespErr)
				priceGetter.On("FilterConfiguredTokens", mock.Anything, mock.Anything).Return(destTokens, tc.filterOutTokens, nil)
			}

			offRampReader := ccipdatamocks.NewOffRampReader(t)
			offRampReader.On("GetTokens", mock.Anything).Return(cciptypes.OffRampTokens{
				DestinationTokens: destTokens,
			}, nil).Maybe()

			destPriceReg := ccipdatamocks.NewPriceRegistryReader(t)
			destPriceReg.On("GetTokensDecimals", mock.Anything, destTokens).Return(destDecimals, nil).Maybe()
			destPriceReg.On("GetFeeTokens", mock.Anything).Return([]cciptypes.Address{destTokens[0]}, nil).Maybe()

			priceService := NewPriceService(
				lggr,
				nil,
				jobId,
				destChainSelector,
				sourceChainSelector,
				"0x123",
				priceGetter,
				offRampReader,
			).(*priceService)
			priceService.destPriceRegistryReader = destPriceReg

			tokenPricesUSD, err := priceService.observeTokenPriceUpdates(context.Background(), lggr)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.True(t, reflect.DeepEqual(tc.expTokenPricesUSD, tokenPricesUSD))
		})
	}
}

func TestPriceService_calculateUsdPer1e18TokenAmount(t *testing.T) {
	testCases := []struct {
		name       string
		price      *big.Int
		decimal    uint8
		wantResult *big.Int
	}{
		{
			name:       "18-decimal token, $6.5 per token",
			price:      big.NewInt(65e17),
			decimal:    18,
			wantResult: big.NewInt(65e17),
		},
		{
			name:       "6-decimal token, $1 per token",
			price:      big.NewInt(1e18),
			decimal:    6,
			wantResult: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e12)), // 1e30
		},
		{
			name:       "0-decimal token, $1 per token",
			price:      big.NewInt(1e18),
			decimal:    0,
			wantResult: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e18)), // 1e36
		},
		{
			name:       "36-decimal token, $1 per token",
			price:      big.NewInt(1e18),
			decimal:    36,
			wantResult: big.NewInt(1),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateUsdPer1e18TokenAmount(tt.price, tt.decimal)
			assert.Equal(t, tt.wantResult, got)
		})
	}
}

func TestPriceService_GetGasAndTokenPrices(t *testing.T) {
	lggr := logger.TestLogger(t)
	jobId := int32(1)
	destChainSelector := uint64(12345)
	sourceChainSelector := uint64(67890)

	token1 := ccipcalc.HexToAddress("0x123")
	token2 := ccipcalc.HexToAddress("0x234")

	gasPrice := big.NewInt(1e18)
	tokenPrices := map[cciptypes.Address]*big.Int{
		token1: big.NewInt(2e18),
		token2: big.NewInt(3e18),
	}

	testCases := []struct {
		name                 string
		ormGasPricesResult   []cciporm.GasPrice
		ormTokenPricesResult []cciporm.TokenPrice

		expectedGasPrices   map[uint64]*big.Int
		expectedTokenPrices map[cciptypes.Address]*big.Int

		gasPriceError   bool
		tokenPriceError bool
		expectedErr     bool
	}{
		{
			name: "ORM called successfully",
			ormGasPricesResult: []cciporm.GasPrice{
				{
					SourceChainSelector: sourceChainSelector,
					GasPrice:            assets.NewWei(gasPrice),
				},
			},
			ormTokenPricesResult: []cciporm.TokenPrice{
				{
					TokenAddr:  string(token1),
					TokenPrice: assets.NewWei(tokenPrices[token1]),
				},
				{
					TokenAddr:  string(token2),
					TokenPrice: assets.NewWei(tokenPrices[token2]),
				},
			},
			expectedGasPrices: map[uint64]*big.Int{
				sourceChainSelector: gasPrice,
			},
			expectedTokenPrices: tokenPrices,
			gasPriceError:       false,
			tokenPriceError:     false,
			expectedErr:         false,
		},
		{
			name: "multiple gas prices with nil token price",
			ormGasPricesResult: []cciporm.GasPrice{
				{
					SourceChainSelector: sourceChainSelector,
					GasPrice:            assets.NewWei(gasPrice),
				},
				{
					SourceChainSelector: sourceChainSelector + 1,
					GasPrice:            assets.NewWei(big.NewInt(200)),
				},
				{
					SourceChainSelector: sourceChainSelector + 2,
					GasPrice:            assets.NewWei(big.NewInt(300)),
				},
			},
			ormTokenPricesResult: nil,
			expectedGasPrices: map[uint64]*big.Int{
				sourceChainSelector:     gasPrice,
				sourceChainSelector + 1: big.NewInt(200),
				sourceChainSelector + 2: big.NewInt(300),
			},
			expectedTokenPrices: map[cciptypes.Address]*big.Int{},
			gasPriceError:       false,
			tokenPriceError:     false,
			expectedErr:         false,
		},
		{
			name:               "multiple token prices with nil gas price",
			ormGasPricesResult: nil,
			ormTokenPricesResult: []cciporm.TokenPrice{
				{
					TokenAddr:  string(token1),
					TokenPrice: assets.NewWei(tokenPrices[token1]),
				},
				{
					TokenAddr:  string(token2),
					TokenPrice: assets.NewWei(tokenPrices[token2]),
				},
			},
			expectedGasPrices:   map[uint64]*big.Int{},
			expectedTokenPrices: tokenPrices,
			gasPriceError:       false,
			tokenPriceError:     false,
			expectedErr:         false,
		},
		{
			name: "nil prices filtered out",
			ormGasPricesResult: []cciporm.GasPrice{
				{
					SourceChainSelector: sourceChainSelector,
					GasPrice:            nil,
				},
				{
					SourceChainSelector: sourceChainSelector + 1,
					GasPrice:            assets.NewWei(gasPrice),
				},
			},
			ormTokenPricesResult: []cciporm.TokenPrice{
				{
					TokenAddr:  string(token1),
					TokenPrice: assets.NewWei(tokenPrices[token1]),
				},
				{
					TokenAddr:  string(token2),
					TokenPrice: nil,
				},
			},
			expectedGasPrices: map[uint64]*big.Int{
				sourceChainSelector + 1: gasPrice,
			},
			expectedTokenPrices: map[cciptypes.Address]*big.Int{
				token1: tokenPrices[token1],
			},
			gasPriceError:   false,
			tokenPriceError: false,
			expectedErr:     false,
		},
		{
			name:            "gasPrice clear failed",
			gasPriceError:   true,
			tokenPriceError: false,
			expectedErr:     true,
		},
		{
			name:            "tokenPrice clear failed",
			gasPriceError:   false,
			tokenPriceError: true,
			expectedErr:     true,
		},
		{
			name:            "both ORM calls failed",
			gasPriceError:   true,
			tokenPriceError: true,
			expectedErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tests.Context(t)

			mockOrm := ccipmocks.NewORM(t)
			if tc.gasPriceError {
				mockOrm.On("GetGasPricesByDestChain", ctx, destChainSelector).Return(nil, fmt.Errorf("gas prices error")).Once()
			} else {
				mockOrm.On("GetGasPricesByDestChain", ctx, destChainSelector).Return(tc.ormGasPricesResult, nil).Once()
			}
			if tc.tokenPriceError {
				mockOrm.On("GetTokenPricesByDestChain", ctx, destChainSelector).Return(nil, fmt.Errorf("token prices error")).Once()
			} else {
				mockOrm.On("GetTokenPricesByDestChain", ctx, destChainSelector).Return(tc.ormTokenPricesResult, nil).Once()
			}

			priceService := NewPriceService(
				lggr,
				mockOrm,
				jobId,
				destChainSelector,
				sourceChainSelector,
				"",
				nil,
				nil,
			).(*priceService)
			gasPricesResult, tokenPricesResult, err := priceService.GetGasAndTokenPrices(ctx, destChainSelector)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedGasPrices, gasPricesResult)
				assert.Equal(t, tc.expectedTokenPrices, tokenPricesResult)
			}
		})
	}
}

func val1e18(val int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(val))
}

func setupORM(t *testing.T) cciporm.ORM {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	orm, err := cciporm.NewORM(db)

	require.NoError(t, err)

	return orm
}

func checkResultLen(t *testing.T, priceService PriceService, destChainSelector uint64, gasCount int, tokenCount int) error {
	ctx := tests.Context(t)
	dbGasResult, dbTokenResult, err := priceService.GetGasAndTokenPrices(ctx, destChainSelector)
	if err != nil {
		return nil
	}
	if len(dbGasResult) != gasCount {
		return fmt.Errorf("expected %d gas prices, got %d", gasCount, len(dbGasResult))
	}
	if len(dbTokenResult) != tokenCount {
		return fmt.Errorf("expected %d token prices, got %d", tokenCount, len(dbTokenResult))
	}
	return nil
}

func TestPriceService_priceWriteAndCleanupInBackground(t *testing.T) {
	lggr := logger.TestLogger(t)
	jobId := int32(1)
	destChainSelector := uint64(12345)
	sourceChainSelector := uint64(67890)
	ctx := tests.Context(t)

	sourceNative := cciptypes.Address("0x123")
	feeTokens := []cciptypes.Address{"0x234"}
	rampTokens := []cciptypes.Address{"0x345", "0x456"}
	rampFilteredTokens := []cciptypes.Address{"0x345"}
	rampFilterOutTokens := []cciptypes.Address{"0x456"}

	laneTokens := []cciptypes.Address{"0x234", "0x345"}
	laneTokenDecimals := []uint8{18, 18}

	tokens := []cciptypes.Address{sourceNative, "0x234", "0x345"}
	tokenPrices := []int64{2, 3, 4}
	gasPrice := big.NewInt(10)

	orm := setupORM(t)

	priceGetter := pricegetter.NewMockPriceGetter(t)
	defer priceGetter.AssertExpectations(t)

	gasPriceEstimator := prices.NewMockGasPriceEstimatorCommit(t)
	defer gasPriceEstimator.AssertExpectations(t)

	priceGetter.On("TokenPricesUSD", mock.Anything, tokens[:1]).Return(map[cciptypes.Address]*big.Int{
		tokens[0]: val1e18(tokenPrices[0]),
	}, nil)
	priceGetter.On("TokenPricesUSD", mock.Anything, tokens[1:]).Return(map[cciptypes.Address]*big.Int{
		tokens[1]: val1e18(tokenPrices[1]),
		tokens[2]: val1e18(tokenPrices[2]),
	}, nil)
	priceGetter.On("FilterConfiguredTokens", mock.Anything, rampTokens).Return(rampFilteredTokens, rampFilterOutTokens, nil)

	offRampReader := ccipdatamocks.NewOffRampReader(t)
	offRampReader.On("GetTokens", mock.Anything).Return(cciptypes.OffRampTokens{
		DestinationTokens: rampTokens,
	}, nil).Maybe()

	gasPriceEstimator.On("GetGasPrice", mock.Anything).Return(gasPrice, nil)
	pUSD := ccipcalc.CalculateUsdPerUnitGas(gasPrice, val1e18(tokenPrices[0]))
	gasPriceEstimator.On("DenoteInUSD", mock.Anything, mock.Anything).Return(pUSD, nil)

	destPriceReg := ccipdatamocks.NewPriceRegistryReader(t)
	destPriceReg.On("GetTokensDecimals", mock.Anything, laneTokens).Return(laneTokenDecimals, nil).Maybe()
	destPriceReg.On("GetFeeTokens", mock.Anything).Return(feeTokens, nil).Maybe()

	priceService := NewPriceService(
		lggr,
		orm,
		jobId,
		destChainSelector,
		sourceChainSelector,
		tokens[0],
		priceGetter,
		offRampReader,
	).(*priceService)

	gasUpdateInterval := 2000 * time.Millisecond
	tokenUpdateInterval := 5000 * time.Millisecond
	cleanupInterval := 3000 * time.Millisecond

	// run gas price task every 2 second
	priceService.gasUpdateInterval = gasUpdateInterval
	// run token price task every 5 second
	priceService.tokenUpdateInterval = tokenUpdateInterval
	// run cleanup every 3 seconds
	priceService.cleanupInterval = cleanupInterval
	// expire all prices during every cleanup
	priceService.priceExpireThreshold = time.Duration(0)

	// initially, db is empty
	assert.NoError(t, checkResultLen(t, priceService, destChainSelector, 0, 0))

	// starts PriceService in the background
	assert.NoError(t, priceService.Start(ctx))

	// setting dynamicConfig triggers initial price update
	err := priceService.UpdateDynamicConfig(ctx, gasPriceEstimator, destPriceReg)
	assert.NoError(t, err)
	assert.NoError(t, checkResultLen(t, priceService, destChainSelector, 1, len(laneTokens)))

	// eventually prices will be cleaned
	assert.Eventually(t, func() bool {
		err := checkResultLen(t, priceService, destChainSelector, 0, 0)
		return err == nil
	}, testutils.WaitTimeout(t), testutils.TestInterval)

	// then prices will be updated again
	assert.Eventually(t, func() bool {
		err := checkResultLen(t, priceService, destChainSelector, 1, len(laneTokens))
		return err == nil
	}, testutils.WaitTimeout(t), testutils.TestInterval)

	assert.NoError(t, priceService.Close())
	assert.NoError(t, priceService.runCleanup(ctx))

	// after stopping PriceService and runCleanup, no more updates are inserted
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		assert.NoError(t, checkResultLen(t, priceService, destChainSelector, 0, 0))
	}
}