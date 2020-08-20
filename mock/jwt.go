package mock

import (
	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"

	jwtapple2 "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	mocket "github.com/selvatico/go-mocket"
)

// ConfigClaims configure claims with user
func ConfigClaims(c *gin.Context, u model.User) {
	config.TestInit()

	claims := jwtapple2.MapClaims{config.IdentityKey: float64(u.ID)}

	c.Set("JWT_PAYLOAD", claims)

	dbResponse := GetMapArrayByUser(u)
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"`).WithReply(dbResponse)
}
