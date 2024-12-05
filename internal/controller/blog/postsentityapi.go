package blog

import (
	"go-blog/internal/config"
	controller "go-blog/internal/controller"
	"go-blog/internal/mvc"

	"go-blog/internal/i18n"
	"go-blog/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BlogPostDTO struct {
	service.BlogPost
}
type PostsEntityDTO struct {
	Input struct {
		Code string `param:"code"`
	}
	Meta struct {
		Status int
	}
	Output struct {
		mvc.ModelBaseDTO
		Data BlogPostDTO `json:"data,omitempty"`
	}
}
type PostsEntityAPIController struct {
	appService service.AppService
	appConfig  *config.AppConfig
	userLang   i18n.UserLang

	IsGET bool

	webCtxt echo.Context // webCtxt

	DTO PostsEntityDTO
}

func (x *PostsEntityAPIController) Handler() error {
	// TODO sign out force

	err := x.validateDTO()
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

// NewAccountController is constructor.
func NewPostsEntityAPIController(appService service.AppService, c echo.Context) *PostsEntityAPIController {

	appConfig := appService.Config()

	return &PostsEntityAPIController{
		appService: appService,
		appConfig:  appConfig,
		userLang:   controller.UserLang(c, appService),
		IsGET:      controller.IsGET(c),

		webCtxt: c,
	}
}

func (x *PostsEntityAPIController) validateDTOFields() (err error) {

	return nil

}

func (x *PostsEntityAPIController) validateDTO() error {

	dto := &x.DTO
	input := &dto.Input

	c := x.webCtxt

	if err := c.Bind(input); err != nil {
		return err
	}

	return x.validateDTOFields()

}
func (x *PostsEntityAPIController) handleGET() (err error) {
	dto := &x.DTO
	input := &dto.Input
	meta := &dto.Meta
	output := &dto.Output
	srv := x.appService.Blog()

	var res *service.BlogPost

	if input.Code != "" { // /code/:code
		res, err = srv.Posts().FindByCode(input.Code)
	}

	if err != nil {
		return err
	}

	if res == nil {
		meta.Status = http.StatusNotFound
	} else {
		output.Data.BlogPost = *res // copy
		// output.Data.ContentMarkdown = ""
	}

	return nil
}

func (x *PostsEntityAPIController) handleDTO() error {

	dto := &x.DTO
	meta := &dto.Meta
	output := &dto.Output

	if meta.Status > 0 {
		return nil // stop processing
	}

	switch {

	case x.IsGET:
		return x.handleGET()
	default:
		{
			meta.Status = http.StatusMethodNotAllowed
			output.Message = "Method action undef"
		}
	}

	return nil
}
func (x *PostsEntityAPIController) responseDTOAsAPI() (err error) {

	dto := &x.DTO
	meta := &dto.Meta
	output := &dto.Output
	c := x.webCtxt
	controller.CsrfToHeader(c)

	if meta.Status == 0 {
		meta.Status = http.StatusOK
	}

	return c.JSON(meta.Status, output)

}

func (x *PostsEntityAPIController) responseDTO() (err error) {

	return x.responseDTOAsAPI()

}
