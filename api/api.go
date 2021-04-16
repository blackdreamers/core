package api

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	et "github.com/go-playground/validator/v10/translations/en"
	zt "github.com/go-playground/validator/v10/translations/zh"
	"golang.org/x/text/language"

	"github.com/blackdreamers/core/constant"
	log "github.com/blackdreamers/go-micro/v3/logger"
)

var (
	enTrans ut.Translator
	zhTrans ut.Translator
)

// Interface api interface
type Interface interface {
	Validator()
	API404(c *gin.Context)
	Verify(c *gin.Context, obj interface{}) (bool, *Response)
	Err(resp *Response, err error)
	Resp(c *gin.Context, r *Response)
}

type API struct{}

type Response struct {
	HttpStatus int
	Code       int
	Message    interface{}
	Data       interface{}
	Err        *Error
}

func (a *API) Validator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		uni := ut.New(en.New(), zh.New())
		enTrans, _ = uni.GetTranslator("en")
		zhTrans, _ = uni.GetTranslator("zh")
		_ = et.RegisterDefaultTranslations(v, enTrans)
		_ = zt.RegisterDefaultTranslations(v, zhTrans)
	}
}

func (a *API) API404(c *gin.Context) {
	a.Resp(c, &Response{HttpStatus: http.StatusNotFound})
}

func (a *API) Verify(c *gin.Context, obj interface{}) (bool, *Response) {
	bind := false
	resp := &Response{}
	if err := c.ShouldBind(obj); err != nil {
		resp.HttpStatus = http.StatusBadRequest
		if vErrors, ok := err.(validator.ValidationErrors); ok {
			var trans ut.Translator
			switch GetLan(c) {
			case "zh":
				trans = zhTrans
			default:
				trans = enTrans
			}
			vErrs := vErrors.Translate(trans)
			errs := make(map[string]string)
			for s := range vErrs {
				errs[strings.ToLower(strings.Split(s, ".")[1])] = vErrs[s]
			}
			resp.Message = errs
		}
		a.Resp(c, resp)
	} else {
		bind = true
	}
	return bind, resp
}

func (a *API) Err(resp *Response, err error) {
	if e, ok := err.(*Error); ok {
		resp.Err = e
	} else {
		resp.HttpStatus = http.StatusInternalServerError

		pc, file, line, ok := runtime.Caller(1)
		if ok {
			log.Fields(
				constant.ErrKey, err,
				"func", runtime.FuncForPC(pc).Name(),
				"file", fmt.Sprintf("%v:%v", file, line),
			).Log(log.ErrorLevel)
		}
	}
}

func (a *API) Resp(c *gin.Context, r *Response) {
	code := http.StatusOK
	if r == nil {
		r = &Response{}
	}
	if r.HttpStatus != 0 {
		code = r.HttpStatus
		r.Code = r.HttpStatus
	}
	if r.Err != nil {
		r.Code = r.Err.ErrCode()
		r.Message = r.Err.ErrMsg(GetLan(c))
	}
	if r.Code == 0 {
		r.Code = code
	}
	if r.Message == nil {
		r.Message = http.StatusText(code)
	}
	c.JSON(code, gin.H{
		"code":      r.Code,
		"msg":       r.Message,
		"data":      r.Data,
		"timestamp": time.Now().Unix(),
	})
}

func GetLan(c *gin.Context) string {
	tags, _, _ := language.ParseAcceptLanguage(c.GetHeader("Accept-Language"))
	if len(tags) > 0 {
		str := strings.ToLower(tags[0].String())
		switch {
		case strings.Contains(str, "zh"):
			return "zh"
		default:
			return "en"
		}
	}
	return "en"
}

func (a *API) Get(c *gin.Context, key interface{}) interface{} {
	session := sessions.Default(c)
	return session.Get(key)
}

func (a *API) Set(c *gin.Context, key interface{}, val interface{}) error {
	session := sessions.Default(c)
	session.Set(key, val)
	if err := session.Save(); err != nil {
		return err
	}
	return nil
}

func (a *API) Delete(c *gin.Context, key interface{}) error {
	session := sessions.Default(c)
	session.Delete(key)
	if err := session.Save(); err != nil {
		return err
	}
	return nil
}
