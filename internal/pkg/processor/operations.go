package processor

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"

	"truetech/internal/pkg/client"
	"truetech/internal/pkg/utils"
)

func (p *Processor) calculate(
	st *client.Statistics,
	mp *client.MachineParameters,
	resources []*client.Resource,
) (int, int) {
	R := float64(st.Requests)
	NL := p.cfg.NeedLoad / 100
	rLen := float64(len(resources))
	cpu := int(math.Ceil((rLen*mp.CPU + mp.RequestCPU*R) / NL))
	ram := int(math.Ceil((rLen*mp.RAM + mp.RequestRAM*R) / 1000 / NL))
	log.Infof("!calculate: reqCPU=%v, reqRAM=%v, req=%v, needLoad=%v, rLen=%v", mp.RequestCPU, mp.RequestRAM, R, NL, rLen)
	log.Infof("!calculate-res: cpu=%v, ram=%v", cpu, ram)
	return cpu, ram
}

func (p *Processor) couldOptimize(
	st *client.Statistics,
	mp *client.MachineParameters,
	resources []*client.Resource,
	sellResources []*client.Resource,
) bool {
	R := float64(st.Requests)
	NL := 0.85
	rLen := float64(len(resources))

	cpu, ram := p.calculateWorkParameters(resources, sellResources)

	cpuLoad := math.Ceil((rLen*mp.CPU + mp.RequestCPU*R) / cpu)
	ramLoad := math.Ceil((rLen*mp.RAM + mp.RequestRAM*R) / 1000 / ram)
	if cpuLoad > NL || ramLoad > NL {
		return false
	}

	log.Infof("!calculate-opt: cpu=%v, ram=%v", cpuLoad, ramLoad)

	return true
}

func (p *Processor) buyOrSell(
	ctx context.Context,
	tp string,
	needCPU int,
	needRAM int,
	resources []*client.Resource,
) {
	currCPU, currRAM := p.countResources(resources)

	if currCPU < needCPU || currRAM < needRAM {
		p.buy(ctx, tp, needCPU, currCPU, needRAM, currRAM)
	} else if currCPU > needCPU || currRAM > needRAM {
		p.sell(ctx, tp, needCPU, currCPU, needRAM, currRAM, resources)
	} else {
		log.Infof("no need to sell or buy for %v | needCPU=%v, currCPU=%v | needRAM=%v, currRAM=%v", tp, needCPU, currCPU, needRAM, currRAM)
	}
}

func (p *Processor) buy(
	ctx context.Context,
	tp string,
	needCPU int, currCPU int,
	needRAM int, currRAM int,
) {
	log.Infof("!!buy: tp=%v | needCPU=%v, currCPU=%v | needRAM=%v, currRAM=%v", tp, needCPU, currCPU, needRAM, currRAM)

	prices := p.getPrices(ctx, tp)
	if len(prices) == 0 {
		return
	}

	deltaCPU, deltaRAM := utils.Max(needCPU-currCPU, 0), utils.Max(needRAM-currRAM, 0)

	type state struct {
		cpu  int
		ram  int
		cost int
	}

	resultPrices := client.Prices{}
	usedState := make(map[state]struct{})

	start := time.Now()

	// TODO: change to multidimensional knapsack
	currentState := state{}
	currentPrices := make(client.Prices, 0, 100)
	minCost := int(1e8)

	var backtrack func()
	backtrack = func() {
		if currentState.cost > minCost {
			return
		}

		if _, ok := usedState[currentState]; ok {
			return
		}

		if currentState.cpu >= deltaCPU && currentState.ram >= deltaRAM {
			if currentState.cost < minCost || (currentState.cost == minCost && len(resultPrices) < len(currentPrices)) {
				minCost = currentState.cost

				resultPrices = make(client.Prices, len(currentPrices))
				copy(resultPrices, currentPrices)
			}

			return
		}

		if time.Since(start).Seconds() >= 20 {
			return
		}

		for i := 0; i < len(prices); i++ {
			currentState.cpu += prices[i].CPU
			currentState.ram += prices[i].RAM
			currentState.cost += prices[i].Cost
			currentPrices = append(currentPrices, prices[i])

			backtrack()
			usedState[currentState] = struct{}{}

			currentState.cpu -= prices[i].CPU
			currentState.ram -= prices[i].RAM
			currentState.cost -= prices[i].Cost
			currentPrices = currentPrices[:len(currentPrices)-1]
		}
	}

	backtrack()

	log.Println("calc time=%v seconds", time.Since(start).Seconds())
	log.Printf("need to buy %v%v, cost=%v", len(resultPrices), tp, minCost)

	for _, rp := range resultPrices {
		res, err := p.api.BuyResource(ctx, client.BuyResourceRequest{
			CPU:  rp.CPU,
			RAM:  rp.RAM,
			Type: rp.Type,
		})

		if p.cfg.UseUpdateAbuse && err == nil {
			p.api.UpdateResource(ctx, res.ID, client.UpdateResourceRequest{
				CPU:  res.CPU,
				RAM:  res.RAM,
				Type: res.Type,
			})
		}
	}

	return
}

