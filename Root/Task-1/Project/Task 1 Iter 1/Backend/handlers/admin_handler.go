package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samin-craftsmen/gin-project/models"
	"github.com/samin-craftsmen/gin-project/utils"
)

// ---------- Admin: Bulk Opt-Out Selected Meals For Everyone ----------
func AdminOptOut(c *gin.Context) {

	role := c.GetString("role")

	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	var request struct {
		Date  string            `json:"date"`
		Meals []models.MealType `json:"meals"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if request.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date required"})
		return
	}

	if len(request.Meals) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one meal type required"})
		return
	}

	users, err := utils.LoadUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load users"})
		return
	}

	data, err := utils.LoadParticipation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load participation"})
		return
	}

	// Default full meal list (implicit opt-in base)
	defaultMeals := []models.MealType{
		"Lunch",
		"Snacks",
		"Iftar",
		"Event Dinner",
		"Optional Dinner",
	}

	updatedCount := 0

	for _, user := range users {

		found := false

		for i, entry := range data {

			if entry.Username == user.Username && entry.Date == request.Date {

				removeSet := make(map[models.MealType]bool)
				for _, m := range request.Meals {
					removeSet[m] = true
				}

				var updatedMeals []models.MealType
				for _, existingMeal := range entry.Meals {
					if !removeSet[existingMeal] {
						updatedMeals = append(updatedMeals, existingMeal)
					}
				}

				data[i].Meals = updatedMeals
				found = true
				updatedCount++
				break
			}
		}

		// If no record → start from default & remove
		if !found {

			removeSet := make(map[models.MealType]bool)
			for _, m := range request.Meals {
				removeSet[m] = true
			}

			var remainingMeals []models.MealType
			for _, m := range defaultMeals {
				if !removeSet[m] {
					remainingMeals = append(remainingMeals, m)
				}
			}

			data = append(data, models.Participation{
				Username: user.Username,
				Date:     request.Date,
				Meals:    remainingMeals,
			})

			updatedCount++
		}
	}

	if err := utils.SaveParticipation(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save participation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Selected meals opted out for everyone",
		"updated_count": updatedCount,
	})
}

// ---------- Admin: Bulk Opt-In Selected Meals For Everyone ----------
func AdminOptIn(c *gin.Context) {

	role := c.GetString("role")

	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	var request struct {
		Date  string            `json:"date"`
		Meals []models.MealType `json:"meals"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if request.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date required"})
		return
	}

	if len(request.Meals) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one meal type required"})
		return
	}

	users, err := utils.LoadUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load users"})
		return
	}

	data, err := utils.LoadParticipation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load participation"})
		return
	}

	updatedCount := 0

	for _, user := range users {

		found := false

		for i, entry := range data {

			if entry.Username == user.Username && entry.Date == request.Date {

				// Merge meals (avoid duplicates)
				existing := make(map[models.MealType]bool)
				for _, m := range entry.Meals {
					existing[m] = true
				}

				for _, newMeal := range request.Meals {
					if !existing[newMeal] {
						data[i].Meals = append(data[i].Meals, newMeal)
					}
				}

				found = true
				updatedCount++
				break
			}
		}

		// If no record → create with requested meals only
		if !found {
			data = append(data, models.Participation{
				Username: user.Username,
				Date:     request.Date,
				Meals:    request.Meals,
			})
			updatedCount++
		}
	}

	if err := utils.SaveParticipation(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save participation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Selected meals opted in for everyone",
		"updated_count": updatedCount,
	})
}

// -------------------- Admin checks participation based on teams --------------------
func GetTeamMealCountsByDate(c *gin.Context) {

	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	date := c.Param("date")

	users, err := utils.LoadUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load users"})
		return
	}

	participation, err := utils.LoadParticipation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load participation"})
		return
	}

	// Default meals
	defaultMeals := []models.MealType{
		"Lunch",
		"Snacks",
		"Iftar",
		"Event Dinner",
		"Optional Dinner",
	}

	// Create result grouped by team
	result := make(map[string][]gin.H)

	for _, user := range users {

		found := false
		var userMeals []models.MealType

		for _, entry := range participation {
			if entry.Username == user.Username && entry.Date == date {
				userMeals = entry.Meals
				found = true
				break
			}
		}

		if !found {
			userMeals = defaultMeals
		}

		result[user.Team] = append(result[user.Team], gin.H{
			"username": user.Username,
			"meals":    userMeals,
		})
	}

	c.JSON(http.StatusOK, result)
}

// -------------------- Admin: Set Special Day (Holiday, Celebration, etc.) --------------------
func SetSpecialDay(c *gin.Context) {

	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	var req models.DayControl

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ------------------ Save Day Control ------------------

	data, err := utils.LoadDayControls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load day controls"})
		return
	}

	updated := false
	for i, d := range data {
		if d.Date == req.Date {
			data[i] = req
			updated = true
			break
		}
	}

	if !updated {
		data = append(data, req)
	}

	// ------------------ If OFFICE_CLOSED ------------------

	if req.Type == "OFFICE_CLOSED" {

		participationData, err := utils.LoadParticipation()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load participation"})
			return
		}

		// Remove all participation entries for that date
		newParticipation := []models.Participation{}

		for _, p := range participationData {
			if p.Date != req.Date {
				newParticipation = append(newParticipation, p)
			}
		}

		// Load all users
		users, err := utils.LoadUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load users"})
			return
		}

		// Add forced opt-out (meals: null) for everyone
		for _, user := range users {
			newParticipation = append(newParticipation, models.Participation{
				Username: user.Username,
				Date:     req.Date,
				Meals:    nil, // 🔥 means opted out of all meals
			})
		}

		err = utils.SaveParticipation(newParticipation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save participation"})
			return
		}
	}

	err = utils.SaveDayControls(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save day controls"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Day control set successfully"})
}

// -------------------- Admin: Get Special Day Status --------------------
func GetDayStatus(c *gin.Context) {

	date := c.Param("date")

	data, _ := utils.LoadDayControls()

	for _, d := range data {
		if d.Date == date {
			c.JSON(http.StatusOK, d)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"date": date,
		"type": "normal_day",
	})
}
