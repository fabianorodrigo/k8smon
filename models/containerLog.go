package models

import "time"

//Log contém uma linha de log com data/hora de sua ocorrência
type Log struct {
	Horario time.Time
	Log     string
}

//ContainerLog contém os logs de um container específico
type ContainerLog struct {
	Namespace     string
	PodName       string
	PodPhase      string
	PodNode       string
	PodIP         string
	ContainerName string
	GetLogsFail   bool
	Logs          []Log
}
