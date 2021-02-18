package logs

import (
	"regexp"

	"github.com/k8smon/colors"
	"github.com/k8smon/config"
	"github.com/k8smon/models"
)

//Expressão regular utilizada para localizar logs de interesse
var regExpPickLog regexp.Regexp

func init() {
	configuracoes := config.CarregaConfiguracoes()
	patternRegExp := ""
	for i, expressao := range configuracoes["expressoes"].([]interface{}) {
		if i > 0 {
			patternRegExp += "|"
		}
		patternRegExp += expressao.(map[string]interface{})["texto"].(string)
	}
	regExpPickLog = *regexp.MustCompile(patternRegExp)
}

//PickLogs Retorna uma instância do ContainerLog apenas com os logs que contenham as expressões procuradas
func PickLogs(containerLog models.ContainerLog, canalLogs chan models.ContainerLog) {
	logs := []models.Log{}

	for _, l := range containerLog.Logs {
		match := regExpPickLog.FindAllString(l.Log, -1)
		if match != nil {
			logs = append(logs, l)
		}
	}
	containerLog.Logs = logs
	canalLogs <- containerLog
}

//PrintLog imprime o log no console/stdout
func PrintLog(containerLog models.ContainerLog) {
	if len(containerLog.Logs) == 0 {
		colors.Grayf("%s.%s.%s\n", containerLog.Namespace, containerLog.PodName, containerLog.ContainerName)
	} else {
		colors.Yellowf("%s.%s.%s [%s]\n", containerLog.Namespace, containerLog.PodName, containerLog.ContainerName, containerLog.PodNode)
	}
}
