package networktraffic

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	"sigs.k8s.io/scheduler-plugins/apis/config"
)

const PluginName = "NetworkTraffic"

var _ = framework.ScorePlugin(&NetworkTraffic{})

/*
NetworkTraffic is a scheduling plugin that favours nodes based
on their network traffic amount. Node with less traffic are favoured.
*/
type NetworkTraffic struct {
	handle     framework.Handle
	prometheus *PrometheusHandle
}

func (n *NetworkTraffic) Name() string {
	return PluginName
}

func New(_ context.Context, obj runtime.Object, h framework.Handle) (framework.Plugin, error) {
	args, ok := obj.(*config.NetworkTrafficArgs)
	if !ok {
		return nil, fmt.Errorf("want args to be of type NetworkTrafficArgs, got %T", obj)
	}
	return &NetworkTraffic{
		handle: h,
		prometheus: NewPrometehus(
			args.Address,
			args.NetworkInterface,
			time.Minute*time.Duration(args.TimeRangeInMinutes),
		),
	}, nil
}

/*
Score will get the bandwidth measure for a given node
and convert that value to int64
*/
func (n *NetworkTraffic) Score(
	ctx context.Context,
	state *framework.CycleState,
	p *v1.Pod,
	nodeName string,
) (int64, *framework.Status) {
	nodeBandwidth, err := n.prometheus.GetNodeBandwidthMeasure(nodeName)
	if err != nil {
		return 0, framework.NewStatus(
			framework.Error,
			fmt.Sprintf("error getting node bandwidth measure: %s", err),
		)
	}

	klog.Infof("[NetworkTraffic] node: %s, bandwidth: %s", nodeName, nodeBandwidth)

	return int64(nodeBandwidth.Value), nil
}

func (n *NetworkTraffic) ScoreExtensions() framework.ScoreExtensions {
	return n
}

/*
We will have returned the total bytes received by each node in a determined period of time.
However, the scheduler framework expects a value from 0 to 100, thus, need to normalize the values to fulfill this requirement.

The NormalizeScore basically will take the highest value returned by prometheus and use it as the highest possible value, corresponding to the framework.MaxNodeScore (100). The other values will be calculated relative to the highest score using the rule of three.

Finally, we will have a list where the nodes with more network traffic have a greater score in the range of [0,100]. If we use it like it is, we would favor nodes that have higher traffic, so, we need to reverse the values. For that, we will simply replace the node score with the result of the rule of three, subtracted by the max score.
*/
func (n *NetworkTraffic) NormalizeScore(
	ctx context.Context,
	state *framework.CycleState,
	pod *v1.Pod,
	scores framework.NodeScoreList,
) *framework.Status {
	var higherScore int64
	for _, node := range scores {
		if higherScore < node.Score {
			higherScore = node.Score
		}
	}

	for i, node := range scores {
		scores[i].Score = framework.MaxNodeScore - (node.Score - framework.MaxNodeScore/higherScore)
	}

	klog.Infof("[NetworkTraffic] Nodes final score: %v", scores)
	return nil
}
