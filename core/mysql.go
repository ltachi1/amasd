// 数据库操作
package core

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"scrapyd-admin/config"
	"github.com/ltachi1/logrus"
)

var DBPool *xorm.EngineGroup

func InitDb() {
	//初始化数据库链接
	masterConfig := config.Conf.Db.Master
	dataSourceNameSlice := make([]string, 0)
	dataSourceNameSlice = append(dataSourceNameSlice, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", masterConfig.Username, masterConfig.Password, masterConfig.Host, masterConfig.Database, masterConfig.Charset))
	//设置从库
	slaveConfig := config.Conf.Db.Slave
	for _, slave := range slaveConfig.List {
		dataSourceNameSlice = append(dataSourceNameSlice, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", slave.Username, slave.Password, slave.Host, slave.Database, slave.Charset))
	}
	var err error
	DBPool, err = xorm.NewEngineGroup("mysql", dataSourceNameSlice)
	//设置主库最大连接数和空闲连接数
	DBPool.Master().SetMaxOpenConns(masterConfig.MaxOpenConns)
	DBPool.Master().SetMaxIdleConns(masterConfig.MaxIdleConns)

	//设置从库最大连接数和空闲连接数
	DBPool.SetMaxOpenConns(masterConfig.MaxOpenConns)
	DBPool.SetMaxIdleConns(masterConfig.MaxIdleConns)
	if err != nil {
		WriteLog(config.LogTypeAdmin, logrus.FatalLevel, nil, "数据库配置初始化失败")
	}
	//设置表前缀
	DBPool.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, config.Conf.Db.TablePrefix))
	//DBPool.ShowSQL(true)
}
