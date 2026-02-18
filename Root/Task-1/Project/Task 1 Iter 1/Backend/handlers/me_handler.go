package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samin-craftsmen/gin-project/models"
	"github.com/samin-craftsmen/gin-project/utils"
)

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

func SetWorkLocation(c *gin.Context) {

	var req models.WorkLocation

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate location
	if req.Location != "Office" && req.Location != "WFH" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work location"})
		return
	}

	// ------------------ Save Work Location ------------------

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

	// ------------------ If WFH → Opt Out Meals ------------------

	if req.Location == "WFH" {

		participationData, err := utils.LoadParticipation()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load participation"})
			return
		}

		newParticipation := []models.Participation{}

		found := false

		for _, p := range participationData {
			if p.Username == req.Username && p.Date == req.Date {
				// Replace with meals: null
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

		// If no participation existed, add one explicitly
		if !found {
			newParticipation = append(newParticipation, models.Participation{
				Username: req.Username,
				Date:     req.Date,
				Meals:    nil,
			})
		}

		err = utils.SaveParticipation(newParticipation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save participation"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work location set successfully"})
}

func GetWorkLocation(c *gin.Context) {

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

	// Default if not found
	c.JSON(http.StatusOK, gin.H{
		"username": req.Username,
		"date":     req.Date,
		"location": "Office",
	})
}
