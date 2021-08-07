// ==========================================================================
// This is auto-generated by gf cli tool. DO NOT EDIT THIS FILE MANUALLY.
// ==========================================================================

package internal

import (
    "github.com/gogf/gf/os/gtime"
)

// ItemCate is the golang structure for table sys_item_cate.
type ItemCate struct {
    Id         int         `orm:"id,primary"  json:"id"`         // ID                     
    Name       string      `orm:"name"        json:"name"`       // 栏目名称               
    Pid        int         `orm:"pid"         json:"pid"`        // 父级ID                 
    ItemId     uint        `orm:"item_id"     json:"itemId"`     // 站点ID                 
    Pinyin     string      `orm:"pinyin"      json:"pinyin"`     // 拼音(全)               
    Code       string      `orm:"code"        json:"code"`       // 拼音(简)               
    IsCover    int         `orm:"is_cover"    json:"isCover"`    // 是否有封面：1是 2否    
    Cover      string      `orm:"cover"       json:"cover"`      // 封面                   
    Status     int         `orm:"status"      json:"status"`     // 状态：1启用 2停用      
    Note       string      `orm:"note"        json:"note"`       // 备注                   
    Sort       uint        `orm:"sort"        json:"sort"`       // 排序                   
    CreateUser int         `orm:"create_user" json:"createUser"` // 添加人                 
    CreateTime *gtime.Time `orm:"create_time" json:"createTime"` // 添加时间               
    UpdateUser int         `orm:"update_user" json:"updateUser"` // 更新人                 
    UpdateTime *gtime.Time `orm:"update_time" json:"updateTime"` // 更新时间               
    Mark       int         `orm:"mark"        json:"mark"`       // 有效标识(1正常 0删除)  
}