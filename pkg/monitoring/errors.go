package monitoring

import (
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func init() {
	Writer = &stdPaymentErrors{
		counter: stdPaymentErrorsCounters,
	}
}

var (
	nameSpace       = "goals_scheduler"
	subsystemErrors = "errors"

	labelLayer    = "layer"
	labelMessage  = "message"
	LabelRedis    = "redis"
	LabelPostgres = "postgres"
	LabelBot      = "bot"
)

var (
	stdPaymentErrorsCounters = kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: nameSpace,
		Subsystem: subsystemErrors,
		Name:      "internal_errors",
		Help:      "Errors that need to be alerted",
	}, []string{labelLayer, labelMessage})

	Writer *stdPaymentErrors
)

type stdPaymentErrors struct {
	counter metrics.Counter
}

func (std *stdPaymentErrors) Incr(layer, msg string) {
	cnt := std.counter.With(
		labelLayer, layer,
		labelMessage, msg,
	)

	cnt.Add(1)
}
