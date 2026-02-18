package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/samin-craftsmen/gin-project/middleware"
	"github.com/samin-craftsmen/gin-project/models"
	"github.com/samin-craftsmen/gin-project/utils"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ================= PUBLIC LOGIN =================
	r.POST("/login", func(c *gin.Context) {
		var loginUser models.User
		if err := c.BindJSON(&loginUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		users, _ := utils.LoadUsers()

		for _, user := range users {
			if user.Username == loginUser.Username && user.Password == loginUser.Password {
				token, err := utils.GenerateToken(user.Username, user.Role, user.Team)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Token error"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"token": token})
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	})

	// ================= PROTECTED ROUTES =================
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/me", func(c *gin.Context) {
			username := c.GetString("username")
			role := c.GetString("role")
			team := c.GetString("team")

			c.JSON(http.StatusOK, gin.H{
				"username": username,
				"role":     role,
				"team":     team,
			})
		})

		// ---------- Get Today's Meals  ----------
		authorized.GET("/meals/today", func(c *gin.Context) {
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
		})

		// ---------- Get Tommorow's Meals  ----------
		authorized.GET("/meals/tomorrow", func(c *gin.Context) {
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
		})

		// ---------- Employee Update (Cutoff Applied) ----------
		authorized.POST("/meals/update", func(c *gin.Context) {
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
		})

		// ---------- Override (Team Lead + Admin, Any Date) ----------
		authorized.POST("/meals/override", func(c *gin.Context) {
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
		})

		// ---------- Headcount By Date (Admin Only) ----------
		authorized.GET("/meals/headcount/:date", func(c *gin.Context) {
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
		})

		authorized.GET("/meals/items/:date", func(c *gin.Context) {

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
		})

		authorized.POST("/meals/items/update", func(c *gin.Context) {

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
		})

		// ---------- Team Participation View(Team Lead) ----------

		authorized.GET("/teams/meals/today", func(c *gin.Context) {
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
		})

		// -------------------- Admin checks participation based on teams --------------------
		authorized.GET("/admin/teams/meals/:date", func(c *gin.Context) {

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
		})

		// ---------- Team Lead: Bulk Opt-Out Selected Meals For Entire Team ----------
		authorized.POST("/teams/meals/optout", func(c *gin.Context) {

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
		})

		// ---------- Team Lead: Bulk Opt-In Selected Meals For Entire Team ----------
		authorized.POST("/teams/meals/optin", func(c *gin.Context) {

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
		})

		// ---------- Admin: Bulk Opt-In Selected Meals For Everyone ----------
		authorized.POST("/admin/meals/optin", func(c *gin.Context) {

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
		})

		// ---------- Admin: Bulk Opt-Out Selected Meals For Everyone ----------
		authorized.POST("/admin/meals/optout", func(c *gin.Context) {

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
		})

	}

	r.Run(":8080")
}
