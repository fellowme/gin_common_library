package token

import (
	"encoding/json"
	"errors"
	gin_config "github.com/fellowme/gin_commom_library/config"
	gin_const "github.com/fellowme/gin_commom_library/const"
	gin_remote_service "github.com/fellowme/gin_commom_library/remote_service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/url"
)

func VerifyUserId(c *gin.Context) error {
	userId := c.Request.Header.Get("user_id")
	if userId == "" {
		userId = c.Request.Header.Get("userid")
	}
	token := c.Request.Header.Get("token")
	if userId == "" || token == "" {
		zap.L().Error(gin_const.VerifyUserFail,
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.Path),
			zap.Any("error", "user_id token 两者都不能为空 "),
			zap.String("user_id", userId),
			zap.String("token", token),
		)
		return errors.New(gin_const.VerifyUserFail)
	}
	values := url.Values{}
	values.Add("uid", userId)
	values.Add("token", token)
	resp, err := gin_remote_service.PostForm(gin_config.ServerConfigSettings.Server.PassportUrl, values, gin_remote_service.WithJaegerContext(c))
	if err != nil {
		return err
	}
	ret := &struct {
		Status int `json:"status"`
	}{}
	err = json.Unmarshal(resp.Data, ret)
	if err != nil {
		return err
	}
	if ret.Status != 0 {
		return errors.New("token_verify status failed")
	}
	return nil
}
