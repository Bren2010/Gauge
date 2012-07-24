package main

import (
    "html/template"
    "log"
    "net/http"
    "sort"
)

var views *template.Template
var files = []string {
    "./html/index.html", 
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

func dispatcher(w http.ResponseWriter, r *http.Request) {
    vals := r.URL.Query()
    
    var (
        data map[string] interface{}
        view string
        ok bool
    )
    
    // Get namespace
    namespace := vals.Get("namespace")
    if namespace == "" { data, view = controller_index(); goto show }
    
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
        "operations" : data,
    }, "view"
    
}


func error_page(err string) (map[string] interface{}, string) {
    return map[string] interface{} {"error" : err}, "error"
}
