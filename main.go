package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync/atomic"
	"time"
	"github.com/doemoor/wed-server/api"
)



func main() {
	mux := http.NewServeMux()
	serverStruct := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  2 * time.Second,
	}
	
	var apiCfg = &api.ApiConfig{
		FileserverHits: atomic.Int32{},
	}

	
	mux.Handle("GET /app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", api.Healthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.MetricsReset)
	mux.HandleFunc("POST /api/validate_chirp", api.ValidateChirp)
	
	
	clearTerminal()
	log.Printf("Serving on port: %s\n", serverStruct.Addr)
	log.Fatal(serverStruct.ListenAndServe())

}

func clearTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}



