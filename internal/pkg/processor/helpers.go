package processor

import (
	"context"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"

	"truetech/internal/pkg/client"
)

func (p *Processor) atLeastOneFailed(resources []*client.Resource, tp string) bool {
	for _, r := range resources {
		if r.Type == tp && r.Failed {
			return true
		}
	}
	return false
}

func (p *Processor) filterAndSortResources(resources []*client.Resource, tp string) (res []*client.Resource) {
	for _, r := range resources {
		if r.Type == tp {
			res = append(res, r)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		// failed first
		if res[i].Failed && !res[j].Failed {
			return true
		}

		// failed first
		if !res[i].Failed && res[j].Failed {
			return false
		}

		return res[i].Cost < res[j].Cost
	})

	return
}

func (p *Processor) countResources(resources []*client.Resource) (currCPU, currRAM int) {
	for _, r := range resources {
		currCPU += r.CPU
		currRAM += r.RAM
	}
	return
}

func (p *Processor) getWithParameters(resources []*client.Resource, cpu int) (int, bool) {
	for _, r := range resources {
		if r.CPU == cpu {
			return r.ID, true
		}
	}
	return 0, false
}

func (p *Processor) getPrices(ctx context.Context, tp string) (res client.Prices) {
	prices, err := p.api.GetPrices(ctx)
	if err != nil {
		return
	}

	for _, pr := range *prices {
		if pr.Type == tp && pr.CPU <= p.cfg.MaxMachineParam && pr.RAM <= p.cfg.MaxMachineParam {
			res = append(res, pr)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Cost > res[j].Cost
	})
	return
}

func (p *Processor) sellAll(ctx context.Context) {
	resources, _ := p.api.ListResources(ctx)
	for _, r := range resources {
		if r.Type == VM {
			err := p.api.SellResource(ctx, r.ID)
			fmt.Println("sell", VM, r.ID, err)
		}

		if r.Type == DB {
			err := p.api.SellResource(ctx, r.ID)
			fmt.Println("sell", DB, r.ID, err)
		}
	}
}

func (p *Processor) calculateWorkParameters(
	resources []*client.Resource,
	sellResources []*client.Resource,
) (cpu, ram float64) {
	for _, r := range resources {
		if !r.Failed {
			cpu += float64(r.CPU)
			ram += float64(r.RAM)
		}
	}

	for _, r := range sellResources {
		if !r.Failed {
			cpu -= float64(r.CPU)
			ram -= float64(r.RAM)
		}
	}

	return
}

func (p *Processor) logStatistics(ctx context.Context) {
	st, err := p.api.GetStatistic(ctx)
	if err != nil {
		return
	}

	log.Info("Availability ", st.Availability)

	log.Info("Requests ", st.Requests)
	log.Info("RequestsTotal ", st.RequestsTotal)
	log.Info("ResponseTime ", st.ResponseTime)

	log.Info("VMCPU ", st.VMCPU)
	log.Info("VMCPULoad ", st.VMCPULoad)
	log.Info("VMRAM ", st.VMRAM)
	log.Info("VMRAMLoad ", st.VMRAMLoad)

	log.Info("DBCPU ", st.DBCPU)
	log.Info("DBCPULoad ", st.DBCPULoad)
	log.Info("DBRAM ", st.DBRAM)
	log.Info("DBRAMLoad ", st.DBRAMLoad)
	log.Println("-----------------")
}
