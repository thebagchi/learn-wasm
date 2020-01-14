package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

const dir = "./www"

func Serve() {
	web := http.FileServer(http.Dir(dir))
	handlers := http.NewServeMux()
	handlers.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		dump, err := httputil.DumpRequest(request, true)
		if nil != err {
			fmt.Println("Error: ", err)
		}
		fmt.Println(string(dump))
		defer request.Body.Close()
		body, _ := ioutil.ReadAll(request.Body)
		if nil == body {
			body = []byte("")
		}
		if strings.HasSuffix(request.URL.Path, ".wasm") {
			writer.Header().Set("content-type", "application/wasm")
		}
		if strings.HasPrefix(request.URL.Path, "/web") {
			http.Redirect(writer, request, "/", http.StatusFound)
			return
		}
		web.ServeHTTP(writer, request)
	})
	log.Fatal(http.ListenAndServe(":8080", handlers))
}

func main() {
	Serve()
}
