package handlers

import (
	"net/http"
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

	if role != "admin" && role != "teamLead" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	var request struct {
		Username string            `json:"username"`
		Meals    []models.MealType `json:"meals"`
		Date     string            `json:"date"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if request.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date required"})
		return
	}

	data, _ := utils.LoadParticipation()

	found := false
	for i, entry := range data {
		if entry.Username == request.Username && entry.Date == request.Date {
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

	c.JSON(http.StatusOK, gin.H{"message": "Override successful"})
}

// ---------- Employee Update (Cutoff Applied) ----------
func UpdateMealSelection(c *gin.Context) {
	if !utils.IsBeforeCutoff() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cutoff time passed for tomorrow"})
		return
	}

	username := c.GetString("username")
	role := c.GetString("role")

	if role != "employee" && role != "teamlead" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	var request struct {
		Meals []models.MealType `json:"meals"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	data, _ := utils.LoadParticipation()

	updated := false
	for i, entry := range data {
		if entry.Username == username && entry.Date == tomorrow {
			data[i].Meals = request.Meals
			updated = true
			break
		}
	}

	if !updated {
		data = append(data, models.Participation{
			Username: username,
			Date:     tomorrow,
			Meals:    request.Meals,
		})
	}

	utils.SaveParticipation(data)
	c.JSON(http.StatusOK, gin.H{"message": "Updated for tomorrow"})
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

	tomorrow := time.Now().Format("2006-01-02")

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
func GetUser(c *gin.Context) {
	username := c.GetString("username")
	role := c.GetString("role")
	team := c.GetString("team")

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"role":     role,
		"team":     team,
	})
}
