package protocol // import "github.com/karmaKiller3352/Xray-core/common/protocol"

import (
	"errors"
)

var ErrProtoNeedMoreData = errors.New("protocol matches, but need more data to complete sniffing")
