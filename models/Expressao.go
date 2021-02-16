package models

//Expressao buscada nos logs de erro
type Expressao struct {
	Texto    string   `json:"texto"`
	Excecoes []string `json:"excecoes"`
}
