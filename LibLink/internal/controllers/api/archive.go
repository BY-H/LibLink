package api

import (
	"context"
	"fmt"
	"liblink/internal/controllers/message"
	"liblink/internal/global"
	"liblink/internal/middleware"
	"liblink/internal/models/archive"
	"liblink/internal/models/user"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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

	type archRequest struct {
		message.RequestMsg
		ContractNo  string `json:"contract_no" form:"contract_no"`
		ArcType     string `json:"arc_type" form:"arc_type"`
		InstNo      string `json:"inst_no" form:"inst_no"`
		BorrowState string `json:"borrow_state" form:"borrow_state"`
	}

	var request archRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数错误", "error": err.Error()})
		return
	}

	db := global.DB.Model(&archive.Archive{}).Where("group_permission = ?", currentUser.PermissionGroup)

	// 筛选字段
	if request.ContractNo != "" {
		db = db.Where("contract_no LIKE ?", "%"+request.ContractNo+"%")
	}

	if request.ArcType != "" {
		db = db.Where("arc_type = ?", request.ArcType)
	}

	if request.InstNo != "" {
		db = db.Where("inst_no = ?", request.InstNo)
	}

	if request.BorrowState != "" {
		db = db.Where("borrow_state = ?", request.BorrowState)
	}

	// 自动分页
	if request.Page <= 0 {
		request.Page = 1
	}

	if request.PageSize <= 0 {
		request.PageSize = 10
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库错误"})
		return
	}

	db = db.Offset((request.Page - 1) * request.PageSize).Limit(request.PageSize)

	// 查询档案列表
	var archives []archive.Archive
	if err := db.Find(&archives).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "获取档案列表成功",
		"page":      request.Page,
		"page_size": request.PageSize,
		"total":     total,
		"data":      archives,
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

	// 后端生成字段
	newArchive.FileNo = archive.GetArcTypeFileNo(global.DB, newArchive.ArcType)
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

// BatchImportArchives 批量导入
// 目前只支持 xlsx 格式，后续有需要则扩展其他格式进行导入
func BatchImportArchives(c *gin.Context) {
	// 获取当前用户
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

	// 解析上传的 Excel 文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "文件上传失败", "error": err.Error()})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "无法打开文件", "error": err.Error()})
		return
	}
	defer f.Close()

	// 使用 excelize 解析
	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取Excel失败", "error": err.Error()})
		return
	}

	// 默认读取第一个sheet
	sheet := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取行失败", "error": err.Error()})
		return
	}

	var archives []archive.Archive
	for i, row := range rows {
		if i == 0 {
			// 第一行是表头
			continue
		}

		if len(row) < 6 {
			continue
		}

		// 档案类型 	合同编号	姓名	身份证号	网点编号	客户经理	合同金额 	存档日期
		a := archive.Archive{
			FileNo:          archive.GetArcTypeFileNo(global.DB, row[0]),
			ArcType:         row[0],
			ContractNo:      row[1],
			Name:            row[2],
			IDCard:          row[3],
			InstNo:          row[4],
			Manager:         row[5],
			Amount:          row[6],
			GroupPermission: currentUser.PermissionGroup,
			CreatorID:       currentUser.Email,
			BorrowState:     "0", // 默认未借阅
		}
		// 存档日期如果为空，则默认为今日
		if len(row) > 7 && row[7] != "" {
			a.StorageDate = row[7]
		} else {
			a.StorageDate = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		}

		archives = append(archives, a)
	}

	if len(archives) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Excel中没有有效数据"})
		return
	}

	// 批量插入
	if err := global.DB.Create(&archives).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "批量导入失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "批量导入成功",
		"count":   len(archives),
	})
}

func BorrowArchive(c *gin.Context) {
	// TODO: 鉴权，该用户是否有对应档案操作权限
	contractNo := c.Query("contract_no")
	if contractNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "缺少合同编号"})
		return
	}

	ctx := context.WithValue(context.Background(), archive.ArchiveOperateUserID, middleware.GetEmail(c))
	if err := operateArchive(contractNo, ctx, "1"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "借阅档案失败",
			"error":   err.Error(),
		})
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
	ctx := context.WithValue(context.Background(), archive.ArchiveOperateUserID, middleware.GetEmail(c))
	if err := operateArchive(contractNo, ctx, "0"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "归还档案失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "归还成功"})
}

// 批量更新档案状态
// 目前只支持 xlsx 格式，后续有需要则扩展其他格式进行导入
func BatchOperateArchives(c *gin.Context) {
	// 获取当前用户
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

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "文件上传失败", "error": err.Error()})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "无法打开文件", "error": err.Error()})
		return
	}
	defer f.Close()

	// 使用 excelize 解析
	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取Excel失败", "error": err.Error()})
		return
	}

	// 默认读取第一个sheet
	sheet := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取行失败", "error": err.Error()})
		return
	}

	// 批量解析借阅
	var message []string
	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) < 2 {
			message = append(message, "第"+strconv.Itoa(i+1)+"行数据格式错误")
			continue
		}

		ctx := context.WithValue(context.Background(), archive.ArchiveOperateUserID, currentUser.Email)
		err := operateArchive(row[0], ctx, row[1])
		if err != nil {
			message = append(message, "第"+strconv.Itoa(i+1)+"行操作失败: "+err.Error())
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "批量操作完成",
		"detail":  message,
	})
}

func operateArchive(contractNo string, ctx context.Context, status string) error {
	var arch archive.Archive
	if err := global.DB.Where("contract_no = ?", contractNo).First(&arch).Error; err != nil {
		return err
	}

	if arch.BorrowState == status {
		return fmt.Errorf("档案状态已是 %s", status)
	}

	arch.BorrowState = status
	if err := global.DB.WithContext(ctx).
		Model(&arch).
		Update("borrow_state", status).Error; err != nil {
		return err
	}

	return nil
}

// UpdateArchive 编辑档案
func UpdateArchive(c *gin.Context) {
	// 获取档案ID
	contractNo := c.Param("id")
	if contractNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "缺少档案ID"})
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

	// 查找档案
	var arc archive.Archive
	if err := global.DB.Where("contract_no = ?", contractNo).First(&arc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "档案不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库错误"})
		return
	}

	// 绑定请求参数
	var req struct {
		Title       string `json:"title"`
		Name        string `json:"name"`
		IDCard      string `json:"id_card"`
		InstNo      string `json:"inst_no"`
		Manager     string `json:"manager"`
		Amount      string `json:"amount"`
		ArcType     string `json:"arc_type"`
		StorageDate string `json:"storage_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数错误", "error": err.Error()})
		return
	}

	// 更新数据
	updates := map[string]interface{}{
		"title":        req.Title,
		"name":         req.Name,
		"id_card":      req.IDCard,
		"inst_no":      req.InstNo,
		"manager":      req.Manager,
		"amount":       req.Amount,
		"arc_type":     req.ArcType,
		"storage_date": req.StorageDate,
	}

	if err := global.DB.Model(&arc).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新档案失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "档案更新成功",
		"data":    arc,
	})
}
