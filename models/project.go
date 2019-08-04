package models

import (
	"mime/multipart"
	"scrapyd-admin/core"
	"strconv"
	"time"
	"github.com/ltachi1/logrus"
	"sync"
	"fmt"
)

type Project struct {
	Base        core.BaseModel `json:"-" xorm:"-"`
	Id          int            `json:"id" xorm:"pk autoincr"`
	Name        string         `json:"name" xorm:"unique" binding:"required"`
	LastVersion string         `json:"last_version"`
	CreateTime  core.Timestamp `json:"create_time" xorm:"created"`
	UpdateTime  core.Timestamp `json:"update_time" xorm:"created updated"`
}

func (p *Project) InsertOne(relation string, serverIds []int, file *multipart.FileHeader) (bool, string, []string) {
	errorServerList := make([]string, 0)
	if count, _ := core.Db.Where("name = ?", p.Name).Table(p).Count(); count > 0 {
		return false, "project_name_repeat", errorServerList
	}
	//不输入版本号，则默认使用当前系统时间戳作为版本号
	if p.LastVersion == "" {
		p.LastVersion = strconv.FormatInt(time.Now().Unix(), 10)
	}
	session := core.Db.NewSession()
	defer session.Close()
	session.Begin()
	if _, error := session.Insert(p); error != nil {
		session.Rollback()
		core.WriteLog(core.LogTypeProject, logrus.ErrorLevel, logrus.Fields{"project_name": p.Name, "version": p.LastVersion}, fmt.Sprintf("项目创建失败:%s", error))
		return false, "add_error", errorServerList
	}
	//创建版本历史记录
	projectHistory := ProjectHistory{
		ProjectId: p.Id,
		Version:   p.LastVersion,
	}
	if _, error := session.Insert(projectHistory); error != nil {
		session.Rollback()
		core.WriteLog(core.LogTypeProject, logrus.ErrorLevel, logrus.Fields{"project_name": p.Name, "version": p.LastVersion}, fmt.Sprintf("项目历史版本创建失败:%s", error))
		return false, "add_error", errorServerList
	}

	//关联现有服务器
	if relation == "yes" {
		server := new(Server)
		servers := server.FindByIds(serverIds)
		if len(servers) != len(serverIds) {
			session.Rollback()
			core.WriteLog(core.LogTypeProject, logrus.ErrorLevel, logrus.Fields{"project_name": p.Name}, "项目服务获取数量不匹配")
			return false, "server_info_error", errorServerList
		}
		ch := make(chan bool)
		for _, serverInfo := range servers {
			go func(serverInfo Server) {
				projectServer := ProjectServer{
					ProjectId: p.Id,
					ServerId:  serverInfo.Id,
				}
				if _, error := session.InsertOne(&projectServer); error != nil {
					errorServerList = append(errorServerList, serverInfo.Host)
					core.WriteLog(core.LogTypeProject, logrus.ErrorLevel, logrus.Fields{"project_name": p.Name, "version": p.LastVersion}, fmt.Sprintf("项目与服务器关联失败:%s", error))
					ch <- false
				} else {
					//上传项目文件到服务器
					scrapyd := Scrapyd{
						Host:     serverInfo.Host,
						Auth:     serverInfo.Auth,
						Username: serverInfo.Username,
						Password: serverInfo.Password,
					}
					if !scrapyd.AddVersion(p, file) {
						errorServerList = append(errorServerList, serverInfo.Host)
						ch <- false
					} else {
						ch <- true
					}
				}

			}(serverInfo)
		}
		isFailure := true
		for i := 0; i < len(servers); i++ {
			if ! <-ch {
				isFailure = false
			}
		}
		close(ch)
		if !isFailure {
			session.Rollback()
			return false, "project_server_relation_error", errorServerList
		}

		//添加项目所包含爬虫列表
		scrapyd := &Scrapyd{Host: servers[0].Host, Auth: servers[0].Auth, Username: servers[0].Username, Password: servers[0].Password}
		err, spiders := scrapyd.ListSpiders(p)
		if err != nil {
			core.WriteLog(core.LogTypeProject, logrus.ErrorLevel, logrus.Fields{"project_name": p.Name, "version": p.LastVersion}, fmt.Sprintf("爬虫列表获取失败:%s", err))
			session.Rollback()
			return false, "project_spider_number_error", errorServerList
		}
		spider := new(Spider)
		if !spider.UpdateProjectSpiders(p, spiders, session) {
			session.Rollback()
			return false, "project_spider_update_error", errorServerList
		}
	}

	session.Commit()
	return true, "", errorServerList

}

