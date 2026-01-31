package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// 字段名称映射（JSON 字段名 -> 中文名）
var fieldNames = map[string]string{
	// 用户相关
	"username":         "用户名",
	"password":         "密码",
	"current_password": "当前密码",
	"new_password":     "新密码",
	"nickname":         "昵称",
	"email":            "邮箱",
	"bio":              "个人简介",
	"avatar_style":     "头像样式",
	"avatar_color":     "头像颜色",
	"confirm":          "确认信息",

	// 笔记本相关
	"name": "名称",

	// 笔记相关
	"notebook_id": "笔记本",
	"title":       "标题",
	"tag_ids":     "标签",
	"note_ids":    "笔记",
	"content":     "内容",
	"action":      "操作类型",

	// 游戏化相关
	"daily_char_goal": "每日目标",
}

// 特殊字段的自定义消息 (字段名 -> tag -> 消息)
var customMessages = map[string]map[string]string{
	"password": {
		"required": "密码不能为空",
		"min":      "密码长度不能少于6位",
		"max":      "密码长度不能超过50位",
	},
	"new_password": {
		"required": "新密码不能为空",
		"min":      "新密码长度不能少于6位",
		"max":      "新密码长度不能超过50位",
	},
	"current_password": {
		"required": "当前密码不能为空",
		"min":      "当前密码长度不能少于6位",
	},
	"username": {
		"required": "用户名不能为空",
		"min":      "用户名长度不能少于3个字符",
		"max":      "用户名长度不能超过100个字符",
	},
	"confirm": {
		"required": "请输入确认信息",
		"eq":       "请输入 DELETE 确认注销",
	},
	"notebook_id": {
		"required": "请选择笔记本",
	},
	"note_ids": {
		"required": "请选择要操作的笔记",
		"min":      "请至少选择一条笔记",
		"max":      "批量操作最多支持100条笔记",
	},
	"tag_ids": {
		"required": "请选择标签",
	},
	"action": {
		"required": "请选择操作类型",
		"oneof":    "操作类型无效",
	},
	"daily_char_goal": {
		"required": "每日目标不能为空",
		"min":      "每日目标不能少于100字符",
		"max":      "每日目标不能超过10000字符",
	},
	"name": {
		"required": "名称不能为空",
		"max":      "名称长度不能超过255个字符",
	},
}

// TranslateValidationError 将验证错误转换为用户友好的中文提示
// 只返回第一个错误
func TranslateValidationError(err error) string {
	// 尝试转换为 validator.ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// 如果不是验证错误，返回原始错误信息
		return "请求参数格式错误"
	}

	if len(validationErrors) == 0 {
		return "请求参数错误"
	}

	// 只处理第一个错误
	fe := validationErrors[0]
	field := strings.ToLower(fe.Field())
	tag := fe.Tag()

	// 1. 先检查是否有自定义消息
	if fieldMsgs, ok := customMessages[field]; ok {
		if msg, ok := fieldMsgs[tag]; ok {
			return msg
		}
	}

	// 2. 获取字段中文名
	fieldName := fieldNames[field]
	if fieldName == "" {
		fieldName = fe.Field()
	}

	// 3. 根据 tag 生成通用消息
	switch tag {
	case "required":
		return fmt.Sprintf("%s不能为空", fieldName)
	case "min":
		return fmt.Sprintf("%s长度不能少于%s", fieldName, fe.Param())
	case "max":
		return fmt.Sprintf("%s长度不能超过%s", fieldName, fe.Param())
	case "eq":
		return fmt.Sprintf("%s必须等于 %s", fieldName, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s必须是以下值之一: %s", fieldName, fe.Param())
	case "email":
		return fmt.Sprintf("%s格式不正确", fieldName)
	case "url":
		return fmt.Sprintf("%s必须是有效的URL", fieldName)
	case "len":
		return fmt.Sprintf("%s长度必须为%s", fieldName, fe.Param())
	case "gt":
		return fmt.Sprintf("%s必须大于%s", fieldName, fe.Param())
	case "gte":
		return fmt.Sprintf("%s必须大于或等于%s", fieldName, fe.Param())
	case "lt":
		return fmt.Sprintf("%s必须小于%s", fieldName, fe.Param())
	case "lte":
		return fmt.Sprintf("%s必须小于或等于%s", fieldName, fe.Param())
	default:
		return fmt.Sprintf("%s验证失败", fieldName)
	}
}
