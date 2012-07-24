package main

import (
    zmq "github.com/alecthomas/gozmq"
    
    "encoding/json"
    "fmt"
    "strconv"
    "time"
)

var database = make(map[string] map[int] [2]float64)

func handler() {
    // Create sockets
    context, err := zmq.NewContext();handleErr(err)
    defer context.Close()

    sink, err := context.NewSocket(zmq.PULL);handleErr(err)
    defer sink.Close()
    sink.Bind(fmt.Sprint(config["listen"]))
    
    // Main work
    for true {
        packet := map[string] interface{} {
            "namespace" : "Unknown",
            "operation" : "Unknown",
            "count" : 0,
            "latency" : 0,
        }
        
        msgBytes, err := sink.Recv(0);handleErr(err)
        err = json.Unmarshal(msgBytes, &packet)
        if err != nil { continue } // Throw out bad packets.
        
        namespace := fmt.Sprint(packet["namespace"])
        operation := fmt.Sprint(packet["operation"])
        count, cErr := strconv.ParseFloat(fmt.Sprint(packet["count"]), 64)
        latency, lErr := strconv.ParseFloat(fmt.Sprint(packet["latency"]), 64)
        key := namespace + "." + operation
        
        _, ok := database[key]
        
        // Throw out packets that have bad information.
        if !ok  || cErr != nil || lErr != nil { continue }
        
        // Save the packet to memory
        minute := int(time.Now().Unix() / 60)
        _, created := database[key][minute]
        
        if created {
            count += database[key][minute][0]
            latency += database[key][minute][1]
        }
        
        database[key][minute] = [2]float64{count, latency}
    }
}

func clean() {
    cutoffSecs := int64((config["cutoff"].(time.Duration)).Seconds())
    cutoff := int((time.Now().Unix() - cutoffSecs) / 60)

    for name, _ := range database {
        for minute, _ := range database[name] {
            if minute <= cutoff {
                delete(database[name], minute)
            }
        }
    }
}
