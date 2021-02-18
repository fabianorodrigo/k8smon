package models

type Node struct {
	Nome         string
	QtPods       int
	QtContainers int
}

//NodeIsIn retorna TRUE se algum elemento de {arr} tem a propriedade
//nome igual a {nomeNo}
func NodeIsIn(arr []Node, nomeNo string) bool {
	for _, a := range arr {
		if a.Nome == nomeNo {
			return true
		}
	}
	return false
}

//getNodeByName retorna nó do array que tenha o nome
//igual a {nomeNo}, ou nil se não houver
func GetNodeByName(arr []Node, nomeNo string) *Node {
	for _, a := range arr {
		if a.Nome == nomeNo {
			return &a
		}
	}
	return nil
}
