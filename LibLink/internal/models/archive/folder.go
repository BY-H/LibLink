package archive

import "gorm.io/gorm"

type Folder struct {
	gorm.Model
	Name     string `gorm:"column:name;comment:'文件夹名称'" json:"name"`
	Path     string `gorm:"column:path;comment:'文件夹路径'" json:"path"`
	ParentID uint   `gorm:"column:parent_id;comment:'父文件夹ID'" json:"parent_id"`
}
