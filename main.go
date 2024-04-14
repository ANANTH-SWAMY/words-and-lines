package main

import(
    "fmt"
    "net/http"
)

const index string = `
<!DOCTYPE html>
<html data-theme="night">
    <head>
        <link href="http://cdn.jsdelivr.net/npm/daisyui@latest/dist/full.min.css" rel="stylesheet"/>
        <script src="https://cdn.tailwindcss.com"></script>
        <script src="https://unpkg.com/axios/dist/axios.min.js"></script>
    </head>
    <body class="h-svh flex flex-col items-center justify-center">
        <div class="flex items-center justify-center gap-4">
            <input type="file" class="file-input file-input-bordered file-input-primary"></input>
            <button class="btn btn-primary">Upload</button>
        </div>
    </body>
</html>
`

func Index(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, index)
}

func main(){
    http.HandleFunc("/", Index)
    fmt.Println("Serving on http://localhost:7000")
    http.ListenAndServe(":7000",nil)
}
