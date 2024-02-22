package main

import "github.com/gin-gonic/gin"

func Route() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("./ui/templates/*")
	r.Static("/static", "./ui/static")

	r.GET("/", requireAuth, getMain)
	r.GET("/product/:id", requireAuth, getProduct)
	r.GET("/product/:id/review", requireAuth, getReview)
	r.POST("/product/:id/review", requireAuth, postReview)
	r.GET("/logout", logout)

	r.GET("/login", getLogin)
	r.GET("/register", getReg)
	r.GET("/search", requireAuth, search)
	r.POST("/register", postReg)
	r.POST("/login", postLogin)

	r.Use(recoverPanic())
	return r
}
