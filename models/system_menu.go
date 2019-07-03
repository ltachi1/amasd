package models

import (
	"fmt"
	"html/template"
	"scrapyd-admin/core"
	"regexp"
	"strings"
)

type SystemMenu struct {
	core.BaseModel      `xorm:"-"`
	Id           int    `xorm:"pk autoincr" json:"id" form:"id"`
	Name         string `json:"name" form:"name"`
	ParentId     int    `json:"parent_id" form:"parent_id"`
	App          string `json:"app" form:"app"`
	Controller   string `json:"controller" form:"controller"`
	Action       string `json:"action" form:"action"`
	Parameter    string `json:"parameter" form:"parameter"`
	Icon         string `json:"icon" form:"icon"`
	Type         string `json:"type" form:"type"`
	Status       int    `json:"status" form:"status"`
	ListingOrder int    `json:"listing_order" form:"listing_order"`
}

const (
	MenuStatusEnable  = 1 //菜单显示状态
	MenuStatusDisable = 2 //菜单禁用状态
)

//更新数据
func (s *SystemMenu) Insert() bool {
	_, error := core.DBPool.Master().InsertOne(s)
	if error != nil {
		return false
	}
	return true
}

func (s *SystemMenu) Update(id int, data core.A) error {
	_, error := core.DBPool.Master().Table(s).ID(id).Update(data)
	return error
}

//删除菜单
func (s *SystemMenu) DeleteById() bool {
	_, error := core.DBPool.Master().Id(s.Id).Delete(s)
	if error != nil {
		return false
	}
	return true
}

//获取子菜单
func (s *SystemMenu) getChildMenuList(parentId int) []SystemMenu {
	accessList := core.GetRbacInstance().GetCurrentUserAccessList()
	menus, overMenus := make([]SystemMenu, 0), make([]SystemMenu, 0)
	if err := core.DBPool.Slave().Where("parent_id = ? and status = ? ", parentId, MenuStatusEnable).OrderBy("listing_order asc, id asc").Find(&menus); err == nil {
		if core.GetPassportInstance().IsAdminRole() {
			return menus
		}
		for i := 0; i < len(menus); i++ {
			if s.hasCompetence(menus[i].App, menus[i].Controller, menus[i].Action, accessList) {
				overMenus = append(overMenus, menus[i])
			}
		}
	}
	return overMenus
}

//验证当前登录的用户是否对此菜单有操作权限
func (s *SystemMenu) hasCompetence(app string, controller string, action string, accessList map[string]map[string]map[string]string) bool {
	actionRegexp := regexp.MustCompile(`^public`)
	if actionRegexp.MatchString(action) {
		return true
	}
	if _, exists := accessList[strings.ToUpper(app)][strings.ToUpper(controller)][strings.ToUpper(action)]; exists {
		return true
	}
	return false
}

//获取后台管理员所有用的菜单
func (s *SystemMenu) GetMenuStr() string {
	menus := s.getTree(0, 1)
	str := ""
	for i := 0; i < len(menus); i++ {
		if items, exists := menus[i]["items"]; exists {
			str += `<li class="kt-menu__section "><h4 class="kt-menu__section-text">` + menus[i]["name"].(string) + `</h4><i class="kt-menu__section-icon fa fa-ellipsis-h"></i></li>`

			subItems := items.([]core.A)
			for n := 0; n < len(subItems); n++ {
				if si, exists := subItems[n]["items"]; exists {
					str += `
			<li class="kt-menu__item  kt-menu__item--submenu"><a href="javascript:;" class="kt-menu__link kt-menu__toggle">
				<span class="kt-menu__link-text">` + subItems[n]["name"].(string) + `</span><i class="kt-menu__ver-arrow la la-angle-right"></i></a>
				<div class="kt-menu__submenu "><span class="kt-menu__arrow"></span>
					<ul class="kt-menu__subnav">
			`
					s.GetSubMenuStr(si.([]core.A), &str)
					str += "</ul></div></li>"
				} else {
					str += `
<li class="kt-menu__item " aria-haspopup="true">
	<a href="` + subItems[n]["url"].(string) + `" class="kt-menu__link " target="mainFrame">
		<span class="kt-menu__link-icon"><i class="` + subItems[n]["icon"].(string) + `"></i></span>
		<span class="kt-menu__link-text">` + subItems[n]["name"].(string) + `</span>
	</a>
</li>`
				}
			}

		} else {
			str += `
<li class="kt-menu__item " aria-haspopup="true">
	<a href="` + menus[i]["url"].(string) + `" class="kt-menu__link " target="mainFrame">
		<span class="kt-menu__link-icon"><i class="` + menus[i]["icon"].(string) + `"></i></span>
		<span class="kt-menu__link-text">` + menus[i]["name"].(string) + `</span>
	</a>
</li>`
		}
	}
	return str
}

//获取后台管理员所有用的菜单（不展开二级菜单)
//func (s *SystemMenu) GetMenuStr() string {
//	menus := s.getTree(0, 1)
//	str := ""
//	for i := 0; i < len(menus); i++ {
//		if items, exists := menus[i]["items"]; exists {
//			//str += "<li><a href=\"javascript:;\"> <i class=\"\"></i> <span class=\"menu-title\">"+menus[i]["name"].(string)+"</span> <i class=\"arrow\"></i></a> <ul class=\"collapse\">"
//			//str += `<li class="kt-menu__section "><h4 class="kt-menu__section-text">` + menus[i]["name"].(string) + `</h4><i class="kt-menu__section-icon flaticon-more-v2"></i></li>`
//			str += `
//			<li class="kt-menu__item  kt-menu__item--submenu"><a href="javascript:;" class="kt-menu__link kt-menu__toggle">
//				<span class="kt-menu__link-text">` + menus[i]["name"].(string) + `</span><i class="kt-menu__ver-arrow la la-angle-right"></i></a>
//				<div class="kt-menu__submenu "><span class="kt-menu__arrow"></span>
//					<ul class="kt-menu__subnav">
//			`
//			s.GetSubMenuStr(items.([]core.A), &str)
//			str += "</ul></div></li>"
//		} else {
//			str += "<li class=\"kt-menu__item \" aria-haspopup=\"true\"><a href=\"" + menus[i]["url"].(string) + "\" class=\"kt-menu__link \" target=\"mainFrame\"><span class=\"kt-menu__link-text\">" + menus[i]["name"].(string) + "</span></a></li>"
//		}
//	}
//	return str
//}

