package models

import (
	"amasd/core"
)

type ProjectHistory struct {
	Base       core.BaseModel `json:"-" xorm:"-"`
	Id         int            `json:"id" xorm:"pk autoincr"`
	ProjectId  int            `json:"project_id" binding:"required"`
	Version    string         `json:"version"`
	CreateTime core.Timestamp `json:"create_time" xorm:"created"`
}

//根据项目id获取历史版本记录
func (p *ProjectHistory) FindByProjectId() []core.B {
	projectHistory := make([]core.B, 0)
	core.Db.Where("project_id = ?", p.ProjectId).Table(p).OrderBy("create_time desc").Find(&projectHistory)
	for i := 0; i < len(projectHistory); i++ {
		projectHistory[i]["create_time"] = core.FormatDateByString(projectHistory[i]["create_time"])
	}
	return projectHistory
}

func (p *ProjectHistory) CountByProjectIdAndVersion() int {
	count, _ := core.Db.Table("project_history").Where("project_id = ? and version = ?", p.ProjectId, p.Version).Count()
	return int(count)
}
