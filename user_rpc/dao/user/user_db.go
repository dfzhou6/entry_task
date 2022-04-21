package user

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	userModel "user_rpc/model/user"
	"user_rpc/pkg/database"
	"user_rpc/pkg/logger"
	"user_rpc/pkg/util"
)

// GetByUsername 根据用户名获取用户信息
func GetByUsername(username string) (user userModel.User, err error) {
	err = database.DB.
		Table(util.GetTableByUsername(username)).
		Where("username = ?", username).
		First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return user, nil
	}
	return
}

// CreateOne 创建用户信息
func CreateOne(username, password, nickname string) (user userModel.User, err error) {
	currentTime := time.Now()
	user.Username = username
	user.Salt = util.GenerateRandomStr()
	user.Password = util.GeneratePwdHash(password, user.Salt)
	user.Nickname = nickname
	user.CreateTime = currentTime
	user.UpdateTime = currentTime
	res := database.DB.
		Table(util.GetTableByUsername(username)).
		Model(&userModel.User{}).
		Select("username", "password", "salt", "nickname").
		Create(&user)
	return user, res.Error
}

// UpdateColumnByUsername 更新用户信息的某个字段
func UpdateColumnByUsername(username string, column string, value string) error {
	// 开启事务，使用乐观锁更新数据
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var user userModel.User

		// 先查询用户的version
		err := tx.
			Table(util.GetTableByUsername(username)).
			Select("version").
			Where("username = ?", username).
			First(&user).Error
		if err != nil {
			return err
		}

		// 在根据version去更新
		currentTime := time.Now()
		newVersion := user.Version + 1
		res := tx.
			Table(util.GetTableByUsername(username)).
			Where("username = ? AND version = ?", username, user.Version).
			Select(column, "update_time", "version").
			Updates(map[string]interface{}{
				column:        value,
				"update_time": currentTime,
				"version":     newVersion,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected <= 0 {
			logger.Debug("userDao", "UpdateColumnByUsername not affected",
				username, "column", column, "value", value)
			return errors.New(fmt.Sprintf("UpdateColumnByUsername not affected, "+
				"username: %s, column: %s, value: %s", username, column, value))
		}

		return nil
	})
}

// UpdateNickNameByUsername 更新用户的Nickname
func UpdateNickNameByUsername(username string, nickname string) error {
	return UpdateColumnByUsername(username, "nickname", nickname)
}

// UpdatePicPathByUsername 更新用户的PicPath
func UpdatePicPathByUsername(username string, picPath string) error {
	return UpdateColumnByUsername(username, "pic_path", picPath)
}
