package main

import (
	"context"
	"flag"
	"github.com/heejin-github/umee-exporter/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
	"strconv"
	"k8s.io/klog/v2"
)

const (
	defaultAddr = "http://localhost:1317"
	defaultPort = ":26670"
	httpTimeout = 5 * time.Second
)
var (
	addr string
	port string
	valoperAddr string
)

type umeeCollector struct {
	rpcClient *rpc.RPCClient
	validatorMissCounter	*prometheus.Desc
	validatorPrevoteSubmit	*prometheus.Desc
	slashWindow		*prometheus.Desc
}
func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "Umee RPC URI Address")
	flag.StringVar(&port, "port", defaultPort, "Listening port")
	flag.StringVar(&valoperAddr, "valoperAddr", "", "Umee validator operator Address")

	flag.Parse()
	klog.InitFlags(nil)
}

func NewUmeeCollector(rpcAddr string) *umeeCollector {
	return &umeeCollector{
		rpcClient: rpc.NewRPCClient(rpcAddr),
		validatorMissCounter: prometheus.NewDesc(
			"umee_validator_oracle_miss_counter",
			"The number of oracle miss of a  validator",
			[]string{"validator"}, nil),
		validatorPrevoteSubmit: prometheus.NewDesc(
			"umee_validator_aggregate_prevote_submit_block",
			"Last submit block number of oracle prevote",
			[]string{"validator"}, nil),
		slashWindow: prometheus.NewDesc(
			"umee_slash_window_progress",
			"slash window progress in current period",
			[]string{"validator"}, nil),
	}
}

func (c *umeeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.validatorMissCounter
	ch <- c.validatorPrevoteSubmit
	ch <- c.slashWindow
}

func (c *umeeCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	miss, err := c.rpcClient.GetMissCount(ctx, valoperAddr)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.validatorMissCounter, err)
	} else {
		if i, err := strconv.Atoi(*miss); err == nil {
			ch <- prometheus.MustNewConstMetric(c.validatorMissCounter, prometheus.GaugeValue,
				float64(i), valoperAddr)
		}
	}

	submit, err := c.rpcClient.GetPrevoteSubmit(ctx, valoperAddr)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.validatorPrevoteSubmit, err)
	} else {
		if i, err := strconv.Atoi(*submit); err == nil {
			ch <- prometheus.MustNewConstMetric(c.validatorPrevoteSubmit, prometheus.GaugeValue,
				float64(i), valoperAddr)
		}
	}

	window, err := c.rpcClient.GetSlashWindow(ctx)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.slashWindow, err)
	} else {
		if i, err := strconv.Atoi(*window); err == nil {
			ch <- prometheus.MustNewConstMetric(c.slashWindow, prometheus.GaugeValue,
				float64(i), valoperAddr)
		}
	}
}

func main() {
	if valoperAddr == "" {
		klog.Fatal("Please specify -valoperAddr")
	}

	collector := NewUmeeCollector(addr)
	prometheus.MustRegister(collector)
	http.Handle("/metrics", promhttp.Handler())

	klog.Infof("listening on %s", port)
	klog.Fatal(http.ListenAndServe(port, nil))
}
