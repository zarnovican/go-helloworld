package main

import (
	"log"
	"log/syslog"
	"net/http"
	"os"
)

type Config struct {
	LOG_TARGET   string
	PORT         string
	SERVICE_NAME string
}

var config *Config

var (
	debug   = log.New(os.Stderr, "DEBUG ", log.LstdFlags|log.Lshortfile)
	info    = log.New(os.Stderr, "INFO ", log.LstdFlags|log.Lshortfile)
	warning = log.New(os.Stderr, "WARNING ", log.LstdFlags|log.Lshortfile)
	err     = log.New(os.Stderr, "ERROR ", log.LstdFlags|log.Lshortfile)
)

func loadConfig() *Config {
	ret := &Config{}
	ret.LOG_TARGET = os.Getenv("LOG_TARGET")
	if ret.LOG_TARGET == "" {
		ret.LOG_TARGET = "console"
	}
	ret.PORT = os.Getenv("PORT")
	if ret.PORT == "" {
		ret.PORT = "8080"
	}
	ret.SERVICE_NAME = os.Getenv("SERVICE_NAME")
	if ret.SERVICE_NAME == "" {
		ret.SERVICE_NAME = "helloworld"
	}
	return ret
}

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
	w.Write([]byte("hello\n"))
}

func main() {
	config = loadConfig()
	if config.LOG_TARGET == "syslog" {
		info.Printf("Logging is redirected to systemd journal. Tail with \"journalctl -t %s -f\"", config.SERVICE_NAME)
		setupSyslog()
	}
	http.HandleFunc("/", root)

	debug.Print("sample Debug message")
	info.Print("sample Info message")
	warning.Print("sample Warning message")
	err.Print("sample Error message\nfoo\n    bar\nfoo2\n    bar2")

	info.Printf("Listening on 0.0.0.0:%s", config.PORT)
	http.ListenAndServe(":"+config.PORT, nil)
}
