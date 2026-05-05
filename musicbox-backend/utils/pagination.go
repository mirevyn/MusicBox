package utils

import (
	"gorm.io/gorm"
)

// Paginate 对 GORM 查询进行分页。
// 此函数通过添加 Offset 和 Limit 子句来修改传入的 gorm.DB 对象
// db: 将被修改的 GORM 查询对象。
// pageIndex: 当前页码。
// pageSize: 每页的项目数量。
// out: 用于存储结果的切片指针。
// 返回（分页前的）总记录数和任何错误
func Paginate(db *gorm.DB, pageIndex, pageSize int, out interface{}) (int64, error) {
	var total int64

	// 在应用分页之前获取总记录数
	err := db.Model(out).Count(&total).Error
	if err != nil {
		return 0, err
	}

	// 计算偏移量
	offset := (pageIndex - 1) * pageSize

	// 应用分页（限制和偏移量）并执行查询
	err = db.Offset(offset).Limit(pageSize).Find(out).Error
	if err != nil {
		return 0, err
	}

	return total, nil
}
