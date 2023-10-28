package client

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type MachineParameters struct {
	MinLoad    float64 `json:"MinLoad"`
	MaxLoad    float64 `json:"MaxLoad"`
	RequestCPU float64 `json:"RequestCPU"`
	RequestRAM float64 `json:"RequestRAM"`
	CPU        float64 `json:"CPU"`
	RAM        float64 `json:"RAM"`
}

type ProcessorConfig struct {
	NeedLoad              float64            `json:"NeedLoad"`
	SleepTimeSeconds      int                `json:"SleepTimeSeconds"`
	ErrorSleepTimeSeconds int                `json:"ErrorSleepTimeSeconds"`
	UseUpdateAbuse        bool               `json:"UseUpdateAbuse"`
	MaxMachineParam       int                `json:"MaxMachineParam"`
	DB                    *MachineParameters `json:"DB"`
	VM                    *MachineParameters `json:"VM"`
}

func DefaultProcessorConfig() *ProcessorConfig {
	return &ProcessorConfig{
		NeedLoad:              70,
		SleepTimeSeconds:      60,
		ErrorSleepTimeSeconds: 7,
		UseUpdateAbuse:        true,
		MaxMachineParam:       16,
		DB: &MachineParameters{
			MinLoad:    65,
			MaxLoad:    75,
			RequestCPU: 0.001,
			RequestRAM: 1.0,
			CPU:        0.05,
			RAM:        500,
		},
		VM: &MachineParameters{
			MinLoad:    65,
			MaxLoad:    75,
			RequestCPU: 0.001,
			RequestRAM: 5.0,
			CPU:        0.05,
			RAM:        300,
		},
	}
}

func (c *Client) GetConfig(ctx context.Context) *ProcessorConfig {
	resp, err := c.cfgClient.Get(ctx, "/config.json")
	if err != nil {
		log.Error("get /config", err)
		return DefaultProcessorConfig()
	}

	if !resp.Ok() {
		log.Warn("get /config resp != 200", resp.Status())
		return DefaultProcessorConfig()
	}

	var res ProcessorConfig
	if err = resp.Unmarshal(&res); err != nil {
		log.Error("get: unmarshal to ProcessorConfig", err)
		return DefaultProcessorConfig()
	}
	return &res
}
