package equipment

import "fmt"

type EquipmentType string

const (
	EquipmentTypePickaxe     EquipmentType = "pickaxe"
	EquipmentTypeVentilation EquipmentType = "ventilation"
	EquipmentTypeTrolleys    EquipmentType = "trolleys"
)

const (
	EquipmentPickaxeCost     int64 = 3000
	EquipmentVentilationCost int64 = 15000
	EquipmentTrolleysCost    int64 = 50000
)

type Equipment struct {
	// Кирки
	pickaxe bool

	// Вентиляция
	ventilation bool

	// Вагонетки
	trolleys bool
}

func NewEquipment() Equipment {
	return Equipment{}
}

func (e *Equipment) BuyPickaxe() {
	e.pickaxe = true

	fmt.Println("pickaxe is purchased")
}

func (e *Equipment) BuyVentilation() {
	e.ventilation = true

	fmt.Println("ventilation is purchased")
}

func (e *Equipment) BuyTrolleys() {
	e.trolleys = true

	fmt.Println("trolleys is purchased")
}

func (e *Equipment) PickaxesPurchased() bool {
	return e.pickaxe
}

func (e *Equipment) VentilationPurchased() bool {
	return e.ventilation
}

func (e *Equipment) TrolleysPurchased() bool {
	return e.trolleys
}

func (e *Equipment) AllBought() bool {
	return e.PickaxesPurchased() && e.VentilationPurchased() && e.TrolleysPurchased()
}
