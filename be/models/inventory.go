package models

import "gorm.io/gorm"

// Inventory represents an inventory item
type Inventory struct {
	gorm.Model
	ItemName    string  `json:"item_name"`
	Quantity    float64 `json:"quantity"`
	Uom         string  `json:"uom"`
	PricePerQty float64 `json:"price_per_qty"`
}
