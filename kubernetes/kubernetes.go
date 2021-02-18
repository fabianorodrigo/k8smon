package kubernetes

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/k8smon/colors"
	"github.com/k8smon/models"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	/*
		kubernetes : This package offers you APIs using which you can connect to your Kubernetes service.
		It provides different kinds of connection APIs.
		It also provides you the access to clientset interface, which can be used to specify and manipulate K8s objects
	*/

	config "github.com/k8smon/config"
	utils "github.com/k8smon/utils"
	kubernetes "k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

const minuto int64 = 60
const layoutHoraLog string = time.RFC3339Nano

var regExpHorarioLog = regexp.MustCompile(`([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{1,2}:[0-9]{2}:[0-9]{2}.[0-9]{1,9}Z)`)

/* ConectarK8S Conecta ao cluster baseando-se, primeiramente, na variável
 * de ambiente KUBECONFIG, caso não exista, buscará no path '/$HOME/.kube/config'.
 * Ao se conectar com sucesso, cria um objeto do tipo clientset e o retorna
 */
func ConectarK8s() *kubernetes.Clientset {
	//prioriza o que estiver na env KUBECONFIG
	configPath, exists := os.LookupEnv("KUBECONFIG")
	if !exists {
		home, exists := os.LookupEnv("HOME")
		if !exists {
			home = "/root"
		}
		configPath = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Panicln("Falha ao criar configurações K8s")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln("Falha ao criar clienteset K8s")
	}

	return clientset
}

//Pods Busca os pods no cluster ao qual o {clientset} está conectado
func Pods(clientset *kubernetes.Clientset, configuracoes config.Dictionary,
	canalPods chan models.Pod, namespaces ...string) {
	clienteCoreV1 := clientset.CoreV1()
	namespaceInterface := clienteCoreV1.Namespaces()
	namespaceList, err := namespaceInterface.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Panicf("Falha ao buscar namespaces ")
	}

	//Percorre os namespaces
	for _, namespace := range (*namespaceList).Items {
		//Se o namespace estiver configurado para não ser analisado, salta para o próximo
		if utils.ArrayContains(configuracoes["ignoraNamespaces"].([]interface{}), namespace.Name) {
			continue
		}
		podInterface := clienteCoreV1.Pods(namespace.Name)

		podList, err := podInterface.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Panicf("Falha ao buscar PODs do namespace %s", namespace.Name)
		}
		// List() returns a pointer to slice, derefernce it, before iterating
		for _, podInfo := range (*podList).Items {
			canalPods <- models.Pod{Nome: podInfo.Name, IP: podInfo.Status.PodIP, Namespace: namespace.Name,
				Fase: string(podInfo.Status.Phase), NomeNode: podInfo.Spec.NodeName, IPNode: podInfo.Status.HostIP, NomeConteineres: getContainerNames(podInfo)}
		}
	}
	close(canalPods)
}

