package company

import (
	"coal_company/company/miners"
	"sync/atomic"
	"time"
)

type CompanyStatistics struct {
	ballance    *atomic.Int64
	totalEarned *atomic.Int64

	totalMinersStatistics map[miners.MinerType]int

	createdTime   time.Time
	completedTime *time.Time
}

func NewCompanyStatistics() *CompanyStatistics {
	return &CompanyStatistics{
		ballance:              &atomic.Int64{},
		totalEarned:           &atomic.Int64{},
		totalMinersStatistics: make(map[miners.MinerType]int),
		createdTime:           time.Now(),
	}
}

func (c CompanyStatistics) Ballance() int64 {
	return c.ballance.Load()
}

func (c CompanyStatistics) TotalEarned() int64 {
	return c.totalEarned.Load()
}

func (c CompanyStatistics) TimeToComple() string {
	if c.completedTime == nil {
		return ""
	}

	return c.completedTime.Sub(c.createdTime).String()
}

func (c CompanyStatistics) TotalMinersStatistics() map[miners.MinerType]int {
	tmp := make(map[miners.MinerType]int)
	for k, v := range c.totalMinersStatistics {
		tmp[k] = v
	}

	return tmp
}
