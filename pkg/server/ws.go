package server

import (
	"context"
	"fmt"
	"github.com/go-crew/group/async"
	"github.com/gorilla/websocket"
	"github.com/king19800105/go-ws-demo/pkg/hardware"
	"github.com/king19800105/go-ws-demo/pkg/hardware/factory"
	"net/http"
)

const (
	buffer = 1024
)

const (
	msgPath      = "访问的路由错误"
	msgWebsocket = "websocket error : %v"
)

type Server struct {
	addr     string
	uri      string
	upgrade  *websocket.Upgrader
	hardware factory.Hardware
}

// 创建websocket服务对象
func NewServer(addr string, uri string) *Server {
	var ws *Server
	{
		ws = new(Server)
		ws.addr = addr
		ws.uri = uri
		ws.upgrade = &websocket.Upgrader{
			ReadBufferSize:  buffer,
			WriteBufferSize: buffer,
			CheckOrigin: func(r *http.Request) bool {
				if r.URL.Path != uri {
					fmt.Println(msgPath)
					return false
				}

				return true
			},
		}
	}

	return ws
}

// 启动websocket服务
// 每个连接都会触发一次全新的http.HandleFunc处理
func (ws *Server) StartBy(app *hardware.App) (err error) {
	http.HandleFunc(ws.uri, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != ws.uri {
			httpCode := http.StatusInternalServerError
			status := http.StatusText(httpCode)
			http.Error(w, status, httpCode)
			fmt.Println(msgPath)
			return
		}

		conn, err := ws.upgrade.Upgrade(w, r, nil)
		if nil != err {
			fmt.Printf(msgWebsocket, err)
			return
		}

		var (
			gp  = async.NewGroup()
			ctx = context.Background()
		)

		var hw = factory.HardwareFactory(app.Name)
		if nil == hw {
			conn.Close()
			return
		}

		hw.BindEvents()
		go hw.Process(ctx, gp, conn)
	})

	return http.ListenAndServe(ws.addr, nil)
}
