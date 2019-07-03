// rbac 权限相关
package core

import (
	"encoding/json"
	"sync"
)

var (
	obj                  *rbac
	once                 sync.Once
	accessListSessionKey = "_ACCESS_LIST"
)

type rbac struct{}

//获取当前对象实例
func GetRbacInstance() *rbac {
	once.Do(func() {
		obj = &rbac{}
	})
	return obj
}

//验证权限
func (r *rbac) CheckAccess() bool {
	return true
}

//把当前用户权限放入session中
func (r *rbac) SaveAccessList(accessList map[string]map[string]map[string]string) bool {
	if len(accessList) > 0 {
		session := GetSession()
		info, error := json.Marshal(accessList)
		if error != nil {
			return false
		}
		session.Set(accessListSessionKey, info)
		session.Save()
	}
	return true
}

//获取当前登录用户所拥有的权限
func (r *rbac) GetCurrentUserAccessList() map[string]map[string]map[string]string {
	session := GetSession()
	accessList := make(map[string]map[string]map[string]string)
	if sessionAccess := session.Get(accessListSessionKey); sessionAccess != nil {
		json.Unmarshal(session.Get(accessListSessionKey).([]byte), &accessList)
	}
	return accessList
}
