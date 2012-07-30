package main

import (
    "fmt"
    "html/template"
    "io/ioutil"
    "log"
    "mime"
    "net/http"
    "sort"
    "strconv"
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
        "equals" : func(a, b int) bool { return a == b },
        "getSizeName" : func(s int) string { return sizeNames[s] },
        "getSizeLetter" : func(s int) string { return sizeLetters[s] },
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
        
        size int
        stringSize string
        
        err error
    )
    
    // Get namespace
    namespace := vals.Get("namespace")
    if namespace == "" {
        for namespace, _ = range (config["namespaces"].(map[string] interface{})) { break } // WTF
    } else {
        //Validate namespace
        _, ok = (config["namespaces"].(map[string] interface{}))[namespace]
        if !ok { data, view = error_page("Invalid namespace."); goto show }
    }
    
    // Get size
    stringSize = vals.Get("size")
    if stringSize == "" {
        size = 1
    } else {
        size, err = strconv.Atoi(stringSize)
        
        //Validate size
        _, ok = sizeNames[size]
        if !ok { data, view = error_page("Invalid size."); goto show }
        if err != nil { data, view = error_page(err.Error()); goto show }
    }
    
    data, view = controller_view(namespace, size)
    
    show:
    namespaces := []string{}
    for name, _ := range config["namespaces"].(map[string] interface{}) {
        namespaces = append(namespaces, name)
    }
    sort.Strings(namespaces)
    
    data["name"] = config["project"].(string)
    data["namespaces"] = namespaces
    err = views.Lookup(view).Execute(w, data)
    handleErr(err)
}

func controller_index() (map[string] interface{}, string) {
    return map[string] interface{} {}, "index"
}

func controller_view(namespace string, size int) (map[string] interface{}, string) {
    data := cache[namespace + "." + fmt.Sprint(size)]
    sort.Sort(Head{ data })
    
    return map[string] interface{} {
        "namespace" : namespace,
        "sizeNames" : sizeNames,
        "size" : size,
        "operations" : data,
    }, "view"
    
}


func error_page(err string) (map[string] interface{}, string) {
    return map[string] interface{} {"error" : err}, "error"
}
