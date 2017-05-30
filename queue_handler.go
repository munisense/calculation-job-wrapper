package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

type QueueHandler struct {
	jobs <-chan amqp.Delivery

	lock     sync.RWMutex
	openJobs map[string]*amqp.Delivery

	amqpman *AMQPManager
	config  *MQConfig
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

	// Panic on connection problem
	go func() {
		err := <-amqpman.closeNotifier
		panic(err)
	}()

	h := &QueueHandler{
		jobs:     jobs,
		config:   config,
		openJobs: make(map[string]*amqp.Delivery),
		amqpman:  amqpman,
	}

	return h, nil
}

func (h *QueueHandler) ServeJob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request")

	job := <-h.jobs

	body := job.Body

	// Unzip contents if needed
	if job.ContentEncoding == "gzip" {
		reader := bytes.NewReader(job.Body)
		gr, _ := gzip.NewReader(reader)
		body, _ = ioutil.ReadAll(gr)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(body)

	h.lock.Lock()
	h.openJobs[job.CorrelationId] = &job
	h.lock.Unlock()

	go h.HandleJobTimeout(job.CorrelationId, &job)

	fmt.Printf("Served %s\n", job.CorrelationId)
}

func (h *QueueHandler) HandleResponse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Received response for " + vars["correlationId"])

	h.lock.RLock()
	job, exists := h.openJobs[vars["correlationId"]]
	h.lock.RUnlock()

	if exists {
		h.lock.Lock()
		delete(h.openJobs, vars["correlationId"])
		h.lock.Unlock()

		job.Ack(false)

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

// Nacks a job if it is not handled within a timeout period
func (h *QueueHandler) HandleJobTimeout(id string, job *amqp.Delivery) {
	<-time.After(2 * time.Minute)

	h.lock.RLock()
	job, exists := h.openJobs[id]
	h.lock.RUnlock()

	if exists {
		h.lock.Lock()
		delete(h.openJobs, id)
		h.lock.Unlock()

		if job.Redelivered {
			fmt.Println("Rejecting job " + id + " because we received no response twice")
			job.Reject(false)
		} else {
			fmt.Println("Requeued job " + id + " because no response has been received within 2 minutes")
			job.Nack(false, true)
		}
	}
}

func (h *QueueHandler) TransferFileToQueue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Posting jobs in " + vars["file"] + " on queue")

	filePath := CURRENT_PATH + string(os.PathSeparator) + "static" + string(os.PathSeparator) + vars["file"]

	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not read file: %v", err)
		return
	}

	for i, jobStr := range strings.Split(string(dat), "\n") {

		// Try to parse the job
		var job map[string]interface{}
		if err := json.Unmarshal([]byte(jobStr), &job); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Could not parse %s:%d: %v", vars["file"], i, err)
		}

		correlationId := job["correlation_id"].(string)

		if err := h.amqpman.channel.Publish(h.config.InputQueue, "", false, false, amqp.Publishing{
			CorrelationId:   correlationId,
			Body:            []byte(jobStr),
			AppId:           getAppId(),
			ContentEncoding: "application/json",
		}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Could not publish %s:%d: %v", vars["file"], i, err)
		}
	}
}
