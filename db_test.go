package database
import (
    "testing"
    "fmt"
    "encoding/json"
   "gorm.io/driver/mysql"
   "gorm.io/gorm"
)
type Users struct {
    Id uint `gorm:"column:id;PRIMARY_KEY;type:int(10)"`
    EnName string `gorm:"column:en_name;type:varchar(256)"`
    Password string `gorm:"column:password;type:varchar(256)"`
    Status int `gorm:"column:status;type:int(10)";default:1`
}
func Test_InitDb(t *testing.T) {
    dsn := "root:**********@tcp(127.0.0.1:3306)/user_info?charset=utf8&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        fmt.Printf("connect mysql failed:%v\n",err)
        return
    }
    if db != nil {
        fmt.Println("Database connected successfully")
    }
}

func Test_QueryRowGorm(t *testing.T) {
    dsn := "root:**********@tcp(127.0.0.1:3306)/user_info?charset=utf8&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        fmt.Printf("connect mysql failed:%v\n",err)
        return
    }
    if db != nil {
        fmt.Println("Database connected successfully")
    }
    dbclient := NewDBOperation(db)
    var user *Users = &Users{}
   result,err := dbclient.QueryRow("users",user,"en_name = ? and status = ?","mingyu","2")
   if err != nil {
       fmt.Printf("delete user info failed:%v\n",err)
       return
   }
   re,_ := json.Marshal(user)
   fmt.Printf("result:%v\n",result)
   fmt.Printf("output:%v\n",string(re))
}
func Test_QueryListGorm(t *testing.T) {
    dsn := "root:**********@tcp(127.0.0.1:3306)/user_info?charset=utf8&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        fmt.Printf("connect mysql failed:%v\n",err)
        return
    }
    if db != nil {
        fmt.Println("Database connected successfully")
    }
    dbclient := NewDBOperation(db)
   var userlist []*Users = make([]*Users,0)

   var page int = 1
   var pagesize int = 2
   count,err := dbclient.QueryList("users","id desc",page,pagesize,&userlist,"en_name LIKE ?","%mingyu%")
   if err != nil {
       fmt.Printf("delete user info failed:%v\n",err)
       return
   }
   fmt.Printf("count:%v\n",count)
   re,_ := json.Marshal(userlist)
   fmt.Printf("output:%v\n",string(re))
}

func Test_InsertSingleRow(t *testing.T) {
    dsn := "root:**********@tcp(127.0.0.1:3306)/user_info?charset=utf8&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        fmt.Printf("connect mysql failed:%v\n",err)
        return
    }
    if db != nil {
        fmt.Println("Database connected successfully")
    }
    dbclient := NewDBOperation(db)

    var user *Users = &Users{}
    user.EnName = "alden"
    user.Password = "alden"
    user.Status = 1
    err = dbclient.Create("users",user)
    if err != nil {
        fmt.Printf("insert user info failed:%v\n",err)
        return
    }
    fmt.Println("insert user successfully")
}

func Test_InsertMutilRow(t *testing.T) {
    dsn := "root:**********@tcp(127.0.0.1:3306)/user_info?charset=utf8&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        fmt.Printf("connect mysql failed:%v\n",err)
        return
    }
    if db != nil {
        fmt.Println("Database connected successfully")
    }
    dbclient := NewDBOperation(db)
    var users []*Users = make([]*Users,0)
    var user1 *Users = &Users{}
    var user2 *Users = &Users{}
    user1.EnName = "summer"
    user1.Password = "alden"
    user1.Status = 1
    user2.EnName = "mingyu"
    user2.Password = "alden"
    user2.Status = 1
    users = append(users,user1)
    users = append(users,user2)
    err = dbclient.CreateInBatches("users",&users,len(users))
    if err != nil {
        fmt.Printf("insert user info failed:%v\n",err)
        return
    }
    fmt.Println("insert user successfully")
}

func Test_DeleteRowGorm(t *testing.T) {
    dsn := "root:**********@tcp(127.0.0.1:3306)/user_info?charset=utf8&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        fmt.Printf("connect mysql failed:%v\n",err)
        return
    }
    if db != nil {
        fmt.Println("Database connected successfully")
    }
   dbclient := NewDBOperation(db)
   var user *Users = &Users{}

   rows,err := dbclient.DeleteRow("users",user,"en_name = ?","alden")
   if err != nil {
       fmt.Printf("delete user info failed:%v\n",err)
       return
   }
   fmt.Printf("rows:%v\n",rows)
   fmt.Printf("delete success")
}

func Test_UpdateRow(t *testing.T) {
    dsn := "root:**********@tcp(127.0.0.1:3306)/user_info?charset=utf8&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        fmt.Printf("connect mysql failed:%v\n",err)
        return
    }
    if db != nil {
        fmt.Println("Database connected successfully")
    }
   dbclient := NewDBOperation(db)
   var user *Users = &Users{}
   user.Status = 2
   rows,err := dbclient.UpdateRow("users",user,"en_name = ?","summer")
   if err != nil {
       fmt.Printf("delete user info failed:%v\n",err)
       return
   }
   fmt.Printf("rows:%v\n",rows)
   fmt.Printf("delete success")
}
