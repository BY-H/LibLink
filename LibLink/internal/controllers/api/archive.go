package api

import (
	"liblink/internal/global"
	"liblink/internal/middleware"
	"liblink/internal/models/archive"
	"liblink/internal/models/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetArchive 获取指定档案
func GetArchive(c *gin.Context) {
	// 获取档案ID
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "缺少档案ID"})
		return
	}
	archiveID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "档案ID无效"})
		return
	}

	// 获取当前用户信息
	email := middleware.GetEmail(c)

	var currentUser user.User
	if err := global.DB.Where("email = ?", email).First(&currentUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库错误"})
		return
	}

	// 查询档案
	var arc archive.Archive
	if err := global.DB.First(&arc, archiveID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "档案不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库错误"})
		return
	}

	// 校验权限
	if !archive.CheckPermission(arc.GroupPermission, currentUser.PermissionGroup) {
		c.JSON(http.StatusForbidden, gin.H{"message": "无权访问该档案"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取档案成功",
		"data":    arc,
	})
}

// CreateArchive 创建档案
func CreateArchive(c *gin.Context) {
	// 获取当前用户信息
	emailVal, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}
	email := emailVal.(string)

	var currentUser user.User
	if err := global.DB.Where("email = ?", email).First(&currentUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库错误"})
		return
	}

	// 绑定请求参数
	var req struct {
		FileNo      string `json:"file_no" binding:"required"`
		Title       string `json:"title" binding:"required"`
		ContractNo  string `json:"contract_no"`
		InstNo      string `json:"inst_no"`
		ArcType     string `json:"arc_type"`
		BorrowState string `json:"borrow_state"`
		FolderID    uint   `json:"folder_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数错误", "error": err.Error()})
		return
	}

	// 校验文件夹权限
	var folder archive.Folder
	if err := global.DB.First(&folder, req.FolderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "文件夹不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库错误"})
		return
	}

	if !archive.CheckPermission(folder.GroupPermission, currentUser.PermissionGroup) {
		c.JSON(http.StatusForbidden, gin.H{"message": "无权在该文件夹下创建档案"})
		return
	}

	// 创建档案
	newArc, err := archive.CreateArchive(
		global.DB,
		req.FileNo,
		req.Title,
		req.ContractNo,
		req.InstNo,
		req.ArcType,
		req.BorrowState,
		currentUser.Email,
		req.FolderID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建档案失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "档案创建成功",
		"data":    newArc,
	})
}

func GetArchives(c *gin.Context) {
	// 获取当前用户信息
	email := middleware.GetEmail(c)

	var currentUser user.User
	if err := global.DB.Where("email = ?", email).First(&currentUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库错误"})
		return
	}

	// 查询档案列表
	var archives []archive.Archive
	if err := global.DB.Where("group_permission = ?", currentUser.PermissionGroup).Find(&archives).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取档案列表成功",
		"data":    archives,
	})
}

func AddArchive(c *gin.Context) {

}
