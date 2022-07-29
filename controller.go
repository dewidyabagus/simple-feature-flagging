package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	role := ctx.Request().Header.Get("role")
	if strings.TrimSpace(role) == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "header not valid"})
	}

	paymentID, err := strconv.Atoi(ctx.Param("payment_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "payment id not valid"})
	}

	result, err := c.service.Get(role, paymentID)
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

// Mocking API Notifier yang digunakan untuk menerima pemberitahuan ketika file config berubah
type FFController struct {
}

func NewFFController() *FFController {
	return &FFController{}
}

func (f *FFController) NotifierHook(ctx echo.Context) error {
	// Model auth yang digunakan menggunakan auth key yang value nya berasal dari
	// body yang di encrypt dengan kunci yang sudah di deskripsikan. Auth yang digunakan HMAC(SHA256)
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

	// Decode bytes payload to map[string]interface
	payload := map[string]interface{}{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}

	fmt.Println(string(body))

	return ctx.JSON(http.StatusOK, echo.Map{"message": "successfully send notifier"})
}
