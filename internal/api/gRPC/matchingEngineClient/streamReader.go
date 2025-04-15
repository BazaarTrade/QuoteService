package mClient

import (
	"context"
	"errors"
	"io"

	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/BazaarTrade/QuoteService/internal/converter"
)

func (c *Client) StartStreamReaders(pair string) {
	go c.readOrderBookSnapshotStream(pair)
	go c.readTradesStream(pair)
}

func (c *Client) readOrderBookSnapshotStream(pair string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.client.StreamOrderBookSnapshot(ctx, &pbM.Pair{Pair: pair})
	if err != nil {
		c.logger.Error("failed to connect to pbM OBS stream", "pair", pair, "error", err)
		return
	}

	c.logger.Info("successfully connected to pbM OBS stream", "pair", pair)

	for {
		pbOBS, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				c.logger.Info("OBS stream closed by server", "pair", pair)
				return
			}
			c.logger.Error("failed to receive pbOBS", "pair", pair, "error", err)
			return
		}

		c.logger.Debug("recieved OBS", "pair", pair)

		go c.service.PreciseOrderBookSnaphot(converter.PbMOBSToModelsOBS(pbOBS))
	}
}

func (c *Client) readTradesStream(pair string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.client.StreamTrades(ctx, &pbM.Pair{Pair: pair})
	if err != nil {
		c.logger.Error("failed to connect to pbM trade stream", "pair", pair, "error", err)
		return
	}

	c.logger.Info("successfully connected to pbM trade stream", "pair", pair)

	for {
		trades, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				c.logger.Info("trades stream reader closed by server", "pair", pair)
				return
			}
			c.logger.Error("failed to receive trades", "pair", pair, "error", err)
			return
		}

		c.logger.Debug("recieved trades", "pair", pair)

		modelsTrades := converter.PbMTradesToModelsTrades(trades)

		go c.service.TickerFormation(modelsTrades)
		go c.service.CandleStickFormation(modelsTrades)
		go c.service.PreciseTrades(modelsTrades)
	}
}
