package auth

import (
	"time"
	"todo-api/internal/model"
	"todo-api/internal/service"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var identityKey = "id"

// JwtMiddleware инициализирует middleware
func JwtMiddleware(userSvc service.UserService) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "notion zone",
		Key:         []byte("your-very-secret-key"), // вынесите в .env!
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			u := data.(*model.User)
			return jwt.MapClaims{
				identityKey: u.ID,
				"name":      u.Name,
			}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.User{
				ID:   uint(claims[identityKey].(float64)),
				Name: claims["email"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals model.UserInput
			if err := c.ShouldBindJSON(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			user, err := userSvc.CreateUser(&loginVals)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// здесь можно проверять роль или ID
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},
		TokenLookup:   "header: Authorization, query: token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
