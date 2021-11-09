package translator

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"go.uber.org/zap"
	"reflect"
	"strings"
)

var translator ut.Translator

// InitTranslator ******** 初始化翻译器********//
func InitTranslator() {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New() // 中文翻译器
		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(zhT, zhT)
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		translator, ok = uni.GetTranslator("zh")
		if !ok {
			msg := fmt.Sprintf("uni.GetTranslator(%s) failed", "zh")
			zap.L().Error("Translator init error", zap.String("error", msg))
		}
		// 注册翻译器
		err := zhTranslations.RegisterDefaultTranslations(v, translator)
		if err != nil {
			zap.L().Error("Translator RegisterDefaultTranslations error", zap.Any("error", err))
		}
	}
	return
}

// RemoveTopStruct ******去除 结构体的名称******//
func removeTopStructToString(fields map[string]string) string {
	res := make([]string, 0)
	for field, err := range fields {
		res = append(res, fmt.Sprintf("%s:%s", field[strings.Index(field, ".")+1:], err))
	}
	return strings.Join(res, ",")
}

// GetErrorMessage ******* 获取错误信息*********//
func GetErrorMessage(err error) string {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return errs.Error()
	} else {
		return removeTopStructToString(errs.Translate(translator))
	}
}
