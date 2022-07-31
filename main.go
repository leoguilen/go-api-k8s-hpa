package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/google/uuid"
)

var (
	httpPort string = os.Getenv("HTTP_PORT")
)

type ContextDetails struct {
	Environment EnvironmentInfo `json:"environment"`
	Request     RequestInfo     `json:"request"`
}

type EnvironmentInfo struct {
	MachineName string `json:"machineName"`
	OSName      string `json:"osName"`
	OSArch      string `json:"osArch"`
	GoVersion   string `json:"goVersion"`
}

type RequestInfo struct {
	Identifier uuid.UUID `json:"identifier"`
	Host       string    `json:"host"`
	Path       string    `json:"path"`
	Scheme     string    `json:"schema"`
	Method     string    `json:"method"`
}

func NewContextInfo(r *http.Request) *ContextDetails {
	h, _ := os.Hostname()
	envInfo := EnvironmentInfo{
		MachineName: h,
		OSName:      runtime.GOOS,
		OSArch:      runtime.GOARCH,
		GoVersion:   runtime.Version(),
	}

	requestInfo := RequestInfo{
		Identifier: uuid.New(),
		Host:       r.Host,
		Path:       r.URL.Path,
		Scheme:     r.URL.Scheme,
		Method:     r.Method,
	}

	return &ContextDetails{
		Environment: envInfo,
		Request:     requestInfo,
	}
}

// func GetOSDetails() ([]string, error) {
// 	output, err := exec.Command("uname", "-smnr").Output()
// 	if err != nil {
// 		log.Printf("[ERROR] - Error get OS details: %v", err.Error())
// 		return nil, err
// 	}
// 	log.Printf("[DEBUG] - Command output: %v", string(output))
// 	outputParts := strings.Split(string(output), " ")
// 	return outputParts, nil
// }

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("[ERROR] - HTTP Method %v is not allowed --- %v", r.Method, http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctxInfo := NewContextInfo(r)

	log.Printf("[INFO] - HTTP %v %v --- %v", r.Host, r.Method, http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ctxInfo)
}

func main() {
	http.HandleFunc("/server-info", Handler)

	log.Printf("[INFO] - About to listen on %v. Go to http://127.0.0.1:%v/", httpPort, httpPort)
	log.Fatalf("[FATAL] - %v", http.ListenAndServe(fmt.Sprintf(":%v", httpPort), nil))
}
