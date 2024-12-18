package metrics

import (
	"fmt"
	"log"
	"net/url"
)

// The MetricSink interface is used to transmit metrics information
// to an external system
type MetricSink interface {
	// A Gauge should retain the last value it is set to
	SetGauge(key []string, val float32)
	SetGaugeWithLabels(key []string, val float32, labels []Label)

	// Should emit a Key/Value pair for each call
	EmitKey(key []string, val float32)

	// Counters should accumulate values
	IncrCounter(key []string, val float32)
	IncrCounterWithLabels(key []string, val float32, labels []Label)

	// Samples are for timing information, where quantiles are used
	AddSample(key []string, val float32)
	AddSampleWithLabels(key []string, val float32, labels []Label)
}

// PrecisionGaugeMetricSink interfae is used to support 64 bit precisions for Sinks, if needed.
type PrecisionGaugeMetricSink interface {
	SetPrecisionGauge(key []string, val float64)
	SetPrecisionGaugeWithLabels(key []string, val float64, labels []Label)
}

type ShutdownSink interface {
	MetricSink

	// Shutdown the metric sink, flush metrics to storage, and cleanup resources.
	// Called immediately prior to application exit. Implementations must block
	// until metrics are flushed to storage.
	Shutdown()
}

// BlackholeSink is used to just blackhole messages
type BlackholeSink struct{}

func (b *BlackholeSink) SetGauge(key []string, val float32) {
	b.noopsWarnings()
}

func (b *BlackholeSink) SetGaugeWithLabels(key []string, val float32, labels []Label) {
	b.noopsWarnings()
}

func (b *BlackholeSink) SetPrecisionGauge(key []string, val float64) {
	b.noopsWarnings()
}

func (b *BlackholeSink) SetPrecisionGaugeWithLabels(key []string, val float64, labels []Label) {
	b.noopsWarnings()
}

func (b *BlackholeSink) EmitKey(key []string, val float32) {
	b.noopsWarnings()
}

func (b *BlackholeSink) IncrCounter(key []string, val float32) {
	b.noopsWarnings()
}

func (b *BlackholeSink) IncrCounterWithLabels(key []string, val float32, labels []Label) {
	b.noopsWarnings()
}

func (b *BlackholeSink) AddSample(key []string, val float32) {
	b.noopsWarnings()
}

func (b *BlackholeSink) AddSampleWithLabels(key []string, val float32, labels []Label) {
	b.noopsWarnings()
}

func (b *BlackholeSink) noopsWarnings() {
	log.Print("[WARN] No default metrics sink was set. Using noops metrics sink as default. Set the default metrics sink to do all functions\n")
}

// FanoutSink is used to sink to fanout values to multiple sinks
type FanoutSink []MetricSink

func (fh FanoutSink) SetGauge(key []string, val float32) {
	fh.SetGaugeWithLabels(key, val, nil)
}

func (fh FanoutSink) SetGaugeWithLabels(key []string, val float32, labels []Label) {
	for _, s := range fh {
		s.SetGaugeWithLabels(key, val, labels)
	}
}

func (fh FanoutSink) SetPrecisionGauge(key []string, val float64) {
	fh.SetPrecisionGaugeWithLabels(key, val, nil)
}

func (fh FanoutSink) SetPrecisionGaugeWithLabels(key []string, val float64, labels []Label) {
	for _, s := range fh {
		// The Sink needs to implement PrecisionGaugeMetricSink, in case it doesn't, the metric value won't be set and ingored instead
		if s64, ok := s.(PrecisionGaugeMetricSink); ok {
			s64.SetPrecisionGaugeWithLabels(key, val, labels)
		}
	}
}

func (fh FanoutSink) EmitKey(key []string, val float32) {
	for _, s := range fh {
		s.EmitKey(key, val)
	}
}

func (fh FanoutSink) IncrCounter(key []string, val float32) {
	fh.IncrCounterWithLabels(key, val, nil)
}

func (fh FanoutSink) IncrCounterWithLabels(key []string, val float32, labels []Label) {
	for _, s := range fh {
		s.IncrCounterWithLabels(key, val, labels)
	}
}

func (fh FanoutSink) AddSample(key []string, val float32) {
	fh.AddSampleWithLabels(key, val, nil)
}

func (fh FanoutSink) AddSampleWithLabels(key []string, val float32, labels []Label) {
	for _, s := range fh {
		s.AddSampleWithLabels(key, val, labels)
	}
}

func (fh FanoutSink) Shutdown() {
	for _, s := range fh {
		if ss, ok := s.(ShutdownSink); ok {
			ss.Shutdown()
		}
	}
}

// sinkURLFactoryFunc is an generic interface around the *SinkFromURL() function provided
// by each sink type
type sinkURLFactoryFunc func(*url.URL) (MetricSink, error)

// sinkRegistry supports the generic NewMetricSink function by mapping URL
// schemes to metric sink factory functions
var sinkRegistry = map[string]sinkURLFactoryFunc{
	"statsd":   NewStatsdSinkFromURL,
	"statsite": NewStatsiteSinkFromURL,
	"inmem":    NewInmemSinkFromURL,
}

// NewMetricSinkFromURL allows a generic URL input to configure any of the
// supported sinks. The scheme of the URL identifies the type of the sink, the
// and query parameters are used to set options.
//
// "statsd://" - Initializes a StatsdSink. The host and port are passed through
// as the "addr" of the sink
//
// "statsite://" - Initializes a StatsiteSink. The host and port become the
// "addr" of the sink
//
// "inmem://" - Initializes an InmemSink. The host and port are ignored. The
// "interval" and "duration" query parameters must be specified with valid
// durations, see NewInmemSink for details.
func NewMetricSinkFromURL(urlStr string) (MetricSink, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	sinkURLFactoryFunc := sinkRegistry[u.Scheme]
	if sinkURLFactoryFunc == nil {
		return nil, fmt.Errorf(
			"cannot create metric sink, unrecognized sink name: %q", u.Scheme)
	}

	return sinkURLFactoryFunc(u)
}
