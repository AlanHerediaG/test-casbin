package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AlanHerediaG/test-casbin/auth"
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/persist"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}

func Authorize(obj string, act string, adapter persist.Adapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "user hasn't logged in yet")
			c.Abort()
			return
		}

		metadata, err := auth.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}

		ok, err := enforce(metadata.UserName, obj, act, adapter)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, "error ocurred when authorizing user")
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, "forbidden")
			return
		}

		c.Next()
	}
}

func enforce(sub string, obj string, act string, adapter persist.Adapter) (bool, error) {
	enforcer := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	err := enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}

	ok := enforcer.Enforce(sub, obj, act)
	return ok, nil
}
