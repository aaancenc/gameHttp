package miners

import (
	"context"

	"github.com/google/uuid"
)

type (
	MinerType string
	Coal      int64
)

type Miner interface {
	Run(ctx context.Context) <-chan Coal
	Info() MinerInfo
}

type MinerInfo struct {
	ID          uuid.UUID
	MinerType   MinerType
	MinerEnergy int64
	MinerPower  int64
}
