package company

import "errors"

var ErrUnknokwnMinerType = errors.New("unknown miner type")
var ErrUnknokwnEquipmentType = errors.New("unknown equipment type")
var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrNotAllEquipmentPurchased = errors.New("not all equipment purchased")
