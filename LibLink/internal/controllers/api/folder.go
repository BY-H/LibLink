package api

import "github.com/gin-gonic/gin"

func GetFolders(c *gin.Context) {
	// 这里可以添加获取文件夹的逻辑
	c.JSON(200, gin.H{
		"message": "获取文件夹列表",
	})
}

func CreateFolder(c *gin.Context) {
	// 这里可以添加创建文件夹的逻辑
	c.JSON(200, gin.H{
		"message": "创建文件夹",
	})
}
