package client

import (
	"context"
	"errors"
	"strconv"

	log "github.com/sirupsen/logrus"

	"truetech/internal/pkg/utils"
)

type Statistics struct {
	Availability  float64 `json:"availability"`
	CostTotal     int     `json:"cost_total"`
	DBCPU         int     `json:"db_cpu"`
	DBCPULoad     float64 `json:"db_cpu_load"`
	DBRAM         int     `json:"db_ram"`
	DBRAMLoad     float64 `json:"db_ram_load"`
	ID            int     `json:"id"`
	Last1         int     `json:"last1"`
	Last15        int     `json:"last15"`
	Last5         int     `json:"last5"`
	LastDay       int     `json:"lastDay"`
	LastHour      int     `json:"lastHour"`
	LastWeek      int     `json:"lastWeek"`
	OfflineTime   int     `json:"offline_time"`
	Online        bool    `json:"online"`
	OnlineTime    int     `json:"online_time"`
	Requests      int     `json:"requests"`
	RequestsTotal int     `json:"requests_total"`
	ResponseTime  int     `json:"response_time"`
	Timestamp     string  `json:"timestamp"`
	UserID        int     `json:"user_id"`
	UserName      string  `json:"user_name"`
	VMCPU         int     `json:"vm_cpu"`
	VMCPULoad     float64 `json:"vm_cpu_load"`
	VMRAM         int     `json:"vm_ram"`
	VMRAMLoad     float64 `json:"vm_ram_load"`
}

func (c *Client) GetStatistic(ctx context.Context) (*Statistics, error) {
	resp, err := c.client.Get(ctx, "/statistic", c.tokenOpt)
	if err != nil {
		log.Error("get /statistic", err)
		return nil, err
	}

	if !resp.Ok() {
		log.Warn("get /statistic resp != 200", resp.Status())
		return nil, errors.New(strconv.Itoa(resp.Status()))
	}

	var res Statistics
	if err = resp.Unmarshal(&res); err != nil {
		log.Error("get: unmarshal to Statistics", err)
		return nil, err
	}

	return &res, nil
}

func (s *Statistics) MaxLoad(tp string) float64 {
	switch tp {
	case "vm":
		return utils.Max(s.VMRAMLoad, s.VMCPULoad)
	case "db":
		return utils.Max(s.DBRAMLoad, s.DBCPULoad)
	default:
		panic("unknown type")
	}
}
