package miners

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	TypeMinerPowerful   MinerType = "powerful"
	PowerfulMinerSalary           = 450
)

type PowerfulMiner struct {
	id       uuid.UUID
	energy   *atomic.Int64
	power    *atomic.Int64
	interval time.Duration
}

func NewPowerfulMiner() *PowerfulMiner {
	const (
		powerfulMinerEnergy   = 60
		powerfulMinerPower    = 10
		powerfulMinerInterval = 1 * time.Second
	)

	energy := &atomic.Int64{}
	power := &atomic.Int64{}

	energy.Add(powerfulMinerEnergy)
	power.Add(powerfulMinerPower)

	return &PowerfulMiner{
		id:       uuid.New(),
		energy:   energy,
		power:    power,
		interval: powerfulMinerInterval,
	}
}

func (m *PowerfulMiner) Run(ctx context.Context) <-chan Coal {
	transferPoint := make(chan Coal)

	go func() {
		defer func() {
			close(transferPoint)
			fmt.Println("powerful miner stopped:", m.id)
		}()

		fmt.Println("powerful miner started:", m.id)

		for {
			select {
			case <-ctx.Done():
				return

			case <-time.After(m.interval):
			}

			select {
			case <-ctx.Done():
				return

			case transferPoint <- Coal(m.power.Load()):
				m.power.Add(3)

				if new := m.energy.Add(-1); new <= 0 {
					return
				}
			}
		}
	}()

	return transferPoint
}

func (m *PowerfulMiner) Info() MinerInfo {
	return MinerInfo{
		ID:          m.id,
		MinerType:   TypeMinerPowerful,
		MinerEnergy: m.energy.Load(),
		MinerPower:  m.power.Load(),
	}
}
