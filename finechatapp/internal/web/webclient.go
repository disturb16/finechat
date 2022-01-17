package web

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed dist
var assetsSource embed.FS

// SetAssetsHandlers sets the http handlers for the static assets.
func SetAssetsHandlers(e *echo.Echo) {
	fsys, err := fs.Sub(assetsSource, "dist")
	if err != nil {
		panic(err)
	}
	assetsHandler := http.FileServer(http.FS(fsys))

	// serves the index.html
	e.GET("/", echo.WrapHandler(assetsHandler))

	// serves other static files
	e.GET("/css/*", echo.WrapHandler(assetsHandler))
	e.GET("/js/*", echo.WrapHandler(assetsHandler))
	e.GET("/img/*", echo.WrapHandler(assetsHandler))
}
