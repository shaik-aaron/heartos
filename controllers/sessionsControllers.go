package controllers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shaik-aaron/fantasy-backend/intializers"
	"github.com/shaik-aaron/fantasy-backend/models"
	"github.com/shaik-aaron/fantasy-backend/utils"
)

func CreateSession(c *gin.Context) {
	var reqBody struct {
		UserID          uint   `json:"userId" binding:"required"`
		CompletedAt     string `json:"completedAt" binding:"required"`
		DurationMinutes int    `json:"durationMinutes"`
		DurationSeconds int    `json:"durationSeconds"`
		SessionType     string `json:"sessionType" binding:"required"`
		Status          string `json:"status" binding:"required"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Verify that the user exists
	var user models.User
	if err := intializers.DB.First(&user, reqBody.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Parse the completedAt timestamp
	completedAt, err := utils.ParseTime(reqBody.CompletedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid completedAt format", "details": err.Error()})
		return
	}

	// Create the session
	session := models.Session{
		UserID:          reqBody.UserID,
		CompletedAt:     completedAt,
		DurationMinutes: reqBody.DurationMinutes,
		DurationSeconds: reqBody.DurationSeconds,
		SessionType:     reqBody.SessionType,
		Status:          reqBody.Status,
	}

	result := intializers.DB.Create(&session)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session", "details": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Session created successfully",
		"session": session,
	})
}

func GetSessions(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse limit from query param, default 10, max 100
	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
			if limit > 100 {
				limit = 100
			}
		}
	}

	var sessions []models.Session
	result := intializers.DB.Where("user_id = ?", userID).
		Order("completed_at DESC").
		Limit(limit).
		Find(&sessions)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions", "details": result.Error.Error()})
		return
	}

	// Total sessions count (all time)
	var totalSessions int64
	intializers.DB.Model(&models.Session{}).Where("user_id = ?", userID).Count(&totalSessions)

	// Total time in seconds (sum of all session durations)
	var totalTimeSeconds int64
	intializers.DB.Model(&models.Session{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(duration_seconds), 0)").
		Scan(&totalTimeSeconds)

	// Session type breakdown (count and percentage per type)
	var typeCounts []struct {
		SessionType string
		Count       int64
	}
	intializers.DB.Model(&models.Session{}).
		Select("session_type, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("session_type").
		Find(&typeCounts)

	sessionTypeBreakdown := make([]gin.H, 0, len(typeCounts))
	for _, tc := range typeCounts {
		var pct float64
		if totalSessions > 0 {
			pct = math.Round(float64(tc.Count)/float64(totalSessions)*1000) / 10
		}
		sessionTypeBreakdown = append(sessionTypeBreakdown, gin.H{
			"type":       tc.SessionType,
			"count":      tc.Count,
			"percentage": pct,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions":             sessions,
		"count":                len(sessions),
		"totalSessions":        totalSessions,
		"totalTimeSeconds":     totalTimeSeconds,
		"sessionTypeBreakdown": sessionTypeBreakdown,
	})
}
