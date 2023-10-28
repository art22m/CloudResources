package client

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	log "github.com/sirupsen/logrus"

	ghc "github.com/bozd4g/go-http-client"
)

func (c *Client) SellResource(ctx context.Context, id int) error {
	log.Warnf("sell id=%v", id)
	resp, err := c.client.Delete(ctx, "/resource/"+strconv.Itoa(id), c.tokenOpt)
	if err != nil {
		log.Error("delete /resource/ ", id, err)
		return err
	}

	if !resp.Ok() {
		log.Warn("delete /resource/ resp != 200 ", resp.Status())
		return errors.New(strconv.Itoa(resp.Status()))
	}

	return nil
}

func (c *Client) GetResource(ctx context.Context, id int) (*Resource, error) {
	log.Infof("get id=%v", id)
	resp, err := c.client.Get(ctx, "/resource/"+strconv.Itoa(id), c.tokenOpt)
	if err != nil {
		log.Error("get /resource/ ", id, err)
		return nil, err
	}

	if !resp.Ok() {
		log.Warn("get /resource/ resp != 200 ", resp.Status())
		return nil, errors.New(strconv.Itoa(resp.Status()))
	}

	var res Resource
	if err = resp.Unmarshal(&res); err != nil {
		log.Error("get: unmarshal to Resource ", err)
		return nil, err
	}

	return &res, nil
}

type UpdateResourceRequest struct {
	CPU  int    `json:"cpu"`
	RAM  int    `json:"ram"`
	Type string `json:"type"`
}

func (c *Client) UpdateResource(ctx context.Context, id int, req UpdateResourceRequest) error {
	log.Warnf("update id=%v", id)
	body, err := json.Marshal(req)
	if err != nil {
		log.Error("marshal UpdateResourceRequest ", err)
		return err
	}

	bodyOpt := ghc.WithBody(body)
	bodyHeaderOpt := ghc.WithHeader("Content-Type", "application/json")
	resp, err := c.client.Put(ctx, "/resource/"+strconv.Itoa(id), c.tokenOpt, bodyHeaderOpt, bodyOpt)
	if err != nil {
		log.Error("put /resource/ ", id, err)
		return err
	}

	if !resp.Ok() {
		log.Warn("put /resource/ resp != 200 ", resp.Status())
		return errors.New(strconv.Itoa(resp.Status()))
	}

	return nil
}
