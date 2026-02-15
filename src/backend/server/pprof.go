package server

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func init() {
	if v, ok := os.LookupEnv("PPROF"); ok && v != "" {
		// http://0.0.0.0:6060/debug/pprof/
		go func() {
			log.Println("pprof server started http://0.0.0.0:6060")
			log.Println(http.ListenAndServe(":6060", nil))
		}()
	}
}
