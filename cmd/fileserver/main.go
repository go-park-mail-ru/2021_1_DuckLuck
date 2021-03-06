package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
)

func main() {
	port := flag.String("p", "8100", "port to serve on")
	flag.Parse()

	http.Handle(configs.UrlToAvatar, http.StripPrefix(configs.UrlToAvatar,
		http.FileServer(http.Dir(configs.PathToUploadAvatar))))

	http.Handle(configs.UrlToProductImg, http.StripPrefix(configs.UrlToProductImg,
		http.FileServer(http.Dir(configs.PathToUploadProductImg))))

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
