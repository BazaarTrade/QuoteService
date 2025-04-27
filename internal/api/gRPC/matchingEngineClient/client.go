package mClient

import (
	"context"
	"log/slog"

	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/BazaarTrade/QuoteService/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client       pbM.MatchingEngineClient
	conn         *grpc.ClientConn
	service      *service.Service
	ctx          context.Context
	cancel       context.CancelFunc
	cancelByPair map[string]context.CancelFunc
	logger       *slog.Logger
}

func New(service *service.Service, logger *slog.Logger) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		service:      service,
		ctx:          ctx,
		cancel:       cancel,
		cancelByPair: make(map[string]context.CancelFunc),
		logger:       logger,
	}
}

func (c *Client) Run(CONN_ADDR_MATCHING_ENGINE string) error {
	var err error
	c.conn, err = grpc.NewClient(CONN_ADDR_MATCHING_ENGINE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.logger.Error("GRPC error connecting to mClient", "error", err)
		return err
	}

	c.client = pbM.NewMatchingEngineClient(c.conn)

	return nil
}

func (c *Client) CloseConnection() {
	if c.cancel != nil {
		c.cancel()
		c.cancelByPair = make(map[string]context.CancelFunc)
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.Error("failed to close mClient connection", "error", err)
		}
	}
}

func (c *Client) StopStreamReadersByPair(pair string) {
	if cancel, ok := c.cancelByPair[pair]; ok {
		cancel()
		delete(c.cancelByPair, pair)
	}
}
