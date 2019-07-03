package models

import (
	"scrapyd-admin/core"
	"strings"
)

type AdminAccess struct {
	core.BaseModel `xorm:"-"`
	Id             int
	RoleId         int
	App            string
	Controller     string
	Action         string
	Status         int
}

//设置用户权限列表
func (a *AdminAccess) SetAccessList(roleIdList []int) bool {
	//如果当前用户角色中包含超级官员角色则不用读取相关权限
	if core.InIntArray(core.SuperAdminRoleId, roleIdList) {
		return true
	}
	accessList := a.getAccessListByRoleIdList(roleIdList)
	accessMap := make(map[string]map[string]map[string]string)
	for _, access := range accessList {
		if _, exists := accessMap[strings.ToUpper(access.App)]; !exists {
			accessMap[strings.ToUpper(access.App)] = make(map[string]map[string]string)
		}
		if _, exists := accessMap[strings.ToUpper(access.App)][strings.ToUpper(access.Controller)]; !exists {
			accessMap[strings.ToUpper(access.App)][strings.ToUpper(access.Controller)] = make(map[string]string)
		}
		accessMap[strings.ToUpper(access.App)][strings.ToUpper(access.Controller)][strings.ToUpper(access.Action)] = strings.ToUpper(access.Action)
	}
	return core.GetRbacInstance().SaveAccessList(accessMap)
}

//获取角色下所有权限
func (a *AdminAccess) getAccessListByRoleIdList(roleIdList []int) []AdminAccess {
	accessList := make([]AdminAccess, 0)
	core.DBPool.Slave().In("role_id", roleIdList).Find(&accessList)
	return accessList
}
