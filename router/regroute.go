package router

import "github.com/gin-gonic/gin"

func RegrouteInit(app *gin.Engine)  {
	rootGroup := app.Group("/")
	{
		rootGroup.Any("/status", func(context *gin.Context) {
			
		})
		
	}
}
