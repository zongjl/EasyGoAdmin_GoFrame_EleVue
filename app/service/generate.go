// +----------------------------------------------------------------------
// | EasyGoAdmin敏捷开发框架 [ 赋能开发者，助力企业发展 ]
// +----------------------------------------------------------------------
// | 版权所有 2019~2022 深圳EasyGoAdmin研发中心
// +----------------------------------------------------------------------
// | Licensed LGPL-3.0 EasyGoAdmin并不是自由软件，未经许可禁止去掉相关版权
// +----------------------------------------------------------------------
// | 官方网站: http://www.easygoadmin.vip
// +----------------------------------------------------------------------
// | Author: @半城风雨 团队荣誉出品
// +----------------------------------------------------------------------
// | 版权和免责声明:
// | 本团队对该软件框架产品拥有知识产权（包括但不限于商标权、专利权、著作权、商业秘密等）
// | 均受到相关法律法规的保护，任何个人、组织和单位不得在未经本团队书面授权的情况下对所授权
// | 软件框架产品本身申请相关的知识产权，禁止用于任何违法、侵害他人合法权益等恶意的行为，禁
// | 止用于任何违反我国法律法规的一切项目研发，任何个人、组织和单位用于项目研发而产生的任何
// | 意外、疏忽、合约毁坏、诽谤、版权或知识产权侵犯及其造成的损失 (包括但不限于直接、间接、
// | 附带或衍生的损失等)，本团队不承担任何法律责任，本软件框架禁止任何单位和个人、组织用于
// | 任何违法、侵害他人合法利益等恶意的行为，如有发现违规、违法的犯罪行为，本团队将无条件配
// | 合公安机关调查取证同时保留一切以法律手段起诉的权利，本软件框架只能用于公司和个人内部的
// | 法律所允许的合法合规的软件产品研发，详细声明内容请阅读《框架免责声明》附件；
// +----------------------------------------------------------------------

/**
 * 代码生成器-服务类
 * @author 半城风雨
 * @since 2021/8/2
 * @File : generate
 */
package service

import (
	"easygoadmin/app/dao"
	"easygoadmin/app/model"
	"easygoadmin/app/utils"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"os"
	"regexp"
	"strings"
)

// 中间件管理服务
var Generate = new(generateService)

type generateService struct{}

func (s *generateService) GetList(req *model.GeneratePageReq) ([]model.GenerateInfo, error) {
	// 查询SQL
	sql := "SHOW TABLE STATUS"
	// 查询条件
	if req != nil {
		// 表名称
		if req.Name != "" {
			sql += " WHERE Name like \"%" + req.Name + "%\""
		}
		// 表描述
		if req.Comment != "" {
			sql += " WHERE Comment like \"%" + req.Comment + "%\""
		}
	}
	// 对象转换
	var list []model.GenerateInfo
	err := g.DB().GetScan(&list, sql)
	if err != nil {
		return nil, err
	}
	// 返回结果
	return list, nil
}

func (s *generateService) Generate(req model.GenerateFileReq, r *ghttp.Request) error {
	if utils.AppDebug() {
		return gerror.New("演示环境，暂无权限操作")
	}
	// 数据表名
	tableName := req.Name
	// 数据表描述
	moduleTitle := req.Comment
	// 替换“表”
	if gstr.Contains(moduleTitle, "表") {
		moduleTitle = gstr.Replace(moduleTitle, "表", "")
	}
	// 替换“管理”
	if gstr.Contains(moduleTitle, "管理") {
		moduleTitle = gstr.Replace(moduleTitle, "管理", "")
	}
	// 模型名称
	moduleName := gstr.Replace(tableName, "sys_", "")
	// 作者名称
	authorName := "半城风雨"

	// 获取字段列表
	columnList, err := GetColumnList(tableName)
	if err != nil {
		return err
	}

	// 生成控制器
	if err := GenerateController(r, columnList, authorName, moduleName, moduleTitle); err != nil {
		return err
	}

	// 生成控制器
	if err := GenerateModel(r, columnList, authorName, moduleName, moduleTitle); err != nil {
		return err
	}

	// 生成服务类
	if err := GenerateService(r, columnList, authorName, moduleName, moduleTitle); err != nil {
		return err
	}

	// 生成模块index.html
	if err := GenerateIndex(r, columnList, moduleName, moduleTitle); err != nil {
		return err
	}

	// 生成模块edit.html
	if err := GenerateEdit(r, columnList, moduleName, moduleTitle); err != nil {
		return err
	}

	// 生成菜单权限
	if err := GeneratePermission(moduleName, moduleTitle, utils.Uid(r)); err != nil {
		return err
	}

	// 生成路由
	if err := GenerateRouter(r, columnList, authorName, moduleName, moduleTitle); err != nil {
		return err
	}

	return nil
}

