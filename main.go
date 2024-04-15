package main

import(
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
    "strconv"
    "sync"
)

const index string = `
<!DOCTYPE html>
<html data-theme="night">
    <head>
        <link href="http://cdn.jsdelivr.net/npm/daisyui@latest/dist/full.min.css" rel="stylesheet"/>
        <script src="https://cdn.tailwindcss.com"></script>
        <script src="https://unpkg.com/axios/dist/axios.min.js"></script>
        <script src="https://unpkg.com/htmx.org@latest"></script>
    </head>
    <body class="h-svh flex flex-col items-center justify-center gap-4">
        <form class="flex items-center justify-center gap-4">
            <input type="file" name="textfile" class="file-input file-input-bordered file-input-primary"></input>
            <button type="submit" class="btn btn-primary" 
            hx-post="/count" hx-encoding="multipart/form-data" hx-target="#new"
            >Upload</button>
        </form>
        <div id="new"></div>
    </body>
</html>
`

const responseString string = `
<div class="stats bg-neutral shadow">
    <div class="stat">
        <div class="stat-title">Lines</div>
        <div class="stat-value text-primary">{{lines}}</div>
    </div>

    <div class="stat">
        <div class="stat-title">Words</div>
        <div class="stat-value text-secondary">{{words}}</div>
    </div>
</div>
`

type Response struct{
    mutex sync.Mutex
    res string
}

func (r *Response) words(s string) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    count := strconv.Itoa(len(strings.Fields(s)))
    r.res = strings.Replace(r.res, "{{words}}", count, -1)
}

func (r *Response) lines(s string) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    count := strconv.Itoa(strings.Count(s, "\n"))
    r.res = strings.Replace(r.res, "{{lines}}", count, -1)
}

func Index(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, index)
}

func Count(w http.ResponseWriter, r *http.Request){
    r.ParseMultipartForm(10<<20)
    file, _, err :=r.FormFile("textfile")
    if(err!=nil){
        fmt.Println(err)
        return
    }
    defer file.Close()

    content,_ := ioutil.ReadAll(file)

    response := Response{
        res: responseString,
    }

    var wg sync.WaitGroup
    countWords := func(){
                response.words(string(content))
                defer wg.Done()
                }

    countLines := func(){
                response.lines(string(content))
                defer wg.Done()
                }

    wg.Add(2)

    go countWords()
    go countLines()

    wg.Wait()

    fmt.Fprintf(w, response.res)

}

func main(){
    http.HandleFunc("/", Index)
    http.HandleFunc("/count", Count)
    fmt.Println("Serving on http://localhost:7000")
    http.ListenAndServe(":7000",nil)
}
