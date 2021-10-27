package core

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/stnc/pongo2gin"
	"os"
	"path/filepath"
)

type server interface {
}

var addr = flag.String("addr", ":8080", "http address to bind to.")

func RunServer() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(gin.Recovery())

	app.HTMLRender = pongo2gin.TemplatePath("templates")

	appDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	assetsDir := appDir + "/assets"
	app.Static("/static", assetsDir+"/common/static")
	router.R

}
