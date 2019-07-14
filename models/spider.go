package models

import (
	"github.com/go-xorm/xorm"
	"scrapyd-admin/core"
	"strconv"
)

type Spider struct {
	Base       core.BaseModel `json:"-" xorm:"-"`
	Id         int            `json:"id" xorm:"pk autoincr"`
	ProjectId  int            `json:"project_id"`
	Name       string         `json:"name"`
	Version    string         `json:"version"`
	Status     uint8          `json:"status"`
	CreateTime core.Timestamp `json:"create_time" xorm:"created"`
}

var (
	SpiderStatusNormal  uint8 = 1 //爬虫状态正常
	ServerStatusDiscard uint8 = 2 //爬虫状态废弃
)

//更新项目下所有爬虫
func (s *Spider) UpdateProjectSpiders(p *Project, spiderNames []string, session *xorm.Session) bool {
	spiders := make([]Spider, 0)
	//查询当前项目下有哪些爬虫
	if core.Db.Where("project_id = ? and version = ?", p.Id, p.LastVersion).Find(&spiders) != nil {
		return false
	}

	updateSpiders := make([]core.B, 0)
	newSpiders, oldSpiderNames := make([]Spider, 0), make([]string, 0)
	for _, spider := range spiders {
		oldSpiderNames = append(oldSpiderNames, spider.Name)
	}

	//处理需要增加的爬虫
	for _, spiderName := range spiderNames {
		if !core.InStringArray(spiderName, oldSpiderNames) {
			newSpiders = append(newSpiders, Spider{
				ProjectId: p.Id,
				Name:      spiderName,
				Version:   p.LastVersion,
				Status:    SpiderStatusNormal,
			})
		}
	}

	//处理需要修改的爬虫
	for _, spiderObj := range spiders {
		if core.InStringArray(spiderObj.Name, spiderNames) {
			updateSpiders = append(updateSpiders, core.B{
				"id":     strconv.Itoa(spiderObj.Id),
				"status": strconv.Itoa(int(SpiderStatusNormal)),
			})
		} else {
			updateSpiders = append(updateSpiders, core.B{
				"id":     strconv.Itoa(spiderObj.Id),
				"status": strconv.Itoa(int(ServerStatusDiscard)),
			})
		}
	}

	if len(newSpiders) > 0 {
		_, error := session.Insert(newSpiders)
		if error != nil {
			return false
		}
	}

	if len(updateSpiders) > 0 {
		_, error := session.Exec(core.JoinBatchUpdateSql("spider", updateSpiders, "id"))
		if error != nil {
			return false
		}
	}
	return true
}

//获取项目下爬虫数量
func (s *Spider) CountByProjectId(projectId int) int {
	count, _ := core.Db.Table("spider").Where("project_id = ?", projectId).Count()
	return int(count)
}

//分页获取爬虫数据
func (s *Spider) FindPageSpiders(projectId int, version string, page int, pageSize int, order string) ([]core.B, int) {
	spiders := make([]core.B, 0)
	countObj := core.Db.Table("spider").Alias("s").Join("INNER", "project as p", "s.project_id = p.id")
	selectObj := core.Db.Select("s.*,p.name as project_name").Table("spider").Alias("s").Join("INNER", "project as p", "s.project_id = p.id")
	if projectId > 0 {
		countObj.Where("s.project_id = ? ", projectId)
		selectObj.Where("s.project_id = ? ", projectId)
	}
	if version != "" {
		countObj.Where("s.version = ? ", version)
		selectObj.Where("s.version = ? ", version)
	}
	totalCount, _ := countObj.Count()
	selectObj.OrderBy("create_time desc").Limit(pageSize, (page-1)*pageSize).Find(&spiders)
	for i := 0; i < len(spiders); i++ {
		spiders[i]["create_time"] = core.FormatDateByString(spiders[i]["create_time"], "2006-01-02 15:04:05")
	}
	return spiders, int(totalCount)
}

func (s *Spider) FindByProjectIdAndVersion(projectId int, version string) []core.B {
	spiders := make([]core.B, 0)
	core.Db.Select("id,name").Table("spider").Where("project_id = ? and version = ? and status = ?", projectId, version, SpiderStatusNormal).Find(&spiders)
	return spiders
}

func (s *Spider) FindBySpiderIds(spiderIds []string) []Spider {
	spiders := make([]Spider, 0)
	core.Db.Select("id,name").Table("spider").In("id", spiderIds).Find(&spiders)
	return spiders
}
