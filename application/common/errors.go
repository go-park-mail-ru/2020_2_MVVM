package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/mailru/easyjson"
	"net/http"
)

const (
	EmptyFieldErr = "Обязательные поля не заполнены."
	SessionErr    = "Ошибка авторизации. Попробуйте авторизоваться повторно."
	AuthRequiredErr = "Необходима авторизация."
	DataBaseErr   = "Что-то пошло не так. Попробуйте позже."
	UserExistErr  = "Пользователь уже существует."
	AuthErr       = "Пользователь с такими данными не зарегистрирован."
	WrongPasswd   = "Неверное имя пользователя или пароль."
	EmpHaveComp   = "Работодатель может являться представителем только одной компании."
	NoRecommendation = "Для этого пользователя еще нет рекомендаций."
)

type Err struct {
	code    int         `json:"code"`
	message string      `json:"message"`
	meta    interface{} `json:"meta"`
}

func (e Err) Code() int         { return e.code }
func (e Err) Error() string     { return e.message }
func (e Err) Meta() interface{} { return e.meta }

func (e Err) MarshalJSON() ([]byte, error) {
	ret := map[string]interface{}{
		"code":    e.code,
		"message": e.message,
	}
	if e.meta != nil {
		ret["meta"] = e.meta
	}
	return json.Marshal(ret)
}

func (e Err) String() string {
	data, _ := e.MarshalJSON()
	return string(data)
}

func NewErr(code int, message string, meta interface{}) Err {
	return Err{
		code:    code,
		message: message,
		meta:    meta,
	}
}

func WriteErrResponse(ctx *gin.Context, code int, message string) {
	resp := models.RespError{
		Err:   message,
	}
	ctx.Status(code)
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(resp, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}
