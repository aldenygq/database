package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"time"
)

type GormConfig struct {
	User            string
	Passwd          string
	Host            string
	Port            int
	Dbcharset       string
	MaxIdleConns    int
	MaxOpenConns    int
	MaxConnLifeTime int
	DBName          string
}
type DBOperation interface {
	Create(table string, value interface{}) error
	CreateInBatches(table string, value interface{}, batchSize int) error
	QueryRow(table string, value, query interface{}, args ...interface{}) (bool, error)
	QueryList(table, order string, page, pagesize int, value, query interface{}, args ...interface{}) (int64, error)
	QueryAll(table, order string, value, query interface{}, args ...interface{}) (int64, error)
	DeleteRow(table string, value, query interface{}, args ...interface{}) (int64, error)
	UpdateRow(table string, value, query interface{}, args ...interface{}) (int64, error)
}

type dbOperation struct {
	db *gorm.DB
}

func NewDBOperation(conf *GormConfig) (DBOperation, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&loc=Local",
		conf.User,
		conf.Passwd,
		conf.Host,
		conf.Port,
		conf.DBName,
		conf.Dbcharset,
	)
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		DisableAutomaticPing:                     true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: NewGormLogger(),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := client.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(conf.MaxConnLifeTime))
	return &dbOperation{db: client}, nil
}

// Create 通用插入数据
func (d *dbOperation) Create(table string, value interface{}) error {
	if value == nil {
		return errors.New("value not be null")
	}
	if reflect.ValueOf(value).Kind() != reflect.Pointer {
		return errors.New("value must be a pointer")
	}
	tx := d.db.Table(table).Begin()
	err := tx.Create(value).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

// 批量插入
func (d *dbOperation) CreateInBatches(table string, value interface{}, batchSize int) error {
	if value == nil {
		return errors.New("value not be null")
	}
	if reflect.ValueOf(value).Kind() != reflect.Pointer {
		return errors.New("value must be a pointer")
	}
	tx := d.db.Table(table).Begin()
	err := tx.CreateInBatches(value, batchSize).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

// 获取单条数据,默认正序第一条
func (d *dbOperation) QueryRow(table string, value, query interface{}, args ...interface{}) (bool, error) {
	if value == nil {
		return false, errors.New("out not be null")
	}
	if reflect.ValueOf(value).Kind() != reflect.Pointer {
		return false, errors.New("value must be a pointer")
	}
	result := d.db.Table(table).Where(query, args...).First(&value)
	if result.Error == gorm.ErrRecordNotFound {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// 分页查询数据列表
func (d *dbOperation) QueryList(table, order string, page, pagesize int, value, query interface{}, args ...interface{}) (int64, error) {

	var count int64
	if value == nil {
		return count, errors.New("out not be null")
	}
	if reflect.ValueOf(value).Kind() != reflect.Pointer {
		return count, errors.New("value must be a pointer")
	}
	db := d.db.Table(table).Where(query, args...)
	//排序
	if order != "" {
		db = db.Order(order)
	}
	//总数
	db = db.Count(&count)
	//分页
	if page <= 0 {
		page = 1
	}
	if pagesize < 0 {
		pagesize = 10
	}
	offset := (page - 1) * pagesize
	db = db.Offset(offset).Limit(pagesize)

	//查询
	if err := db.Find(value).Error; err != nil {
		return count, err
	}
	return count, nil
}

// 查询所有数据
func (d *dbOperation) QueryAll(table, order string, value, query interface{}, args ...interface{}) (int64, error) {

	var count int64
	if value == nil {
		return count, errors.New("out not be null")
	}
	if reflect.ValueOf(value).Kind() != reflect.Pointer {
		return count, errors.New("value must be a pointer")
	}
	db := d.db.Table(table).Where(query, args...)
	//排序
	if order != "" {
		db = db.Order(order)
	}
	//总数
	db = db.Count(&count)
	//查询
	if err := db.Find(value).Error; err != nil {
		return count, err
	}
	return count, nil
}

// 删除数据
func (d *dbOperation) DeleteRow(table string, value, query interface{}, args ...interface{}) (int64, error) {
	var count int64
	var err error
	if value == nil {
		return count, errors.New("value not be null")
	}
	if reflect.ValueOf(value).Kind() != reflect.Pointer {
		return count, errors.New("value must be a pointer")
	}
	tx := d.db.Table(table).Begin()
	result := tx.Where(query, args...).Delete(value)
	if err = result.Error; err != nil {
		tx.Rollback()
		return count, err
	}
	err = tx.Commit().Error
	if err != nil {
		return count, err
	}
	count = result.RowsAffected
	return count, nil
}

// 更新数据
func (d *dbOperation) UpdateRow(table string, value, query interface{}, args ...interface{}) (int64, error) {
	var count int64
	var err error
	if value == nil {
		return count, errors.New("value not be null")
	}
	if reflect.ValueOf(value).Kind() != reflect.Pointer {
		return count, errors.New("value must be a pointer")
	}
	tx := d.db.Table(table).Begin()
	result := tx.Where(query, args...).Updates(value)
	if err = result.Error; err != nil {
		tx.Rollback()
		return count, err
	}

	err = tx.Commit().Error
	if err != nil {
		return count, err
	}
	count = result.RowsAffected
	return count, nil
}
