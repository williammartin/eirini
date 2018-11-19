package metrics

import (
	"strconv"

	loggregator "code.cloudfoundry.org/go-loggregator"
)

type LoggregatorForwarder struct {
	client *loggregator.IngressClient
}

func NewLoggregatorForwarder(client *loggregator.IngressClient) *LoggregatorForwarder {
	return &LoggregatorForwarder{
		client: client,
	}
}

func (l *LoggregatorForwarder) Forward(msg Message) {
	index, _ := strconv.Atoi(msg.IndexID)
	l.client.EmitGauge(
		loggregator.WithGaugeSourceInfo(msg.AppID, msg.IndexID),
		loggregator.WithGaugeAppInfo(msg.AppID, index),
		loggregator.WithGaugeValue("cpu", msg.CPU, "percentage"),
		loggregator.WithGaugeValue("memory", msg.Memory, "bytes"),
		loggregator.WithGaugeValue("disk", msg.Disk, "bytes"),
		loggregator.WithGaugeValue("memory_quota", msg.MemoryQuota, "bytes"),
		loggregator.WithGaugeValue("disk_quota", msg.DiskQuota, "bytes"),
	)
}