//更新项目文件
func (p *Project) UpdateVersion(useHistoryVersion string, version string, file *multipart.FileHeader) (bool, string) {
	//运行中的项目不允许更新项目文件,包含定时任务
	task := new(Task)
	if task.HaveRunning(p.Id) || new(SchedulesTask).HaveEnabled(p.Id) {
		return false, "task_running_error"
	}
	session := core.Db.NewSession()
	defer session.Close()
	session.Begin()
	if useHistoryVersion == "no" {
		if version == "" {
			version = strconv.FormatInt(time.Now().Unix(), 10)
		}
		p.LastVersion = version
		//判断版本号是否有重复的
		projectHistory := &ProjectHistory{
			ProjectId: p.Id,
			Version:   version,
		}
		if projectHistory.CountByProjectIdAndVersion() > 0 {
			return false, "project_version_repeat"
		}
		if _, error := session.ID(p.Id).Cols("last_version").Update(p); error != nil {
			session.Rollback()
			return false, "update_error"
		}
		//增加版本更新历史记录
		if _, error := session.Insert(projectHistory); error != nil {
			session.Rollback()
			return false, "update_error"
		}
		if error := session.Commit(); error != nil {
			session.Rollback()
			return false, "update_error"
		}
	} else {
		if _, error := session.Exec("update project set update_time = ? where id = ?", core.Timestamp(time.Now().Unix()), p.Id); error != nil {
			session.Rollback()
			return false, "update_error"
		}
	}
	if !p.Get(p.Id) {
		session.Rollback()
		return false, "parameter_error"
	}
	//查询当前项目下可用服务器
	server := Server{
		Status: ServerStatusNormal,
	}
	servers := server.FindByProjectId(p.Id)
	if len(servers) == 0 {
		session.Rollback()
		return false, "project_no_server"
	}
	for _, s := range servers {
		//上传项目文件到服务器
		scrapyd := &Scrapyd{Host: s.Host, Auth: s.Auth, Username: s.Username, Password: s.Password}
		if !scrapyd.AddVersion(p, file) {
			session.Rollback()
			return false, "project_update_version_error"
		}
	}
	//当所有服务器都更新成功再更新项目所包含爬虫
	scrapyd := &Scrapyd{Host: servers[0].Host, Auth: servers[0].Auth, Username: servers[0].Username, Password: servers[0].Password}
	err, spiders := scrapyd.ListSpiders(p)
	if err != nil {
		session.Rollback()
		core.WriteLog(core.LogTypeProject, logrus.ErrorLevel, logrus.Fields{"project_name": p.Name, "version": p.LastVersion}, fmt.Sprintf("爬虫列表获取失败:%s", err))
		return false, "project_spider_number_error"
	}
	spider := new(Spider)
	if !spider.UpdateProjectSpiders(p, spiders, session) {
		session.Rollback()
		return false, "project_spider_update_error"
	}
	session.Commit()
	return true, ""
}

