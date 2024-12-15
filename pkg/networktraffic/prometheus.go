package networktraffic

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"k8s.io/klog/v2"
)

/*
This file uses prometheus to get a metric value for a given node of the sum of received bytes
in a given timeRange.
*/

/*
Sum of the received bytes in a time range per node in a specific network interface.The kubernetes_node filter will query the metrics for the node provided, as described by the query below.

The device filter will query the metrics on the provided network interface, and the last value between [%s] defines the time range taken into account.

`sum_over_time` will sum all the values in the provided time range
*/
const (
	nodeMeasureQueryTemplate = "sum_over_time(node_network_receive_bytes_total{kubernetes_node=\"%s\",device=\"%s\"}[%s])"
)

type PrometheusHandle struct {
	networkInterface string
	timeRange        time.Duration
	address          string
	api              v1.API
}

func NewPrometehus(address, networkInterface string, timeRange time.Duration) *PrometheusHandle {
	// Create prometheus client
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		klog.Fatalf("[NetworkTraffic] Error creating prometheus client: %v", err)
	}

	return &PrometheusHandle{
		networkInterface: networkInterface,
		timeRange:        timeRange,
		address:          address,
		api:              v1.NewAPI(client),
	}
}

func (p *PrometheusHandle) GetNodeBandwidthMeasure(node string) (*model.Sample, error) {
	// model.Sample returns a sample pair associated with a metric
	query := getNodeBandwidthQuery(node, p.networkInterface, p.timeRange)
	result, err := p.query(query)
	if err != nil {
		return nil, fmt.Errorf("[NetworkTraffic] Error: %v", err)
	}

	nodeMeasure := result.(model.Vector)
	if len(nodeMeasure) != 1 {
		return nil, fmt.Errorf(
			"[NetworkTraffic] Error: Invalid response, expected 1 value, got %d",
			len(nodeMeasure),
		)
	}

	return nodeMeasure[0], nil
}

func getNodeBandwidthQuery(node, networkInterface string, timeRange time.Duration) string {
	return fmt.Sprintf(nodeMeasureQueryTemplate, node, networkInterface, timeRange)
}

func (p *PrometheusHandle) query(query string) (model.Value, error) {
	results, warnings, err := p.api.Query(context.Background(), query, time.Now())
	if len(warnings) > 0 {
		klog.Warningf("[NetworkTraffic] Warnings: %v\n", warnings)
	}

	return results, err
}
