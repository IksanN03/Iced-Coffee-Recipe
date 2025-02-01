package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model
	SKU          string         `json:"sku"`
	NumberOfCups int            `json:"number_of_cups"`
	Ingredients  datatypes.JSON `json:"ingredients"` // Using GORM's datatypes.JSON
	COGS         float64        `json:"cogs"`
}

type Measurement struct {
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"` // g, ml, pcs
}

type RecipeInput struct {
	NumberOfCups int                    `json:"number_of_cups"`
	Ingredients  map[string]Measurement `json:"ingredients"`
}
