package api

import "github.com/gin-gonic/gin"

// GetArchive 获取制定档案
func GetArchive(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "获取指定档案",
	})
}

func CreateArchive(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "创建档案",
	})
}

func DeleteArchive(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "删除档案",
	})
}
