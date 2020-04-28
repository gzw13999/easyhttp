package easyhttp

import (
	"fmt"
	"io"

	// "net"
	"net/http"
	"path"

	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"
)

type Easyhttp struct {
	server *http.Server
	Routes map[string]func(http.ResponseWriter, *http.Request)
}

func (ehttp *Easyhttp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", " Server ")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	visiturl := strings.ToLower(path.Clean(r.URL.Path))

	//返回类型判断
	w.Header().Set("Content-Type", Ext2ct(path.Ext(visiturl)))

	if h, ok := ehttp.Routes[visiturl]; ok {
		h(w, r)
		return
	} else {
		for k, v := range ehttp.Routes {
			if reg, err := regexp.Compile(k); err == nil {
				if reg.MatchString(visiturl) {
					v(w, r)
					return
				}

			}

		}
		// for k, v := range ehttp.RERoutes {

		// 	if k.pattern.MatchString(visiturl) {
		// 		// route.handler.ServeHTTP(w, r)
		// 		// fmt.Println("匹配", r.URL.Path)
		// 		v(w, r)
		// 		return
		// 	}
		// }
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "   404 Not Found!\n\nUrl:"+visiturl+"\npath:"+r.URL.Path+"\nurl:"+r.URL.String())
	}
}

func New() *Easyhttp {
	ehttp := new(Easyhttp)
	ehttp.Routes = make(map[string]func(http.ResponseWriter, *http.Request))

	return ehttp

}

func (ehttp *Easyhttp) Run(addr string) {

	ehttp.server = &http.Server{
		Addr:           addr,
		Handler:        ehttp,
		ReadTimeout:    10 * time.Second, // 读超时设置  读取clent超时 不可更改，否则客户会提示 io timeout
		WriteTimeout:   10 * time.Second, // 写超时设置  给client写数据超时
		MaxHeaderBytes: 1 << 20,
	}

	// 一个通知退出的chan
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		// 接收退出信号
		quit := make(chan os.Signal)
		<-quit
		if err := ehttp.server.Close(); err != nil {
			fmt.Println("Close server:", err)
		}
	}()

	fmt.Println("开始运行", addr)
	//监听
	if err := ehttp.server.ListenAndServe(); err != nil {
		// 正常退出
		if err == http.ErrServerClosed {
			fmt.Println("Server closed under request")
		} else {
			fmt.Println("Server closed unexpected", err)
		}
	}
}

func (http *Easyhttp) Shutdown() {
	err := http.server.Shutdown(nil)
	if err != nil {
		fmt.Println([]byte("shutdown the server err"))
	}
}

func Ext2ct(ext string) string {

	var ct string
	switch ext {
	case ".jpg":
		ct = "image/jpeg"
		break
	case ".gif":
		ct = "image/gif"
		break
	case ".png":
		ct = "image/png"
		break
	case ".bmp":
		ct = "image/bmp"
		break
	case ".html", ".htm":
		ct = "text/html; charset=utf-8"
	case ".css":
		ct = "text/css"
	case ".js":
		ct = "application/javascript"
	case ".xml":
		ct = "text/xml; charset=utf-8"
	case ".json":
		ct = "text/json; charset=utf-8"
	case ".txt":
		ct = "text/plain; charset=utf-8"
	default:
		ct = "text/plain; charset=utf-8"
		break
	}
	return ct
}
