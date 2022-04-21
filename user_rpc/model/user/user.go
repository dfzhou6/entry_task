package user

import (
	"time"
)

// User 用户模型
type User struct {
	ID         int64     `gorm:"column:id"`
	Username   string    `gorm:"column:username"`
	Password   string    `gorm:"column:password"`
	Nickname   string    `gorm:"column:nickname"`
	PicPath    string    `gorm:"column:pic_path"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
	Salt       string    `gorm:"column:salt"`
	Version    int       `gorm:"column:version"`
}
