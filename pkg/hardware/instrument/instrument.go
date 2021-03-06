package instrument

import (
	"context"
	"fmt"
	"github.com/go-crew/group/async"
	"github.com/gorilla/websocket"
	"github.com/king19800105/go-ws-demo/pkg/hardware"
)

// 设备仪器
type Instrument struct {
	hardware.DataHandler
}

// 创建仪器对象
func NewInstrument() (s *Instrument) {
	s = new(Instrument)
	s.Events = make(map[string]hardware.EventHandler)
	s.DataCh = make(chan []byte, 1)
	return
}

// 通道对接后处理
func (s *Instrument) Process(ctx context.Context, gp *async.Async, conn *websocket.Conn) {
	defer conn.Close()
	// 读取数据
	gp.Add(func(ctx context.Context, params ...interface{}) error {
		return hardware.ReadMessage(conn, s.Events)
	}, func(cancel context.CancelFunc, err error) {
		fmt.Println("Exec read cancel, error is:", err)
		cancel()
	})

	// 写数据
	gp.Add(func(ctx context.Context, params ...interface{}) error {
		return hardware.WriteMessage(ctx, conn, s.DataCh)
	}, func(cancel context.CancelFunc, err error) {
		fmt.Println("Exec write cancel, error is:", err)
		cancel()
	})

	if err := gp.Run(ctx, -1); nil != err {
		fmt.Println("Run:", err)
	}

	fmt.Println("Run over")
}

// 事件绑定操作
func (s *Instrument) BindEventListener() {
	// 登入事件处理
	s.On("login", func(event *hardware.Event) {
		s.login(event)
	})
}

// 注册事件
func (s *Instrument) On(evt string, action hardware.EventHandler) {
	s.Events[evt] = action
}

// 登入事件的处理逻辑
func (s *Instrument) login(event *hardware.Event) {
	// 其他逻辑处理...
	s.DataCh <- hardware.ToRaw(&hardware.Response{
		Code:    0,
		Message: "success",
		Data:    "{}",
	})
}

// todo... 其他事件定义
