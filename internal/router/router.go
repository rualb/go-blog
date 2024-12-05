package router

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"

	"go-blog/internal/config/consts"
	"go-blog/internal/controller/blog"

	"go-blog/internal/service"
	xlog "go-blog/internal/util/utillog"
	webfs "go-blog/web"

	"github.com/labstack/echo/v4/middleware"
)

func Init(e *echo.Echo, appService service.AppService) {

	e.Renderer = mustNewRenderer()

	initBlogController(e, appService)

	initDebugController(e, appService)

	initSys(e, appService)
}

func initSys(e *echo.Echo, appService service.AppService) {

	// !!! DANGER for private(non-public) services only
	// or use non-public port via echo.New()

	appConfig := appService.Config()

	listen := appConfig.HTTPServer.Listen
	listenSys := appConfig.HTTPServer.ListenSys
	sysMetrics := appConfig.HTTPServer.SysMetrics
	hasAnyService := sysMetrics
	sysAPIKey := appConfig.HTTPServer.SysAPIKey
	hasAPIKey := sysAPIKey != ""
	hasListenSys := listenSys != ""
	startNewListener := listenSys != listen

	if !hasListenSys {
		return
	}

	if !hasAnyService {
		return
	}

	if !hasAPIKey {
		xlog.Panic("Sys api key is empty")
		return
	}

	if startNewListener {

		e = echo.New() // overwrite override

		e.Use(middleware.Recover())
		// e.Use(middleware.Logger())
	} else {
		xlog.Warn("Sys api serve in main listener: %v", listen)
	}

	sysAPIAccessAuthMW := middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:api-key,header:Authorization",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == sysAPIKey, nil
		},
	})

	if sysMetrics {
		// may be eSys := echo.New() // this Echo will run on separate port
		e.GET(
			consts.PathSysMetricsAPI,
			echoprometheus.NewHandler(),
			sysAPIAccessAuthMW,
		) // adds route to serve gathered metrics

	}

	if startNewListener {

		// start as async task
		go func() {
			xlog.Info("Sys api serve on: %v main: %v", listenSys, listen)

			if err := e.Start(listenSys); err != nil {
				if err != http.ErrServerClosed {
					xlog.Error("%v", err)
				} else {
					xlog.Info("shutting down the server")
				}
			}
		}()

	} else {
		xlog.Info("Sys api server serve on main listener: %v", listen)
	}

}

type tmplRenderer struct {
	// viewsMvc       echo.Renderer
	blogIndex *template.Template
}

func (x *tmplRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {

	if name == "blog.html" {

		return x.blogIndex.ExecuteTemplate(w, name, data)
	}

	return fmt.Errorf("undefined tmpl")

}

func mustNewRenderer() echo.Renderer {

	blogIndex, err := template.New("blog.html").Parse(webfs.MustBlogIndexHTML())

	if err != nil {
		panic(err)
	}

	//	err := t.templates.ExecuteTemplate(w, "layout_header", data)

	handler := &tmplRenderer{
		// viewsMvc:       mvc.NewTemplateRenderer(controller.ViewsFs(), "views/blog/*.html"),
		blogIndex: blogIndex,
	}

	return handler

}

func initDebugController(e *echo.Echo, _ service.AppService) {
	e.GET(consts.PathBlogPingDebugAPI, func(c echo.Context) error { return c.String(http.StatusOK, "pong") })

}

func initBlogController(e *echo.Echo, appService service.AppService) {

	xlog.Info("Adding blog controllers")

	prefix := consts.PathBlog
	group := e.Group(prefix)

	path := func(s string) string {
		xlog.Info("Route: %s", s)
		return strings.TrimPrefix(s, prefix)
	}

	{

		{
			group.GET(path(consts.PathBlogStatusAPI), func(c echo.Context) error {
				ctrl := blog.NewStatusAPIController(appService, c)
				return ctrl.Handler()
			})
			group.GET(path(consts.PathBlogConfigAPI), func(c echo.Context) error {
				ctrl := blog.NewConfigAPIController(appService, c)
				return ctrl.Handler()
			})
		}

		{
			// return UI
			handler := func(c echo.Context) error {
				ctrl := blog.NewBlogIndexController(appService, c)
				return ctrl.Handler()
			}
			group.GET(path(consts.PathBlogPostsEntity), handler)
		}

		{

			{
				handler := func(c echo.Context) error {
					ctrl := blog.NewPostsEntityAPIController(appService, c)
					return ctrl.Handler()
				}

				group.GET(path(consts.PathBlogPostsEntityByCodeAPI), handler)
			}

		}
	}

}

/////////////////////////////////////////////////////
