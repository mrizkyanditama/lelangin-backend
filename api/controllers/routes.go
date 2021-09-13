package controllers

import (
	"github.com/mrizkyanditama/lelangin/api/middlewares"
)

func (s *Server) initializeRoutes() {

	v1 := s.Router.Group("/api/v1")
	{
		// Login Route
		v1.POST("/login", s.Login)

		// Reset password:
		v1.POST("/password/forgot", s.ForgotPassword)
		v1.POST("/password/reset", s.ResetPassword)

		//Users routes
		v1.POST("/users", s.CreateUser)
		// The user of the app have no business getting all the users.
		// v1.GET("/users", s.GetUsers)
		// v1.GET("/users/:id", s.GetUser)
		v1.PUT("/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateUser)
		v1.PUT("/avatar/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateAvatar)
		v1.DELETE("/users/:id", middlewares.TokenAuthMiddleware(), s.DeleteUser)

		//Product routes
		v1.GET("/products", s.GetProducts)
		v1.GET("/products/:id", s.GetProduct)
		v1.POST("/products", middlewares.TokenAuthMiddleware(), s.CreateAuction)

		v1.GET("/auctions", s.GetAuctions)

		v1.GET("/bid/:id", s.HandleBids)

		//Category routes
		v1.GET("/categories", s.GetCategories)

		//Tag routes
		v1.GET("/tags", s.GetTags)
	}
}
