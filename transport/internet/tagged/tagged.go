package tagged

import (
	"context"

	"github.com/karmaKiller3352/Xray-core/common/net"
	"github.com/karmaKiller3352/Xray-core/features/routing"
)

type DialFunc func(ctx context.Context, dispatcher routing.Dispatcher, dest net.Destination, tag string) (net.Conn, error)

var Dialer DialFunc
