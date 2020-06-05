package factory

import (
	"context"
	"github.com/go-crew/group/async"
	"github.com/gorilla/websocket"
	"github.com/king19800105/go-ws-demo/pkg/hardware/instrument"
)

type (
	// 硬件业务抽象
	Hardware interface {
		Process(ctx context.Context, gp *async.Async, conn *websocket.Conn)
		BindEventListener()
	}

	// todo... 其他业务抽象
)

func HardwareFactory(hw string) Hardware {
	switch hw {
	case "instrument":
		return instrument.NewInstrument()
	default:
		return nil
	}
}
