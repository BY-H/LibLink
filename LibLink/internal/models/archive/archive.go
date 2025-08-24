package archive

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Archive struct {
	gorm.Model
	FileNo          string `gorm:"column:file_no;comment:'档案编号'" json:"file_no"`
	ContractNo      string `gorm:"column:contract_no;comment:'合同编号'" json:"contract_no"`
	Title           string `gorm:"column:title;comment:'档案标题'" json:"title"`
	Name            string `gorm:"column:name;comment:'姓名'" json:"name"`
	IDCard          string `gorm:"column:id_card;comment:'身份证号'" json:"id_card"`
	InstNo          string `gorm:"column:inst_no;comment:'网点编号'" json:"inst_no"`
	Manager         string `gorm:"column:manager;comment:'管户客户经理'" json:"manager"`
	Amount          string `gorm:"column:amount;comment:'合同金额'" json:"amount"`
	ArcType         string `gorm:"column:arc_type;comment:'文献类型'" json:"arc_type"`
	BorrowState     string `gorm:"column:borrow_state;comment:'借阅状态'" json:"borrow_state"`
	FolderID        uint   `gorm:"column:folder_id;comment:'文件夹ID'" json:"folder_id"`
	CreatorID       string `gorm:"column:creator_id;comment:'创建者ID'" json:"creator_id"`
	StorageDate     string `gorm:"column:storage_date;comment:'入库日期'" json:"storage_date"`
	GroupPermission string `gorm:"column:group_permission;comment:'用户组权限,自动继承父文件夹权限,需要有其中所有权限才能够访问该档案'" json:"group_permission"`
}

type ArchiveOperateUserKey string

const ArchiveOperateUserID ArchiveOperateUserKey = "UserID"

