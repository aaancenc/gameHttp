package output_dto

import (
	"coal_company/company"
	"coal_company/company/equipment"
	"coal_company/company/miners"

	"github.com/google/uuid"
)

type EquipmenetDTO struct {
	Pickaxe bool `json:"pickaxe"`

	// Вентиляция
	Ventilation bool `json:"ventilation"`

	// Вагонетки
	Trolleys bool `json:"trolleys"`
}

func NewEquipmentDTO(equipment equipment.Equipment) EquipmenetDTO {
	return EquipmenetDTO{
		Pickaxe:     equipment.PickaxesPurchased(),
		Ventilation: equipment.VentilationPurchased(),
		Trolleys:    equipment.TrolleysPurchased(),
	}
}

type MinersByTypeDTO map[uuid.UUID]miners.MinerInfo

func NewMinersByTypeDTO(m map[uuid.UUID]miners.Miner) MinersByTypeDTO {
	info := make(map[uuid.UUID]miners.MinerInfo)

	for id, miner := range m {
		info[id] = miner.Info()
	}

	return info
}

type AllMinersDTO map[miners.MinerType]MinersByTypeDTO

func NewAllMinersDTO(m map[miners.MinerType]map[uuid.UUID]miners.Miner) AllMinersDTO {
	info := make(map[miners.MinerType]MinersByTypeDTO)

	for minersType, _ := range m {
		minersByTypeOutputDTO := NewMinersByTypeDTO(m[minersType])
		info[minersType] = minersByTypeOutputDTO
	}

	return info
}

type CompanyStatisticsDTO struct {
	Ballance              int64
	TotalEarned           int64
	TotalMinersStatistics map[miners.MinerType]int
	TotalTime             string
}

func NewCompanyStatisticsDTO(companyStatistics company.CompanyStatistics) CompanyStatisticsDTO {
	return CompanyStatisticsDTO{
		Ballance:              companyStatistics.Ballance(),
		TotalEarned:           companyStatistics.TotalEarned(),
		TotalMinersStatistics: companyStatistics.TotalMinersStatistics(),
		TotalTime:             companyStatistics.TimeToComple(),
	}
}

type MinersSalariesDTO struct {
	LittleMinerSalary   int64 `json:"little_miner_price"`
	BasicMinerSalary    int64 `json:"basic_miner_price"`
	PowerfulMinerSalary int64 `json:"powerful_miner_price"`
}

func NewMinersSalariesDTO(
	littleMinerSalary int64,
	basicMinerSalary int64,
	powerfulMinerSalary int64,
) MinersSalariesDTO {
	return MinersSalariesDTO{
		LittleMinerSalary:   littleMinerSalary,
		BasicMinerSalary:    basicMinerSalary,
		PowerfulMinerSalary: powerfulMinerSalary,
	}
}

type EquipmentPricesDTO struct {
	PickaxePrice     int64 `json:"pickaxe_price"`
	VentilationPrice int64 `json:"ventilation_price"`
	TrolleysPrice    int64 `json:"trolleys_price"`
}

func NewEquipmentPricesDTO(
	pickaxePrice int64,
	ventilationPrice int64,
	trolleysPrice int64,
) EquipmentPricesDTO {
	return EquipmentPricesDTO{
		PickaxePrice:     pickaxePrice,
		VentilationPrice: ventilationPrice,
		TrolleysPrice:    trolleysPrice,
	}
}
