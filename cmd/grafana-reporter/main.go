/*
   Copyright 2016 Vastech SA (PTY) LTD

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/IzakMarais/reporter/grafana"
	"github.com/IzakMarais/reporter/report"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var proto = flag.String("proto", "http://", "Grafana Protocol. Change to 'https://' if Grafana is using https. Reporter will still serve http.")
var ip = flag.String("ip", "localhost:3000", "Grafana IP and port")
var port = flag.String("port", ":8686", "Port to serve on")
var templateDir = flag.String("templates", "templates/", "Directory for custom TeX templates")
var logFile = flag.String("logfile", "stdout", "File to save log")
var logLevel = flag.String("loglevel", "INFO", "Set log level use one of ERROR,WARNING,INFO,DEBUG")

func main() {
	flag.Parse()
	log.SetFormatter(&PlainTextFormatter{})
	if *logFile == "stdout" {
		log.SetOutput(os.Stdout)
	} else {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}
	if *logLevel == "ERROR" {
		log.SetLevel(log.ErrorLevel)
		log.Println("Log Level: Error")
	} else if *logLevel == "WARNING" {
		log.SetLevel(log.WarnLevel)
		log.Println("Log Level: Warning")
	} else if *logLevel == "DEBUG" {
		log.SetLevel(log.DebugLevel)
		log.Println("Log Level: Debug")
	} else {
		log.SetLevel(log.InfoLevel)
		log.Println("Log Level: Info")
	}
	//'generated*'' variables injected from build.gradle: task 'injectGoVersion()'
	log.Info(fmt.Sprintf("grafana reporter, version: %s.%s-%s hash: %s", generatedMajor, generatedMinor, generatedRelease, generatedGitHash))
	log.Info(fmt.Sprintf("serving at %s and using grafana at %s", *port, *ip))

	router := mux.NewRouter()
	RegisterHandlers(
		router,
		ServeReportHandler{grafana.NewV4Client, report.New},
		ServeReportHandler{grafana.NewV5Client, report.New},
	)

	log.Fatal(http.ListenAndServe(*port, router))
}
