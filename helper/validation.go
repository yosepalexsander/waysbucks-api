package helper

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	ErrorInvalidFileExtension error = errors.New("invalid file extension")
)

func Validate(value interface{}) (bool, string)  {
	v := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	addTranslation(v, trans, "email", "{0} must be a valid email")
	addTranslation(v, trans, "min", "{0} must be at least {1} char length")
	addTranslation(v, trans, "max", "{0} must be max {1} char length")
	addTranslation(v, trans, "required", "{0} is a required field")
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	err := v.Struct(value)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		msgErr := validationErrors[0].Translate(trans)
		return false, msgErr
	}
	return true, ""
}

func addTranslation(v *validator.Validate, trans ut.Translator, tag string, errMessage string) {
	registerFn := func(ut ut.Translator) error {
		 return ut.Add(tag, errMessage, false)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		 tag := fe.Tag()

		 t, err := ut.T(tag, fe.Field(), param)
		 if err != nil {
				return fe.(error).Error()
		 }
		 return t
	}
	_ = v.RegisterTranslation(tag, trans, registerFn, transFn)
}

func ValidateImageFile(filename string) error {
	regex, _ := regexp.Compile(`\.(jpg|JPEG|png|PNG|svg|SVG)$`)

	// Check if file extension match the regex or not
	if isMatch := regex.MatchString(filename); !isMatch {
		return ErrorInvalidFileExtension
	}

	return nil
}