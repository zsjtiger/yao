package server

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/yaoapp/xiang/global"
)

// Guards 服务中间件
var Guards = map[string]gin.HandlerFunc{
	"bearer-jwt": bearerJWT, // JWT 权限校验
}

// JwtClaims JWT claims
type JwtClaims struct {
	ID     int
	Type   string
	Mobile string
	Name   string
	jwt.StandardClaims
}

func bearerJWT(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(403, gin.H{"code": 403, "message": "无权访问该页面"})
		c.Abort()
		return
	}

	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return global.Conf.JWT.Secret, nil
	})

	if err != nil {
		c.JSON(403, gin.H{"code": 403, "message": fmt.Sprintf("登录已过期或令牌失效(%s)", err)})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		c.Set("id", claims.Subject)
		c.Set("type", claims.Type)
		c.Set("name", claims.Name)
		c.Set("mobile", claims.Mobile)
		c.Next()
		return
	}

	// fmt.Println("bearer-JWT", token.Claims.Valid())
	c.JSON(403, gin.H{"code": 403, "message": "无权访问该页面"})
	c.Abort()
}
