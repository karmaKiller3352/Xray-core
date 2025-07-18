package log_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/karmaKiller3352/Xray-core/common/log"
	"github.com/karmaKiller3352/Xray-core/common/net"
)

type testLogger struct {
	value string
}

func (l *testLogger) Handle(msg log.Message) {
	l.value = msg.String()
}

func TestLogRecord(t *testing.T) {
	var logger testLogger
	log.RegisterHandler(&logger)

	ip := "8.8.8.8"
	log.Record(&log.GeneralMessage{
		Severity: log.Severity_Error,
		Content:  net.ParseAddress(ip),
	})

	if diff := cmp.Diff("[Error] "+ip, logger.value); diff != "" {
		t.Error(diff)
	}
}
