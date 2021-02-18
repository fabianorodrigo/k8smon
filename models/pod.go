package models

type Pod struct {
	Nome         string
	QtContainers int
}

//PodIsIn retorna TRUE se algum elemento de {arr} tem a
//nome igual a {nomePod}
func PodIsIn(arr []Pod, nomePod string) bool {
	for _, a := range arr {
		if a.Nome == nomePod {
			return true
		}
	}
	return false
}

//getPodByName retorna nó do array que tenha o nome
//igual a {nomeNo}, ou nil se não houver
func GetPodByName(arr []Pod, nomePod string) *Pod {
	for _, a := range arr {
		if a.Nome == nomePod {
			return &a
		}
	}
	return nil
}
