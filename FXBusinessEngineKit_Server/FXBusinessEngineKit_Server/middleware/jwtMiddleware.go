package middleware

import (
	"FXBusinessEngineKit_Server/configuration"
	"FXBusinessEngineKit_Server/log"
	"FXBusinessEngineKit_Server/tool"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

func InitToken(userId string) (tokenString string, err error) {
	claim := configuration.MyClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "FXBusinessEngine_project",
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(configuration.JWTExp)),
			ID:        userId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err = token.SignedString([]byte(configuration.JWTSecret))
	if err != nil {
		fmt.Println("用户创建token失败", err)
		log.RecordLog(log.Err, fmt.Sprintf("[JWT]token创建失败 %s", err))
	}
	return tokenString, err
}

func VerifyToken(tokenString string) (*configuration.MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &configuration.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configuration.JWTSecret), nil
	})
	if err != nil {
		fmt.Printf("token验证失败基础层 %s", err)
		return nil, err
	}

	claim, ok := token.Claims.(*configuration.MyClaims)
	if ok && token.Valid {
		fmt.Printf("Token验证成功 %s", claim)
		return claim, nil
	}
	return nil, errors.New("token验证失败")
}

func VerifyProductId(productId int64) (tool.OrganizationModel, error) {
	o := tool.OrganizationModel{}
	result := tool.FXGorm.Where("productId = ?", productId).First(&o)
	if err := result.Error; err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[鉴权] productId 验证失败 %s", err))
	}
	log.RecordLog(log.Msg, fmt.Sprintf("[鉴权] productId %s 验证成功", o.ProductId))
	return o, nil
}

// JWTMiddleware token验证中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 从请求头中获取
		tokenStr := context.Request.Header.Get("token")
		uuidStr := context.Request.Header.Get("uuid")
		productId := tool.ProductIdFormHeader(context)

		// 没有token
		if tokenStr == "" || strings.HasPrefix(tokenStr, "Bearer") == false {
			context.JSON(http.StatusOK, gin.H{"code": 401, "msg": "鉴权失败"})
			log.RecordLog(log.Err, fmt.Sprintf("[JWT]鉴权失败没有 token uuid %s", uuidStr))
			// 阻止
			context.Abort()
			return
		}

		// 提取token的有效部分
		tokenStr = tokenStr[6:]
		// 验证
		_, err := VerifyToken(tokenStr)
		// token验证失败
		if err != nil {
			context.JSON(http.StatusOK, gin.H{"code": 0, "msg": "token验证失败"})
			log.RecordLog(log.Err, fmt.Sprintf("[JWT]鉴权验证失败 uuid %s", uuidStr))
			context.Abort()
			return
		}
		// 验证 productId
		_, err = VerifyProductId(productId)
		if err != nil {
			tool.ResponseFail(context, 401, "productId验证失败", "")
			context.Abort()
			return
		}

		log.RecordLog(log.Msg, fmt.Sprintf("[JWT]鉴权验证通过 uuid %s", uuidStr))
		// 将用户信息放入上下文中
		context.Set("uuid", uuidStr)
		context.Next()
	}
}

// JWTUserLogin token获取中间件
func JWTUserLogin() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 从请求头中获取
		tokenStr := tool.TokenFromHeader(context)
		uuid := tool.UUIDFromHeader(context)
		productId := tool.ProductIdFormHeader(context)
		platform := tool.PlatformFormHeader(context)
		if uuid == "" {
			tool.ResponseFail(context, http.StatusUnauthorized, "uuid为空", "")
			context.Abort()
			return
		}
		// 没有productId
		if productId == 000000000 {
			tool.ResponseFail(context, http.StatusUnauthorized, "productId 为空", "")
			context.Abort()
			return
		}
		// 验证组织
		_, err := VerifyProductId(productId)
		if err != nil {
			tool.ResponseFail(context, http.StatusUnauthorized, "productId 验证失败", "")
			context.Abort()
			return
		}

		if platform == "" {
			tool.ResponseFail(context, http.StatusUnauthorized, "platform为空", "")
			context.Abort()
			return
		}

		// 没有token
		if tokenStr == "" {
			token, err := InitToken(uuid)
			if err != nil {
				context.JSON(http.StatusUnauthorized, gin.H{"code": 200, "msg": "新用户注册获取token失败", "error": err})
				log.RecordLog(log.Err, fmt.Sprintf("鉴权失败 uuid %s err %s", uuid, err))
				context.Abort()
				return
			}

			// 返回token
			context.Set("uuid", uuid)
			context.Set("token", token)
			context.Set("productId", productId)
			context.Set("platform", platform)

			// 放行
			context.Next()
			log.RecordLog(log.Msg, fmt.Sprintf("[JWT]用户获取token成功, uuid%s", uuid))
			return
		}

		context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "鉴权失败"})
		log.RecordLog(log.Err, fmt.Sprintf("[JWT]鉴权失败 uuid %s", uuid))
		context.Abort()
	}
}
