package databases

import (
	"fmt"
	"time"

	"github.com/brkelkar/common_utils/logger"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	cr "github.com/brkelkar/common_utils/configreader"
)

//DB used as cursor for database connection
var DB map[string]*gorm.DB

// DBConfig represents db configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

// BuildDBMsSQLConfig Create required config format
func BuildDBMsSQLConfig(host string, port int, user string, dbName string, password string) *DBConfig {
	dbConfig := DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		DBName:   dbName,
		Password: password,
	}
	logger.Debug(
		fmt.Sprintf("Connecting to Host_name: %s, at port %v, user_name %s, database name  %s",
			dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.DBName))
	return &dbConfig
}

// DbMsSQLURL Create database connetion url
func DbMsSQLURL(dbConfig *DBConfig) string {
	return fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}

// DbMySQLURL Create database connetion url
func DbMySQLURL(dbConfig *DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}

//DbObj Dumy obect
type DbObj struct{}

//GetConnection creates connection for give database name
func (d *DbObj) GetConnection(dbName string, cfg cr.Config) (dbPtr *gorm.DB, err error) {

	// Get configured name for dbs
	var dataBaseVar string
	switch dbName {
	case "awacs_smart":
		dataBaseVar = cfg.DatabaseName.AwacsSmartDBName
	case "smartdb":
		dataBaseVar = cfg.DatabaseName.SmartStockistDBName
	case "awacs":
		dataBaseVar = cfg.DatabaseName.AwacsDBName

	}

	// set data base configuration
	dbConfig := BuildDBMsSQLConfig(cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.Username,
		dataBaseVar,
		cfg.DatabaseConfig.Password,
	)

	dbPtr, err = gorm.Open(sqlserver.Open(DbMsSQLURL(dbConfig)), &gorm.Config{})
	return
}

// CreateConnectionPool create connection pool with idle connections and Max Connection Life time settings
func (d *DbObj) CreateConnectionPool(dbPrt *gorm.DB, maxIdelConnections int, maxConnetions int, timeOutInMinutes time.Duration) {

	sqlDB, err := dbPrt.DB()
	if err != nil {
		logger.Error("Cannot creat connection pool", err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(timeOutInMinutes * time.Minute)

}
