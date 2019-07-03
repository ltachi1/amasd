package models

import (
	"scrapyd-admin/core"
	"scrapyd-admin/config"
	"github.com/ltachi1/logrus"
)

type Admin struct {
	core.BaseModel            `xorm:"-"`
	Id         int            `json:"id" xorm:"pk autoincr"`
	Email      string         `json:"email" form:"email" binding:"required" xorm:"email"`
	Password   string         `json:"-" form:"password" binding:"required" xorm:"password"`
	RealName   string         `json:"real_name" xorm:"real_name"`
	Status     string         `json:"status" xorm:"status"`
	CreateTime core.Timestamp `json:"create_time" xorm:"created"`
}

const (
	AdminStatusNormal   = "enabled"
	AdminStatusDisabled = "disabled"
)

func (a *Admin) Login() (bool, string) {
	a.Password = core.Md5(a.Password)
	ok, _ := core.DBPool.Slave().Get(a)
	if !ok {
		return false, "login_password_error"
	} else {
		if a.Status == AdminStatusDisabled {
			return false, "login_user_disable"
		}
	}
	//获取用户所有角色
	adminRole := new(AdminRole)
	ok, roleIdList := adminRole.getRoleIdList(a.Id)
	if !ok {
		return false, "login_user_disable"
	}

	//注册登录信息
	ok, error := core.GetPassportInstance().Login(core.A{
		"id":          a.Id,
		"email":       a.Email,
		"real_name":   a.RealName,
		"status":      a.Status,
		"create_time": a.CreateTime,
		"role_list":   roleIdList,
	})
	if error != nil {
		return false, "login_user_disable"
	}

	//设置用户权限
	adminAccess := new(AdminAccess)
	adminAccess.SetAccessList(roleIdList)

	return true, ""
}

//分页列表
func (a *Admin) PageList(page int, pageSize int) ([]Admin, int) {
	items := make([]Admin, 0)
	totalCount, _ := core.DBPool.Slave().Table(a).Count()
	core.DBPool.Slave().Table(a).Limit(pageSize, (page-1)*pageSize).Find(&items)
	return items, int(totalCount)
}

func (a *Admin) Update(id int, data core.B) error {
	_, error := core.DBPool.Master().Table(a).ID(id).Update(data)
	return error
}

func (a *Admin) Create() (int, string) {
	//查询邮箱是否重复
	count, _ := core.DBPool.Slave().Table(a).Where("email = ?", a.Email).Count()
	if count > 0 {
		return 0, "system_email_repeat_error"
	}
	session := core.DBPool.Master().NewSession()
	defer session.Close()
	session.Begin()
	if _, error := session.InsertOne(a); error != nil {
		core.WriteLog(config.LogTypeAdmin, logrus.ErrorLevel, logrus.Fields{"email": a.Email, "real_name": a.RealName}, error)
		session.Rollback()
		return 0, "add_error"
	}
	//添加角色
	role := AdminRole{
		AdminId: a.Id,
		RoleId:  core.SuperAdminRoleId,
	}
	if _, error := session.InsertOne(&role); error != nil {
		core.WriteLog(config.LogTypeAdmin, logrus.ErrorLevel, logrus.Fields{"email": a.Email, "real_name": a.RealName}, error)
		session.Rollback()
		return 0, "add_error"
	}
	session.Commit()
	return a.Id, ""
}

func (a *Admin) Get(id int) bool {
	ok, _ := core.DBPool.Slave().Id(id).Get(a)
	return ok
}

func (a *Admin) Delete(id int) error {
	session := core.DBPool.Master().NewSession()
	defer session.Close()
	session.Begin()
	if _, error := session.Where("id = ?", id).NoAutoCondition().Delete(a); error != nil {
		session.Rollback()
		return error
	}
	//删除角色
	if _, error := session.Where("admin_id = ?", id).NoAutoCondition().Delete(AdminRole{}); error != nil {
		session.Rollback()
		return error
	}
	session.Commit()
	return nil
}
