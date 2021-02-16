package main

import (
	"flag"
	"fmt"

	"github.com/k8smon/config"
	"github.com/k8smon/kubernetes"
	"github.com/k8smon/logs"
)

func main() {
	jobName := flag.String("jobname", "test-job", "The name of the job")
	containerImage := flag.String("image", "ubuntu:latest", "Name of the container image")
	entryCommand := flag.String("command", "ls", "The command to run inside the container")

	flag.Parse()

	fmt.Printf("Args : %s %s %s\n", *jobName, *containerImage, *entryCommand)

	configuracoes := config.CarregaConfiguracoes()

	clientSet := kubernetes.ConectarK8s()
	containerLogs := kubernetes.Logs(clientSet, configuracoes)

	for _, cl := range containerLogs {
		clResultado := logs.PickLogs(cl)
		logs.PrintLog(clResultado)
	}

	//launchK8sJob(clientSet, jobName, containerImage, entryCommand)
}
