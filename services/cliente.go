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

    clienteQuery := `
        SELECT saldo, limite
        FROM clientes
        WHERE id = $1
    `

    clienteRow := conn.QueryRow(clienteQuery, id)
    var cliente models.Cliente

    if err := clienteRow.Scan(&cliente.Saldo, &cliente.Limite); err != nil {
        return extrato, err
    }

    transacoesQuery := `
        SELECT valor, tipo, descricao, realizada_em
        FROM transacoes
        WHERE cliente_id = $1
        ORDER BY realizada_em DESC
        LIMIT 10
    `

    rows, err := conn.Query(transacoesQuery, id)
    if err != nil {
        return
    }
    defer rows.Close()

    var transacoes []models.TransacaoExtrato

    for rows.Next() {
        var valor *int
        var tipo models.TipoTransacao
        var descricao string
        var realizadaEm time.Time

        if err := rows.Scan(&valor, &tipo, &descricao, &realizadaEm); err != nil {
            return extrato, err
        }

        var valorInt int
        if valor != nil {
            valorInt = *valor
        }

        transacao := models.TransacaoExtrato{
            Tipo:        tipo,
            Valor:       valorInt,
            Descricao:   descricao,
            RealizadaEm: realizadaEm,
        }

        transacoes = append(transacoes, transacao)
    }

    extrato.Saldo.DataExtrato = time.Now()
    extrato.Saldo.Total = cliente.Saldo
    extrato.Saldo.Limite = cliente.Limite
	if transacoes != nil {
		extrato.UltimasTransacoes = transacoes
	} else {
		extrato.UltimasTransacoes = []models.TransacaoExtrato{}
	}


    return
}




func NewTransaction(id int, transacao models.Transacao) (saldo int, limite int, err error) {
	conn, err := db.OpenConn()
	if err != nil {
		return
	}

	tx, err := conn.Begin()

	defer func() {
		conn.Close()
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	if err != nil {
		return
	}

	var saldoAtual int
	var limiteAtual int
	err = tx.QueryRow("SELECT saldo, limite FROM clientes WHERE id = $1", id).Scan(&saldoAtual, &limiteAtual)
	if err != nil {
		return
	}

	if !((saldoAtual-transacao.Valor) >= limiteAtual * (-1)) {
		err = errors.New("verification failed")
		return
	}

	_, err = tx.Exec("UPDATE clientes SET saldo = saldo - $1 WHERE id = $2", transacao.Valor, id)
	if err != nil {
		return
	}

	_, err = tx.Exec(`INSERT INTO public.transacoes (cliente_id, tipo, valor, descricao, realizada_em)
		VALUES ($1, $2, $3, $4, $5);`, id, transacao.Tipo, transacao.Valor, transacao.Descricao, time.Now())
	if err != nil {
		return
	}

	err = tx.QueryRow("SELECT saldo, limite FROM clientes WHERE id = $1", id).Scan(&saldo, &limite)

	return
}
