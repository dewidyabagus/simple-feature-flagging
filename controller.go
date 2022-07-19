package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type FFController struct {
}

func NewFFController() *FFController {
	return &FFController{}
}

func (f *FFController) NotifierHook(ctx echo.Context) error {
	payload := map[string]interface{}{}

	signature := strings.ReplaceAll(ctx.Request().Header.Get("X-Hub-Signature-256"), "sha256=", "")

	buf := new(bytes.Buffer)
	buf.ReadFrom(ctx.Request().Body)
	defer ctx.Request().Body.Close()

	body := buf.Bytes()
	if valid, err := CompareSignatures(signature, GetSignature([]byte("secret"), body)); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})

	} else if !valid {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"message": "access unauthorized"})

	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}

	fmt.Println(string(body))

	return ctx.JSON(http.StatusOK, echo.Map{"message": "successfully send notifier"})
}
