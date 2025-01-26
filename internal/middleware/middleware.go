package middleware

import (
	"go-blog/internal/config/consts"
	"go-blog/internal/service"
	xweb "go-blog/internal/web"
	"io/fs"

	xlog "go-blog/internal/util/utillog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo-contrib/echoprometheus"
)

func InitBasic(e *echo.Echo, appService service.AppService) {

}
func Init(e *echo.Echo, appService service.AppService) {

	appConfig := appService.Config()

	e.HTTPErrorHandler = newHTTPErrorHandler(appService)

	e.Use(middleware.Recover()) //!!!

	if appConfig.HTTPServer.AccessLog {
		e.Use(middleware.Logger())
	}

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level:     5,
		MinLength: 500,
	}))
	//
	e.Use(xweb.UserLangMiddleware(appService))
	e.Use(xweb.TokenParserMiddleware(appService))
	//

	//
	// e.Use(xweb.CsrfMiddleware(appService))

	initSys(e, appService)
}

func initSys(e *echo.Echo, appService service.AppService) {

	appConfig := appService.Config()

	// name := "" // appConfig.Name // name as var

	if appConfig.HTTPServer.SysMetrics {

		// collect metrics (singleton)
		e.Use(echoprometheus.NewMiddlewareWithConfig(

			echoprometheus.MiddlewareConfig{
				// each 404 has own metric (not good)
				DoNotUseRequestPathFor404: true,
			},
		))
	}
}

func newHTTPErrorHandler(_ service.AppService) echo.HTTPErrorHandler {

	return func(err error, c echo.Context) {

		c.Echo().DefaultHTTPErrorHandler(err, c)

	}

}

func AssetsContentsMiddleware(e *echo.Echo, appService service.AppService,
	assetsBlogFiles fs.FS,

) {

	xlog.Info("start serving embedded static content")

	{
		grp := e.Group(consts.PathBlogAssets, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				// c.Response().Before()
				c.Response().Header().Add(`Cache-Control`, "public,max-age=31536000,immutable")
				return next(c)
			}
		})

		grp.StaticFS("/", assetsBlogFiles)
	}

	// admin

}
