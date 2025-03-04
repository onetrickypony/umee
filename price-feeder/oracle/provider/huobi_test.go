package provider

import (
	"context"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/umee-network/umee/price-feeder/oracle/types"
)

func TestHuobiProvider_GetTickerPrices(t *testing.T) {
	p, err := NewHuobiProvider(context.TODO(), zerolog.Nop(), types.CurrencyPair{Base: "ATOM", Quote: "USDT"})
	require.NoError(t, err)

	t.Run("valid_request_single_ticker", func(t *testing.T) {
		lastPrice := 34.69000000
		volume := 2396974.02000000

		tickerMap := map[string]HuobiTicker{}
		tickerMap["market.atomusdt.ticker"] = HuobiTicker{
			CH: "market.atomusdt.ticker",
			Tick: HuobiTick{
				LastPrice: lastPrice,
				Vol:       volume,
			},
		}

		p.tickers = tickerMap

		prices, err := p.GetTickerPrices(types.CurrencyPair{Base: "ATOM", Quote: "USDT"})
		require.NoError(t, err)
		require.Len(t, prices, 1)
		require.Equal(t, sdk.MustNewDecFromStr(strconv.FormatFloat(lastPrice, 'f', -1, 64)), prices["ATOMUSDT"].Price)
		require.Equal(t, sdk.MustNewDecFromStr(strconv.FormatFloat(volume, 'f', -1, 64)), prices["ATOMUSDT"].Volume)
	})

	t.Run("valid_request_multi_ticker", func(t *testing.T) {
		lastPriceAtom := 34.69000000
		lastPriceLuna := 41.35000000
		volume := 2396974.02000000

		tickerMap := map[string]HuobiTicker{}
		tickerMap["market.atomusdt.ticker"] = HuobiTicker{
			CH: "market.atomusdt.ticker",
			Tick: HuobiTick{
				LastPrice: lastPriceAtom,
				Vol:       volume,
			},
		}

		tickerMap["market.lunausdt.ticker"] = HuobiTicker{
			CH: "market.lunausdt.ticker",
			Tick: HuobiTick{
				LastPrice: lastPriceLuna,
				Vol:       volume,
			},
		}

		p.tickers = tickerMap
		prices, err := p.GetTickerPrices(
			types.CurrencyPair{Base: "ATOM", Quote: "USDT"},
			types.CurrencyPair{Base: "LUNA", Quote: "USDT"},
		)
		require.NoError(t, err)
		require.Len(t, prices, 2)
		require.Equal(t, sdk.MustNewDecFromStr(strconv.FormatFloat(lastPriceAtom, 'f', -1, 64)), prices["ATOMUSDT"].Price)
		require.Equal(t, sdk.MustNewDecFromStr(strconv.FormatFloat(volume, 'f', -1, 64)), prices["ATOMUSDT"].Volume)
		require.Equal(t, sdk.MustNewDecFromStr(strconv.FormatFloat(lastPriceLuna, 'f', -1, 64)), prices["LUNAUSDT"].Price)
		require.Equal(t, sdk.MustNewDecFromStr(strconv.FormatFloat(volume, 'f', -1, 64)), prices["LUNAUSDT"].Volume)
	})

	t.Run("invalid_request_invalid_ticker", func(t *testing.T) {
		prices, err := p.GetTickerPrices(types.CurrencyPair{Base: "FOO", Quote: "BAR"})
		require.Error(t, err)
		require.Equal(t, "failed to get ticker price for FOOBAR", err.Error())
		require.Nil(t, prices)
	})
}

func TestHuobiCurrencyPairToHuobiPair(t *testing.T) {
	cp := types.CurrencyPair{Base: "ATOM", Quote: "USDT"}
	binanceSymbol := currencyPairToHuobiTickerPair(cp)
	require.Equal(t, binanceSymbol, "market.atomusdt.ticker")
}
