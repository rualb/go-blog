package blog

import (
	controller "go-blog/internal/controller"
	"go-blog/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ConfigDTO struct {
	Input struct{}
	Meta  struct {
		Status int
	}
	Output struct{}
}

type ConfigAPIController struct {
	appService service.AppService
	webCtxt    echo.Context // webCtxt
}

// NewAccountController is constructor.
func NewConfigAPIController(appService service.AppService, c echo.Context) *ConfigAPIController {
	return &ConfigAPIController{
		appService: appService,
		webCtxt:    c,
	}
}

func (x *ConfigAPIController) Handler() (err error) {

	dto := &ConfigDTO{}
	//
	meta := &dto.Meta
	output := &dto.Output
	c := x.webCtxt
	//

	//
	if meta.Status == 0 {
		meta.Status = http.StatusOK
	}
	return c.JSON(meta.Status, output)

}

type StatusDTO struct {
	Input struct{}
	Meta  struct {
		Status int
	}
	Output struct{}
}

type StatusAPIController struct {
	appService service.AppService
	webCtxt    echo.Context // webCtxt
}

// NewAccountController is constructor.
func NewStatusAPIController(appService service.AppService, c echo.Context) *StatusAPIController {
	return &StatusAPIController{
		appService: appService,
		webCtxt:    c,
	}
}
func (x *StatusAPIController) Handler() (err error) {

	dto := &StatusDTO{}
	//
	meta := &dto.Meta
	output := &dto.Output
	c := x.webCtxt
	//
	controller.CsrfToHeader(c)
	//
	if meta.Status == 0 {
		meta.Status = http.StatusOK
	}
	return c.JSON(meta.Status, output)

}
