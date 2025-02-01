package handler

import (
	"be-test/database"
	"be-test/helpers"
	"be-test/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetInventory(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	search := c.DefaultQuery("search", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	var inventory []models.Inventory
	var totalItems int64

	query := database.DB.Model(&models.Inventory{})
	if search != "" {
		query = query.Where("item_name ILIKE ?", "%"+search+"%")
	}

	query.Count(&totalItems)
	query.Offset(offset).Limit(limit).Find(&inventory)

	helpers.NewAPIResponse(c, gin.H{
		"page":        page,
		"limit":       limit,
		"total_items": totalItems,
		"total_pages": (totalItems + int64(limit) - 1) / int64(limit),
		"inventory":   inventory,
	}, nil, "", 0, "Inventory retrieved successfully")
}

func AddInventory(c *gin.Context) {
	var input models.Inventory
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.NewAPIResponse(c, nil, err, "binding", 0, "Invalid input")
		return
	}

	if err := database.DB.Create(&input).Error; err != nil {
		helpers.NewAPIResponse(c, nil, err, "db", 0, "Failed to create inventory item")
		return
	}

	helpers.NewAPIResponse(c, gin.H{"inventory": input}, nil, "", 0, "Inventory item added successfully")
}

func UpdateInventory(c *gin.Context) {
	var input models.Inventory
	id := c.Param("id")

	if err := database.DB.First(&input, id).Error; err != nil {
		helpers.NewAPIResponse(c, nil, err, "db", 0, "Inventory item not found")
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.NewAPIResponse(c, nil, err, "binding", 0, "Invalid input")
		return
	}

	if err := database.DB.Save(&input).Error; err != nil {
		helpers.NewAPIResponse(c, nil, err, "db", 0, "Failed to update inventory item")
		return
	}

	helpers.NewAPIResponse(c, gin.H{"inventory": input}, nil, "", 0, "Inventory item updated successfully")
}

func DeleteInventory(c *gin.Context) {
	id := c.Param("id")
	var inventory models.Inventory

	if err := database.DB.First(&inventory, id).Error; err != nil {
		helpers.NewAPIResponse(c, nil, err, "db", 0, "Inventory item not found")
		return
	}

	if err := database.DB.Delete(&inventory).Error; err != nil {
		helpers.NewAPIResponse(c, nil, err, "db", 0, "Failed to delete inventory item")
		return
	}

	helpers.NewAPIResponse(c, nil, nil, "", 0, "Inventory item deleted successfully")
}
