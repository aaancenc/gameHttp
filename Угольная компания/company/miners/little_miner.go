package miners

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	TypeMinerLittle   MinerType = "little"
	LittleMinerSalary           = 5
)

type LittleMiner struct {
	id       uuid.UUID
	energy   *atomic.Int64
	power    *atomic.Int64
	interval time.Duration
}

func NewLittleMiner() *LittleMiner {
	const (
		littleMinerEnergy   = 30
		littleMinerPower    = 1
		littleMinerInterval = 3 * time.Second
	)

	energy := &atomic.Int64{}
	power := &atomic.Int64{}

	energy.Add(littleMinerEnergy)
	power.Add(littleMinerPower)

	return &LittleMiner{
		id:       uuid.New(),
		energy:   energy,
		power:    power,
		interval: littleMinerInterval,
	}
}

func (m *LittleMiner) Run(ctx context.Context) <-chan Coal {
	transferPoint := make(chan Coal)

	go func() {
		defer func() {
			close(transferPoint)
			fmt.Println("little miner stopped:", m.id)
		}()

		fmt.Println("little miner started:", m.id)

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

func (m *LittleMiner) Info() MinerInfo {
	return MinerInfo{
		ID:          m.id,
		MinerType:   TypeMinerLittle,
		MinerEnergy: m.energy.Load(),
		MinerPower:  m.power.Load(),
	}
}
