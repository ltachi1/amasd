package models

import (
	"scrapyd-admin/core"
	"github.com/ltachi1/logrus"
)

type Admin struct {
	core.BaseModel             `xorm:"-"`
	Id          int            `json:"id" xorm:"pk autoincr"`
	Username    string         `json:"username" form:"username" binding:"required" xorm:"username"`
	Password    string         `json:"-" form:"password" binding:"required" xorm:"password"`
	Email       string         `json:"email" form:"email" xorm:"email"`
	DisplayName string         `json:"display_name" xorm:"display_name"`
	Status      string         `json:"status" xorm:"status"`
	CreateTime  core.Timestamp `json:"create_time" xorm:"created"`
}

const (
	AdminStatusNormal   = "enabled"
	AdminStatusDisabled = "disabled"
)

func (a *Admin) Login() (bool, string) {
	a.Password = core.Md5(a.Password)
	ok, _ := core.Db.Get(a)
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
		"id":           a.Id,
		"email":        a.Email,
		"display_name": a.DisplayName,
		"status":       a.Status,
		"create_time":  a.CreateTime,
		"role_list":    roleIdList,
	})
	if error != nil {
		return false, "login_user_disable"
	}

	//设置用户权限
	adminAccess := new(Access)
	adminAccess.SetAccessList(roleIdList)

	return true, ""
}

//分页列表
func (a *Admin) PageList(page int, pageSize int) ([]Admin, int) {
	items := make([]Admin, 0)
	totalCount, _ := core.Db.Table(a).Count()
	core.Db.Table(a).Limit(pageSize, (page-1)*pageSize).Find(&items)
	return items, int(totalCount)
}

func (a *Admin) Update(id int, data core.B) error {
	_, error := core.Db.Table(a).ID(id).Update(data)
	return error
}

func (a *Admin) Create() (int, string) {
	//查询用户名是否重复
	count, _ := core.Db.Table(a).Where("username = ?", a.Username).Count()
	if count > 0 {
		return 0, "system_username_repeat_error"
	}
	session := core.Db.NewSession()
	defer session.Close()
	session.Begin()
	if _, error := session.InsertOne(a); error != nil {
		core.WriteLog(core.LogTypeAdmin, logrus.ErrorLevel, logrus.Fields{"username": a.Username, "display_name": a.DisplayName}, error)
		session.Rollback()
		return 0, "add_error"
	}
	//添加角色
	role := AdminRole{
		AdminId: a.Id,
		RoleId:  core.SuperAdminRoleId,
	}
	if _, error := session.InsertOne(&role); error != nil {
		core.WriteLog(core.LogTypeAdmin, logrus.ErrorLevel, logrus.Fields{"username": a.Username, "display_name": a.DisplayName}, error)
		session.Rollback()
		return 0, "add_error"
	}
	session.Commit()
	return a.Id, ""
}

func (a *Admin) Get(id int) bool {
	ok, _ := core.Db.Id(id).Get(a)
	return ok
}

func (a *Admin) Delete(id int) error {
	session := core.Db.NewSession()
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
