package models

import (
	"amasd/core"
)

type AdminRole struct {
	core.BaseModel `xorm:"-"`
	Id      int    `json:"id" xorm:"pk autoincr"`
	AdminId int
	RoleId  int
}

//获取当前用户所有角色id列表
func (a *AdminRole) getRoleIdList(adminId int) (bool, []int) {
	roles := make([]AdminRole, 0)
	roleIdList := make([]int, 0)
	error := core.Db.Where("admin_id = ?", adminId).Find(&roles)
	if error != nil {
		return false, roleIdList
	}
	for _, role := range roles {
		roleIdList = append(roleIdList, role.RoleId)
	}
	return true, roleIdList
}
