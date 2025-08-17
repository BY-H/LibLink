package api

import (
	"fmt"
	"liblink/internal/controllers/message"
	"liblink/internal/global"
	"liblink/internal/middleware"
	"liblink/internal/models/system"
	"liblink/internal/models/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Notifications(c *gin.Context) {
	var notifications []system.Notification
	global.DB.Find(&notifications)

	c.JSON(http.StatusOK, gin.H{
		"total": len(notifications),
		"list":  notifications,
	})
}

func AddNotification(c *gin.Context) {
	email := middleware.GetEmail(c)
	err := checkRole(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	msg := &message.AddNotificationMsg{}
	err = c.BindJSON(&msg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	n := system.Notification{
		Type:    msg.Type,
		Title:   msg.Title,
		Content: msg.Content,
	}

	global.DB.Save(&n)

	c.JSON(http.StatusOK, gin.H{
		"msg": "add successfully",
	})
}

func checkRole(email string) error {
	var u user.User
	global.DB.Where("email = ?", email).First(&u)
	if u.Email != "" && u.Role == "admin" {
		return nil
	}
	return fmt.Errorf("this user is not admin user")
}
