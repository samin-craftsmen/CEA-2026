package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

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

	if req.Type == "office_closed" {

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

// ------------------ Admin: Remove special day -------------------------- //
func RemoveSpecialDay(c *gin.Context) {

	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	date := c.Param("date")

	// ---------------- Remove from day_controls.json ----------------
	data, err := utils.LoadDayControls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load day controls"})
		return
	}

	newDayControls := []models.DayControl{}
	for _, d := range data {
		if d.Date != date {
			newDayControls = append(newDayControls, d)
		}
	}

	err = utils.SaveDayControls(newDayControls)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save day controls"})
		return
	}

	// ---------------- Remove participation entries for that date ----------------
	participationData, err := utils.LoadParticipation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load participation"})
		return
	}

	newParticipation := []models.Participation{}
	for _, p := range participationData {
		if p.Date != date {
			newParticipation = append(newParticipation, p)
		}
	}

	err = utils.SaveParticipation(newParticipation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save participation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Special day removed and participation reset"})
}
func UpdateWorkLocationByAdmin(c *gin.Context) {

	role := c.GetString("role")

	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only admin can modify work locations",
		})
		return
	}

	var req struct {
		Username string `json:"username"`
		Date     string `json:"date"`
		Location string `json:"location"`
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

	if req.Location != "Office" && req.Location != "WFH" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid work location",
		})
		return
	}

	// Verify user exists
	users, err := utils.LoadUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load users",
		})
		return
	}

	userExists := false
	for _, u := range users {
		if u.Username == req.Username {
			userExists = true
			break
		}
	}

	if !userExists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// ------------------ Update Work Location ------------------

	workData, err := utils.LoadWorkLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load work locations",
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save work location",
		})
		return
	}

	// ------------------ If WFH → Opt Out Meals ------------------

	if req.Location == "WFH" {

		participationData, err := utils.LoadParticipation()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to load participation",
			})
			return
		}

		newParticipation := []models.Participation{}
		found := false

		for _, p := range participationData {
			if p.Username == req.Username && p.Date == req.Date {
				newParticipation = append(newParticipation, models.Participation{
					Username: req.Username,
					Date:     req.Date,
					Meals:    nil,
				})
				found = true
			} else {
				newParticipation = append(newParticipation, p)
			}
		}

		if !found {
			newParticipation = append(newParticipation, models.Participation{
				Username: req.Username,
				Date:     req.Date,
				Meals:    nil,
			})
		}

		err = utils.SaveParticipation(newParticipation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save participation",
			})
			return
		}
	} else if req.Location == "Office" {
		// opt back in: remove the null override for that date
		participationData, err := utils.LoadParticipation()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load participation"})
			return
		}

		newParticipation := []models.Participation{}
		for _, p := range participationData {
			if p.Username == req.Username && p.Date == req.Date {
				// Remove override (restore default by not including Meals=nil)
				continue
			}
			newParticipation = append(newParticipation, p)
		}

		err = utils.SaveParticipation(newParticipation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save participation"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Work location updated successfully by admin",
	})
}

func GetWorkLocationByAdmin(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can access this"})
		return
	}

	username := c.Query("username")
	date := c.Query("date")

	if username == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and date are required"})
		return
	}

	//  Optional: verify user exists
	users, err := utils.LoadUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load users"})
		return
	}

	userExists := false
	for _, u := range users {
		if u.Username == username {
			userExists = true
			break
		}
	}

	if !userExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	//  Load work locations
	workData, err := utils.LoadWorkLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load work locations"})
		return
	}

	for _, w := range workData {
		if w.Username == username && w.Date == date {
			c.JSON(http.StatusOK, w)
			return
		}
	}

	// Default if not set
	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"date":     date,
		"location": "Office",
	})
}

