package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/vnlab/makeshop-payment/docs/api/generated"
)

// SetupSwaggerUI thiết lập Swagger UI cho API documentation
func SetupSwaggerUI(e *echo.Echo) {
	swagger, err := generated.GetSwagger()
	if err != nil {
		panic(err)
	}
	swagger.Servers = nil
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler())
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	e.GET("/openapi.json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, swagger)
	})
}