//更新关联服务器
func (p *Project) UpdateServers(serverIds []int, file *multipart.FileHeader) (bool, string) {
	server := new(Server)
	var (
		projectServer     = new(ProjectServer)
		beforeServerIds   []int
		cutBackServers    []Server
		cutBackServerIds  []int
		increaseServerIds []int
	)
	relatedServers := server.FindByProjectId(p.Id)
	task := new(Task)
	schedulesTask := new(SchedulesTask)
	//减少的服务器
	for _, rs := range relatedServers {
		if !core.InIntArray(rs.Id, serverIds) {
			//如果项目下有正在运行的定时任务或者普通任务则不允许删除
			if task.HaveRunning(p.Id) || schedulesTask.HaveEnabled(p.Id) {
				return false, "server_cutback_task_running_error"
			}
			cutBackServers = append(cutBackServers, rs)
			cutBackServerIds = append(cutBackServerIds, rs.Id)
		}
		beforeServerIds = append(beforeServerIds, rs.Id)
	}
	//增加的服务器id
	for _, id := range serverIds {
		if !core.InIntArray(id, beforeServerIds) {
			increaseServerIds = append(increaseServerIds, id)
		}
	}
	//有新增服务器则必须上传项目文件
	if len(increaseServerIds) > 0 && file == nil {
		return false, "file_upload_error"
	}
	session := core.Db.NewSession()
	defer session.Close()
	session.Begin()
	//处理减少的服务器
	if len(cutBackServers) > 0 {
		ch := make(chan bool)
		for _, cbs := range cutBackServers {
			go func(cbs Server) {
				s := Scrapyd{
					Host:     cbs.Host,
					Auth:     cbs.Auth,
					Username: cbs.Username,
					Password: cbs.Password,
				}
				ch <- s.DelProject(p.Name)
			}(cbs)
		}
		isFailure := true
		for i := 0; i < len(cutBackServers); i++ {
			if ! <-ch {
				isFailure = false
			}
		}
		close(ch)
		if !isFailure {
			session.Rollback()
			return false, "project_server_relation_all_error"
		}

		if error := projectServer.DelProjectServers(p.Id, cutBackServerIds, session); error != nil {
			core.WriteLog(core.LogTypeProject, logrus.ErrorLevel, logrus.Fields{"project_name": p.Name}, error)
			session.Rollback()
			return false, "update_error"
		}

		//删除所有定时任务
		if _, error := session.Where("project_id = ?", p.Id).In("server_id", cutBackServerIds).NoAutoCondition().Delete(schedulesTask); error != nil {
			return false, "update_error"
		}
		//删除所有任务
		if _, error := session.Where("project_id = ?", p.Id).In("server_id", cutBackServerIds).NoAutoCondition().Delete(task); error != nil {
			return false, "update_error"
		}
	}
	//处理增加的服务器
	if len(increaseServerIds) > 0 {
		//获取新增加的服务器信息
		increaseServers := server.FindByIds(increaseServerIds)
		if len(increaseServers) != len(increaseServerIds) {
			session.Rollback()
			return false, "server_info_error"
		}
		ch := make(chan bool)
		for _, s := range increaseServers {
			go func(s Server) {
				scrapyd := Scrapyd{
					Host:     s.Host,
					Auth:     s.Auth,
					Username: s.Username,
					Password: s.Password,
				}
				ch <- scrapyd.AddVersion(p, file)
			}(s)
		}
		isFailure := true
		for i := 0; i < len(increaseServers); i++ {
			if ! <-ch {
				isFailure = false
			}
		}
		close(ch)
		if !isFailure {
			session.Rollback()
			return false, "project_server_relation_all_error"
		}

		if error := projectServer.InsertProjectServers(p.Id, increaseServerIds, session); error != nil {
			core.WriteLog(core.LogTypeProject, logrus.ErrorLevel, logrus.Fields{"project_name": p.Name}, error)
			session.Rollback()
			return false, "update_error"
		}
		//如果没有关联过服务器，则需要更新项目所包含爬虫列表
		scrapyd := &Scrapyd{Host: increaseServers[0].Host, Auth: increaseServers[0].Auth, Username: increaseServers[0].Username, Password: increaseServers[0].Password}
		err, spiders := scrapyd.ListSpiders(p)
		if err != nil {
			session.Rollback()
			core.WriteLog(core.LogTypeProject, logrus.ErrorLevel, logrus.Fields{"project_name": p.Name, "version": p.LastVersion}, fmt.Sprintf("爬虫列表获取失败:%s", err))
			return false, "project_spider_number_error"
		}
		spider := new(Spider)
		if !spider.UpdateProjectSpiders(p, spiders, session) {
			session.Rollback()
			return false, "project_spider_update_error"
		}
	}
	session.Commit()
	return true, ""
}

