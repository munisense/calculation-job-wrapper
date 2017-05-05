package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type FileHandler struct {
	fileName string
	jobs     []string
	index    int
}

func NewFileHandler(fileName string) (*FileHandler, error) {
	filePath := CURRENT_PATH + string(os.PathSeparator) + "static" + string(os.PathSeparator) + fileName

	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	h := &FileHandler{
		jobs:     strings.Split(string(dat), "\n"),
		fileName: fileName,
	}

	return h, nil
}

func (h *FileHandler) ServeJob(w http.ResponseWriter, r *http.Request) {

	if len(h.jobs) <= h.index {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(h.jobs[h.index]))
	fmt.Printf("Served job %s#%d\n", h.fileName, h.index+1)
	h.index++
}

func (h *FileHandler) HandleResponse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Received response for " + vars["correlationId"])

	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Could not read response body %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		outputFile := CURRENT_PATH + string(os.PathSeparator) + "output" + string(os.PathSeparator) + vars["correlationId"] + ".json"
		fmt.Println("Writing response to " + outputFile)
		err := ioutil.WriteFile(outputFile, response, 0644)
		if err != nil {
			fmt.Printf("Could not write response body to file %v\n", err)
		}
	}
}
