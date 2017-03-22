package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DOCKER_TASK_SLOT string `envconfig:"DOCKER_TASK_SLOT" default:"1"`
	LOG_TARGET       string `envconfig:"LOG_TARGET" default:"console"`
	PORT             string `envconfig:"PORT" default:"8080"`
	SERVICE_NAME     string `envconfig:"SERVICE_NAME" default:"helloworld"`
}

var config Config

// Set at build time.
var GitDescribe = "snapshot"

var (
	debug   = log.New(os.Stderr, "DEBUG ", log.LstdFlags|log.Lshortfile)
	info    = log.New(os.Stderr, "INFO ", log.LstdFlags|log.Lshortfile)
	warning = log.New(os.Stderr, "WARNING ", log.LstdFlags|log.Lshortfile)
	err     = log.New(os.Stderr, "ERROR ", log.LstdFlags|log.Lshortfile)
)

func setupSyslog() {
	var logger *syslog.Writer
	var e error

	logger, e = syslog.New(syslog.LOG_USER|syslog.LOG_DEBUG, config.SERVICE_NAME)
	if e == nil {
		debug.SetOutput(logger)
		debug.SetFlags(log.Lshortfile)
		debug.SetPrefix("")
	}
	logger, e = syslog.New(syslog.LOG_USER|syslog.LOG_INFO, config.SERVICE_NAME)
	if e == nil {
		info.SetOutput(logger)
		info.SetFlags(log.Lshortfile)
		info.SetPrefix("")
	}
	logger, e = syslog.New(syslog.LOG_USER|syslog.LOG_WARNING, config.SERVICE_NAME)
	if e == nil {
		warning.SetOutput(logger)
		warning.SetFlags(log.Lshortfile)
		warning.SetPrefix("")
	}
	logger, e = syslog.New(syslog.LOG_USER|syslog.LOG_ERR, config.SERVICE_NAME)
	if e == nil {
		err.SetOutput(logger)
		err.SetFlags(log.Lshortfile)
		err.SetPrefix("")
	}
}

func root(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello!\n"))
}

func get_info(w http.ResponseWriter, req *http.Request) {
	iam := config.SERVICE_NAME
	if config.DOCKER_TASK_SLOT != "" {
		iam = iam + "." + config.DOCKER_TASK_SLOT
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "<unknown>"
	}
	w.Write([]byte(fmt.Sprintf("Go %s (%s) on %s: your IP %s\n", iam, GitDescribe, hostname, req.RemoteAddr)))
}

func log_sample(w http.ResponseWriter, req *http.Request) {
	debug.Print("called /log/<foo> endpoint")

	if strings.HasSuffix(req.URL.Path, "/info") {
		info.Print("sample Info message")
	} else if strings.HasSuffix(req.URL.Path, "/warning") {
		warning.Print("sample Warning message")
	} else if strings.HasSuffix(req.URL.Path, "/error") {
		err.Print("sample Error message\nfoo\n    bar\nfoo2\n    bar2")
	} else {
		err.Printf("path %s not found", req.URL.Path)
		w.Write([]byte("not found\n"))
		return
	}
	w.Write([]byte("ok\n"))
}

func main() {
	e := envconfig.Process("helloworld", &config)
	if e != nil {
		log.Fatal(e.Error())
	}
	if config.LOG_TARGET == "syslog" {
		info.Printf("Logging is redirected to systemd journal. Tail with \"journalctl -t %s -f\"", config.SERVICE_NAME)
		setupSyslog()
	}
	http.HandleFunc("/", root)
	http.HandleFunc("/info", get_info)
	http.HandleFunc("/log/", log_sample)

	info.Printf("Listening on 0.0.0.0:%s", config.PORT)
	e = http.ListenAndServe(":"+config.PORT, nil)
	if e != nil {
		err.Print(e.Error())
		log.Fatal(e.Error())
	}
}
