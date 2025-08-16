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

// FileFolders 用于查询文件夹和档案的结构体
type FileFolders struct {
	FileFolders []FileFolders `json:"file_folders"` // 文件夹列表
	Archives    []Archive     `json:"archives"`     // 当前文件夹下的档案列表
}

func GetFoldersAndFilesByParentID(DB *gorm.DB, parentID uint, permission string) (*FileFolders, error) {
	var fileFolders FileFolders
	var archives []Archive
	// 查询当前文件夹底下的所有文件
	err := DB.Model(&Archive{}).Where("parent_id = ? AND FIND_IN_SET(?, group_permission)", parentID, permission).Find(&archives).Error
	if err != nil {
		return nil, err
	}
	// 查询当前文件夹底下的所有文件夹
	err = DB.Model(&Folder{}).Where("parent_id = ? AND FIND_IN_SET(?, group_permission)", parentID, permission).Find(&fileFolders).Error
	if err != nil {
		return nil, err
	}
	fileFolders.Archives = archives
	return &fileFolders, nil
}
