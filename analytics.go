package main

import (
    "strings"
    "fmt"
    "time"
)

var cache = make(map[string] Namespace)
var sizes = []int{1}

type Head struct { Namespace }
type Namespace []*Operation
type Operation struct {
    Name string
    Backlog []*Backlog
    Score float64
}
type Backlog struct {
    Operations float64 // Total operations
    Latency float64 // Average latency
}

func (n Namespace) Len() int { return len(n) }
func (n Namespace) Less(i, j int) bool { return n[i].Score > n[j].Score }
func (n Namespace) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

func compile() {
    var key string
    for _, size := range sizes {
        for name, _ := range database {
            // Declare/extract variables.
            split := strings.Split(name, ".")
            namespace, operation := split[0], split[1]
            key = namespace + "." + fmt.Sprint(size)
            
            // Create cache entry if does not exist.
            _, ok := cache[key]
            if !ok { cache[key] = []*Operation{} }
            
            // Calculate first entry and total score.
            operations, latency := countUp(name, size, 1)
            score := genScore(operations, latency)
            
            // Create operation entry.
            oper := new(Operation)
            oper.Name = operation
            oper.Backlog = []*Backlog{{operations, latency}}
            oper.Score = score
            cache[key] = append(cache[key], oper)
            where := len(cache[key]) - 1
            
            // Generate the rest of the backlog entries.
            for n := 2;n <= 10;n++ {
                operations, latency = countUp(name, size, n)
                
                backlog := new(Backlog)
                backlog.Operations, backlog.Latency = operations, latency
                cache[key][where].Backlog = append(cache[key][where].Backlog, backlog)
            }
        }
    }
}

func countUp(name string, size, place int) (float64, float64) {
    var operations, latency float64 = 0, 0
    offset := (place - 1) * size
    now := int(time.Now().Unix() / 60)
    
    for i := offset;i < size + offset;i++ {
        min := now - offset
        _, ok := database[name][min]
        if !ok { continue }
        
        operations += database[name][min][0]
        latency += database[name][min][1]
    }
    
    if operations == 0 { return 0, 0 } // DIVISION BY ZERO!!! :O
    latency = latency / operations
    return operations, latency
}

func genScore(operations float64, latency float64) float64 {
    return (operations * latency)
}
