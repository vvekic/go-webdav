package main

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
)

const (
	logLevelNone = "NONE"
	logLevelAll  = "ALL"
	envPrefix    = "WEBDAV_PREFIX"
	envLoglevel  = "WEBDAV_LOGLEVEL"
	envInMemory  = "WEBDAV_INMEMORY"
	pathRoot     = "/webdav/root"
	pathLog      = "/webdav/log/webdav.log"
)

func logger() func(*http.Request, error) {
	switch os.Getenv(envLoglevel) {
	case logLevelNone:
		return nil
	case logLevelAll:
		return func(r *http.Request, err error) {
			log.Printf("REQUEST %s %s length:%d %s %s\n", r.Method, r.URL, r.ContentLength, r.RemoteAddr, r.UserAgent())
		}
	default:
		return func(r *http.Request, err error) {
			if err != nil {
				log.Printf("ERROR %v\n", err)
			}
		}
	}
}

func filesystem() webdav.FileSystem {
	switch os.Getenv(envInMemory) {
	case "true":
		log.Printf("INFO using in-memory filesystem")
		return webdav.NewMemFS()
	default:
		if err := os.Mkdir(pathRoot, os.ModePerm); !os.IsExist(err) {
			log.Fatalf("FATAL %v", err)
		}
		log.Printf("INFO using local filesystem at %s", pathRoot)
		return webdav.Dir(pathRoot)
	}
}

func main() {
	logFile, err := os.OpenFile(pathLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("FATAL %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	h := &webdav.Handler{
		Prefix:     os.Getenv(envPrefix),
		FileSystem: filesystem(),
		LockSystem: webdav.NewMemLS(),
		Logger:     logger(),
	}

	http.HandleFunc("/", h.ServeHTTP)
	http.ListenAndServe("", h)
}
