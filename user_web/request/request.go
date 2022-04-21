package request

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
	"user_web/pkg/captcha"
	"user_web/response"
)

func init() {
	// 最大中文字符规则
	govalidator.AddCustomRule("max_cn", func(field, rule, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "max_cn:"))
		if valLength > l {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度不能超过 %d 个字", l)
		}
		return nil
	})

	// 最小中文字符规则
	govalidator.AddCustomRule("min_cn", func(field, rule, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "min_cn:"))
		if valLength < l {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度不能少于 %d 个字", l)
		}
		return nil
	})

	// 合法字符规则(防xss sql注入字符)
	govalidator.AddCustomRule("legal_char", func(field, rule, message string, value interface{}) error {
		match, err := regexp.Match("^[^`~!#$%^&*+=\\\\|{};:\"',/<>?]*$", []byte(value.(string)))
		if err != nil {
			return err
		}
		if match {
			return nil
		}
		return fmt.Errorf("非法字符")
	})
}

// 参数校验逻辑类型声明
type validRequestFunc func(interface{}, *gin.Context) map[string][]string

// RegisterRequest 注册参数对象
type RegisterRequest struct {
	Username   string `json:"username" valid:"username"`
	Password   string `json:"password" valid:"password"`
	Nickname   string `json:"nickname" valid:"nickname"`
	CaptchaID  string `json:"captcha_id" valid:"captcha_id"`
	CaptchaAns string `json:"captcha_ans" valid:"captcha_ans"`
}

// RegisterRequestValid 注册参数校验逻辑
func RegisterRequestValid(data interface{}, ctx *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"username":    []string{"required", "min:3", "max:64", "alpha_num"},
		"password":    []string{"required", "min:6", "max:20"},
		"nickname":    []string{"min_cn:3", "max_cn:20", "legal_char"},
		"captcha_id":  []string{"required", "max:32", "alpha_num"},
		"captcha_ans": []string{"required", "digits:6"},
	}

	messages := govalidator.MapData{
		"username": []string{
			"required:用户名 必填",
			"min:用户名 长度最小为3",
			"max:用户名 长度最大为64",
			"alpha_num: 用户名 必须是字母或数字",
		},
		"password": []string{
			"required:密码 必填",
			"min:密码 长度最小为6",
			"max:密码 长度最大为20",
		},
		"nickname": []string{
			"min_cn:昵称 长度最小为3",
			"max_cn:昵称 长度最大为20",
			"legal_char:昵称 必须是合法字符",
		},
		"captcha_id": []string{
			"required:captcha_id 必填",
			"max:captcha_id 长度最大为32",
			"alpha_num:captcha_id 必须是字母或数字",
		},
		"captcha_ans": []string{
			"required:验证码 必填",
			"digits:验证码 长度为6",
		},
	}

	errs := validateStruct(data, rules, messages)
	_data := data.(*RegisterRequest)

	if ok := captcha.NewCaptcha().VerifyCaptcha(_data.CaptchaID, _data.CaptchaAns, true); !ok {
		errs["captcha_ans"] = append(errs["captcha_ans"], "图片验证码错误")
	}

	return errs
}

// LoginRequest 登录参数对象
type LoginRequest struct {
	Username   string `json:"username" valid:"username"`
	Password   string `json:"password" valid:"password"`
	CaptchaID  string `json:"captcha_id" valid:"captcha_id"`
	CaptchaAns string `json:"captcha_ans" valid:"captcha_ans"`
}

// LoginRequestValid 登录参数校验逻辑
func LoginRequestValid(data interface{}, ctx *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"username":    []string{"required", "min:3", "max:64", "alpha_num"},
		"password":    []string{"required", "min:6", "max:20"},
		"captcha_id":  []string{"required", "max:32", "alpha_num"},
		"captcha_ans": []string{"required", "digits:6"},
	}

	messages := govalidator.MapData{
		"username": []string{
			"required:用户名 必填",
			"min:用户名 长度最小为3",
			"max:用户名 长度最大为64",
			"alpha_num: 用户名 必须是字母或数字",
		},
		"password": []string{
			"required:密码 必填",
			"min:密码 长度最小为6",
			"max:密码 长度最大为20",
		},
		"captcha_id": []string{
			"required:captcha_id 必填",
			"max:captcha_id 长度最大为32",
			"alpha_num:captcha_id 必须是字母或数字",
		},
		"captcha_ans": []string{
			"required:验证码 必填",
			"digits:验证码 长度为6",
		},
	}

	errs := validateStruct(data, rules, messages)
	_data := data.(*LoginRequest)

	if ok := captcha.NewCaptcha().VerifyCaptcha(_data.CaptchaID, _data.CaptchaAns, true); !ok {
		errs["captcha_ans"] = append(errs["captcha_ans"], "图片验证码错误")
	}

	return errs
}

// EditUserProfileRequest 编辑用户信息参数对象
type EditUserProfileRequest struct {
	Nickname string `json:"nickname" valid:"nickname"`
}

// EditUserProfileRequestValid 编辑用户信息参数校验逻辑
func EditUserProfileRequestValid(data interface{}, ctx *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"nickname": []string{"min_cn:3", "max_cn:20", "legal_char"},
	}

	messages := govalidator.MapData{
		"nickname": []string{
			"min_cn:昵称 长度最小为3",
			"max_cn:昵称 长度最大为20",
			"legal_char:昵称 必须是合法字符",
		},
	}

	return validateStruct(data, rules, messages)
}

// UploadPicRequest 上传图片参数对象
type UploadPicRequest struct {
	Avatar *multipart.FileHeader `form:"avatar" valid:"avatar"`
}

// UploadPicValid 上传图片参数校验逻辑
func UploadPicValid(data interface{}, ctx *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"file:avatar": []string{"required", "ext:png,jpg,jpeg", "size:20971520"},
	}

	messages := govalidator.MapData{
		"file:avatar": []string{
			"ext:ext头像只能上传 png, jpg, jpeg 任意一种的图片",
			"size:头像文件最大不能超过 20MB",
			"required:必须上传图片",
		},
	}

	return validateFile(ctx, rules, messages)
}

// 校验参数对象
func validateStruct(data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: "valid",
	}
	return govalidator.New(opts).ValidateStruct()
}

// 校验文件对象
func validateFile(ctx *gin.Context, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	opts := govalidator.Options{
		Request:       ctx.Request,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: "valid",
	}
	return govalidator.New(opts).Validate()
}

// Validate 参数校验入口
func Validate(ctx *gin.Context, data interface{}, handler validRequestFunc) error {
	if err := ctx.ShouldBind(data); err != nil {
		response.BadRequestRsp(ctx, err)
		return err
	}

	errMap := handler(data, ctx)
	if len(errMap) == 0 {
		return nil
	}

	var errList []string
	for key, val := range errMap {
		errList = append(errList, fmt.Sprintf("%s(%s)", key, strings.Join(val, ";")))
	}
	err := errors.New(strings.Join(errList, " "))

	response.ValidErrorRsp(ctx, err)
	return err
}
