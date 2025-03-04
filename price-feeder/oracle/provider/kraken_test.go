package provider

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/umee-network/umee/price-feeder/oracle/types"
)

func TestKrakenProvider_GetTickerPrices(t *testing.T) {
	p, err := NewKrakenProvider(context.TODO(), zerolog.Nop(), types.CurrencyPair{Base: "BTC", Quote: "USDT"})
	require.NoError(t, err)

	t.Run("valid_request_single_ticker", func(t *testing.T) {
		lastPrice := sdk.MustNewDecFromStr("34.69000000")
		volume := sdk.MustNewDecFromStr("2396974.02000000")

		tickerMap := map[string]TickerPrice{}
		tickerMap["ATOMUSDT"] = TickerPrice{
			Price:  lastPrice,
			Volume: volume,
		}

		p.tickers = tickerMap

		prices, err := p.GetTickerPrices(types.CurrencyPair{Base: "ATOM", Quote: "USDT"})
		require.NoError(t, err)
		require.Len(t, prices, 1)
		require.Equal(t, lastPrice, prices["ATOMUSDT"].Price)
		require.Equal(t, volume, prices["ATOMUSDT"].Volume)
	})

	t.Run("valid_request_multi_ticker", func(t *testing.T) {
		lastPriceAtom := sdk.MustNewDecFromStr("34.69000000")
		lastPriceLuna := sdk.MustNewDecFromStr("41.35000000")
		volume := sdk.MustNewDecFromStr("2396974.02000000")

		tickerMap := map[string]TickerPrice{}
		tickerMap["ATOMUSDT"] = TickerPrice{
			Price:  lastPriceAtom,
			Volume: volume,
		}

		tickerMap["LUNAUSDT"] = TickerPrice{
			Price:  lastPriceLuna,
			Volume: volume,
		}

		p.tickers = tickerMap
		prices, err := p.GetTickerPrices(
			types.CurrencyPair{Base: "ATOM", Quote: "USDT"},
			types.CurrencyPair{Base: "LUNA", Quote: "USDT"},
		)
		require.NoError(t, err)
		require.Len(t, prices, 2)
		require.Equal(t, lastPriceAtom, prices["ATOMUSDT"].Price)
		require.Equal(t, volume, prices["ATOMUSDT"].Volume)
		require.Equal(t, lastPriceLuna, prices["LUNAUSDT"].Price)
		require.Equal(t, volume, prices["LUNAUSDT"].Volume)
	})

	t.Run("invalid_request_invalid_ticker", func(t *testing.T) {
		prices, err := p.GetTickerPrices(types.CurrencyPair{Base: "FOO", Quote: "BAR"})
		require.Error(t, err)
		require.Equal(t, "failed to get ticker price for FOOBAR", err.Error())
		require.Nil(t, prices)
	})
}

func TestKrakenPairToCurrencyPairSymbol(t *testing.T) {
	cp := types.CurrencyPair{Base: "ATOM", Quote: "USDT"}
	currencyPairSymbol := krakenPairToCurrencyPairSymbol("ATOM/USDT")
	require.Equal(t, cp.String(), currencyPairSymbol)
}

func TestKrakenCurrencyPairToKrakenPair(t *testing.T) {
	cp := types.CurrencyPair{Base: "ATOM", Quote: "USDT"}
	krakenSymbol := currencyPairToKrakenPair(cp)
	require.Equal(t, krakenSymbol, "ATOM/USDT")
}

func TestNormalizeKrakenBTCPair(t *testing.T) {
	btcSymbol := normalizeKrakenBTCPair("XBT/USDT")
	require.Equal(t, btcSymbol, "BTC/USDT")

	atomSymbol := normalizeKrakenBTCPair("ATOM/USDT")
	require.Equal(t, atomSymbol, "ATOM/USDT")
}
