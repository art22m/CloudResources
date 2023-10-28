package client

import (
	"context"
	"errors"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Prices []struct {
	Cost int    `json:"cost"`
	CPU  int    `json:"cpu"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	RAM  int    `json:"ram"`
	Type string `json:"type"`
}

func (c *Client) GetPrices(ctx context.Context) (*Prices, error) {
	resp, err := c.client.Get(ctx, "/price")
	if err != nil {
		log.Error("get /price ", err)
		return nil, err
	}

	if !resp.Ok() {
		log.Warn("get /price resp != 200 ", resp.Status())
		return nil, errors.New(strconv.Itoa(resp.Status()))
	}

	var res Prices
	if err = resp.Unmarshal(&res); err != nil {
		log.Error("get: unmarshal to Prices ", err)
		return nil, err
	}

	return &res, nil
}
