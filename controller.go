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

// To Do: Mengembalikan rincian transaksi sesuai dengan ID pembayaran pada endpoint ini akan mengembalikan
//        laporan dengan format berbeda untuk user tertentu karena menerapkan feature flagging pada bagian
//        business (core) dari endpoint ini.
func (c *Controller) Get(ctx echo.Context) error {
	userIdentity := ctx.Request().Header.Get("user_identity")
	if strings.TrimSpace(userIdentity) == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "header not valid"})
	}

	paymentID := ctx.Param("payment_id")
	if strings.TrimSpace(paymentID) == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "payment id cannot be empty"})
	}

	result, err := c.service.Get(userIdentity, paymentID)
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
