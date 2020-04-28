package main

import (
	"net/http"

	"github.com/gzw13999/easyhttp"
)

func main() {

	easyhttp := easyhttp.New()
	easyhttp.Routes["/"] = func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("欢迎访问."))
	}
	easyhttp.RERoutes[`/static/(.*)`] = func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("欢迎访问正则路由."))
	}
	easyhttp.Run(":80")
}
