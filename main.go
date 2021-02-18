package main

import (
	"fmt"

	"github.com/k8smon/config"
	"github.com/k8smon/kubernetes"
	"github.com/k8smon/logs"
	"github.com/k8smon/models"
	"github.com/k8smon/ui"
)

func main() {
	//jobName := flag.String("jobname", "test-job", "The name of the job")
	//containerImage := flag.String("image", "ubuntu:latest", "Name of the container image")
	//entryCommand := flag.String("command", "ls", "The command to run inside the container")
	//namespace := flag.String("namespace", "", "Namespace de pesquisa")

	//flag.Parse()

	//fmt.Printf("Args : %s %s %s\n", *jobName, *containerImage, *entryCommand)

	configuracoes := config.CarregaConfiguracoes()

	clientSet := kubernetes.ConectarK8s()
	//criando um channel para receber os container logs
	canalLogs := make(chan models.ContainerLog)
	go kubernetes.Logs(clientSet, configuracoes, logs.PickLogs, canalLogs)
	fmt.Printf("passou do goroutine")
	ui.Inicializa(canalLogs)

	/*for _, cl := range containerLogs {
		clResultado := logs.PickLogs(cl)
		logs.PrintLog(clResultado)
	}*/

	//launchK8sJob(clientSet, jobName, containerImage, entryCommand)
}
