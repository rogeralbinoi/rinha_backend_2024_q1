package services

import (
	"errors"
	"rinha-backend-2024-q1/lib/db"
	"rinha-backend-2024-q1/models"
	"time"
)

func GetExtrato(id int) (extrato models.Extrato, err error) {
	conn, err := db.OpenConn()
	if err != nil {
		return
	}
	defer conn.Close()

	clientesResult := make(chan models.Cliente)
	transacoesResult := make(chan []models.TransacaoExtrato)

	go func() {
		clientesRows, err := conn.Query(`SELECT saldo, limite  FROM clientes where id = $1`, id)
		if err != nil {
			clientesResult <- models.Cliente{}
			return
		}
		defer clientesRows.Close()

		var cliente models.Cliente
		for clientesRows.Next() {
			err = clientesRows.Scan(&cliente.Saldo, &cliente.Limite)
			if err != nil {
				continue
			}
		}
		clientesResult <- cliente
	}()

	go func() {
		transacoesRows, err := conn.Query(`SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE cliente_id = $1 LIMIT 10`, id)
		if err != nil {
			transacoesResult <- []models.TransacaoExtrato{}
			return
		}
		defer transacoesRows.Close()

		var transacoes []models.TransacaoExtrato
		for transacoesRows.Next() {
			var transacao models.TransacaoExtrato
			err = transacoesRows.Scan(&transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadaEm)
			if err != nil {
				continue
			}
			transacoes = append(transacoes, transacao)
		}
		transacoesResult <- transacoes
	}()

	cliente := <-clientesResult
	extrato.UltimasTransacoes = <-transacoesResult

	extrato.DataExtrato = time.Now()
	extrato.Total = cliente.Saldo
	extrato.Limite = cliente.Limite

	return
}

func NewTransaction(id int, transacao models.Transacao) (saldo int, limite int, err error) {
	conn, err := db.OpenConn()
	if err != nil {
		return
	}
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		return
	}

	var saldoAtual int
	var limiteAtual int
	err = tx.QueryRow("SELECT saldo, limite FROM clientes WHERE id = $1", id).Scan(&saldoAtual, &limiteAtual)
	if err != nil {
		tx.Rollback()
		return
	}

	if !(saldoAtual-transacao.Valor >= limiteAtual*-1) {
		tx.Rollback()
		err = errors.New("verification failed")
		return
	}

	_, err = tx.Exec("UPDATE clientes SET saldo = saldo - $1 WHERE id = $2", transacao.Valor, id)
	if err != nil {
		tx.Rollback()
		return
	}

	_, err = tx.Exec(`INSERT INTO public.transacoes (cliente_id, tipo, valor, descricao, realizada_em)
		VALUES ($1, $2, $3, $4, $5);`, id, transacao.Tipo, transacao.Valor, transacao.Descricao, time.Now())
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.QueryRow("SELECT saldo, limite FROM clientes WHERE id = $1", id).Scan(&saldo, &limite)
	if err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		return 0, 0, errors.New("Fatal")
	}
	return
}
