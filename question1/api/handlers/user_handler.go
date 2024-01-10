package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"otpapi/db"
)

func CreateUser(dbInstance *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user db.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := dbInstance.CreateUser(context.Background(), user.Name, user.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user_id": userID})
	}
}

func GenerateOTP(db *db.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            PhoneNumber string `json:"phone_number"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        user, err := db.GetUserByPhoneNumber(context.Background(), req.PhoneNumber)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if user == nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        otp := generateRandomOTP()
        expirationTime := time.Now().Add(time.Minute)
        _, err = db.GenerateOTP(context.Background(), otp, expirationTime, req.PhoneNumber)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "OTP generated successfully"})
    }
}

func VerifyOTP(db *db.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            PhoneNumber string `json:"phone_number"`
            OTP          string `json:"otp"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        user, err := db.VerifyOTP(context.Background(), req.PhoneNumber)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if user == nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        if user.OTP != req.OTP || time.Now().After(user.OTPExpirationTime) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP or expired"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
    }
}

func generateRandomOTP() string {
    return "1234"
}