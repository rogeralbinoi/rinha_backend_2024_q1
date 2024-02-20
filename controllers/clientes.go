package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"rinha-backend-2024-q1/models"
	"rinha-backend-2024-q1/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NewTransactionResponse struct {
	Limite int
	Saldo  int
}

func Extrato(context *gin.Context) {
	id := context.Param("id")
	cliente_id, err := strconv.Atoi(id)
	extrato, err := services.GetExtrato(cliente_id)

	if err != nil {
		log.Fatal(err)
	}

	context.JSON(http.StatusOK, extrato)
}

func NewTransaction(context *gin.Context) {
	id := context.Param("id")
	cliente_id, err := strconv.Atoi(id)

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, "Unprocessable Entity")
		return
	}

	transacao := models.Transacao{}
	err = context.ShouldBind(&transacao)

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, "Unprocessable Entity")
		return
	}

	saldo, limite, err := services.NewTransaction(cliente_id, transacao)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, "Not found")
			return
		}
		context.JSON(http.StatusUnprocessableEntity, "Unprocessable Entity")
		return
	}

	context.JSON(http.StatusOK, NewTransactionResponse{Saldo: saldo, Limite: limite})
}
