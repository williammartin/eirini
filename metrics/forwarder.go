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
		loggregator.WithGaugeValue("cpu", convertCPU(msg.CPU), "percentage"),
		loggregator.WithGaugeValue("memory", convertMemory(msg.Memory), "bytes"),
		loggregator.WithGaugeValue("disk", msg.Disk, "bytes"),
		loggregator.WithGaugeValue("memory_quota", msg.MemoryQuota, "bytes"),
		loggregator.WithGaugeValue("disk_quota", msg.DiskQuota, "bytes"),
	)
}

// The kubernetes metrics api returns the cpu usage in "millicores"
// which is described here: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#meaning-of-cpu
// This means that one millicore - 1m is equal to 0.1% of a CPU
func convertCPU(cpuUsage float64) float64 {
	return cpuUsage / 1000
}

// We assume that the kubernetes metrics api returns memory in Ki.
func convertMemory(memoryUsage float64) float64 {
	return memoryUsage * 1024
}
