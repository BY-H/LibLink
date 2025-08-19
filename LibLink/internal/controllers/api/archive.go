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
func GetArchiveByID(c *gin.Context) {
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

// CreateArchive 创建档案(文件夹层级)
func CreateArchive(c *gin.Context) {
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

// GetArchives 获取当前用户的档案列表
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

	// 按合同号搜索
	contractNo := c.Query("contract_no")
	db := global.DB.Where("group_permission = ?", currentUser.PermissionGroup)
	if contractNo != "" {
		db = db.Where("contract_no LIKE ?", "%" + contractNo + "%")
	}

	// 查询档案列表
	var archives []archive.Archive
	if err := db.Find(&archives).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库错误"})
        return
    }

	c.JSON(http.StatusOK, gin.H{
		"message": "获取档案列表成功",
		"data":    archives,
	})
}

// AddArchive 新增档案(不管文件夹层级)
func AddArchive(c *gin.Context) {
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

	newArchive := &archive.Archive{}
	if err := c.ShouldBindJSON(&newArchive); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数错误", "error": err.Error()})
		return
	}
	newArchive.GroupPermission = currentUser.PermissionGroup
	newArchive.CreatorID = currentUser.Email

	if err := global.DB.Create(newArchive).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建档案失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "档案创建成功",
		"data":    newArchive,
	})
}

func BorrowArchive(c *gin.Context) {
	// TODO: 鉴权，该用户是否有对应档案操作权限
	contractNo := c.Query("contract_no")
	if contractNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "缺少合同编号"})
		return
	}

	var archive archive.Archive
	if err := global.DB.Where("contract_no = ?", contractNo).First(&archive).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "档案不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库错误"})
		return
	}

	// 检查档案是否已被借出
	if archive.BorrowState == "1" {
		c.JSON(http.StatusConflict, gin.H{"message": "档案已被借出"})
		return
	}

	archive.BorrowState = "1" // 设置为已借出状态
	if err := global.DB.Save(&archive).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新档案状态失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "借阅成功"})
}

func ReturnArchive(c *gin.Context) {
	// TODO: 鉴权，该用户是否有对应档案操作权限
	contractNo := c.Query("contract_no")
	if contractNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "缺少合同编号"})
		return
	}

	var archive archive.Archive
	if err := global.DB.Where("contract_no = ?", contractNo).First(&archive).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "档案不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库错误"})
		return
	}

	// 检查档案是否已被借出
	if archive.BorrowState == "0" {
		c.JSON(http.StatusConflict, gin.H{"message": "档案未被借出"})
		return
	}

	archive.BorrowState = "0" // 设置为已借出状态
	if err := global.DB.Save(&archive).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新档案状态失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "归还成功"})
}