//Logs Busca os logs no cluster ao qual o {clientset} está conectado
func Logs(clientset *kubernetes.Clientset, configuracoes config.Dictionary,
	funcaoPick func(models.ContainerLog, chan models.ContainerLog),
	canalLogs chan models.ContainerLog, namespaces ...string) []models.ContainerLog {
	retorno := []models.ContainerLog{}
	clienteCoreV1 := clientset.CoreV1()
	namespaceInterface := clienteCoreV1.Namespaces()
	namespaceList, err := namespaceInterface.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Panicf("Falha ao buscar namespaces ")
	}

	var configMinutos int64 = int64(configuracoes["janelaMinutos"].(float64)) * minuto
	//Percorre os namespaces
	for _, namespace := range (*namespaceList).Items {
		//Se o namespace estiver configurado para não ser analisado, salta para o próximo
		if utils.ArrayContains(configuracoes["ignoraNamespaces"].([]interface{}), namespace.Name) {
			continue
		}
		podInterface := clienteCoreV1.Pods(namespace.Name)

		/*options := metav1.ListOptions{
			LabelSelector: "app=minio",
		}*/
		podList, err := podInterface.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Panicf("Falha ao buscar PODs do namespace %s", namespace.Name)
		}
		// List() returns a pointer to slice, derefernce it, before iterating
		for _, podInfo := range (*podList).Items {
			/*colors.Redf("pods-name=%v.%v\n", podInfo.Namespace, podInfo.Name)
			fmt.Printf("pods-namespace=%v\n", podInfo.Namespace)
			fmt.Printf("pods-status=%v\n", podInfo.Status.Phase)
			fmt.Printf("pods-message=%v\n", podInfo.Status.Message)
			fmt.Printf("pods-reason=%v\n", podInfo.Status.Reason)
			fmt.Printf("pods-hostIP=%v\n", podInfo.Status.HostIP)
			fmt.Printf("pods-podIP=%v\n", podInfo.Status.PodIP)
			fmt.Printf("pods-node=%v,%v\n", podInfo.Status.NominatedNodeName, podInfo.Spec.NodeName)
			fmt.Printf("pods-condition=%v\n", podInfo.Status.Conditions)*/

			for _, container := range getContainerNames(podInfo) {
				podLogOpts := v1.PodLogOptions{Container: container, Timestamps: true, SinceSeconds: &configMinutos}

				containerLog := models.ContainerLog{Namespace: namespace.Name, PodName: podInfo.Name, PodPhase: string(podInfo.Status.Phase), PodNode: podInfo.Spec.NodeName, PodIP: podInfo.Status.PodIP, ContainerName: container}

				logsReq := podInterface.GetLogs(podInfo.Name, &podLogOpts)
				logs, err := logsReq.Stream(context.TODO())
				if err != nil {
					panic("Falha ao abrir stream de logs do pod " + err.Error())
				}
				defer logs.Close()
				buf := new(bytes.Buffer)
				_, err = io.Copy(buf, logs)
				if err != nil {
					log.Panicf("Falha na cópia dos logs para o buffer. Pod: %s", podInfo.Name)
				}
				for _, linha := range strings.Split(buf.String(), "\n") {
					if len(strings.TrimSpace(linha)) > 0 {
						//o serviço pode retornar esse tipo de mensagem: unable to retrieve container logs for docker://8c30969187f7de8fc9dc
						if strings.HasPrefix(linha, "unable to retrieve") {
							containerLog.GetLogsFail = true
						} else {
							hora := regExpHorarioLog.FindAllString(linha, -1)
							if hora == nil {
								log.Panicf("Horário não encontrado no log: %s", linha)
							}
							horarioLog, err := time.Parse(layoutHoraLog, hora[0])
							if err != nil {
								colors.Redln(`Falha ao converter horário do log: {{hora[0]}}`)
							}
							containerLog.Logs = append(containerLog.Logs, models.Log{Horario: horarioLog, Log: linha})
						}
					}
				}
				funcaoPick(containerLog, canalLogs)
				retorno = append(retorno, containerLog)
			}

		}
	}
	close(canalLogs)
	return retorno
}

//Pega os nomes dos containeres do POD
func getContainerNames(pod v1.Pod) []string {
	names := []string{}
	//Init containers
	for _, c := range pod.Spec.InitContainers {
		names = append(names, c.Name)
	}
	//containers
	for _, c := range pod.Spec.Containers {
		names = append(names, c.Name)
	}

	return names
}

func launchK8sJob(clientset *kubernetes.Clientset, jobName *string, image *string, cmd *string) {
	jobs := clientset.BatchV1().Jobs("default")
	var backOffLimit int32 = 0

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *jobName,
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    *jobName,
							Image:   *image,
							Command: strings.Split(*cmd, " "),
						},
					},
					//Note that, we are setting RestartPolicy as Never and BackOffLimit to zero,
					//it means, if the job fails, it will never be restarted by the job controller
					//because of BackOffLimit and the Pod created by the job will die gracefully
					//because of RestartPolicy being set to Never
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln("Falha na criação de job K8s.")
	}

	//print job details
	log.Println("Job K8s criado com sucesso")
}
