package archive

import (
	"gorm.io/gorm"
)

type Archive struct {
	gorm.Model
	FileNo      string `gorm:"column:file_no;comment:'文献编号'" json:"file_no"`
	Title       string `gorm:"column:title;comment:'文献标题'" json:"title"`
	ContractNo  string `gorm:"column:contract_no;comment:'合同号'" json:"contract_no"`
	InstNo      string `gorm:"column:inst_no;comment:'机构号'" json:"inst_no"`
	ArcType     string `gorm:"column:arc_type;comment:'文献类型'" json:"arc_type"`
	BorrowState string `gorm:"column:borrow_state;comment:'借阅状态'" json:"borrow_state"`
	FolderID    uint   `gorm:"column:folder_id;comment:'文件夹ID'" json:"folder_id"`
}

type Folder struct {
	gorm.Model
	Name            string `gorm:"column:name;comment:'文件夹名称'" json:"name"`
	Path            string `gorm:"column:path;comment:'文件夹路径'" json:"path"`
	ParentID        uint   `gorm:"column:parent_id;index:idx_parent_id;comment:'父文件夹ID'" json:"parent_id"`
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

	// 查询当前层级的所有文件夹
	var folders []Folder
	if err := DB.Where("parent_id = ? AND FIND_IN_SET(?, group_permission)", parentID, permission).Find(&folders).Error; err != nil {
		return nil, err
	}

	// 遍历每个文件夹，递归查询
	for _, f := range folders {
		node := FileFolders{
			Folder: f,
		}

		// 查询该文件夹下的档案
		var archives []Archive
		if err := DB.Where("folder_id = ? AND FIND_IN_SET(?, ?)", f.ID, permission, f.GroupPermission).Find(&archives).Error; err != nil {
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
