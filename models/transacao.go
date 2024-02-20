package models

import "time"

type TipoTransacao string

const (
	Debito  TipoTransacao = "d"
	Credito TipoTransacao = "c"
)

type Transacao struct {
	ID          int           `json:"id"`
	ClienteID   int           `json:"cliente_id"`
	Tipo        TipoTransacao `json:"tipo"`
	Valor       int           `json:"valor"`
	Descricao   string        `json:"descricao"`
	RealizadaEm time.Time     `json:"realizada_em"`
}

type TransacaoExtrato struct {
	Tipo        TipoTransacao `json:"tipo"`
	Valor       int           `json:"valor"`
	Descricao   string        `json:"descricao"`
	RealizadaEm time.Time     `json:"realizada_em"`
}

type Extrato struct {
	Total             int                `json:"total"`
	DataExtrato       time.Time          `json:"data_extrato"`
	Limite            int                `json:"limite"`
	UltimasTransacoes []TransacaoExtrato `json:"ultimas_transacoes"`
}
