package client

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	ghc "github.com/bozd4g/go-http-client"
	log "github.com/sirupsen/logrus"
)

type Resource struct {
	Cost        int     `json:"cost"`
	CPU         int     `json:"cpu"`
	CPULoad     float64 `json:"cpu_load"`
	Failed      bool    `json:"failed"`
	FailedUntil string  `json:"failed_until"`
	ID          int     `json:"id"`
	RAM         int     `json:"ram"`
	RAMLoad     float64 `json:"ram_load"`
	Type        string  `json:"type"`
}

func (c *Resource) GetValue(param int) int {
	switch param {
	case 0:
		return c.CPU
	case 1:
		return c.RAM
	default:
		panic("unknown type")
	}
}

func (c *Client) ListResources(ctx context.Context) ([]*Resource, error) {
	resp, err := c.client.Get(ctx, "/resource", c.tokenOpt)
	if err != nil {
		log.Error("list /resource ", err)
		return nil, err
	}

	if !resp.Ok() {
		log.Warn("list /resource resp != 200 ", resp.Status())
		return nil, errors.New(strconv.Itoa(resp.Status()))
	}

	var res []*Resource
	if err = resp.Unmarshal(&res); err != nil {
		log.Error("list: unmarshal to []*Resource ", err)
		return nil, err
	}

	return res, nil
}

type BuyResourceRequest struct {
	CPU  int    `json:"cpu"`
	RAM  int    `json:"ram"`
	Type string `json:"type"`
}

func (c *Client) BuyResource(ctx context.Context, req BuyResourceRequest) (*Resource, error) {
	log.Warnf("buy resource of type=%v, cpu=%v, ram=%v", req.Type, req.CPU, req.RAM)

	body, err := json.Marshal(req)
	if err != nil {
		log.Error("marshal BuyResourceRequest", err)
		return nil, err
	}

	bodyOpt := ghc.WithBody(body)
	bodyHeaderOpt := ghc.WithHeader("Content-Type", "application/json")
	resp, err := c.client.Post(ctx, "/resource", c.tokenOpt, bodyHeaderOpt, bodyOpt)
	if err != nil {
		log.Error("buy /resource ", err)
		return nil, err
	}

	if !resp.Ok() {
		log.Warn("buy /resource resp != 200 ", resp.Status())
		return nil, errors.New(strconv.Itoa(resp.Status()))
	}

	var res Resource
	if err = resp.Unmarshal(&res); err != nil {
		log.Error("buy: unmarshal to Resource ", err)
		return nil, err
	}

	return &res, nil
}
