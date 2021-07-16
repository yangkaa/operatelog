package db

import (
	"fmt"
	"gorm.io/gorm/schema"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//DatabaseType db type
type DatabaseType string

//MysqlDB mysql db
var MysqlDB DatabaseType = "mysql"

//InstanceConfig db config
type InstanceConfig struct {
	DBType               DatabaseType
	User                 string
	Passwd               string
	Net                  string
	Addr                 string
	DBName               string
	AllowNativePasswords bool
	ParseTime            bool
	DebugMode            bool
	Loc                  *time.Location
	Models               []interface{}
	TablePrefix          string
}

//Instance db instance
type Instance struct {
	config *InstanceConfig
	db     *gorm.DB
}

//NewDBInstance new db Instance
func NewDBInstance(config *InstanceConfig) (*Instance, error) {
	logrus.Infof("start create db instance")
	var dsn string
	var tableOptions string
	switch config.DBType {
	case MysqlDB:
		mySQLConfig := &mysql.Config{
			User:                 config.User,
			Passwd:               config.Passwd,
			Net:                  "tcp",
			Addr:                 config.Addr,
			DBName:               config.DBName,
			AllowNativePasswords: config.AllowNativePasswords,
			ParseTime:            config.ParseTime,
			Params:               map[string]string{"charset": "utf8"},
			Loc:                  time.Local,
			Timeout:              time.Second * 5,
		}
		dsn = mySQLConfig.FormatDSN()
		tableOptions = "ENGINE=InnoDB DEFAULT CHARSET=utf8"
	default:
		return nil, fmt.Errorf("%s database type is not support", config.DBType)
	}
	var instance = Instance{
		config: config,
	}

	for {
		db, err := gorm.Open(gmysql.Open(dsn), &gorm.Config{
			NamingStrategy: &schema.NamingStrategy{
				TablePrefix: config.TablePrefix,
			},
		})
		// db.Statement.RaiseErrorOnNotFound = true

		if err != nil {
			logrus.Errorf("open db connection failure %s, will retry", err.Error())
			time.Sleep(time.Second * 3)
			continue
		}
		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("get *sql.DB: %v", err)
		}
		sqlDB.SetConnMaxLifetime(time.Second)
		if config.DebugMode {
			db = db.Debug()
		}
		instance.db = db
		break
	}
	//gorm
	for _, model := range config.Models {
		if err := instance.db.Set("gorm:table_options", tableOptions).AutoMigrate(model); err != nil {
			return nil, fmt.Errorf("auto migrate table %+v failure %v", model, err)
		}
	}
	return &instance, nil
}

//DB get db
func (i *Instance) DB() *gorm.DB {
	return i.db
}

//RunInitFunc run db init function
func (i *Instance) RunInitFunc(funs ...func(db *gorm.DB)) {
	for _, f := range funs {
		f(i.DB())
	}
}
