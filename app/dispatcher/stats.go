package dispatcher

import (
	"github.com/karmaKiller3352/Xray-core/common"
	"github.com/karmaKiller3352/Xray-core/common/buf"
	"github.com/karmaKiller3352/Xray-core/features/stats"
)

type SizeStatWriter struct {
	Counter stats.Counter
	Writer  buf.Writer
}

func (w *SizeStatWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	w.Counter.Add(int64(mb.Len()))
	return w.Writer.WriteMultiBuffer(mb)
}

func (w *SizeStatWriter) Close() error {
	return common.Close(w.Writer)
}

func (w *SizeStatWriter) Interrupt() {
	common.Interrupt(w.Writer)
}
