package main

import (
	"flag"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"log"
	"net/http"
)

func main() {
	port := flag.String("p", "8100", "port to serve on")
	flag.Parse()

	http.Handle("/avatar/", http.StripPrefix("/avatar/", http.FileServer(http.Dir(configs.PathToUploadAvatar))))
	http.Handle("/product/", http.StripPrefix("/product/", http.FileServer(http.Dir(configs.PathToUploadProductImg))))

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
