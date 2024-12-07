package blog

import (
	"fmt"
	"go-blog/internal/config"
	controller "go-blog/internal/controller"
	"go-blog/internal/service"
	"time"

	"go-blog/internal/i18n"
	"go-blog/internal/mvc"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BlogIndexController struct {
	appService service.AppService
	appConfig  *config.AppConfig
	userLang   i18n.UserLang

	IsGET  bool
	IsPOST bool

	webCtxt echo.Context // webCtxt

	DTO struct {
		Input struct {
		}
		Meta struct {
			IsFragment bool `json:"-"`
		}
		Output struct {
			mvc.ModelBaseDTO
			LangCode  string
			AppConfig struct {
				AppTitle string `json:"app_title,omitempty"`
				TmTitle  string `json:"tm_title,omitempty"`
			}
			Title     string
			LangWords map[string]string
		}
	}
}

func (x *BlogIndexController) Handler() error {

	err := x.createDTO()
	if err != nil {
		return err
	}

	err = x.handleDTO()
	if err != nil {
		return err
	}

	err = x.responseDTO()
	if err != nil {
		return err
	}

	return nil
}

func NewBlogIndexController(appService service.AppService, c echo.Context) *BlogIndexController {

	return &BlogIndexController{
		appService: appService,
		appConfig:  appService.Config(),
		userLang:   controller.UserLang(c, appService),
		IsGET:      controller.IsGET(c),
		IsPOST:     controller.IsPOST(c),
		webCtxt:    c,
	}
}

func (x *BlogIndexController) validateFields() {

}

func (x *BlogIndexController) createDTO() error {

	dto := &x.DTO
	c := x.webCtxt

	if err := c.Bind(dto); err != nil {
		return err
	}

	x.validateFields()

	return nil
}

func (x *BlogIndexController) handleDTO() error {

	dto := &x.DTO
	// input := &dto.Input
	output := &dto.Output
	// meta := &dto.Meta
	// c := x.webCtxt

	userLang := x.userLang
	output.LangCode = userLang.LangCode()
	output.Title = userLang.Lang("Blog") // TODO /*Lang*/
	output.LangWords = userLang.LangWords()

	cfg := &output.AppConfig

	cfg.AppTitle = x.appConfig.Title
	cfg.TmTitle = fmt.Sprintf("© %v %s", time.Now().Year(), x.appConfig.Title)

	return nil
}

func (x *BlogIndexController) responseDTOAsMvc() (err error) {

	dto := &x.DTO
	meta := &dto.Meta
	output := &dto.Output
	appConfig := x.appConfig
	lang := x.userLang
	c := x.webCtxt

	data, err := mvc.NewModelWrap(c, output, meta.IsFragment, "Blog" /*Lang*/, appConfig, lang)
	if err != nil {
		return err
	}

	err = c.Render(http.StatusOK, "index.html", data)

	if err != nil {
		return err
	}

	return nil
}

func (x *BlogIndexController) responseDTO() (err error) {

	return x.responseDTOAsMvc()

}
