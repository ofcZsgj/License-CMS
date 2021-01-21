package licensedb

import (
	"fmt"
	"lisence/pkg/dao/dbpool"
)

// HasTable 检查表和表模型是否存在/对应 table传入MySQL中表的名字，module传入类似&LicenseInfo{}的映射数据库的结构体变量
func HasTable(table string, module interface{}) (string, string) {
	db := dbpool.Pool().DB

	flag := db.HasTable(table)
	var isExist, isMatch string
	if flag {
		isExist = fmt.Sprintf("table %s is exist!\n", table)
	} else {
		isExist = fmt.Sprintf("table %s is no exist, you can create table\n", table)
	}

	flag = db.HasTable(module)
	if !flag {
		isMatch = fmt.Sprintf("moudle %s is no match table!\n", module)
	} else {
		isMatch = fmt.Sprintf("moudle %s is match table, you can use gorm\n", module)
	}

	return isExist, isMatch
}