//获取下级菜单字符串
func (s *SystemMenu) GetSubMenuStr(items []core.A, str *string) {
	for i := 0; i < len(items); i++ {
		if _, exists := items[i]["items"]; exists {
			*str += `
<li class="kt-menu__item kt-menu__item--submenu" aria-haspopup="true" data-ktmenu-submenu-toggle="hover">
	<a href="javascript:;" class="kt-menu__link kt-menu__toggle">
		<i class="kt-menu__link-bullet kt-menu__link-bullet--dot">
			<span></span>
		</i>
		<span class="kt-menu__link-text">Tabs</span>
		<i class="kt-menu__ver-arrow la la-angle-right"></i>
	</a>
`
		} else {
			*str += `
<li class="kt-menu__item" aria-haspopup="true">
	<a href="` + items[i]["url"].(string) + `" class="kt-menu__link" target="mainFrame">
		<i class="kt-menu__link-bullet kt-menu__link-bullet--dot">
			<span></span>
		</i>
		<span class="kt-menu__link-text">` + items[i]["name"].(string) + `</span>
	</a>
`
		}

		if _, exists := items[i]["items"]; exists {
			*str += `
<div class="kt-menu__submenu ">
	<span class="kt-menu__arrow"></span>
	<ul class="kt-menu__subnav">
`
			s.GetSubMenuStr(items[i]["name"].([]core.A), str)
			*str += "</ul></div>"
		}
		*str += "</li>"
	}
}

func (s *SystemMenu) getTree(parentId int, level int) []core.A {
	menus := s.getChildMenuList(parentId)
	ret := make([]core.A, 0)
	level++
	if len(menus) > 0 {
		for i := 0; i < len(menus); i++ {
			menu := core.A{}
			menu["id"] = menus[i].Id
			menu["app"] = menus[i].App
			menu["name"] = menus[i].Name
			menu["controller"] = menus[i].Controller
			menu["action"] = menus[i].Action
			menu["icon"] = menus[i].Icon
			parameter := fmt.Sprintf("?menuid=%d", menus[i].Id)
			if menus[i].Parameter != "" {
				parameter = fmt.Sprintf("?%s&menuid=%d", menus[i].Parameter, menus[i].Id)
			}
			menu["url"] = fmt.Sprintf("/%s/%s/%s%s", menus[i].App, menus[i].Controller, menus[i].Action, parameter)

			children := s.getTree(menus[i].Id, level)
			if len(children) > 0 && level <= 4 {
				menu["items"] = children
			}
			ret = append(ret, menu)
		}
	}
	return ret
}

//获取所有菜单
func (s *SystemMenu) TreeMenus() []core.A {
	menus := make([]SystemMenu, 0)
	core.DBPool.Slave().OrderBy("listing_order asc, id asc").Find(&menus)
	menuTree := GetMenuTreeInstance(menus)
	menuTree.GetTree(0, "")
	return menuTree.FinalMenus
}

//获取菜单信息
func (s *SystemMenu) Get(id int) bool {
	ok, _ := core.DBPool.Slave().Where("id = ?", id).Get(s)
	return ok
}

//属性菜单
type MenuTree struct {
	icon       []string
	nbsp       string
	menus      []SystemMenu
	FinalMenus []core.A
}

func GetMenuTreeInstance(menus []SystemMenu) MenuTree {
	return MenuTree{
		icon:  []string{"&nbsp;&nbsp;&nbsp;│ ", "&nbsp;&nbsp;&nbsp;├─ ", "&nbsp;&nbsp;&nbsp;└─ "},
		nbsp:  "&nbsp;&nbsp;&nbsp;",
		menus: menus,
	}
}

//获取下级所有菜单
func (m *MenuTree) getChildren(parentId int) []SystemMenu {
	children := make([]SystemMenu, 0)
	for i := 0; i < len(m.menus); i++ {
		if m.menus[i].ParentId == parentId {
			children = append(children, m.menus[i])
		}
	}
	return children
}

//获取树形列表
//parentId 父节点id
//addSpacer 增加的间隔字符
func (m *MenuTree) GetTree(parentId int, addSpacer string) {
	children := m.getChildren(parentId)
	total := len(children)
	if total > 0 {
		for i := 0; i < len(children); i++ {
			j, k := "", ""
			if (i + 1) == total {
				j += m.icon[2]
			} else {
				j += m.icon[1]
				if addSpacer == "" {
					k = ""
				} else {
					k = m.icon[0]
				}
			}
			spacer := ""
			if addSpacer != "" {
				spacer = addSpacer + j
			}
			m.FinalMenus = append(m.FinalMenus, core.A{
				"id":         children[i].Id,
				"name":       template.HTML(spacer + children[i].Name),
				"app":        children[i].App,
				"controller": children[i].Controller,
				"action":     children[i].Action,
				"status":     children[i].Status,
			})
			m.GetTree(children[i].Id, addSpacer+k+m.nbsp)
		}
	}
}
