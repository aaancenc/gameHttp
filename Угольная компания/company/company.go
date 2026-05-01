package company

import (
	"coal_company/company/equipment"
	"coal_company/company/miners"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Company struct {
	incomeCh chan miners.Coal

	mtx    sync.RWMutex
	miners map[miners.MinerType]map[uuid.UUID]miners.Miner

	equipment equipment.Equipment

	companyCtx  context.Context
	companyStop context.CancelFunc

	// -----
	statisctics *CompanyStatistics
}

func NewCompany(ctx context.Context) *Company {
	companyCtx, companyStop := context.WithCancel(ctx)

	c := &Company{
		incomeCh:    make(chan miners.Coal),
		miners:      make(map[miners.MinerType]map[uuid.UUID]miners.Miner),
		equipment:   equipment.NewEquipment(),
		companyCtx:  companyCtx,
		companyStop: companyStop,

		// -----
		statisctics: NewCompanyStatistics(),
	}

	go c.baseIncome()
	go c.collectIncome()

	return c
}

func (c *Company) HireMiner(minerType miners.MinerType) (miners.Miner, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var miner miners.Miner

	switch minerType {
	case miners.TypeMinerBasic:
		if c.statisctics.ballance.Load() >= miners.BasicMinerSalary {
			miner = miners.NewBasicMiner()
			c.statisctics.ballance.Add(-miners.BasicMinerSalary)
		} else {
			return nil, ErrInsufficientFunds
		}

	case miners.TypeMinerLittle:
		if c.statisctics.ballance.Load() >= miners.LittleMinerSalary {
			miner = miners.NewLittleMiner()
			c.statisctics.ballance.Add(-miners.LittleMinerSalary)
		} else {
			return nil, ErrInsufficientFunds
		}

	case miners.TypeMinerPowerful:
		if c.statisctics.ballance.Load() >= miners.PowerfulMinerSalary {
			miner = miners.NewPowerfulMiner()
			c.statisctics.ballance.Add(-miners.PowerfulMinerSalary)
		} else {
			return nil, ErrInsufficientFunds
		}

	default:
		return nil, ErrUnknokwnMinerType
	}

	info := miner.Info()

	if c.miners[info.MinerType] == nil {
		c.miners[info.MinerType] = make(map[uuid.UUID]miners.Miner)
	}

	c.miners[info.MinerType][info.ID] = miner

	coalCh := miner.Run(c.companyCtx)

	go func() {
		for v := range coalCh {
			c.incomeCh <- v
		}

		c.mtx.Lock()
		delete(c.miners[info.MinerType], info.ID)

		if len(c.miners[minerType]) == 0 {
			delete(c.miners, info.MinerType)
		}
		c.mtx.Unlock()
	}()

	c.statisctics.totalMinersStatistics[minerType]++

	return miner, nil
}

func (c *Company) GetAllMiners() map[miners.MinerType]map[uuid.UUID]miners.Miner {
	tmp := make(map[miners.MinerType]map[uuid.UUID]miners.Miner)
	for minerType, minersMap := range c.miners {
		tmp[minerType] = make(map[uuid.UUID]miners.Miner)

		for k, v := range minersMap {
			tmp[minerType][k] = v
		}
	}

	return tmp
}

func (c *Company) GetMinersByType(minerType miners.MinerType) map[uuid.UUID]miners.Miner {
	tmp := make(map[uuid.UUID]miners.Miner)
	for k, v := range c.miners[minerType] {
		tmp[k] = v
	}

	return tmp
}

func (c *Company) BuyEquipment(equipmentType equipment.EquipmentType) (equipment.Equipment, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if equipmentType == equipment.EquipmentTypePickaxe {
		if c.statisctics.ballance.Load() >= equipment.EquipmentPickaxeCost {
			c.equipment.BuyPickaxe()
			c.statisctics.ballance.Add(-equipment.EquipmentPickaxeCost)
		} else {
			return equipment.Equipment{}, ErrInsufficientFunds
		}
	} else if equipmentType == equipment.EquipmentTypeVentilation {
		if c.statisctics.ballance.Load() >= equipment.EquipmentVentilationCost {
			c.equipment.BuyVentilation()
			c.statisctics.ballance.Add(-equipment.EquipmentVentilationCost)
		} else {
			return equipment.Equipment{}, ErrInsufficientFunds
		}
	} else if equipmentType == equipment.EquipmentTypeTrolleys {
		if c.statisctics.ballance.Load() >= equipment.EquipmentTrolleysCost {
			c.equipment.BuyTrolleys()
			c.statisctics.ballance.Add(-equipment.EquipmentTrolleysCost)
		} else {
			return equipment.Equipment{}, ErrInsufficientFunds
		}
	} else {
		return equipment.Equipment{}, ErrUnknokwnEquipmentType
	}

	return c.equipment, nil
}

func (c *Company) GetEquipment() equipment.Equipment {
	return c.equipment
}

func (c *Company) Complete() (CompanyStatistics, error) {
	if !c.equipment.AllBought() {
		return CompanyStatistics{}, ErrNotAllEquipmentPurchased
	}

	c.companyStop()

	now := time.Now()
	c.statisctics.completedTime = &now

	return *c.statisctics, nil
}

func (c *Company) GetStatistics() CompanyStatistics {
	return *c.statisctics
}

func (c *Company) collectIncome() {
	for {
		select {
		case <-c.companyCtx.Done():
			return
		case income := <-c.incomeCh:
			c.statisctics.ballance.Add(int64(income))
			fmt.Println("ballance:", c.statisctics.ballance.Load())

			c.statisctics.totalEarned.Add(int64(income))
		}
	}
}

func (c *Company) baseIncome() {
	for {
		select {
		case <-c.companyCtx.Done():
			return

		case <-time.After(1 * time.Second):
		}

		select {
		case <-c.companyCtx.Done():
			return
		case c.incomeCh <- 1:
		}
	}
}
