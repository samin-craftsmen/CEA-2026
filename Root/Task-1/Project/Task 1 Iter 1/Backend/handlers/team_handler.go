package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samin-craftsmen/gin-project/models"
	"github.com/samin-craftsmen/gin-project/utils"
)

// ---------- Team Participation View(Team Lead) ----------

func GetTodayTeamMeals(c *gin.Context) {
	role := c.GetString("role")
	teamLead_team := c.GetString("team")
	if role != "teamLead" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	today := time.Now().Format("2006-01-02")

	users, _ := utils.LoadUsers()
	participation, _ := utils.LoadParticipation()

	result := []models.Participation{}

	for _, user := range users {

		if user.Team != teamLead_team {
			continue
		}

		found := false

		for _, entry := range participation {
			if entry.Username == user.Username && entry.Date == today {
				result = append(result, entry)
				found = true
				break
			}
		}

		if !found {
			// Default opt-in
			result = append(result, models.Participation{
				Username: user.Username,
				Date:     today,
				Meals: []models.MealType{
					"Lunch",
					"Snacks",
					"Iftar",
					"Event Dinner",
					"Optional Dinner",
				},
			})
		}
	}

	c.JSON(http.StatusOK, result)
}

// ---------- Team Lead: Bulk Opt-Out Selected Meals For Entire Team ----------
func TeamBulkOptOut(c *gin.Context) {

	role := c.GetString("role")
	team := c.GetString("team")

	if role != "teamLead" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Team Lead only"})
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

		if user.Team != team {
			continue
		}

		found := false

		// Try to find existing participation record
		for i, entry := range data {

			if entry.Username == user.Username && entry.Date == request.Date {

				// Remove selected meals from existing
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

		// 🔥 If no record exists → start from default meals
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
		"message":       "Selected meals opted out successfully",
		"updated_count": updatedCount,
	})
}

// ---------- Team Lead: Bulk Opt-In Selected Meals For Entire Team ----------
func TeamBulkOptIn(c *gin.Context) {

	role := c.GetString("role")
	team := c.GetString("team")

	if role != "teamLead" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Team Lead only"})
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

		if user.Team != team {
			continue
		}

		found := false

		for i, entry := range data {
			if entry.Username == user.Username && entry.Date == request.Date {

				// Merge meals (avoid duplicates)
				existingMeals := make(map[models.MealType]bool)
				for _, m := range entry.Meals {
					existingMeals[m] = true
				}

				for _, newMeal := range request.Meals {
					if !existingMeals[newMeal] {
						data[i].Meals = append(data[i].Meals, newMeal)
					}
				}

				found = true
				updatedCount++
				break
			}
		}

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
		"message":       "Selected meals opted in successfully",
		"updated_count": updatedCount,
	})
}

func UpdateTeamMemberWorkLocation(c *gin.Context) {
	role := c.GetString("role")
	leadTeam := c.GetString("team")

	if role != "teamLead" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only team leads can modify team work locations"})
		return
	}

	var req struct {
		Username string `json:"username"`
		Date     string `json:"date"`
		Location string `json:"location"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.Username == "" || req.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and date are required"})
		return
	}

	if req.Location != "Office" && req.Location != "WFH" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work location"})
		return
	}

	users, err := utils.LoadUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load users"})
		return
	}

	validMember := false
	for _, u := range users {
		if u.Username == req.Username && u.Team == leadTeam {
			validMember = true
			break
		}
	}

	if !validMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only modify members of your own team"})
		return
	}

	// Update work location
	workData, err := utils.LoadWorkLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load work locations"})
		return
	}

	updated := false
	for i, w := range workData {
		if w.Username == req.Username && w.Date == req.Date {
			workData[i].Location = req.Location
			updated = true
			break
		}
	}

	if !updated {
		workData = append(workData, models.WorkLocation{
			Username: req.Username,
			Date:     req.Date,
			Location: req.Location,
		})
	}

	err = utils.SaveWorkLocations(workData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save work location"})
		return
	}

	// Handle meals
	participationData, err := utils.LoadParticipation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load participation"})
		return
	}

	defaultMeals := []string{"Lunch", "Snacks", "Iftar", "Event Dinner", "Optional Dinner"}
	mealTypes := []models.MealType{}
	for _, m := range defaultMeals {
		mealTypes = append(mealTypes, models.MealType(m))
	}

	newParticipation := []models.Participation{}
	found := false

	for _, p := range participationData {
		if p.Username == req.Username && p.Date == req.Date {
			found = true
			if req.Location == "WFH" {
				newParticipation = append(newParticipation, models.Participation{
					Username: req.Username,
					Date:     req.Date,
					Meals:    nil,
				})
			} else {
				newParticipation = append(newParticipation, models.Participation{
					Username: req.Username,
					Date:     req.Date,
					Meals:    mealTypes,
				})
			}
		} else {
			newParticipation = append(newParticipation, p)
		}
	}

	if !found {
		if req.Location == "WFH" {
			newParticipation = append(newParticipation, models.Participation{
				Username: req.Username,
				Date:     req.Date,
				Meals:    nil,
			})
		} else {
			newParticipation = append(newParticipation, models.Participation{
				Username: req.Username,
				Date:     req.Date,
				Meals:    mealTypes,
			})
		}
	}

	err = utils.SaveParticipation(newParticipation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save participation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work location and meals updated successfully"})
}
func GetTeamMemberWorkLocation(c *gin.Context) {

	// 🔐 Role check
	role := c.GetString("role")
	leadTeam := c.GetString("team")

	if role != "teamLead" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only team leads can access this",
		})
		return
	}

	var req struct {
		Username string `json:"username"`
		Date     string `json:"date"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	if req.Username == "" || req.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username and date are required",
		})
		return
	}

	// 🔍 Verify team membership
	users, err := utils.LoadUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load users",
		})
		return
	}

	isMember := false
	for _, u := range users {
		if u.Username == req.Username && u.Team == leadTeam {
			isMember = true
			break
		}
	}

	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only view members of your own team",
		})
		return
	}

	// 🔎 Load work locations
	workData, err := utils.LoadWorkLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load work locations",
		})
		return
	}

	for _, w := range workData {
		if w.Username == req.Username && w.Date == req.Date {
			c.JSON(http.StatusOK, w)
			return
		}
	}

	// Default if not set
	c.JSON(http.StatusOK, gin.H{
		"username": req.Username,
		"date":     req.Date,
		"location": "Office",
	})
}
