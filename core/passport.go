// 通行证
package core

import (
	"encoding/json"
	"sync"
)

var (
	passportObj        *passport
	passportOnce       sync.Once
	SessionUserInfoKey = "session_user_info"
)

type passport struct {
	Id         int    `json:"id"`
	Email      string `json:"email"`
	RealName   string `json:"real_name"`
	Mobile     string `json:"mobile"`
	Status     int8   `json:"status"`
	CreateTime int    `json:"create_time"`
	RoleList   []int  `json:"role_list"`
}

const SuperAdminRoleId = 1 //超级管理员角色id

//获取当前对象实例
func GetPassportInstance() *passport {
	passportOnce.Do(func() {
		passportObj = &passport{}
		passportObj.getSessionUserInfo()
	})
	return passportObj
}

func (p *passport) Login(userInfo map[string]interface{}) (bool, error) {
	//保存用户信息到session中
	session := GetSession()
	info, error := json.Marshal(userInfo)
	if error != nil {
		return false, error
	}
	session.Set(SessionUserInfoKey, info)
	session.Save()
	json.Unmarshal(info, &p)
	return true, nil
}

func (p *passport) getSessionUserInfo() {
	session := GetSession()
	if info := session.Get(SessionUserInfoKey); info != nil {
		json.Unmarshal(info.([]byte), &p)
	}

}

func (p *passport) UserInfo() {

}

//检查用户是否已经登录
func (p *passport) CheckLogin() bool {
	session := GetSession()
	userInfo := session.Get(SessionUserInfoKey)
	if userInfo == nil {
		return false
	}
	return true
}

//判断当前登录用户是否拥有管理员权限
func (p *passport) IsAdminRole() bool {
	return InIntArray(SuperAdminRoleId, p.RoleList)
}
