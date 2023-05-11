package main

import (
	"encoding/json"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"os/exec"
)

type Req struct {
	Curl string `json:"curl_command"`
}

func main() {
	http.Handle("/", enableCORS(http.HandlerFunc(rootHandler)))
	http.Handle("/run-curl", enableCORS(http.HandlerFunc(runCurlHandler)))
	http.ListenAndServe(":8088", nil)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		if r.URL.Query().Get("seed") != "" {
			processRequest(w, r)
		} else {
			http.ServeFile(w, r, "index.html")
		}
	} else {
		http.Error(w, "Use the GET method.", http.StatusMethodNotAllowed)
	}
}

func runCurlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Use the POST method.", http.StatusMethodNotAllowed)
		return
	}

	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// curlCommand := r.FormValue("curl_command")
	// if curlCommand == "" {
	// 	http.Error(w, "Please send the curl command.", http.StatusBadRequest)
	// 	return
	// }

	cmd := exec.Command("bash", "-c", req.Curl)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to run command: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "The command was successfully executed. Output:\n\n%s", string(output))
}

func processRequest(w http.ResponseWriter, r *http.Request) {
	seedStr := r.URL.Query().Get("seed")
	if seedStr == "" {
		http.Error(w, "Please enter the seed value.", http.StatusBadRequest)
		return
	}

	seed, err := strconv.ParseInt(seedStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid seed value.", http.StatusBadRequest)
		return
	}

	rand.Seed(seed)

	nID := generateRandHexEncodedNamespaceID()
	msg := generateRandMessage()

	fmt.Fprintf(w, "My hex-encoded namespace ID: %s\n\nMy hex-encoded message: %s", nID, msg)
}

func generateRandHexEncodedNamespaceID() string {
	nID := make([]byte, 8)
	_, err := rand.Read(nID)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(nID)
}

func generateRandMessage() string {
	lenMsg := rand.Intn(100)
	msg := make([]byte, lenMsg)
	_, err := rand.Read(msg)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(msg)
}
