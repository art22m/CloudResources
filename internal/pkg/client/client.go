package client

import ghc "github.com/bozd4g/go-http-client"

type Client struct {
	client    *ghc.Client
	cfgClient *ghc.Client
	tokenOpt  ghc.Option
}

func NewClient(baseUrl string) *Client {
	client := ghc.New(baseUrl)
	tokenOpt := ghc.WithQuery(
		"token",
		"some-token",
	)

	cfgClient := ghc.New("http://storage.yandexcloud.net/cloud-resources")

	return &Client{
		client:    client,
		cfgClient: cfgClient,
		tokenOpt:  tokenOpt,
	}
}