// -------------------- Admin: Set Company-wide WFH Range --------------------
func SetCompanyWFH(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	var req struct {
		StartDate string  `json:"start_date"`
		EndDate   string  `json:"end_date"`
		Note      *string `json:"note,omitempty"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.StartDate == "" || req.EndDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date required"})
		return
	}

	start, err1 := time.Parse("2006-01-02", req.StartDate)
	end, err2 := time.Parse("2006-01-02", req.EndDate)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	if start.After(end) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be <= end_date"})
		return
	}

	users, err := utils.LoadUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load users"})
		return
	}

	participationData, err := utils.LoadParticipation()
	if err != nil {
		// initialize empty if file missing
		participationData = []models.Participation{}
	}

	workData, err := utils.LoadWorkLocations()
	if err != nil {
		workData = []models.WorkLocation{}
	}

	dayControls, err := utils.LoadDayControls()
	if err != nil {
		dayControls = []models.DayControl{}
	}

	// Iterate dates in range
	daysAffected := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")

		// Participation: opt everyone out (meals: nil)
		for _, user := range users {
			found := false
			for i, p := range participationData {
				if p.Username == user.Username && p.Date == dateStr {
					participationData[i].Meals = nil
					found = true
					break
				}
			}
			if !found {
				participationData = append(participationData, models.Participation{
					Username: user.Username,
					Date:     dateStr,
					Meals:    nil,
				})
			}
		}

		// Work locations: set everyone to WFH
		for _, user := range users {
			updated := false
			for i, w := range workData {
				if w.Username == user.Username && w.Date == dateStr {
					workData[i].Location = "WFH"
					updated = true
					break
				}
			}
			if !updated {
				workData = append(workData, models.WorkLocation{
					Username: user.Username,
					Date:     dateStr,
					Location: "WFH",
				})
			}
		}

		// Day controls: add/update entry for this date
		updatedDC := false
		for i, dc := range dayControls {
			if dc.Date == dateStr {
				dayControls[i].Type = "COMPANY_WFH"
				dayControls[i].Note = req.Note
				updatedDC = true
				break
			}
		}
		if !updatedDC {
			dayControls = append(dayControls, models.DayControl{
				Date: dateStr,
				Type: "COMPANY_WFH",
				Note: req.Note,
			})
		}

		daysAffected++
	}

	if err := utils.SaveParticipation(participationData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save participation"})
		return
	}

	if err := utils.SaveWorkLocations(workData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save work locations"})
		return
	}

	if err := utils.SaveDayControls(dayControls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save day controls"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Company WFH range applied",
		"days_affected": daysAffected,
	})
}

// ---------- Headcount By Team ----------
// ---------- Headcount By Team ----------
func HeadcountByTeam(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	date := c.Param("date")
	participation, _ := utils.LoadParticipation()
	users, _ := utils.LoadUsers()

	teamMeals := map[string]map[string]int{}

	// Initialize all teams with all meals
	for _, user := range users {
		if _, exists := teamMeals[user.Team]; !exists {
			teamMeals[user.Team] = map[string]int{
				"Lunch":           0,
				"Snacks":          0,
				"Iftar":           0,
				"Event Dinner":    0,
				"Optional Dinner": 0,
			}
		}
	}

	// Track which users have entries for this date
	userHasEntry := make(map[string]bool)

	for _, p := range participation {
		if p.Date == date {
			userHasEntry[p.Username] = true
			// Only count if meals is not nil (nil means opted out)
			if p.Meals != nil {
				// Find user's team
				for _, user := range users {
					if user.Username == p.Username {
						for _, meal := range p.Meals {
							teamMeals[user.Team][string(meal)]++
						}
						break
					}
				}
			}
		}
	}

	// For users with no entry for this date, they default to all meals
	for _, user := range users {
		if !userHasEntry[user.Username] {
			// Default: all meals opted in
			teamMeals[user.Team]["Lunch"]++
			teamMeals[user.Team]["Snacks"]++
			teamMeals[user.Team]["Iftar"]++
			teamMeals[user.Team]["Event Dinner"]++
			teamMeals[user.Team]["Optional Dinner"]++
		}
	}

	c.JSON(http.StatusOK, teamMeals)
}

// ---------- Headcount By Location (Office vs WFH) ----------
func HeadcountByLocation(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	date := c.Param("date")
	workLocations, _ := utils.LoadWorkLocations()
	users, _ := utils.LoadUsers()

	officeCount := 0
	wfhCount := 0

	for _, user := range users {
		location := "Office" // default
		for _, wl := range workLocations {
			if wl.Username == user.Username && wl.Date == date {
				location = wl.Location
				break
			}
		}

		if location == "WFH" {
			wfhCount++
		} else {
			officeCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"office": officeCount,
		"wfh":    wfhCount,
	})
}

// ---------- Overall Headcount Summary ----------
func HeadcountSummary(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	date := c.Param("date")

	// --- Load Participation locally using a local struct ---
	type LocalParticipation struct {
		Username string   `json:"username"`
		Date     string   `json:"date"`
		Meals    []string `json:"meals"`
	}

	participationFile, _ := os.ReadFile("participation.json")
	var participation []LocalParticipation
	json.Unmarshal(participationFile, &participation)

	// --- Load other required data ---
	workLocations, _ := utils.LoadWorkLocations()
	users, _ := utils.LoadUsers()
	dayControls, _ := utils.LoadDayControls() // Load day controls

	// Determine day status
	dayStatus := "normal_day"
	var dayNote *string
	isCompanyWideWFH := false
	for _, d := range dayControls {
		if d.Date == date {
			dayStatus = d.Type
			dayNote = d.Note
			if d.Type == "COMPANY_WFH" {
				isCompanyWideWFH = true
			}
			break
		}
	}

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
		// Get work location (default Office)
		location := "Office"
		for _, wl := range workLocations {
			if wl.Username == user.Username && wl.Date == date {
				location = wl.Location
				break
			}
		}

		// Override if company-wide WFH
		if isCompanyWideWFH {
			location = "WFH"
		}

		// Find participation for this user on this date
		hasParticipation := false
		for _, p := range participation {
			if p.Username == user.Username && p.Date == date {
				hasParticipation = true
				if p.Meals == nil {
					optedOut++
				} else {
					totalParticipants++
					if location == "WFH" {
						wfhCount++
					} else {
						officeCount++
					}
					for _, meal := range p.Meals {
						mealCount[meal]++
					}
				}
				break
			}
		}

		if !hasParticipation {
			// No participation record = default Office, opted in
			totalParticipants++
			if location == "WFH" {
				wfhCount++
			} else {
				officeCount++
			}
			for meal := range mealCount {
				mealCount[meal]++
			}
		}
	}

	resp := gin.H{
		"total_participants": totalParticipants,
		"office":             officeCount,
		"wfh":                wfhCount,
		"opted_out":          optedOut,
		"by_meal":            mealCount,
		"day_status":         dayStatus,
	}
	if dayNote != nil {
		resp["day_note"] = *dayNote
	}

	c.JSON(http.StatusOK, resp)
}