//根据id获取项目信息
func (p *Project) Get(id int) bool {
	ok, _ := core.Db.Id(id).NoAutoCondition().Get(p)
	return ok && p.Id > 0
}

//获取所有项目
func (p *Project) Find() []Project {
	projects := make([]Project, 0)
	core.Db.OrderBy("id asc").Find(&projects)
	return projects
}

//分页获取项目数据
func (p *Project) GetPageProjects(serverId int, page int, pageSize int) ([]Project, int) {
	projects := make([]Project, 0)
	var totalCount int64 = 0
	if serverId == 0 {
		totalCount, _ = core.Db.Table("project").Count()
		core.Db.Limit(pageSize, (page-1)*pageSize).Find(&projects)
	} else {
		totalCount, _ = core.Db.Table("project").Alias("p").Join("INNER", "project_server as ps", "ps.project_id = p.id").Where("ps.server_id = ?", serverId).Count()
		core.Db.Select("p.*").Table("project").Alias("p").Join("INNER", "project_server as ps", "ps.project_id = p.id").Where("ps.server_id = ?", serverId).Limit(pageSize, (page-1)*pageSize).Find(&projects)
	}
	return projects, int(totalCount)
}

func (p *Project) Del(id int) (bool, string) {
	if !p.Get(id) {
		return false, "project_del_error"
	}
	//如果项目下有正在运行的定时任务或者普通任务则不允许删除
	if new(Task).HaveRunning(p.Id) || new(SchedulesTask).HaveEnabled(p.Id) {
		return false, "project_task_running_error"
	}
	var wg sync.WaitGroup
	//查询所有关联服务器
	server := new(Server)
	successServerIds, failureServerIds := make([]int, 0), make([]int, 0)
	for _, s := range server.FindByProjectId(id) {
		wg.Add(1)
		go func(s Server) {
			scrapyd := Scrapyd{
				Host:     s.Host,
				Auth:     s.Auth,
				Username: s.Username,
				Password: s.Password,
			}
			if scrapyd.DelProject(p.Name) {
				successServerIds = append(successServerIds, s.Id)
			} else {
				failureServerIds = append(failureServerIds, s.Id)
			}
			wg.Done()
		}(s)
	}
	wg.Wait()
	if len(successServerIds) > 0 {
		session := core.Db.NewSession()
		defer session.Close()
		session.Begin()
		//删除关联服务器
		if _, err := session.Where("project_id = ?", id).In("server_id", successServerIds).NoAutoCondition().Delete(&ProjectServer{}); err != nil {
			return false, "project_del_error"
		}
		//删除所有定时任务
		if _, err := session.Where("project_id = ?", id).In("server_id", successServerIds).NoAutoCondition().Delete(&SchedulesTask{}); err != nil {
			return false, "project_del_error"
		}
		//删除所有任务
		if _, err := session.Where("project_id = ?", id).In("server_id", successServerIds).NoAutoCondition().Delete(&Task{}); err != nil {
			return false, "project_del_error"
		}
		if len(failureServerIds) == 0 {
			//删除项目历史
			if _, err := session.Where("project_id = ?", id).NoAutoCondition().Delete(&ProjectHistory{}); err != nil {
				return false, "project_del_error"
			}
			//删除项目
			if _, err := session.Where("id = ?", id).NoAutoCondition().Delete(&Project{}); err != nil {
				return false, "project_del_error"
			}
		}
		session.Commit()
	} else if len(successServerIds) == 0 && len(failureServerIds) == 0 {
		if _, err := core.Db.Where("id = ?", id).NoAutoCondition().Delete(&Project{}); err != nil {
			return false, "del_error"
		}
	}

	return true, "project_del_error"
}
