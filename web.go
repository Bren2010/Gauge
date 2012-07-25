package main

import (
    "fmt"
    "html/template"
    "io/ioutil"
    "log"
    "mime"
    "net/http"
    "sort"
    "strings"
)

var views *template.Template
var files = []string {
    "./html/view.html",
    "./html/error.html", 
}

func server() {
    funcs := map[string] interface{} {
        "plus1" : func(i int) int { return i + 1 },
    }
    
    var err error
    views, err = template.ParseFiles("./html/global.html");handleErr(err)
    views.Funcs(funcs)
    views.ParseFiles(files...)
    
    log.Fatal(http.ListenAndServe(config["server"].(string), nil))
}

func httpCdn(w http.ResponseWriter, r *http.Request) {
    uri := "." + r.RequestURI
    data, err := ioutil.ReadFile(uri)
    
    if err == nil {
        l := strings.LastIndex(uri, ".")
        ext := uri[l:]
        w.Header().Add("Content-Type", mime.TypeByExtension(ext))
        
        fmt.Fprintf(w, "%s", data)
    } else {
        fmt.Fprint(w, err)
    }
}

func dispatcher(w http.ResponseWriter, r *http.Request) {
    vals := r.URL.Query()
    
    var (
        data map[string] interface{}
        view string
        ok bool
    )
    
    // Get namespace
    namespace := vals.Get("namespace")
    if namespace == "" {
        var ns string
        for ns, _ = range (config["namespaces"].(map[string] interface{})) { break } // WTF
        data, view = controller_minute(ns)
        goto show
    }
    
    //Validate namespace
    _, ok = (config["namespaces"].(map[string] interface{}))[namespace]
    if !ok { data, view = error_page("Invalid namespace."); goto show }
    
    data, view = controller_minute(namespace)
    
    show:
    namespaces := []string{}
    for name, _ := range config["namespaces"].(map[string] interface{}) {
        namespaces = append(namespaces, name)
    }
    sort.Strings(namespaces)
    
    data["name"] = config["project"].(string)
    data["namespaces"] = namespaces
    err := views.Lookup(view).Execute(w, data)
    handleErr(err)
}

func controller_index() (map[string] interface{}, string) {
    return map[string] interface{} {}, "index"
}

func controller_minute(namespace string) (map[string] interface{}, string) {
    data := cache[namespace + ".1"]
    sort.Sort(Head{ data })
    
    return map[string] interface{} {
        "namespace" : namespace,
        "operations" : data,
    }, "view"
    
}


func error_page(err string) (map[string] interface{}, string) {
    return map[string] interface{} {"error" : err}, "error"
}
