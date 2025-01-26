package controller

import (
	"go-blog/internal/i18n"
	"go-blog/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UserLang(c echo.Context, appLang i18n.AppLang) i18n.UserLang {

	lang, _ := c.Get("lang_code").(string)
	return appLang.UserLang(lang)
}

func IsGET(c echo.Context) bool {
	return c.Request().Method == http.MethodGet
}

func IsPOST(c echo.Context) bool {
	return c.Request().Method == http.MethodPost
}

func IsPUT(c echo.Context) bool {
	return c.Request().Method == http.MethodPut
}

func IsDELETE(c echo.Context) bool {
	return c.Request().Method == http.MethodDelete
}

// func CsrfToHeader(c echo.Context) {
// 	csrf, _ := c.Get("_csrf").(string)
// 	c.Response().Header().Set("X-CSRF-Token", csrf)
// }

func GetAccount(c echo.Context) *service.UserAccount {
	acc, _ := c.Get("user_account").(*service.UserAccount) // cached by middleware
	return acc
}
