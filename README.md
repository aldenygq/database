# database
gorm mysql general lib

# 官方文档
https://gorm.io/zh_CN/docs/generic_interface.html
# 开始使用

## 下载database

```SQL
 go get github.com/aldenygq/database
```

## 实现逻辑
```Go
package main
import (
    "github.com/aldenygq/database"
    "fmt"
)
type Users struct {
    Id uint `gorm:"column:id;PRIMARY_KEY;type:int(10)"`
    EnName string `gorm:"column:en_name;type:varchar(256)"`
    Password string `gorm:"column:password;type:varchar(256)"`
    Status int `gorm:"column:status;type:int(10)";default:1`
}

func main() {
   var conf *database.GormConfig  = &database.GormConfig{}
   conf.User = "root"
   conf.Passwd = "qiang19940114**"
   conf.Host = "127.0.0.1"
   conf.Port = 3306
   conf.Dbcharset = "utf8"
   conf.MaxIdleConns = 10
   conf.MaxOpenConns = 100
   conf.MaxConnLifeTime = 600
   conf.DBName = "user_info"
   dbclient,err := database.NewDBOperation(conf)
   if err != nil {
       fmt.Printf("connect mysql failed:%v\n",err)
       return
   }
   if dbclient != nil {
       fmt.Println("Database connected successfully")
   }
   var user *Users = &Users{}
   user.Status = 1
   rows,err := dbclient.UpdateRow("users",user,"en_name = ?","summer")
   if err != nil {
       fmt.Printf("delete user info failed:%v\n",err)
       return
   }
   fmt.Printf("rows:%v\n",rows)
   fmt.Printf("delete success")
}

```
## 控制台输出：

```Go
Database connected successfully
rows:11
delete success%
```

## 日志输出
```
time="2024-12-07 19:58:02" level=info msg="Execute SQL: UPDATE `users` SET `status`=1 WHERE en_name = 'summer', rows affected: 11, duration: 11.040208ms"
time="2024-12-07 20:23:24" level=info msg="Execute SQL: UPDATE `users` SET `status`=2 WHERE en_name = 'summer', rows affected: 11, duration: 17.329792ms"
```
