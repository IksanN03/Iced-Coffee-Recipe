package handler

import (
	"be-test/database"
	"be-test/helpers"
	"be-test/models"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func AddRecipe(c *gin.Context) {
	var input models.RecipeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.NewAPIResponse(c, nil, err, "binding", 0, "Invalid input")
		return
	}

	// Convert ingredients map to JSON
	ingredientsJSON, err := json.Marshal(input.Ingredients)
	if err != nil {
		helpers.NewAPIResponse(c, nil, err, "json", 0, "Failed to marshal ingredients")
		return
	}

	recipe := models.Recipe{
		NumberOfCups: input.NumberOfCups,
		Ingredients:  datatypes.JSON(ingredientsJSON),
	}

	// Calculate COGS
	var inventoryItems []models.Inventory
	if err := database.DB.Find(&inventoryItems).Error; err != nil {
		helpers.NewAPIResponse(c, nil, err, "db", 0, "Failed to fetch inventory")
		return
	}

	var totalCOGS float64
	ingredients := make(map[string]models.Measurement)
	json.Unmarshal(ingredientsJSON, &ingredients)

	for itemName, measurement := range ingredients {
		// Find matching inventory item
		var inventoryItem models.Inventory
		found := false
		for _, item := range inventoryItems {
			if item.ItemName == itemName {
				inventoryItem = item
				found = true
				break
			}
		}

		if !found {
			helpers.NewAPIResponse(c, nil, fmt.Errorf("item not found"), "validation", 0, fmt.Sprintf("Item %s not found in inventory", itemName))
			return
		}

		// Convert to base unit and calculate cost
		var baseAmount float64
		switch strings.ToLower(measurement.Unit) {
		case "g":
			baseAmount = measurement.Amount / 1000 // to kg
		case "kg":
			baseAmount = measurement.Amount // already in kg
		case "ml":
			baseAmount = measurement.Amount / 1000 // to l
		case "liter":
			baseAmount = measurement.Amount // already in l
		case "pcs":
			baseAmount = measurement.Amount
		default:
			helpers.NewAPIResponse(c, nil, fmt.Errorf("invalid unit"), "validation", 0, "Invalid unit")
			return
		}

		totalCOGS += baseAmount * inventoryItem.PricePerQty * float64(input.NumberOfCups)
	}

	recipe.COGS = totalCOGS

	// Generate SKU
	currentTime := time.Now()
	var lastRecipe models.Recipe
	database.DB.Order("created_at desc").First(&lastRecipe)

	sequence := 1
	if !lastRecipe.CreatedAt.IsZero() && lastRecipe.CreatedAt.Format("20060102") == currentTime.Format("20060102") {
		fmt.Sscanf(lastRecipe.SKU, "IC-%8s-%03d", new(string), &sequence)
		sequence++
	}

	recipe.SKU = fmt.Sprintf("IC-%s-%03d", currentTime.Format("20060102"), sequence)

	if err := database.DB.Create(&recipe).Error; err != nil {
		helpers.NewAPIResponse(c, nil, err, "db", 0, "Failed to save recipe")
		return
	}

	helpers.NewAPIResponse(c, gin.H{
		"sku":            recipe.SKU,
		"cogs":           recipe.COGS,
		"number_of_cups": recipe.NumberOfCups,
	}, nil, "", 0, "Recipe added successfully")
}

func GetRecipe(c *gin.Context) {
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
	var recipes []models.Recipe
	var totalItems int64

	query := database.DB.Model(&models.Recipe{})
	if search != "" {
		query = query.Where("sku ILIKE ?", "%"+search+"%")
	}

	query.Count(&totalItems)
	query.Offset(offset).Limit(limit).Find(&recipes)

	helpers.NewAPIResponse(c, gin.H{
		"page":        page,
		"limit":       limit,
		"total_items": totalItems,
		"total_pages": (totalItems + int64(limit) - 1) / int64(limit),
		"recipes":     recipes,
	}, nil, "", 0, "Recipes retrieved successfully")
}

func UpdateRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	var input models.RecipeInput

	if err := database.DB.First(&recipe, id).Error; err != nil {
		helpers.NewAPIResponse(c, nil, err, "db", 0, "Recipe not found")
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.NewAPIResponse(c, nil, err, "binding", 0, "Invalid input")
		return
	}

	// Convert ingredients map to JSON
	ingredientsJSON, err := json.Marshal(input.Ingredients)
	if err != nil {
		helpers.NewAPIResponse(c, nil, err, "json", 0, "Failed to marshal ingredients")
		return
	}

	// Recalculate COGS
	var inventoryItems []models.Inventory
	if err := database.DB.Find(&inventoryItems).Error; err != nil {
		helpers.NewAPIResponse(c, nil, err, "db", 0, "Failed to fetch inventory")
		return
	}

	var totalCOGS float64
	ingredients := make(map[string]models.Measurement)
	json.Unmarshal(ingredientsJSON, &ingredients)

	for itemName, measurement := range ingredients {
		var inventoryItem models.Inventory
		found := false
		for _, item := range inventoryItems {
			if item.ItemName == itemName {
				inventoryItem = item
				found = true
				break
			}
		}

		if !found {
			helpers.NewAPIResponse(c, nil, fmt.Errorf("item not found"), "validation", 0, fmt.Sprintf("Item %s not found in inventory", itemName))
			return
		}

		var baseAmount float64
		switch strings.ToLower(measurement.Unit) {
		case "g":
			baseAmount = measurement.Amount / 1000
		case "kg":
			baseAmount = measurement.Amount
		case "ml":
			baseAmount = measurement.Amount / 1000
		case "liter":
			baseAmount = measurement.Amount
		case "pcs":
			baseAmount = measurement.Amount
		default:
			helpers.NewAPIResponse(c, nil, fmt.Errorf("invalid unit"), "validation", 0, "Invalid unit")
			return
		}

		totalCOGS += baseAmount * inventoryItem.PricePerQty * float64(input.NumberOfCups)
	}

	recipe.NumberOfCups = input.NumberOfCups
	recipe.Ingredients = datatypes.JSON(ingredientsJSON)
	recipe.COGS = totalCOGS

	if err := database.DB.Save(&recipe).Error; err != nil {
		helpers.NewAPIResponse(c, nil, err, "db", 0, "Failed to update recipe")
		return
	}

	helpers.NewAPIResponse(c, gin.H{
		"sku":            recipe.SKU,
		"cogs":           recipe.COGS,
		"number_of_cups": recipe.NumberOfCups,
	}, nil, "", 0, "Recipe updated successfully")
}
