package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
)

type Packet struct {
    namespace string
    operation string
    count string
    latency string
}

var config map[string] interface{}

func main() {
    log.Print("Starting up...")
    defer log.Print("Shutting down.")
    
    // Read config.
    configBytes, err := ioutil.ReadFile("./config.json");handleErr(err)
    err = json.Unmarshal(configBytes, &config);handleErr(err)
    
    // Validate config
    _, pOk := config["project"]
    _, nOk := config["namespaces"]
    _, lOk := config["listen"]
    _, sOk := config["server"]
    _, cOk := config["cutoff"]
    
    wrong := ""
    if !pOk { wrong = "project" }
    if !nOk { wrong = "namespaces" }
    if !lOk { wrong = "listen" }
    if !sOk { wrong = "server" }
    if !cOk { wrong = "cutoff" }
    if wrong != "" { log.Fatal("Required config value:  " + wrong) }
    
    config["cutoff"], err = time.ParseDuration(config["cutoff"].(string))
    handleErr(err)
    
    // Build data structures
    log.Print("Building data structures...")
    
    for name, ops := range config["namespaces"].(map[string] interface{}) {
        for _, op := range ops.([]interface{}) {
            database[name + "." + fmt.Sprint(op)] = make(map[int] [2]float64)
        }
    }
    
    // Start real work.
    http.HandleFunc("/", dispatcher)
    http.HandleFunc("/style/", httpCdn)
    
    go handler()
    go server()
    
    log.Print("Started.")
    
    // Watch for Ctrl+C's, and handle them appropriately so defers are run.
    signals := make(chan os.Signal)
    signal.Notify(signals, os.Interrupt)
    
    // Compile and clean the datastore every minute.
    ticker := time.Tick(time.Minute)
    
    compile()
    for true {
        select {
            case <- signals: // Die
                return
            case <- ticker: // Clean
                cache = make(map[string] Namespace)
                compile()
                clean()
        }
    }
}

func handleErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
