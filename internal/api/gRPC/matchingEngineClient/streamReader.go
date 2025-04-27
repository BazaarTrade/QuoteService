package mClient

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/BazaarTrade/QuoteService/internal/converter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var retry = time.Second * 5

func (c *Client) StartStreamReaders(pair string) {
	ctx, cancel := context.WithCancel(c.ctx)
	c.cancelByPair[pair] = cancel
	go c.readOrderBookSnapshotStream(ctx, pair)
	go c.readTradesStream(ctx, pair)
}

func (c *Client) readOrderBookSnapshotStream(ctx context.Context, pair string) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("stopped OBS stream reader", "pair", pair)
			return
		default:
		}

		stream, err := c.client.StreamOrderBookSnapshot(ctx, &pbM.Pair{Pair: pair})
		if err != nil {
			c.logger.Warn("failed to connect to pbM OBS stream", "pair", pair)
			time.Sleep(retry)
			continue
		}

		c.logger.Info("successfully connected to pbM OBS stream", "pair", pair)

		for {
			select {
			case <-ctx.Done():
				c.logger.Info("stopping OBS stream reader", "pair", pair)
				return
			default:
			}

			OBS, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					c.logger.Info("OBS stream closed by server", "pair", pair)
					time.Sleep(retry)
					break
				}

				if status.Code(err) == codes.Canceled {
					c.logger.Info("stopped OBS stream reader", "pair", pair)
					return
				}

				c.logger.Error("failed to receive OBS", "pair", pair, "error", err, "status", status.Code(err))
				time.Sleep(retry)
				break
			}

			c.logger.Debug("recieved OBS", "pair", pair)

			go c.service.PreciseOrderBookSnaphot(converter.PbMOBSToModelsOBS(OBS))
		}
	}
}

func (c *Client) readTradesStream(ctx context.Context, pair string) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("stopped trades stream reader", "pair", pair)
			return
		default:
		}

		stream, err := c.client.StreamTrades(ctx, &pbM.Pair{Pair: pair})
		if err != nil {
			c.logger.Warn("failed to connect to pbM trade stream", "pair", pair)
			time.Sleep(retry)
			continue
		}

		c.logger.Info("successfully connected to pbM trade stream", "pair", pair)

		for {
			select {
			case <-ctx.Done():
				c.logger.Info("stopped trades stream reader", "pair", pair)
				return
			default:
			}

			trades, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					c.logger.Info("trades stream closed by server", "pair", pair)
					time.Sleep(retry)
					break
				}

				if status.Code(err) == codes.Canceled {
					c.logger.Info("stopped trades stream reader", "pair", pair)
					return
				}

				c.logger.Error("failed to receive trades", "pair", pair, "error", err)
				time.Sleep(retry)
				break
			}

			c.logger.Debug("recieved trades", "pair", pair)

			modelsTrades := converter.PbMTradesToModelsTrades(trades)

			go c.service.TickerFormation(modelsTrades)
			go c.service.CandleStickFormation(modelsTrades)
			go c.service.PreciseTrades(modelsTrades)
		}
	}
}
