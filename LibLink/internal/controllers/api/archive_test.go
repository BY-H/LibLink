package api

import (
	"encoding/json"
	"fmt"
	"liblink/internal/db"
	"liblink/internal/models/archive"
	"testing"
)

// 辅助函数：打印 JSON 结果
func printJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

func TestGetFoldersAndFilesByParentID(t *testing.T) {
	// 初始化数据库连接（这里替换为你自己的配置）
	DB, err := db.InitDB("undefiner.cn:3306", "liblink", "liblink")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	// 清空旧数据
	DB.Exec("DELETE FROM archives")
	DB.Exec("DELETE FROM folders")

	// 插入模拟文件夹
	root := archive.Folder{Name: "根目录", Path: "/", ParentID: 0, GroupPermission: "admin,user"}
	DB.Create(&root)

	sub1 := archive.Folder{Name: "子目录1", Path: "/子目录1", ParentID: root.ID, GroupPermission: "admin,user"}
	sub2 := archive.Folder{Name: "子目录2", Path: "/子目录2", ParentID: root.ID, GroupPermission: "admin"}
	DB.Create(&sub1)
	DB.Create(&sub2)

	sub1_1 := archive.Folder{Name: "子目录1-1", Path: "/子目录1/子目录1-1", ParentID: sub1.ID, GroupPermission: "admin,user"}
	DB.Create(&sub1_1)

	// 插入档案
	DB.Create(&archive.Archive{FileNo: "A001", Title: "档案1", FolderID: sub1.ID, BorrowState: "未借阅", GroupPermission: "admin,user"})
	DB.Create(&archive.Archive{FileNo: "A002", Title: "档案2", FolderID: sub1_1.ID, BorrowState: "未借阅", GroupPermission: "admin,user"})
	DB.Create(&archive.Archive{FileNo: "A003", Title: "档案3", FolderID: sub2.ID, BorrowState: "未借阅", GroupPermission: "admin,user"})

	// 查询树形结构（admin 权限）
	tree, err := archive.GetFoldersAndFilesByParentID(DB, 0, "admin")
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}

	printJSON(tree)
}
