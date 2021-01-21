package dbpool

import (
	"lisence/pkg/config"
	"lisence/pkg/iocgo"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var dbpool *TpaasDbPool

type TpaasDbPool struct {
	DB *gorm.DB
}

// dbpool Reload时重新读取配置文件中的信息
type dbpoolCnf struct {
	dbname      string
	maxopen     int
	maxidle     int
	maxlifetime int
	maxidletime int
}

func Pool() *TpaasDbPool {
	if dbpool == nil {
		dbpool = new(TpaasDbPool)
	}
	return dbpool
}

func (d *TpaasDbPool) Init(cfg *config.TpaasConfig) error {
	var err error
	d.DB, err = gorm.Open("mysql", cfg.GetString("mysql.datasource"))
	if err != nil {
		return err
	}
	maxopen := cfg.GetInt("mysql.max_open_conn")
	if maxopen == 0 {
		maxopen = 50
	}
	d.DB.DB().SetMaxOpenConns(maxopen)

	maxidle := cfg.GetInt("mysql.max_idle_conn")
	if maxidle == 0 {
		maxidle = 5
	}
	d.DB.DB().SetMaxIdleConns(maxidle)

	maxlifetime := cfg.GetInt("mysql.max_conn_lifetime")
	if maxlifetime != 0 {
		d.DB.DB().SetConnMaxLifetime(time.Duration(maxlifetime) * time.Minute)
	}

	maxidletime := cfg.GetInt("mysql.max_conn_idletime")
	if maxidletime != 0 {
		d.DB.DB().SetConnMaxIdleTime(time.Duration(maxidletime) * time.Minute)
	}
	return nil
}

func (d *TpaasDbPool) Close() error {
	return d.DB.Close()
}

func (d *TpaasDbPool) Reload(cfg *config.TpaasConfig) error {
	var err error

	poolcnf := new(dbpoolCnf)
	poolcnf.dbname = cfg.GetString("mysql.datasource")
	poolcnf.maxopen = cfg.GetInt("mysql.max_open_conn")
	poolcnf.maxidle = cfg.GetInt("mysql.max_idle_conn")
	poolcnf.maxidletime = cfg.GetInt("mysql.max_conn_lifetime")

	d.DB, err = gorm.Open("mysql", poolcnf.dbname)
	if err != nil {
		return err
	}
	maxopen := poolcnf.maxopen
	if maxopen == 0 {
		maxopen = 50
	}
	d.DB.DB().SetMaxOpenConns(maxopen)

	maxidle := poolcnf.maxidle
	if maxidle == 0 {
		maxidle = 5
	}
	d.DB.DB().SetMaxIdleConns(maxidle)

	maxlifetime := poolcnf.maxlifetime
	if maxlifetime != 0 {
		d.DB.DB().SetConnMaxLifetime(time.Duration(maxlifetime) * time.Minute)
	}

	maxidletime := poolcnf.maxidletime
	if maxidletime != 0 {
		d.DB.DB().SetConnMaxIdleTime(time.Duration(maxidletime) * time.Minute)
	}

	return nil
}

// 注册到iocgo 帮助统一初始化和销毁
func init() {
	iocgo.Register("TpaasDbPool", Pool())
}

func NewDbPoolWithDsn(dsn string) (*TpaasDbPool, error) {
	dbpool := new(TpaasDbPool)
	var err error
	dbpool.DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return dbpool, nil

}
