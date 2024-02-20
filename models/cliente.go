package models

type Cliente struct {
	ID     int `json:"-"`
	Limite int `json:"limite"`
	Saldo  int `json:"saldo"`
}
