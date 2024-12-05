package web

/*
validate jwt by middleware
rotate jwt by middleware
set jwt (auth jwt)
validate ExpiresAt,Issuer,Audience
*/
import (
	"fmt"
	"go-blog/internal/config/consts"
	"go-blog/internal/service"
	xtoken "go-blog/internal/token"
	"go-blog/internal/util/utilstring"
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	JwtKey = "_auth" // string value "auth"
)

func CsrfMiddleware(appService service.AppService) echo.MiddlewareFunc {

	csrfConfig := middleware.CSRFConfig{
		Skipper: assetsReqSkipper,

		TokenLookup: "header:X-CSRF-Token,form:_csrf",
		CookiePath:  "/",
		// CookieDomain:   "example.com",
		// CookieSecure:   true, // https only
		CookieHTTPOnly: true,
		CookieName:     "_csrf",
		ContextKey:     "_csrf",
		CookieSameSite: http.SameSiteDefaultMode,
	}

	return middleware.CSRFWithConfig(csrfConfig)

}

func UserLangMiddleware(appService service.AppService) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			var lang string

			// Check the _lang query parameter
			lang1 := c.QueryParam("_lang")
			if appService.HasLang(lang1) {
				lang = lang1
			} else {
				// Check the _lang cookie
				lang2, err := c.Cookie("_lang")
				if err == nil && lang2 != nil && lang2.Value != "" && appService.HasLang(lang2.Value) {
					lang = lang2.Value
				} else {
					// Fallback to the Accept-Language header
					lang3 := c.Request().Header.Get("Accept-Language")
					if len(lang3) > 2 {
						lang3 = lang3[:2]
						if appService.HasLang(lang3) {
							lang = lang3
						}
					}
				}
			}

			c.Set("lang_code", lang)

			c.Response().Header().Set("Content-Language", lang)

			return next(c)
		}
	}
}

func TokenParserMiddleware(appService service.AppService) echo.MiddlewareFunc {

	vaultKeyScopeAuth := appService.Vault().KeyScopeAuth()

	appConfig := appService.Config()
	jwtMd := echojwt.WithConfig(echojwt.Config{
		Skipper:    assetsReqSkipper,
		ContextKey: JwtKey,
		// SigningMethod:          echojwt.AlgorithmHS256, // jwt.SigningMethodHS256
		KeyFunc: func(t *jwt.Token) (any, error) {

			issuer, err := t.Claims.GetIssuer()
			if err != nil {
				return nil, err
			}

			// protect from invalid issuer
			if issuer != appConfig.Identity.AuthTokenIssuer {
				return nil, fmt.Errorf("token issuer not for auth")
			}

			return xtoken.JwtSecretSearch(t, vaultKeyScopeAuth)
		},
		SuccessHandler:         jwtParseSuccessHandler,
		ErrorHandler:           jwtParseErrorHandler,
		ContinueOnIgnoredError: true,
		TokenLookup:            fmt.Sprintf("cookie:%s,header:Authorization:Bearer ", JwtKey), // "Authorization:Bearer jwt"
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(xtoken.TokenClaimsDTO)
		},
		// Validator: // configToken.AuthTokenIssuer
	})

	return jwtMd

}

func AuthorizeMiddleware(appService service.AppService, reddirect bool) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			if IsSignedIn(c) {
				// ok
			} else {

				if reddirect {
					reqURI := c.Request().RequestURI // "/dashboard?view=weekly"
					redirectURL := utilstring.AppendURL(consts.PathAuthSignin, "return_url", reqURI)
					return c.Redirect(http.StatusFound /*302*/, redirectURL)
				} else {
					return c.NoContent(http.StatusUnauthorized) // 401
				}

			}

			return next(c)
		}
	}
}

func assetsReqSkipper(c echo.Context) bool {
	path := c.Request().URL.Path
	prefixes := []string{consts.PathBlogAssets}
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			// Skip the middleware
			return true
		}
	}
	return false
}

func jwtParseSuccessHandler(c echo.Context) {

	// user, _ := c.Get(JwtKey).(*jwt.Token)
	// claims, _ := user.Claims.(*xtoken.TokenClaimsDTO)

	//

}

func jwtParseErrorHandler(c echo.Context, err error) error {
	return nil
}

func IsSignedIn(c echo.Context) bool {
	claims := AuthTokenClaims(c)
	return claims != nil && claims.IsSignedIn()
}

func GetAccount(c echo.Context, srv service.AppService) (*service.UserAccount, error) {

	if srv == nil {
		return nil, fmt.Errorf("arg is nil: service")
	}

	acc, _ := c.Get("user_account").(*service.UserAccount)

	if acc != nil {
		return acc, nil
	}

	userID := UserID(c)

	acc, err := srv.Account().FindByID(userID)
	if err != nil {
		return nil, err
	}

	c.Set("user_account", acc)

	return acc, nil
}

func AuthTokenClaims(c echo.Context) *xtoken.TokenClaimsDTO {

	jwtToken, ok := c.Get(JwtKey).(*jwt.Token)
	if ok && jwtToken != nil && jwtToken.Valid {

		claims, _ := jwtToken.Claims.(*xtoken.TokenClaimsDTO)
		if claims != nil && claims.IsValid() {
			// if claims.HasScope(ScopeAuth) { // check token has scope auth
			return claims
			//}
		}

	}

	return nil
	//
}

func UserID(c echo.Context) string {
	claims := AuthTokenClaims(c)
	if claims != nil /*&& claims.HasScope(ScopeAuth)*/ {
		return claims.UserID
	}
	return ""
}
