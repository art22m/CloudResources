package main

import (
	"fmt"

	"truetech/internal/pkg/utils"
)

type Price struct {
	CPU  int
	RAM  int
	Cost int
	Name string
}

func main() {
	needCPU, currCPU := 8, 0
	needRAM, currRAM := 32, 0

	prices := []Price{
		{
			CPU:  1,
			RAM:  1,
			Cost: 2,
			Name: "2",
		},
		{
			CPU:  1,
			RAM:  2,
			Cost: 3,
			Name: "3",
		},
		{
			CPU:  2,
			RAM:  2,
			Cost: 4,
			Name: "4",
		},
		{
			CPU:  2,
			RAM:  4,
			Cost: 6,
			Name: "5",
		},
		{
			CPU:  2,
			RAM:  4,
			Cost: 6,
			Name: "6",
		}, {
			CPU:  2,
			RAM:  8,
			Cost: 10,
			Name: "7",
		},
		{
			CPU:  4,
			RAM:  4,
			Cost: 8,
			Name: "8",
		},
		{
			CPU:  4,
			RAM:  8,
			Cost: 12,
			Name: "9",
		},
		{
			CPU:  4,
			RAM:  16,
			Cost: 20,
			Name: "10",
		},
		{
			CPU:  8,
			RAM:  8,
			Cost: 16,
			Name: "11",
		},
		{
			CPU:  8,
			RAM:  16,
			Cost: 24,
			Name: "12",
		},
		{
			CPU:  8,
			RAM:  32,
			Cost: 40,
			Name: "13",
		},
	}

	deltaCPU, deltaRAM := utils.Max(needCPU-currCPU, 0), utils.Max(needRAM-currRAM, 0)

	type state struct {
		cpu  int
		ram  int
		cost int
	}

	resultPrices := make([]Price, 0, 100)
	used := make(map[state]struct{})

	// TODO: change to multidimensional knapsack
	currentState := state{}
	currentPrices := make([]Price, 0, 100)
	minCost := int(1e8)

	//start := time.Now()
	var backtrack func()
	backtrack = func() {
		if currentState.cost >= minCost {
			return
		}

		//if _, ok := used[currentState]; ok {
		//	return
		//}

		if currentState.cpu >= deltaCPU && currentState.ram >= deltaRAM {
			if currentState.cost < minCost {
				minCost = currentState.cost
				resultPrices = make([]Price, len(currentPrices))
				copy(resultPrices, currentPrices)
				fmt.Println(len(currentPrices))
			}

			fmt.Println(currentState.cost)
			return
		}

		for i := 0; i < len(prices); i++ {
			currentState.cpu += prices[i].CPU
			currentState.ram += prices[i].RAM
			currentState.cost += prices[i].Cost
			currentPrices = append(currentPrices, prices[i])

			backtrack()
			used[currentState] = struct{}{}

			currentState.cpu -= prices[i].CPU
			currentState.ram -= prices[i].RAM
			currentState.cost -= prices[i].Cost
			currentPrices = currentPrices[:len(currentPrices)-1]
		}
	}

	backtrack()

	fmt.Println("----")
	fmt.Println(minCost)
	fmt.Println(len(resultPrices))
	for _, rp := range resultPrices {
		fmt.Println(rp.CPU, rp.RAM, rp.Cost, rp.Name)
	}
}
