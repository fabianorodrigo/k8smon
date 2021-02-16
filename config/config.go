package config

import (
	"encoding/json"
	"io/ioutil"
)

//Dictionary é um map[string]interface{}
type Dictionary map[string]interface{}

//CarregaConfiguracoes carrega as configurações do arquivo expresssoesBuscadas.json
func CarregaConfiguracoes() Dictionary {
	// read file
	data, err := ioutil.ReadFile("./config/logConfig.json")
	if err != nil {
		panic(err.Error())
	}
	var config Dictionary
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err.Error())
	}
	return config
}
