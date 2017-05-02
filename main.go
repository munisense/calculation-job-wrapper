package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	SERVICE_NAME = "calculation-job-wrapper"
	PORT         = 8765
)

var (
	// These can be injected at build time -ldflags "-X main.VERSION=dev main.BUILD_TIME=201610251410"
	VERSION    = "Undefined"
	BUILD_TIME = "Undefined"
)

func main() {
	fmt.Printf("Starting %s\n", getAppId())

	// Read flags
	var fileName string
	flag.StringVar(&fileName, "file", "", "Reads job input from file instead of from queue")
	flag.Parse()

	// Read Config
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	go func() {
		fmt.Printf("* Request new jobs at http://localhost:%d/input\n", PORT)
		fmt.Printf("* Post job results to http://localhost:%d/output/{correlationId}\n", PORT)

		r := mux.NewRouter()

		// When serving static files
		if fileName != "" {
			handler, err := NewFileHandler(fileName)
			if err != nil {
				panic(err)
			}

			r.HandleFunc("/input", handler.ServeJob)
			r.HandleFunc("/output/{correlationId}", handler.HandleResponse)
		} else {
			// When serving from the queue
			handler, err := NewQueueHandler(config.MQ)
			if err != nil {
				panic(err)
			}

			r.HandleFunc("/input", handler.ServeJob)
			r.HandleFunc("/output/{correlationId}", handler.HandleResponse)
		}

		err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
		if err != nil {
			panic(err)
		}
	}()

	// Make sure we dont terminate
	forever := make(chan bool)
	<-forever
}

// Returns the appId which is used to identify this instance
func getAppId() string {
	hostname, _ := os.Hostname()
	pid := os.Getpid()

	return fmt.Sprintf("%s[%s-%s]@%s-%d", SERVICE_NAME, VERSION, BUILD_TIME, hostname, pid)
}