// BeforeUpdate 更新了文献前，需要记录对应日志
func (a *Archive) BeforeUpdate(tx *gorm.DB) (err error) {
	changes := make(map[string]interface{})

	// 检查借阅状态，日后如果需要添加其他的变更，则加入对应的监听
	if tx.Statement.Changed("borrow_state") {
		var old Archive
		tx.Model(&Archive{}).Select("borrow_state").Where("id = ?", a.ContractNo).Take(&old)
		changes["borrow_state"] = map[string]interface{}{
			"old": old.BorrowState,
			"new": a.BorrowState,
		}
	}

	if len(changes) > 0 {
		// 通过上下文传递
		if v := tx.Statement.Context.Value(ArchiveOperateUserID); v != nil {
			operatorID := v.(string)

			log := ArchiveRecord{
				ContractNo:  a.ContractNo,
				CreatorID:   operatorID,
				OperateType: a.BorrowState, // 直接使用新的借阅状态作为操作类型
				OperateDate: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			}

			if err := tx.Save(&log).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

type Folder struct {
	gorm.Model
	Name            string `gorm:"column:name;comment:'文件夹名称'" json:"name"`
	Path            string `gorm:"column:path;comment:'文件夹路径'" json:"path"`
	ParentID        uint   `gorm:"column:parent_id;index:idx_parent_id;comment:'父文件夹ID'" json:"parent_id"`
	CreatorID       string `gorm:"column:creator_id;comment:'创建者ID'" json:"creator_id"`
	GroupPermission string `gorm:"column:group_permission;comment:'用户组权限,自动继承父文件夹权限,需要有其中所有权限才能够访问该文件夹'" json:"group_permission"`
}

// FileFolders 用于查询文件夹和档案的树形结构
type FileFolders struct {
	Folder   Folder        `json:"folder"`   // 当前文件夹信息
	Children []FileFolders `json:"children"` // 子文件夹列表
	Archives []Archive     `json:"archives"` // 当前文件夹下的档案列表
}

// CheckPermission 检查用户组是否包含资源所需的权限
// resourcePerm 是资源的 GroupPermission (逗号分隔的权限组)
// userPerm 是用户的 PermissionGroup (逗号分隔的权限组)
func CheckPermission(resourcePerm string, userPerm string) bool {
	if resourcePerm == "" {
		return true // 没有限制
	}

	resourceGroups := strings.Split(resourcePerm, ",")
	userGroups := strings.Split(userPerm, ",")

	groupMap := make(map[string]struct{})
	for _, g := range userGroups {
		groupMap[strings.TrimSpace(g)] = struct{}{}
	}

	// 要求 resourceGroups 里的每个权限组都在 userGroups 中
	for _, rg := range resourceGroups {
		if _, ok := groupMap[strings.TrimSpace(rg)]; !ok {
			return false
		}
	}

	return true
}

func GetFoldersAndFilesByParentID(DB *gorm.DB, parentID uint, permission string) ([]FileFolders, error) {
	var result []FileFolders

	// 查询当前层级的所有文件夹，确保用户组权限包含传入的 permission
	var folders []Folder
	if err := DB.Where("parent_id = ? AND FIND_IN_SET(?, group_permission)", parentID, permission).
		Find(&folders).Error; err != nil {
		return nil, err
	}

	for _, f := range folders {
		node := FileFolders{
			Folder: f,
		}

		// 查询该文件夹下的档案（不加“创建者限制”）
		var archives []Archive
		if err := DB.Where("folder_id = ? AND FIND_IN_SET(?, group_permission)", f.ID, permission).
			Find(&archives).Error; err != nil {
			return nil, err
		}
		node.Archives = archives

		// 递归查询子文件夹
		children, err := GetFoldersAndFilesByParentID(DB, f.ID, permission)
		if err != nil {
			return nil, err
		}
		node.Children = children

		result = append(result, node)
	}

	return result, nil
}

// CreateFolder 新增文件夹，自动继承父文件夹的权限
func CreateFolder(DB *gorm.DB, name, path string, parentID uint, creatorID string, groupPermission string) (*Folder, error) {
	newFolder := &Folder{
		Name:      name,
		Path:      path,
		ParentID:  parentID,
		CreatorID: creatorID,
	}

	// 如果有父文件夹，则继承父文件夹的权限
	if parentID != 0 {
		var parent Folder
		if err := DB.First(&parent, parentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("父文件夹不存在")
			}
			return nil, err
		}

		// 继承父文件夹权限
		newFolder.GroupPermission = parent.GroupPermission
	} else {
		// 顶层文件夹，可以直接设置权限
		newFolder.GroupPermission = groupPermission
	}

	// 保存到数据库
	if err := DB.Create(newFolder).Error; err != nil {
		return nil, err
	}

	return newFolder, nil
}

// CreateArchive 新增文献，自动继承所属文件夹的权限
func CreateArchive(DB *gorm.DB, fileNo, title, contractNo, instNo, arcType, borrowState, creatorID string, folderID uint) (*Archive, error) {
	if folderID == 0 {
		return nil, errors.New("文献必须归属在一个文件夹下")
	}

	// 查找文件夹
	var folder Folder
	if err := DB.First(&folder, folderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("所属文件夹不存在")
		}
		return nil, err
	}

	// 创建文献，继承文件夹权限
	newArchive := &Archive{
		FileNo:          fileNo,
		Title:           title,
		ContractNo:      contractNo,
		InstNo:          instNo,
		ArcType:         arcType,
		BorrowState:     borrowState,
		FolderID:        folder.ID,
		CreatorID:       creatorID,
		GroupPermission: folder.GroupPermission, // 继承
	}

	// 保存
	if err := DB.Create(newArchive).Error; err != nil {
		return nil, err
	}

	return newArchive, nil
}

func GetArcTypeFileNo(DB *gorm.DB, arcType string) string {
	var count int64
	if err := DB.Model(&Archive{}).Where("arc_type = ?", arcType).Count(&count).Error; err != nil {
		return ""
	}
	return strconv.Itoa(int(count + 1))
}

type ArchiveRecord struct {
	gorm.Model
	ContractNo  string `gorm:"column:contract_no;comment:'合同编号'" json:"contract_no"`
	CreatorID   string `gorm:"column:creator_id;comment:'借阅人ID'" json:"creator_id"`
	OperateType string `gorm:"column:operate_type;comment:'操作类型，借阅或归还'" json:"operate_type"`
	OperateDate string `gorm:"column:operate_date;comment:'操作日期'" json:"operate_date"`
}