// 生成控制器
func GenerateController(r *ghttp.Request, dataList *garray.Array, authorName string, moduleName string, moduleTitle string) error {
	// 初始化表单数组
	columnList := make([]map[string]interface{}, 0)
	for i := 0; i < dataList.Len(); i++ {
		// 当前元素
		data, _ := dataList.Get(i)
		// 类型转换
		item := data.(map[string]interface{})
		// 字段列名
		columnName := gconv.String(item["columnName"])
		// 移除部分非表单字段
		if columnName == "id" ||
			columnName == "create_user" ||
			columnName == "create_time" ||
			columnName == "update_user" ||
			columnName == "update_time" ||
			columnName == "mark" {
			continue
		}
		// 加入数组
		columnList = append(columnList, item)
	}

	// 读取HTML模板并绑定数据
	if filePath, err := r.Response.ParseTpl("controller.html", g.Map{
		"author":      authorName,
		"since":       gtime.Now().Format("Y/m/d"),
		"moduleName":  moduleName,
		"entityName":  gstr.UcWords(moduleName),
		"moduleTitle": moduleTitle,
		"columnList":  columnList,
	}); err == nil {
		// 获取项目根目录
		curDir, err := os.Getwd()
		if err != nil {
			return err
		}
		// 文件路径
		fileName := strings.Join([]string{curDir, "/app/controller/", moduleName, ".go"}, "")
		// 删除现有文件
		if err := gfile.Remove(fileName); err != nil {
			return err
		}
		// 写入文件
		if !gfile.Exists(fileName) {
			f, err := gfile.Create(fileName)
			if err == nil {
				_, err := f.WriteString(filePath)
				if err != nil {
					return err
				}
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// 生成模型
func GenerateModel(r *ghttp.Request, dataList *garray.Array, authorName string, moduleName string, moduleTitle string) error {
	// 初始化查询条件
	queryList := make([]map[string]interface{}, 0)
	// 初始化表单数组
	columnList := make([]map[string]interface{}, 0)
	for i := 0; i < dataList.Len(); i++ {
		// 当前元素
		data, _ := dataList.Get(i)
		// 类型转换
		item := data.(map[string]interface{})
		// 字段列名
		columnName := gconv.String(item["columnName"])
		// 移除部分非表单字段
		if columnName == "id" ||
			columnName == "create_user" ||
			columnName == "create_time" ||
			columnName == "update_user" ||
			columnName == "update_time" ||
			columnName == "mark" {
			continue
		}
		// 常规字段查询条件
		if columnName == "name" || columnName == "title" {
			// 加入查询条件数组
			queryList = append(queryList, item)
		}
		// 常规下拉选择查询条件
		if _, ok := item["columnValue"]; ok {
			queryList = append(queryList, item)
		}
		// 加入数组
		columnList = append(columnList, item)
	}

	// 读取HTML模板并绑定数据
	if filePath, err := r.Response.ParseTpl("model.html", g.Map{
		"author":      authorName,
		"since":       gtime.Now().Format("Y/m/d"),
		"moduleName":  moduleName,
		"entityName":  gstr.UcWords(moduleName),
		"moduleTitle": moduleTitle,
		"columnList":  columnList,
		"queryList":   queryList,
	}); err == nil {
		// 获取项目根目录
		curDir, err := os.Getwd()
		if err != nil {
			return err
		}
		// 文件路径
		fileName := strings.Join([]string{curDir, "/app/model/", moduleName, ".go"}, "")
		// 删除现有文件
		if err := gfile.Remove(fileName); err != nil {
			return err
		}
		// 写入文件
		if !gfile.Exists(fileName) {
			f, err := gfile.Create(fileName)
			if err == nil {
				_, err := f.WriteString(filePath)
				if err != nil {
					return err
				}
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// 生成服务类
func GenerateService(r *ghttp.Request, dataList *garray.Array, authorName string, moduleName string, moduleTitle string) error {
	// 初始化查询条件
	queryList := make([]map[string]interface{}, 0)
	// 初始化表单数组
	columnList := make([]map[string]interface{}, 0)
	for i := 0; i < dataList.Len(); i++ {
		// 当前元素
		data, _ := dataList.Get(i)
		// 类型转换
		item := data.(map[string]interface{})
		// 字段列名
		columnName := gconv.String(item["columnName"])
		// 移除部分非表单字段
		if columnName == "id" ||
			columnName == "create_user" ||
			columnName == "create_time" ||
			columnName == "update_user" ||
			columnName == "update_time" ||
			columnName == "mark" {
			continue
		}
		// 常规字段查询条件
		if columnName == "name" || columnName == "title" {
			// 加入查询条件数组
			queryList = append(queryList, item)
		}
		// 常规下拉选择查询条件
		if _, ok := item["columnValue"]; ok {
			queryList = append(queryList, item)
		}
		// 加入数组
		columnList = append(columnList, item)
	}
	// 读取HTML模板并绑定数据
	if filePath, err := r.Response.ParseTpl("service.html", g.Map{
		"author":      authorName,
		"since":       gtime.Now().Format("Y/m/d"),
		"moduleName":  moduleName,
		"entityName":  gstr.UcWords(moduleName),
		"moduleTitle": moduleTitle,
		"columnList":  columnList,
		"queryList":   queryList,
	}); err == nil {
		// 获取项目根目录
		curDir, err := os.Getwd()
		if err != nil {
			return err
		}
		// 文件路径
		fileName := strings.Join([]string{curDir, "/app/service/", moduleName, ".go"}, "")
		// 删除现有文件
		if err := gfile.Remove(fileName); err != nil {
			return err
		}
		// 写入文件
		if !gfile.Exists(fileName) {
			f, err := gfile.Create(fileName)
			if err == nil {
				_, err := f.WriteString(filePath)
				if err != nil {
					return err
				}
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// 生成列表页
func GenerateIndex(r *ghttp.Request, dataList *garray.Array, moduleName string, moduleTitle string) error {
	// 初始化查询条件
	queryList := make([]map[string]interface{}, 0)
	// 初始化表单数组
	columnList := make([]map[string]interface{}, 0)
	for i := 0; i < dataList.Len(); i++ {
		// 当前元素
		data, _ := dataList.Get(i)
		// 类型转换
		item := data.(map[string]interface{})
		// 字段列名
		columnName := item["columnName"]
		if columnName == "name" || columnName == "title" {
			// 加入查询条件数组
			queryList = append(queryList, item)
		}
		// 判断是否有columnValue键值
		if _, ok := item["columnValue"]; ok {
			// 加入查询条件数组
			queryList = append(queryList, item)
		}
		// 移除部分非表单字段
		if columnName == "id" ||
			columnName == "create_user" ||
			columnName == "update_user" ||
			columnName == "mark" {
			continue
		}
		// 加入数组
		columnList = append(columnList, item)
	}

	// 读取HTML模板并绑定数据
	if filePath, err := r.Response.ParseTpl("index.html", g.Map{
		"moduleName":  moduleName,
		"entityName":  gstr.UcWords(moduleName),
		"moduleTitle": moduleTitle,
		"columnList":  columnList,
		"queryList":   queryList,
	}); err == nil {
		// 获取项目根目录
		curDir, err := os.Getwd()
		if err != nil {
			return err
		}
		// 文件路径
		fileName := strings.Join([]string{curDir, "/evui/src/views/tool/example//", moduleName, "/index.vue"}, "")
		// 删除现有文件
		if err := gfile.Remove(fileName); err != nil {
			return err
		}
		// 写入文件
		if !gfile.Exists(fileName) {
			f, err := gfile.Create(fileName)
			if err == nil {
				_, err := f.WriteString(filePath)
				if err != nil {
					return err
				}
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// 生成编辑表单
func GenerateEdit(r *ghttp.Request, dataList *garray.Array, moduleName string, moduleTitle string) error {
	// 初始化表单数组
	formList := make([]map[string]interface{}, 0)
	// 初始化图片数组
	imageList := make([]map[string]interface{}, 0)
	// 初始化多行数组
	rowsList := make([]map[string]interface{}, 0)
	for i := 0; i < dataList.Len(); i++ {
		// 当前元素
		data, _ := dataList.Get(i)
		// 类型转换
		item := data.(map[string]interface{})
		// 字段列名
		columnName := gconv.String(item["columnName"])
		// 移除部分非表单字段
		if columnName == "id" ||
			columnName == "create_user" ||
			columnName == "create_time" ||
			columnName == "update_user" ||
			columnName == "update_time" ||
			columnName == "mark" {
			continue
		}
		// 图片上传
		if _, isImage := item["columnImage"]; isImage {
			// 加入数组
			imageList = append(imageList, item)
			continue
		}

		// 多行文本输入
		if _, isText := item["columnText"]; isText {
			// 加入数组
			rowsList = append(rowsList, item)
			continue
		}
		// 加入数组
		formList = append(formList, item)
	}

	// 初始化数据列数组
	columnList := make([]map[string]interface{}, 0)

	// 根据控制的个数实行分列显示(一行两列)
	if len(formList)+len(imageList)+len(rowsList) > 50 {
		// 一行两列排列
	} else {
		// 单行排列
		columnList = formList
		// 图片
		if len(imageList) > 0 {
			// 遍历
			for _, v := range imageList {
				columnList = append(columnList, v)
			}
		}
		// 多行文本
		if len(rowsList) > 0 {
			// 遍历
			for _, v := range rowsList {
				columnList = append(columnList, v)
			}
		}
	}

	// 读取HTML模板并绑定数据
	if filePath, err := r.Response.ParseTpl("edit.html", g.Map{
		"moduleName":  moduleName,
		"entityName":  gstr.UcWords(moduleName),
		"moduleTitle": moduleTitle,
		"columnList":  columnList,
	}); err == nil {
		// 获取项目根目录
		curDir, err := os.Getwd()
		if err != nil {
			return err
		}
		// 文件路径
		fileName := strings.Join([]string{curDir, "/evui/src/views/tool/example/", moduleName, "/" + moduleName + "-edit.vue"}, "")
		// 删除现有文件
		if err := gfile.Remove(fileName); err != nil {
			return err
		}
		// 写入文件
		if !gfile.Exists(fileName) {
			f, err := gfile.Create(fileName)
			if err == nil {
				_, err := f.WriteString(filePath)
				if err != nil {
					return err
				}
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// 获取表字段列表
func GetColumnList(tableName string) (*garray.Array, error) {
	// 获取数据库名
	DbName, err := utils.GetDatabase()
	if err != nil {
		return nil, err
	}

	// 获取字段列表
	data, err := g.DB().GetAll("SELECT COLUMN_NAME,COLUMN_DEFAULT,DATA_TYPE,COLUMN_TYPE,COLUMN_COMMENT FROM information_schema.`COLUMNS` where TABLE_SCHEMA = ? AND TABLE_NAME = ?", DbName, tableName)
	if err != nil {
		return nil, err
	}

	// 初始化数组
	result := garray.NewArray(true)
	for _, v := range data {
		// 初始化Map
		item := make(map[string]interface{})
		// 字段列名
		columnName := v["COLUMN_NAME"].String()
		// 系统常规字段直接跳过
		if columnName == "create_user" ||
			columnName == "create_time" ||
			columnName == "update_user" ||
			columnName == "update_time" ||
			columnName == "mark" {
			continue
		}
		item["columnName"] = columnName
		// 字段名称驼峰格式一
		columnName2 := gstr.UcWords(columnName)
		if gstr.Contains(columnName, "_") {
			nameArr := gstr.Split(columnName, "_")
			columnName2 = gstr.UcWords(nameArr[0]) + gstr.UcWords(nameArr[1])
		}
		item["columnName2"] = columnName2

		// 字段名称驼峰格式二
		columnName3 := columnName
		if gstr.Contains(columnName, "_") {
			nameArr := gstr.Split(columnName, "_")
			columnName3 = nameArr[0] + gstr.UcWords(nameArr[1])
		}
		item["columnName3"] = columnName3

		// 字段默认值
		item["columnDefault"] = v["COLUMN_DEFAULT"].String()
		// 字段数据类型
		dataType := v["DATA_TYPE"].String()
		item["dataType"] = dataType
		if dataType == "int" || dataType == "tinyint" || dataType == "smallint" {
			// 整形
			item["columnType"] = "int"
		} else if dataType == "bigint" {
			item["columnType"] = "int64"
		} else if dataType == "datetime" {
			item["columnType"] = "*gtime.Time"
		} else {
			// 字符串类型
			item["columnType"] = "string"
		}

		// 字段描述
		columnComment := v["COLUMN_COMMENT"].String()
		item["columnComment"] = columnComment
		// 判断是否有规则描述
		if gstr.Contains(columnComment, ":") || gstr.Contains(columnComment, "：") {
			// 正则根据冒号分裂字符串
			re := regexp.MustCompile("[：；]")
			commentItem := gstr.Split(re.ReplaceAllString(columnComment, "|"), "|")
			// 字段标题
			item["columnTitle"] = commentItem[0]

			// 字段描述数据处理
			commentStr := gstr.Replace(commentItem[1], " ", "|")
			commentArr := gstr.Split(commentStr, "|")

			// 实例化字段描述参数数组
			columnValue := make([]string, 0)
			// 参数值Map列表
			columnValueList := make(map[int]string)
			// 实例化字段描述文字数组
			columnSwitchValue := make([]string, 0)
			for _, v := range commentArr {
				// 正则提取数字键
				regexp := regexp.MustCompile(`[0-9]+`)
				vItem := regexp.FindStringSubmatch(v)
				// 键
				key := vItem[0]
				// 值
				value := gstr.Replace(v, vItem[0], "")
				// 加入数组
				columnValue = append(columnValue, key+"="+value)
				// 参数值Map
				columnValueList[gconv.Int(key)] = value
				// 开关专用参数值
				columnSwitchValue = append(columnSwitchValue, value)
			}
			// 字符串逗号拼接
			item["columnValue"] = gstr.Join(columnValue, ",")
			item["columnValueList"] = columnValueList

			// 开关判断
			if columnName == "status" || gstr.SubStr(columnName, 0, 3) == "is_" {
				item["columnSwitch"] = true
				item["columnSwitchValue"] = gstr.Join(columnSwitchValue, "|")
				// 方法名处理
				columnSwitchName := ""
				if gstr.Contains(columnName, "_") {
					switchArr := gstr.Split(columnName, "_")
					columnSwitchName = "set" + gstr.UcWords(switchArr[0]) + gstr.UcWords(switchArr[1])
				} else {
					columnSwitchName = "set" + gstr.UcWords(columnName)
				}
				item["columnSwitchName"] = columnSwitchName
			}
		} else {
			// 字段标题
			item["columnTitle"] = columnComment
		}

		// 判断是否是图片
		if gstr.Contains(columnName, "cover") ||
			gstr.Contains(columnName, "avatar") ||
			gstr.Contains(columnName, "image") ||
			gstr.Contains(columnName, "logo") ||
			gstr.Contains(columnName, "pic") {
			item["columnImage"] = true
		}

		// 判断是否多行文本或富文本
		if gstr.Contains(columnName, "note") ||
			gstr.Contains(columnName, "remark") ||
			gstr.Contains(columnName, "content") ||
			gstr.Contains(columnName, "description") ||
			gstr.Contains(columnName, "intro") {
			item["columnText"] = true
		}

		// 加入数组
		result.Append(item)
	}
	return result, nil
}

// 生成菜单和权限
func GeneratePermission(modelName string, modelTitle string, userId int) error {
	// 查询记录
	info, err := dao.Menu.Where("permission", "sys:"+modelName+":view").FindOne()
	if err != nil {
		return err
	}
	// 创建菜单
	var entity model.Menu
	entity.Title = modelTitle
	entity.Icon = "el-icon-setting"
	entity.Path = "/tool/example/" + modelName
	entity.Component = entity.Path
	entity.ParentId = 164
	entity.Type = 0
	entity.Permission = "sys:" + modelName + ":view"
	entity.Status = 1
	entity.Target = "_self"
	entity.Sort = 10
	entity.CreateUser = userId
	entity.CreateTime = gtime.Now()
	entity.UpdateUser = userId
	entity.UpdateTime = gtime.Now()
	entity.Mark = 1
	// 记录ID
	menuId := 0
	// 插入或更新记录
	if info == nil {
		// 创建菜单
		// 插入记录
		result, err := dao.Menu.Insert(entity)
		if err != nil {
			return err
		}
		// 获取插入ID
		id, err := result.LastInsertId()
		if err != nil || id == 0 {
			return err
		}
		// 菜单ID
		menuId = gconv.Int(id)
	} else {
		// 更新菜单
		entity.Id = info.Id
		// 更新记录
		result, err := dao.Menu.Save(entity)
		if err != nil {
			return err
		}
		// 获取受影响行数
		rows, err := result.RowsAffected()
		if err != nil || rows == 0 {
			return err
		}
		// 菜单ID
		menuId = info.Id
	}

	// 删除现有节点
	result, err := dao.Menu.Delete("parent_id=?", menuId)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows < 0 {
		return err
	}

	// 创建节点
	funcList := []int{1, 5, 10, 15, 20, 25, 30}
	for _, v := range funcList {
		// 实例化对象
		var item model.Menu
		item.ParentId = menuId
		item.Type = 1
		item.Status = 1
		item.Target = "_self"
		item.Sort = v
		item.CreateUser = userId
		item.CreateTime = gtime.Now()
		item.Mark = 1

		// 权限节点
		if v == 1 {
			// 列表
			item.Title = "查询" + modelTitle
			item.Path = "/" + modelName + "/list"
			item.Permission = "sys:" + modelName + ":list"
			item.Method = "GET"
		} else if v == 5 {
			// 添加
			item.Title = "添加" + modelTitle
			item.Path = "/" + modelName + "/add"
			item.Permission = "sys:" + modelName + ":add"
			item.Method = "POST"
		} else if v == 10 {
			// 修改
			item.Title = "修改" + modelTitle
			item.Path = "/" + modelName + "/update"
			item.Permission = "sys:" + modelName + ":update"
			item.Method = "PUT"
		} else if v == 15 {
			// 删除
			item.Title = "删除" + modelTitle
			item.Path = "/" + modelName + "/delete"
			item.Permission = "sys:" + modelName + ":delete"
			item.Method = "DELETE"
		} else if v == 20 {
			// 详情
			item.Title = modelTitle + "详情"
			item.Path = "/" + modelName + "/detail"
			item.Permission = "sys:" + modelName + ":detail"
			item.Method = "GET"
		} else if v == 25 {
			// 状态
			item.Title = "设置状态"
			item.Path = "/" + modelName + "/status"
			item.Permission = "sys:" + modelName + ":status"
			item.Method = "PUT"
		} else if v == 30 {
			// 批量删除
			item.Title = "批量删除"
			item.Path = "/" + modelName + "/dall"
			item.Permission = "sys:" + modelName + ":dall"
			item.Method = "DELETE"
		}

		// 插入数据
		_, err := dao.Menu.Insert(item)
		if err != nil {
			break
		}
	}
	return nil
}

// 生成路由文件
func GenerateRouter(r *ghttp.Request, dataList *garray.Array, authorName string, moduleName string, moduleTitle string) error {
	// 初始化表单数组
	columnList := make([]map[string]interface{}, 0)
	for i := 0; i < dataList.Len(); i++ {
		// 当前元素
		data, _ := dataList.Get(i)
		// 类型转换
		item := data.(map[string]interface{})
		// 字段列名
		columnName := gconv.String(item["columnName"])
		// 移除部分非表单字段
		if columnName == "id" ||
			columnName == "create_user" ||
			columnName == "create_time" ||
			columnName == "update_user" ||
			columnName == "update_time" ||
			columnName == "mark" {
			continue
		}
		// 加入数组
		columnList = append(columnList, item)
	}
	// 读取HTML模板并绑定数据
	if filePath, err := r.Response.ParseTpl("router.html", g.Map{
		"author":      authorName,
		"since":       gtime.Now().Format("Y/m/d"),
		"moduleName":  moduleName,
		"entityName":  gstr.UcWords(moduleName),
		"moduleTitle": moduleTitle,
		"columnList":  columnList,
	}); err == nil {
		// 获取项目根目录
		curDir, err := os.Getwd()
		if err != nil {
			return err
		}
		// 文件路径
		fileName := strings.Join([]string{curDir, "/router/", moduleName, ".go"}, "")
		// 删除现有文件
		if err := gfile.Remove(fileName); err != nil {
			return err
		}
		// 写入文件
		if !gfile.Exists(fileName) {
			f, err := gfile.Create(fileName)
			if err == nil {
				_, err := f.WriteString(filePath)
				if err != nil {
					return err
				}
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}
