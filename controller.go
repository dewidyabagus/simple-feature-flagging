package main

import (
	"net/http"
	"strings"

	echo "github.com/labstack/echo/v4"
)

type Controller struct {
	service Servicer
}

func NewController(service Servicer) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Pay(ctx echo.Context) error {
	return ctx.JSON(http.StatusCreated, echo.Map{"message": "successful payment"})
}

func (c *Controller) Get(ctx echo.Context) error {
	exampleUserID := ctx.Request().Header.Get("example-user-id")
	if strings.TrimSpace(exampleUserID) == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "header not valid"})
	}

	orderID := ctx.Param("order_id")
	if strings.TrimSpace(orderID) == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "order id cannot be empty"})
	}

	result, err := c.service.Get(exampleUserID, orderID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success", "data": result})
}

func (c *Controller) Generate(ctx echo.Context) error {
	results, err := c.service.Generate()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success", "data": results})
}
