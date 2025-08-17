package archive

import (
	"errors"

	"gorm.io/gorm"
)

type Archive struct {
	gorm.Model
	FileNo          string `gorm:"column:file_no;comment:'文献编号'" json:"file_no"`
	Title           string `gorm:"column:title;comment:'文献标题'" json:"title"`
	ContractNo      string `gorm:"column:contract_no;comment:'合同号'" json:"contract_no"`
	InstNo          string `gorm:"column:inst_no;comment:'机构号'" json:"inst_no"`
	ArcType         string `gorm:"column:arc_type;comment:'文献类型'" json:"arc_type"`
	BorrowState     string `gorm:"column:borrow_state;comment:'借阅状态'" json:"borrow_state"`
	FolderID        uint   `gorm:"column:folder_id;comment:'文件夹ID'" json:"folder_id"`
	CreatorID       string `gorm:"column:creator_id;comment:'创建者ID'" json:"creator_id"`
	GroupPermission string `gorm:"column:group_permission;comment:'用户组权限,自动继承父文件夹权限,需要有其中所有权限才能够访问该档案'" json:"group_permission"`
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
