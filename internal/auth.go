package auth

import (
	"os"
	"time"
	"todo-api/internal/model"
	"todo-api/internal/service"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var identityKey = "id"

func JwtMiddleware(userSvc service.UserService) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "notion zone",
		Key:         []byte(os.Getenv("KEY")),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if u, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: u.ID,
					"name":      u.Name,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			idAny, ok := claims[identityKey]
			var id uint = 0
			if ok {
				switch v := idAny.(type) {
				case float64:
					id = uint(v)
				case int:
					id = uint(v)
				}
			}
			name, _ := claims["name"].(string)
			return &model.User{
				ID:   id,
				Name: name,
			}
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals model.LoginInput
			if err := c.ShouldBindJSON(&loginVals); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			user, err := userSvc.AuthenticateUser(loginVals.Name, loginVals.Password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return user, nil
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			// Здесь можно реализовать проверки ролей и прав
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
