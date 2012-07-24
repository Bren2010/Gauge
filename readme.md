Gauge
=====
A real-time profiler that's not worth your time.  Dependencies:
```go get github.com/alecthomas/gozmq```

Data is split up into namespaces and operations (Namespace:  Website, 
Operations:  Create, Read, Update, Delete).  Send data to the configured 
port (localhost:5553 by default) through ZeroMQ in the form of:
```
{
    namespace: "Website",
    operation: "Read", 
    count: 1, 
    latency: 5
}
```
Count is the total number of times the operation was run after sending 
the last packet pertaining to this operation.  Latency is the total 
number of milliseconds the batch of operations took.  The compiled 
results are shown on a web frontend, sorted by score.


Todo
----
1.  Styled HTML
2.  Viewing data at different increments.
