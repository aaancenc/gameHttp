package miners

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	TypeMinerBasic   MinerType = "normal"
	BasicMinerSalary           = 50
)

type NormalMiner struct {
	id uuid.UUID

	energy   *atomic.Int64
	power    *atomic.Int64
	interval time.Duration
}

func NewBasicMiner() *NormalMiner {
	const (
		normalMinerEnergy    = 45
		normalMinerPower     = 3
		normalMinerlInterval = 2 * time.Second
	)

	energy := &atomic.Int64{}
	power := &atomic.Int64{}

	energy.Add(normalMinerEnergy)
	power.Add(normalMinerPower)

	return &NormalMiner{
		id:     uuid.New(),
		energy: energy,
		power:  power,
	}
}

func (m *NormalMiner) Run(ctx context.Context) <-chan Coal {
	transferPoint := make(chan Coal)

	go func() {
		defer func() {
			close(transferPoint)
			fmt.Println("normal miner stopped:", m.id)
		}()

		fmt.Println("normal miner started:", m.id)

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
				if new := m.energy.Add(-1); new <= 0 {
					return
				}
			}
		}
	}()

	return transferPoint
}

func (m *NormalMiner) Info() MinerInfo {
	return MinerInfo{
		ID:          m.id,
		MinerType:   TypeMinerBasic,
		MinerEnergy: m.energy.Load(),
		MinerPower:  m.power.Load(),
	}
}
