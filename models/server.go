package models

import (
	"scrapyd-admin/core"
	"regexp"
	"scrapyd-admin/config"
	"github.com/ltachi1/logrus"
)

type Server struct {
	Base       core.BaseModel `json:"-" xorm:"-"`
	Id         int            `json:"id" xorm:"pk autoincr"`
	Alias      string         `json:"alias"`
	Host       string         `json:"host" binding:"required"`
	Username   string         `json:"-"`
	Password   string         `json:"-"`
	Auth       uint8          `json:"auth"`
	Status     uint8          `json:"status"`
	CreateTime core.Timestamp `json:"create_time" xorm:"created"`
	UpdateTime core.Timestamp `json:"update_time" xorm:"created"`
}

var (
	ServerStatusNormal uint8 = 1 //服务器状态正常
	ServerStatusFault  uint8 = 2 //服务器状态故障
	ServerAuthClose    uint8 = 1 //服务器验证关闭
	ServerAuthOpen     uint8 = 2 //服务器验证开启
)

//添加服务器
func (s *Server) InsertOne() (bool, string) {
	ipPortReg := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+:\d+$`)
	if !ipPortReg.MatchString(s.Host) {
		return false, "host_error"
	}
	//校验服务器是否已经添加过
	if count, _ := core.DBPool.Where("host = ?", s.Host).Table(s).Count(); count > 0 {
		return false, "host_repeat_error"
	}
	//检验服务器是否可用
	scrapyd := Scrapyd{Host: s.Host}
	if !scrapyd.DaemonStatus() {
		return false, "scrapyd_server_error"
	}

	if _, error := core.DBPool.Master().InsertOne(s); error != nil {
		core.WriteLog(config.LogTypeServer, logrus.ErrorLevel, logrus.Fields{"host": s.Host}, error)
		return false, "add_error"
	}
	return true, ""
}

//获取所有服务器
func (s *Server) Find() []Server {
	servers := make([]Server, 0)
	core.DBPool.Slave().OrderBy("id asc").Find(&servers, s)
	return servers
}

func (s *Server) FindByIds(ids []int) []Server {
	servers := make([]Server, 0)
	core.DBPool.Slave().In("id", ids).OrderBy("id asc").Find(&servers)
	return servers
}

//根据项目id获取所有服务器
func (s *Server) FindByProjectId(projectId int) []Server {
	servers := make([]Server, 0)
	core.DBPool.Slave().Table("server").Alias("s").Join("INNER", "project_server as ps", "ps.server_id = s.id").Where("ps.project_id = ?", projectId).Find(&servers, s)
	return servers
}

//根据项目id获取所有未拥有此项目的服务器
func (s *Server) FindByProjectIdNotProject(projectId int) []Server {
	servers := make([]Server, 0)
	core.DBPool.Slave().Table("server").Alias("s").Join("INNER", "project_server as ps", "ps.server_id = s.id").Where("ps.project_id != ?", projectId).Find(&servers, s)
	return servers
}

//分页获取服务器数据
func (s *Server) PageList(projectId int, page int, pageSize int) ([]Server, int) {
	servers := make([]Server, 0)
	var totalCount int64 = 0
	if projectId == 0 {
		totalCount, _ = core.DBPool.Slave().Table("server").Count()
		core.DBPool.Slave().OrderBy("update_time desc").Limit(pageSize, (page-1)*pageSize).Find(&servers)
	} else {
		totalCount, _ = core.DBPool.Slave().Table("server").Alias("s").Join("INNER", "project_server as ps", "ps.server_id = s.id").Where("ps.project_id = ?", projectId).Count()
		core.DBPool.Slave().Select("s.*").Table("server").Alias("s").Join("INNER", "project_server as ps", "ps.server_id = s.id").Where("ps.project_id = ?", projectId).Limit(pageSize, (page-1)*pageSize).OrderBy("s.update_time desc").Find(&servers)
	}
	return servers, int(totalCount)
}

func (s *Server) Get(id int) bool {
	ok, _ := core.DBPool.Slave().Id(id).NoAutoCondition().Get(s)
	return ok
}

//检测服务器状态并更新
func (s *Server) DetectionStatus() {
	serverList := s.Find()
	for _, server := range serverList {
		go func(server Server) {
			scrapyd := Scrapyd{
				Host:     server.Host,
				Auth:     server.Auth,
				Username: server.Username,
				Password: server.Password,
			}
			if scrapyd.DaemonStatus() {
				if server.Status == ServerStatusFault {
					if _, error := core.DBPool.Master().Id(server.Id).Update(&Server{Status: ServerStatusNormal}); error != nil {
						core.WriteLog(config.LogTypeServer, logrus.ErrorLevel, logrus.Fields{"host": server.Host}, error)
					}
				}
			} else {
				if server.Status == ServerStatusNormal {
					if _, error := core.DBPool.Master().Id(server.Id).Update(&Server{Status: ServerStatusFault}); error != nil {
						core.WriteLog(config.LogTypeServer, logrus.ErrorLevel, logrus.Fields{"host": server.Host}, error)
					}
				}
			}
		}(server)
	}
}
