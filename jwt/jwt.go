package jwt

import (
	"errors"
	gin_config "github.com/fellowme/gin_common_library/config"
	gin_const "github.com/fellowme/gin_common_library/const"
	gin_util "github.com/fellowme/gin_common_library/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

/*
	JwtAuth gin中间件验证jwt
*/
var (
	TokenExpired     = errors.New(gin_const.TokenExpiredTip)
	TokenNotValidYet = errors.New(gin_const.TokenNotValidYetTip)
	TokenMalformed   = errors.New(gin_const.TokenMalformedTip)
	TokenInvalid     = errors.New(gin_const.TokenInvalidTip)
)

/*
	JwtAuth gin中间件验证jwt
*/
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := c.GetQuery("token")
		if !ok {
			token = c.Request.Header.Get("token")
		}
		if token == "" {
			gin_util.ReturnResponse(http.StatusOK, gin_util.FailCode, gin_const.TokenNotEmptyTip, nil, c)
			c.Abort()
			return
		}
		newJwt := NewJwt()
		claims, err := newJwt.ParseJwtToken(token)
		if err != nil {
			gin_util.ReturnResponse(http.StatusOK, gin_util.FailCode, err.Error(), nil, c)
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserId)
	}
}

type Jwt struct {
	SignKey []byte
}

type CustomClaims struct {
	UserId int `json:"userId"`
	*jwt.StandardClaims
}

func NewJwt() *Jwt {
	signKey := gin_config.ServerConfigSettings.Server.SignKey
	return &Jwt{
		SignKey: []byte(signKey),
	}
}

/*
	CreateJwtToken 生成token
*/
func (j *Jwt) CreateJwtToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SignKey)
}

/*
	ParseToken 解析token
*/
func (j *Jwt) ParseJwtToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

/*
	RefreshJwtToken 更新token
*/
func (j *Jwt) RefreshJwtToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(gin_const.DefaultJwtExpiresAt).Unix()
		return j.CreateJwtToken(*claims)
	}
	return "", TokenInvalid
}
