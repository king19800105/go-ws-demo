package hardware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"net"
	"time"
)

// 超时时间
const pongWait = 62 * time.Second

const (
	msgRead  = "读取数据时发生错误"
	msgWrite = "写入数据时发生错误"
)

var errChanClosed = errors.New("管道被关闭")

type EventHandler func(*Event)

// 事件对象
// {"event":"login", "timestamp": 1591254792, "data":"{}"}
type Event struct {
	Event     string      `json:"event"`
	Timestamp int         `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// 统一响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// 通用数据对象
type DataHandler struct {
	DataCh chan []byte
	Events map[string]EventHandler
}

// 应用对象
type App struct {
	Name string
}

func CreateApp(name string) *App {
	return &App{Name: name}
}

// 从管道中读取信息
func ReadMessage(conn *websocket.Conn, evts map[string]EventHandler) (err error) {
	defer func() {
		fmt.Println("defer ReadMessage ip:", conn.RemoteAddr().String())
	}()
	for {
		// 设置超时时间
		if err = conn.SetReadDeadline(time.Now().Add(pongWait)); nil != err {
			err = errors.Wrap(err, msgRead)
			break
		}

		// 读取信息
		_, msg, readErr := conn.ReadMessage()
		// 超时错误
		if netErr, ok := readErr.(net.Error); ok {
			if netErr.Timeout() {
				err = errors.Wrap(readErr, msgRead)
				break
			}
		}

		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
			break
		}

		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			break
		}

		// 解析
		evt, formatErr := NewEventForRaw(msg)
		if nil != formatErr {
			err = errors.Wrap(formatErr, msgRead)
			break
		}

		// 触发事件回调
		if action, ok := evts[evt.Event]; ok {
			action(evt)
		}
	}

	return
}

// 写数据
func WriteMessage(ctx context.Context, conn *websocket.Conn, ch chan []byte) (err error) {
	defer func() {
		fmt.Println("defer WriteMessage ip:", conn.RemoteAddr().String())
	}()
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				err = conn.WriteMessage(websocket.CloseMessage, make([]byte, 0))
				if nil == err {
					err = errChanClosed
				}

				return errors.Wrap(err, msgWrite)
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if nil != err {
				return errors.Wrap(err, msgWrite)
			}

			if _, err = w.Write(msg); nil != err {
				return errors.Wrap(err, msgWrite)
			}

			if err := w.Close(); nil != err {
				return err
			}
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), msgWrite)
		}
	}
}

// 事件数据转换为Event对象类型
func NewEventForRaw(raw []byte) (evt *Event, err error) {
	evt = new(Event)
	err = json.Unmarshal(raw, evt)
	return
}

// 响应数据转换为byte类型
func ToRaw(evt *Response) []byte {
	resp, _ := json.Marshal(evt)
	return resp
}
