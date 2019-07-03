package models

import (
	"scrapyd-admin/core"
	"github.com/go-xorm/xorm"
)

type ProjectServer struct {
	Base       core.BaseModel `json:"-" xorm:"-"`
	Id         int            `json:"id" xorm:"pk autoincr"`
	ProjectId  int            `json:"project_id" binding:"required"`
	ServerId   int            `json:"server_id" binding:"required"`
	CreateTime core.Timestamp `json:"create_time" xorm:"created"`
}

//删除项目下制定的服务器列表
func (p *ProjectServer) DelProjectServers(projectId int, serverIds []int, session *xorm.Session) error {
	_, error := session.Where("project_id = ?", projectId).In("server_id", serverIds).NoAutoCondition().Delete(p)
	return error
}

//批量增加项目下的服务器
func (p *ProjectServer) InsertProjectServers(projectId int, serverIds []int, session *xorm.Session) error {
	projectServers := make([]ProjectServer, 0)
	for _, id := range serverIds {
		projectServers = append(projectServers, ProjectServer{
			ProjectId: projectId,
			ServerId:  id,
		})
	}
	_, error := core.DBPool.Master().Insert(projectServers)
	return error
}
