package archive

import "gorm.io/gorm"

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
