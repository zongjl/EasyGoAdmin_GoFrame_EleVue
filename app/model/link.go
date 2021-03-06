// ==========================================================================
// This is auto-generated by gf cli tool. Fill this file as you wish.
// ==========================================================================

package model

import (
	"easygoadmin/app/model/internal"
)

// Link is the golang structure for table sys_link.
type Link internal.Link

// Fill with you ideas below.

// 分页查询条件
type LinkPageReq struct {
	Name     string `p:"name"`     // 友链名称
	Type     int    `p:"type"`     // 友链类型
	Platform int    `p:"platform"` // 投放平台
	Page     int    `p:"page"`     // 页码
	Limit    int    `p:"limit"`    // 每页数
}

// 添加友链
type LinkAddReq struct {
	Name     string `p:"name"        v:"required#友链名称不能为空"` // 友链名称
	Type     int    `p:"type"        v:"required#请选择友链类型"`  // 类型：1友情链接 2合作伙伴
	Url      string `p:"url"`                               // 友链地址
	ItemId   int    `p:"item_id"`                           // 站点ID
	CateId   int    `p:"cate_id"`                           // 栏目ID
	Platform int    `p:"platform"    v:"required#请选择投放平台"`  // 平台：1PC站 2WAP站 3微信小程序 4APP应用
	Form     int    `p:"form"        v:"required#请选择友链形式"`  // 友链形式：1文字链接 2图片链接
	Image    string `p:"image"`                             // 友链图片
	Status   int    `p:"status"      v:"required#请选择友链状态"`  // 状态：1在用 2停用
	Sort     int    `p:"sort"        v:"required#排序号不能为空"`  // 显示顺序
	Note     string `p:"note"`                              // 备注
}

// 修改友链
type LinkUpdateReq struct {
	Id       int    `p:"id" v:"required#主键ID不能为空"`
	Name     string `p:"name"        v:"required#友链名称不能为空"` // 友链名称
	Type     int    `p:"type"        v:"required#请选择友链类型"`  // 类型：1友情链接 2合作伙伴
	Url      string `p:"url"`                               // 友链地址
	ItemId   int    `p:"item_id"`                           // 站点ID
	CateId   int    `p:"cate_id"`                           // 栏目ID
	Platform int    `p:"platform"    v:"required#请选择投放平台"`  // 平台：1PC站 2WAP站 3微信小程序 4APP应用
	Form     int    `p:"form"        v:"required#请选择友链形式"`  // 友链形式：1文字链接 2图片链接
	Image    string `p:"image"`                             // 友链图片
	Status   int    `p:"status"      v:"required#请选择友链状态"`  // 状态：1在用 2停用
	Sort     int    `p:"sort"        v:"required#排序号不能为空"`  // 显示顺序
	Note     string `p:"note"`                              // 备注
}

// 删除友链
type LinkDeleteReq struct {
	Ids string `p:"ids"  v:"required#请选择要删除的数据记录"`
}

// 设置状态
type LinkStatusReq struct {
	Id     int `p:"id" v:"required#主键ID不能为空"`
	Status int `p:"status"    v:"required#状态不能为空"`
}

// 友链信息
type LinkInfoVo struct {
	Link
	TypeName     string `json:"typeName"`     // 友链类型
	FormName     string `json:"formName"`     // 友链形式
	PlatformName string `json:"platformName"` // 投放平台
}
