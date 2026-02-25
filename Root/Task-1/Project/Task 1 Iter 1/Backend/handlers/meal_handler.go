package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samin-craftsmen/gin-project/models"
	"github.com/samin-craftsmen/gin-project/utils"
)

// ---------- Meals: Update Meals Admin only ----------
func UpdateMeals(c *gin.Context) {

	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	var request struct {
		Date  string                       `json:"date"`
		Items map[models.MealType][]string `json:"items"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	data, _ := utils.LoadMealItems()

	found := false

	for i, entry := range data {
		if entry.Date == request.Date {
			data[i].Items = request.Items
			found = true
			break
		}
	}

	if !found {
		data = append(data, models.MealItems{
			Date:  request.Date,
			Items: request.Items,
		})
	}

	utils.SaveMealItems(data)

	c.JSON(http.StatusOK, gin.H{"message": "Meal items updated"})
}

// ---------- Headcount By Date (Admin Only) ----------
func Headcount(c *gin.Context) {
	role := c.GetString("role")

	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	date := c.Param("date")

	data, _ := utils.LoadParticipation()
	users, _ := utils.LoadUsers()

	count := map[string]int{
		"Lunch":           0,
		"Snacks":          0,
		"Iftar":           0,
		"Event Dinner":    0,
		"Optional Dinner": 0,
	}

	for _, user := range users {

		found := false

		for _, entry := range data {
			if entry.Username == user.Username && entry.Date == date {
				for _, meal := range entry.Meals {
					count[string(meal)]++
				}
				found = true
				break
			}
		}

		if !found {
			count["Lunch"]++
			count["Snacks"]++
			count["Iftar"]++
			count["Event Dinner"]++
			count["Optional Dinner"]++
		}
	}

	c.JSON(http.StatusOK, count)
}

// ---------- Daily Announcement Draft (Admin only) ----------
func DailyAnnouncement(c *gin.Context) {
	role := c.GetString("role")

	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	date := c.Param("date")

	participation, _ := utils.LoadParticipation()
	users, _ := utils.LoadUsers()
	dayControls, _ := utils.LoadDayControls()

	count := map[string]int{
		"Lunch":           0,
		"Snacks":          0,
		"Iftar":           0,
		"Event Dinner":    0,
		"Optional Dinner": 0,
	}

	for _, user := range users {
		found := false

		for _, entry := range participation {
			if entry.Username == user.Username && entry.Date == date {
				if entry.Meals != nil {
					for _, meal := range entry.Meals {
						count[string(meal)]++
					}
				}
				found = true
				break
			}
		}

		if !found {
			// default: no entry means user is counted for all meals
			count["Lunch"]++
			count["Snacks"]++
			count["Iftar"]++
			count["Event Dinner"]++
			count["Optional Dinner"]++
		}
	}

	// Build announcement text
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Announcement for %s:\n\n", date))
	b.WriteString(fmt.Sprintf("Lunch: %d\n", count["Lunch"]))
	b.WriteString(fmt.Sprintf("Snacks: %d\n", count["Snacks"]))
	b.WriteString(fmt.Sprintf("Iftar: %d\n", count["Iftar"]))
	b.WriteString(fmt.Sprintf("Event Dinner: %d\n", count["Event Dinner"]))
	b.WriteString(fmt.Sprintf("Optional Dinner: %d\n", count["Optional Dinner"]))

	// Check for special day control
	for _, dc := range dayControls {
		if dc.Date == date {
			t := strings.ToUpper(dc.Type)
			switch t {
			case "OFFICE_CLOSED", "OFFICE-CLOSED":
				b.WriteString("\nNote: Office closed on this date.\n")
			case "GOVERNMENT_HOLIDAY", "GOVERNMENT-HOLIDAY":
				b.WriteString("\nNote: Government holiday.\n")
			case "SPECIAL_CELEBRATION", "SPECIAL-CELEBRATION":
				b.WriteString("\nNote: Special celebration.\n")
			default:
				b.WriteString(fmt.Sprintf("\nNote: %s\n", dc.Type))
			}

			if dc.Note != nil && *dc.Note != "" {
				b.WriteString(fmt.Sprintf("Details: %s\n", *dc.Note))
			}
			break
		}
	}

	c.String(http.StatusOK, b.String())
}

// ---------- Get Meal Items By Date  ----------
func GetMealItemsByDate(c *gin.Context) {

	date := c.Param("date")

	data, _ := utils.LoadMealItems()

	for _, entry := range data {
		if entry.Date == date {
			c.JSON(http.StatusOK, entry)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"date":  date,
		"items": map[string][]string{},
	})
}

// ---------- Override (Team Lead + Admin, Any Date) ----------
func OverrideMealSelection(c *gin.Context) {
	role := c.GetString("role")
	username := c.GetString("username")
	ipAddress := c.ClientIP()

	if role != "admin" && role != "teamLead" {
		utils.LogActionWithError("MEAL_SELECTION_OVERRIDE", username, role, "", "", ipAddress, "Unauthorized access attempt")
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	var request struct {
		Username string            `json:"username"`
		Meals    []models.MealType `json:"meals"`
		Date     string            `json:"date"`
	}

	if err := c.BindJSON(&request); err != nil {
		utils.LogActionWithError("MEAL_SELECTION_OVERRIDE", username, role, request.Username, request.Date, ipAddress, "Invalid input")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if request.Date == "" {
		utils.LogActionWithError("MEAL_SELECTION_OVERRIDE", username, role, request.Username, "", ipAddress, "Date required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date required"})
		return
	}

	data, _ := utils.LoadParticipation()

	var previousMeals []models.MealType
	found := false
	for i, entry := range data {
		if entry.Username == request.Username && entry.Date == request.Date {
			previousMeals = entry.Meals
			data[i].Meals = request.Meals
			found = true
			break
		}
	}

	if !found {
		data = append(data, models.Participation{
			Username: request.Username,
			Date:     request.Date,
			Meals:    request.Meals,
		})
	}

	utils.SaveParticipation(data)

	// Log the override action
	changeDetails := map[string]interface{}{
		"previous_meals": previousMeals,
		"new_meals":      request.Meals,
	}
	utils.LogAction("MEAL_SELECTION_OVERRIDE", username, role, request.Username, request.Date, ipAddress, changeDetails, "success")

	c.JSON(http.StatusOK, gin.H{"message": "Override successful"})
}

// ---------- Employee Update (Up to 7 Days Ahead) ----------
func UpdateMealSelection(c *gin.Context) {
	username := c.GetString("username")
	role := c.GetString("role")
	ipAddress := c.ClientIP()

	var req struct {
		Date  string            `json:"date"`
		Meals []models.MealType `json:"meals"`
	}

	if err := c.BindJSON(&req); err != nil {
		utils.LogActionWithError("MEAL_SELECTION_UPDATE", username, role, username, req.Date, ipAddress, "Invalid input")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate date is within 7 days from today
	if !isValidFutureDate(req.Date) {
		utils.LogActionWithError("MEAL_SELECTION_UPDATE", username, role, username, req.Date, ipAddress, "Date outside allowed range")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date must be within 7 days from today"})
		return
	}

	// Load current participation data
	participation, err := utils.LoadParticipation()
	if err != nil {
		utils.LogActionWithError("MEAL_SELECTION_UPDATE", username, role, username, req.Date, ipAddress, "Failed to load participation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load participation"})
		return
	}

	// Track previous meals for audit
	var previousMeals []models.MealType
	found := false
	for i, p := range participation {
		if p.Username == username && p.Date == req.Date {
			previousMeals = p.Meals
			participation[i].Meals = req.Meals
			found = true
			break
		}
	}

	// If not found, create new entry
	if !found {
		participation = append(participation, models.Participation{
			Username: username,
			Date:     req.Date,
			Meals:    req.Meals,
		})
	}

	// Save updated participation
	if err := utils.SaveParticipation(participation); err != nil {
		utils.LogActionWithError("MEAL_SELECTION_UPDATE", username, role, username, req.Date, ipAddress, "Failed to save participation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save participation"})
		return
	}

	// Log the update action
	changeDetails := map[string]interface{}{
		"previous_meals": previousMeals,
		"new_meals":      req.Meals,
	}
	utils.LogAction("MEAL_SELECTION_UPDATE", username, role, username, req.Date, ipAddress, changeDetails, "success")

	//  Broadcast update to all connected admin clients
	workLocations, _ := utils.LoadWorkLocations()
	users, _ := utils.LoadUsers()

	totalParticipants := 0
	officeCount := 0
	wfhCount := 0
	optedOut := 0

	mealCount := map[string]int{
		"Lunch":           0,
		"Snacks":          0,
		"Iftar":           0,
		"Event Dinner":    0,
		"Optional Dinner": 0,
	}

	for _, user := range users {
		location := "Office"
		for _, wl := range workLocations {
			if wl.Username == user.Username && wl.Date == req.Date {
				location = wl.Location
				break
			}
		}

		for _, p := range participation {
			if p.Username == user.Username && p.Date == req.Date {
				if p.Meals == nil || len(p.Meals) == 0 {
					optedOut++
				} else {
					totalParticipants++
					if location == "WFH" {
						wfhCount++
					} else {
						officeCount++
					}
					for _, meal := range p.Meals {
						mealCount[string(meal)]++
					}
				}
				break
			}
		}
	}

	headcount := gin.H{
		"total_participants": totalParticipants,
		"office":             officeCount,
		"wfh":                wfhCount,
		"opted_out":          optedOut,
		"by_meal":            mealCount,
	}

	BroadcastHeadcountUpdate(req.Date, headcount)

	c.JSON(http.StatusOK, gin.H{
		"message": "Meal selection updated",
		"date":    req.Date,
		"meals":   req.Meals,
	})
}

// Helper function to validate future dates (up to 7 days)
func isValidFutureDate(dateStr string) bool {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	today := time.Now().Truncate(24 * time.Hour)
	sevenDaysFromNow := today.AddDate(0, 0, 7)

	return !date.Before(today) && date.Before(sevenDaysFromNow.AddDate(0, 0, 1))
}

// ---------- Get Tommorow's Meals  ----------
func GetTomorrowMeals(c *gin.Context) {
	username := c.GetString("username")

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	data, _ := utils.LoadParticipation()

	for _, entry := range data {
		if entry.Username == username && entry.Date == tomorrow {
			c.JSON(http.StatusOK, entry)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"date":     tomorrow,
		"meals": []string{
			"Lunch",
			"Snacks",
			"Iftar",
			"Event Dinner",
			"Optional Dinner",
		},
	})
}

func GetTodayMeals(c *gin.Context) {
	username := c.GetString("username")
	today := time.Now().Format("2006-01-02")

	data, _ := utils.LoadParticipation()

	for _, entry := range data {
		if entry.Username == username && entry.Date == today {
			// User has specific meals selected
			c.JSON(http.StatusOK, entry.Meals)
			return
		}
	}

	// Default: all meals if no entry found
	c.JSON(http.StatusOK, []string{
		"Lunch",
		"Snacks",
		"Iftar",
		"Event Dinner",
		"Optional Dinner",
	})
}
