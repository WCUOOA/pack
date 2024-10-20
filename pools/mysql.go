package pools

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type MySQL struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

func ConnectionPools(sql *MySQL, MaxLink int, MaxIdle int, MaxTime time.Duration) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", sql.User, sql.Password, sql.Host, sql.Port, sql.Database)
	var once sync.Once
	once.Do(func() {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Println("SQL连接出错(There was an error with the SQL connection)", err)
		}
		//配置连接池参数
		SQL, err := db.DB()
		if err != nil {
			log.Println("无法获取服务器实例(failed to get DB instance)", err)
		}
		//设置最大连接数
		SQL.SetMaxOpenConns(MaxLink)
		//设置最大空闲数
		SQL.SetMaxIdleConns(MaxIdle)
		//设置最大生存时间
		SQL.SetConnMaxLifetime(MaxTime)
	})
	log.Println("MySQL连接成功")
	return db, nil
}