func (p *Processor) sell(
	ctx context.Context,
	tp string,
	needCPU int, currCPU int,
	needRAM int, currRAM int,
	resources []*client.Resource,
) {
	log.Infof("!!sell: tp=%v | needCPU=%v, currCPU=%v | needRAM=%v, currRAM=%v", tp, needCPU, currCPU, needRAM, currRAM)

	deltaCPU, deltaRAM := utils.Max(currCPU-needCPU, 0), utils.Max(currRAM-needRAM, 0)

	sortedResources := make([]*client.Resource, len(resources))
	copy(sortedResources, resources)
	sort.Slice(sortedResources, func(i, j int) bool {
		// failed first
		if sortedResources[i].Failed && !sortedResources[j].Failed {
			return true
		}

		// failed first
		if !sortedResources[i].Failed && sortedResources[j].Failed {
			return false
		}
		// 17 -- 25
		return sortedResources[i].Cost > sortedResources[j].Cost
	})

	// TODO: change to multidimensional knapsack
	for i, r := range sortedResources {
		if deltaCPU-r.CPU < -1 || deltaRAM-r.RAM < -1 {
			continue
		}

		if i+1 == len(resources) {
			log.Warn("DONT SELL LAST RESOURCE")
			break
		}

		err := p.api.SellResource(ctx, r.ID)
		if err != nil {
			fmt.Println("selled !!!", r.Failed)
			deltaCPU -= r.CPU
			deltaRAM -= r.RAM
		}
	}

	log.Infof("!!deltaCPU=%v, deltaRAM=%v", deltaCPU, deltaRAM)
	return
}

func (p *Processor) optimizeV2(ctx context.Context, tp string) {
	//log.Info("!!!OPTIMIZE V2 ", tp)
	//resources, err := p.api.ListResources(ctx)
	//if err != nil {
	//	return
	//}
	//
	//type state struct {
	//	cpu  int
	//	ram  int
	//	cost int
	//}
}

func (p *Processor) tryExchange(ctx context.Context, tp string) {

}

func (p *Processor) optimizeResources(ctx context.Context, tp string) {
	log.Info("!!!OPTIMIZE ", tp)
	resources, err := p.api.ListResources(ctx)
	if err != nil {
		return
	}

	type state struct {
		cpu  int
		ram  int
		cost int
	}

	prices := p.getPrices(ctx, tp)
	filteredResources := p.filterAndSortResources(resources, tp)

	states := make(map[state][]*client.Resource, 10)
	for _, r := range filteredResources {
		states[state{cpu: r.CPU, ram: r.RAM, cost: r.Cost}] = append(
			states[state{cpu: r.CPU, ram: r.RAM, cost: r.Cost}], r,
		)
	}

	for st, rsrcs := range states {
		for _, pr := range prices {
			if st.cpu == pr.CPU && st.ram == pr.RAM && st.cost == pr.Cost {
				continue
			}

			if st.cpu*len(rsrcs) < pr.CPU || st.ram*len(rsrcs) < pr.RAM {
				continue
			}

			if pr.CPU%st.cpu != 0 || pr.RAM%st.ram != 0 {
				continue
			}

			if pr.CPU/st.cpu != pr.RAM/st.ram {
				continue
			}

			exchangeCnt := pr.CPU / st.cpu
			if exchangeCnt*st.cost <= pr.Cost {
				continue
			}

			log.Warnf(
				"!!find optimization %v (cpu=%v, ram=%v, cost=%v) to 1 (cpu=%v, ram=%v, cost=%v)",
				exchangeCnt, st.cpu, st.ram, st.cost, pr.CPU, pr.RAM, pr.Cost,
			)
			for i := 0; i < len(rsrcs); i++ {
				if i%exchangeCnt == 0 && i+exchangeCnt-1 >= len(rsrcs) {
					break
				}

				if i%exchangeCnt == 0 {
					err = p.api.UpdateResource(ctx, rsrcs[i].ID, client.UpdateResourceRequest{
						CPU:  pr.CPU,
						RAM:  pr.RAM,
						Type: pr.Type,
					})
					if err != nil {
						return
					}
				} else {
					p.api.SellResource(ctx, rsrcs[i].ID)
				}
			}
			return
		}
	}
	return
}
