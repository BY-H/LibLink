package api

import (
	"fmt"
	"liblink/internal/global"
	"liblink/internal/models/archive"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取文件夹列表
func GetFolders(c *gin.Context) {
	parentID := c.Query("parent_id") // 父文件夹ID，可为空表示顶层

	var pid uint
	if parentID != "" {
		// 转换成 uint
		var tmp uint
		_, err := fmt.Sscan(parentID, &tmp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parent_id 参数无效"})
			return
		}
		pid = tmp
	}

	// TODO: 这里需要替换成实际获取用户组的逻辑
	userGroups := "[]string{c.GetHeader(X-User-Group)}"

	// 调用 archive 层方法
	folders, err := archive.GetFoldersAndFilesByParentID(global.DB, pid, userGroups)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": folders})
}

// 创建文件夹
func CreateFolder(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		ParentID uint   `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 这里需要替换成实际获取用户组的逻辑
	userGroups := "X-User-Group"

	// TODO: 这里需要替换成实际获取用户ID的逻辑
	folder, err := archive.CreateFolder(global.DB, req.Name, "pathj", req.ParentID, "", userGroups)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": folder})
}
