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
	"gorm.io/gorm"
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
	ingredients := make(map[string]models.Measurement)
	json.Unmarshal(ingredientsJSON, &ingredients)

	recipe.COGS, err = calculateCOGSWithSubquery(ingredients, input.NumberOfCups, database.DB)
	if err != nil {
		helpers.NewAPIResponse(c, nil, err, "db", 0, "Failed to calculate COGS")
		return
	}
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
	query.Offset(offset).Limit(limit).Order("id desc").Find(&recipes)

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
	ingredients := make(map[string]models.Measurement)
	json.Unmarshal(ingredientsJSON, &ingredients)
	recipe.NumberOfCups = input.NumberOfCups
	recipe.Ingredients = datatypes.JSON(ingredientsJSON)
	recipe.COGS, err = calculateCOGSWithSubquery(ingredients, input.NumberOfCups, database.DB)
	if err != nil {
		helpers.NewAPIResponse(c, nil, err, "db", 0, "Failed to calculate COGS")
		return
	}

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

func calculateCOGSWithSubquery(ingredients map[string]models.Measurement, numberOfCups int, db *gorm.DB) (float64, error) {
	var totalCOGS float64

	// Process each ingredient with individual optimized query
	for itemName, measurement := range ingredients {
		var itemCost float64
		unit := strings.ToLower(measurement.Unit)

		query := db.Model(&models.Inventory{}).
			Select("CASE ? "+
				"WHEN 'g' THEN (? * price_per_qty) / (quantity * 1000) "+
				"WHEN 'kg' THEN ? * price_per_qty "+
				"WHEN 'ml' THEN (? * price_per_qty) / (quantity * 1000) "+
				"WHEN 'liter' THEN ? * price_per_qty "+
				"WHEN 'pcs' THEN CASE WHEN quantity > 0 THEN ? * (price_per_qty / quantity) ELSE 0 END "+
				"ELSE 0 END as cost",
				unit, measurement.Amount, measurement.Amount,
				measurement.Amount, measurement.Amount, measurement.Amount).
			Where("item_name = ?", itemName)

		if err := query.Scan(&itemCost).Error; err != nil {
			return 0, fmt.Errorf("item %s not found or invalid", itemName)
		}

		if itemCost == 0 {
			return 0, fmt.Errorf("invalid unit or calculation for %s", itemName)
		}

		totalCOGS += itemCost * float64(numberOfCups)
	}

	return totalCOGS, nil
}
