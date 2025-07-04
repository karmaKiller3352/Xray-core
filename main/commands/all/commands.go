package all

import (
	"github.com/karmaKiller3352/Xray-core/main/commands/all/api"
	"github.com/karmaKiller3352/Xray-core/main/commands/all/convert"
	"github.com/karmaKiller3352/Xray-core/main/commands/all/tls"
	"github.com/karmaKiller3352/Xray-core/main/commands/base"
)

func init() {
	base.RootCommand.Commands = append(
		base.RootCommand.Commands,
		api.CmdAPI,
		convert.CmdConvert,
		tls.CmdTLS,
		cmdUUID,
		cmdX25519,
		cmdWG,
	)
}
