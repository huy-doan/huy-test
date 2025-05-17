package router

import (
	"net/http"

	generated "github.com/huydq/test/internal/pkg/api/generated"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func SetupSwaggerUI(e *echo.Echo) {
	swagger, err := generated.GetSwagger()
	if err != nil {
		panic(err)
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	e.GET("/swagger/doc.json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, swagger)
	})
}
