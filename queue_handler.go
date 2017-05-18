package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

type QueueHandler struct {
	jobs     <-chan amqp.Delivery
	openJobs map[string]*amqp.Delivery
	amqpman  *AMQPManager
	config   *MQConfig
}

func NewQueueHandler(config *MQConfig) (*QueueHandler, error) {
	amqpman, err := NewAMQPManager(&MQConfig{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
	}, getAppId())

	if err != nil {
		return nil, err
	}

	fmt.Printf("(Jobs are served from %s/%s)\n\n", config.Host, config.InputQueue)

	jobs, err := amqpman.consume(config.InputQueue)
	if err != nil {
		return nil, err
	}

	h := &QueueHandler{
		jobs:     jobs,
		config:   config,
		openJobs: make(map[string]*amqp.Delivery),
	}

	return h, nil
}

func (h *QueueHandler) ServeJob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request")

	job := <-h.jobs

	reader := bytes.NewReader(job.Body)
	gr, _ := gzip.NewReader(reader)

	body, _ := ioutil.ReadAll(gr)

	w.Header().Add("Content-Type", "application/json")
	w.Write(body)

	h.openJobs[job.CorrelationId] = &job

	fmt.Printf("Served %s\n", job.CorrelationId)
}

func (h *QueueHandler) HandleResponse(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	fmt.Println("Received response for " + vars["correlationId"])
	job, exists := h.openJobs[vars["correlationId"]]
	if exists {
		job.Ack(true)
		delete(h.openJobs, vars["correlationId"])
		w.WriteHeader(http.StatusOK)

		// First see if we have something in the URL
		response, err := parseResponseInQuery(r)
		if err != nil {
			fmt.Printf("Could not read response body %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if response == nil {
			response, err = ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
		}

		h.amqpman.channel.Publish(h.config.OutputExchange, "", false, false, amqp.Publishing{
			CorrelationId:   vars["correlationId"],
			Body:            response,
			AppId:           getAppId(),
			ContentEncoding: "application/json",
		})

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "CorrelationId not found or already acked")
	}
}
