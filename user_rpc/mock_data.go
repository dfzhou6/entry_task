package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	userDao "user_rpc/dao/user"
	userModel "user_rpc/model/user"
	"user_rpc/pkg/config"
	"user_rpc/pkg/database"
	"user_rpc/pkg/logger"
	"user_rpc/pkg/redis"
	"user_rpc/pkg/util"
)

func init() {
	config.SetupConfig()
	logger.SetupLogger()
	database.SetupDatabase()
	redis.SetupRedis()
}

func main() {
	//mockLogin(2000)
	//mockUserToken(2000)
	//mockEditUser(2000)
}

// 写数据到文件
func writeDataToFile(filename string, data string) {
	if err := ioutil.WriteFile(filename, []byte(data), 0644); err != nil {
		panic(err)
	}
}

// mock登录接口参数
func mockLogin(num int) {
	var err error

	users := make([]userModel.User, num)
	if err = database.DB.Table("users_0").Limit(num).Find(&users).Error; err != nil {
		panic(err)
	}

	var builder strings.Builder
	for _, user := range users {
		builder.WriteString(fmt.Sprintf("%s,%s\n", user.Username, "123456"))
	}

	filename := fmt.Sprintf("users_%d_登录_固定.txt", num)
	writeDataToFile(filename, builder.String())

	fmt.Println(fmt.Sprintf("writer %s success\n", filename))
}

// mock获取用户信息接口参数
func mockUserToken(num int) {
	var err error
	var token string

	users := make([]userModel.User, num)
	if err = database.DB.Table("users_0").Limit(num).Find(&users).Error; err != nil {
		panic(err)
	}

	var builder strings.Builder
	for _, user := range users {
		// 生成并写入token
		token = util.GenerateToken(user.Username)
		if err = userDao.SetTokenCache(token, user.Username); err != nil {
			panic(err)
		}

		// 写入user cache
		if err = userDao.SetUserCache(user.Username, user); err != nil {
			panic(err)
		}

		builder.WriteString(fmt.Sprintf("%s\n", token))
	}

	filename := fmt.Sprintf("users_%d_获取用户信息_固定.csv", num)
	writeDataToFile(filename, builder.String())

	fmt.Println(fmt.Sprintf("writer %s success\n", filename))
}

// mock编辑用户信息接口参数
func mockEditUser(num int) {
	var err error
	var token string

	users := make([]userModel.User, num)
	if err = database.DB.Table("users_0").Limit(num).Find(&users).Error; err != nil {
		panic(err)
	}

	var builder strings.Builder
	for _, user := range users {
		// 生成并写入token
		token = util.GenerateToken(user.Username)
		if err = userDao.SetTokenCache(token, user.Username); err != nil {
			panic(err)
		}

		// 写入user cache
		if err = userDao.SetUserCache(user.Username, user); err != nil {
			panic(err)
		}

		builder.WriteString(fmt.Sprintf("%s,%s\n", token, "ccccccc"))
	}

	filename := fmt.Sprintf("users_%d_编辑用户信息_固定.csv", num)
	writeDataToFile(filename, builder.String())

	fmt.Println(fmt.Sprintf("writer %s success\n", filename))
}
