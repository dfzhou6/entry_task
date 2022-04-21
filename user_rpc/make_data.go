package main

import (
	"fmt"
	"time"
	userModel "user_rpc/model/user"
	"user_rpc/pkg/config"
	"user_rpc/pkg/database"
	"user_rpc/pkg/logger"
	"user_rpc/pkg/util"
)

// 用户数据量
const userCount = 1000 * 10000

// 开启50个协程
const processCount = 50

// 协程单次批量插入数据量
const insertBatchCount = 20000

func init() {
	config.SetupConfig()
	logger.SetupLogger()
	database.SetupDatabase()
}

func main() {
	createTable()
	insertData()
}

// 创建表结构
func createTable() {
	tableTemplate := `
CREATE TABLE users_0 (
   id bigint unsigned NOT NULL AUTO_INCREMENT,
   username varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
   password varchar(32) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
   salt varchar(14) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
   nickname varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
   pic_path varchar(255) CHARACTER SET ascii COLLATE ascii_general_ci DEFAULT NULL,
   create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
   update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
   version int unsigned NOT NULL DEFAULT '1',
   PRIMARY KEY (id),
   UNIQUE KEY uniq_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
`
	database.DB.Exec(tableTemplate)
	for i := 1; i < config.GetInt("DB_USER_TABLE_COUNT"); i++ {
		tableName := fmt.Sprintf("users_%d", i)
		database.DB.Exec(fmt.Sprintf("CREATE TABLE %s LIKE users_0", tableName))
	}
}

// 插入数据
func insertData() {
	ch := make(chan int, processCount)
	quit := make(chan struct{}, processCount)
	startTime := time.Now().Unix()

	for i := 0; i < processCount; i++ {
		go func(ch chan int, quit chan struct{}) {
			var uid int
			var tableName, username, nickname, salt, password string
			j := 0
			insertMap := make(map[string][]userModel.User)
			for {
				uid = <-ch
				if uid == 0 {
					quit <- struct{}{}
					break
				}

				username = fmt.Sprintf("user%d", uid)
				tableName = util.GetTableByUsername(username)
				nickname = fmt.Sprintf("nick%d", uid)
				salt = util.GenerateRandomStr()
				password = util.GeneratePwdHash("123456", salt)
				insertMap[tableName] = append(insertMap[tableName], userModel.User{
					Username: username,
					Nickname: nickname,
					Password: password,
					Salt:     salt,
				})
				j++

				if j == insertBatchCount {
					for table, users := range insertMap {
						if err := database.DB.Table(table).
							Select("username", "nickname", "password", "salt").
							CreateInBatches(users, 1000).Error; err != nil {
							panic(err)
						}
					}
					j = 0
					insertMap = make(map[string][]userModel.User)
				}
			}
		}(ch, quit)
	}

	fmt.Println("Start to create user data,Please wait...")

	totalNum := userCount
	for i := 1; i <= totalNum; i++ {
		if i%20000 == 0 {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"),
				fmt.Sprintf("Completed %.1f%%", float64(i*100)/float64(totalNum)))
		}
		ch <- i
	}

	for i := 0; i < processCount; i++ {
		ch <- 0
	}
	for i := 0; i < processCount; i++ {
		<-quit
	}

	endTime := time.Now().Unix()
	fmt.Println("Done.Cost", endTime-startTime, "s.")
}
