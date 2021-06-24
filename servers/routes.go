package servers

import (
	"github.com/AlanHerediaG/test-casbin/controllers"
	"github.com/AlanHerediaG/test-casbin/middleware"

	"github.com/gin-gonic/gin"
)

func (s *Server) InitializeRoutes() {
	s.Router.POST("/login", controllers.Login)

	authorized := s.Router.Group("/").
		Use(gin.Logger()).
		Use(gin.Recovery()).
		Use(middleware.TokenAuthMiddleware())
	{
		authorized.POST("/api/todo", middleware.Authorize("resource", "write", s.FileAdapter), controllers.CreateTodo)
		authorized.GET("/api/todo", middleware.Authorize("resource", "read", s.FileAdapter), controllers.GetTodo)
		authorized.POST("/logout", controllers.Logout)
		authorized.POST("/refresh", controllers.Refresh)
	}
}
