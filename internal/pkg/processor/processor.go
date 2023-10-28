package processor

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"truetech/internal/pkg/client"
)

const (
	VM string = "vm"
	DB        = "db"
)

const (
	CPU int = 0
	RAM     = 1
)

type Processor struct {
	api *client.Client
	cfg *client.ProcessorConfig

	tick time.Duration
}

func NewProcessor(cl *client.Client) *Processor {
	return &Processor{
		api:  cl,
		tick: 60 * time.Second,
	}
}

func (p *Processor) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		p.updateConfig(ctx)

		statistics, err := p.api.GetStatistic(ctx)
		if err != nil {
			sleep(p.cfg.ErrorSleepTimeSeconds)
			continue
		}

		if statistics.Requests <= 1200 {
			log.Info("low requests, skip")
			continue
		}

		p.buyOrSellIfNeeded(ctx, p.cfg.VM, VM, statistics)
		fmt.Println()
		p.buyOrSellIfNeeded(ctx, p.cfg.DB, DB, statistics)

		fmt.Println()

		//p.optimizeResources(ctx, VM)
		fmt.Println()
		//p.optimizeResources(ctx, DB)

		sleep(p.cfg.SleepTimeSeconds)
	}
}

func (p *Processor) updateConfig(ctx context.Context) {
	p.cfg = p.api.GetConfig(ctx)
}

func (p *Processor) buyOrSellIfNeeded(ctx context.Context, mp *client.MachineParameters, tp string, st *client.Statistics) {
	if mp.MinLoad <= st.MaxLoad(tp) && st.MaxLoad(tp) <= mp.MaxLoad && st.ResponseTime < 400 {
		log.Infof("no need to buy or sell for %s, maxLoad=%v", tp, st.MaxLoad(tp))
		return
	}

	log.Infof("need to smth for %s, needLoad=%v, currLoad=%v", tp, p.cfg.NeedLoad, st.MaxLoad(tp))

	resources, err := p.api.ListResources(ctx)
	if err != nil {
		return
	}

	filteredResources := p.filterAndSortResources(resources, tp)

	needCPU, needRAM := p.calculate(st, mp, filteredResources)
	p.buyOrSell(ctx, tp, needCPU, needRAM, filteredResources)
}

func sleep(seconds int) {
	log.Infof("sleep for %v seconds", seconds)
	fmt.Println()
	time.Sleep(time.Duration(seconds) * time.Second)
}
