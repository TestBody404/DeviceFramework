package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	pcieBwLabel = append(cardLabel, pcieBwType)
)

func pcieBwLabelVal(cardLabels []string, pcieBwType string) []string {
	return append(cardLabels, pcieBwType)
}

func metricWithPcieBw(labelsVal []string, metrics *prometheus.Desc, val float64, valType string) prometheus.Metric {
	return prometheus.MustNewConstMetric(metrics, prometheus.GaugeValue, val, pcieBwLabelVal(labelsVal, valType)...)
}

func updatePcieBwInfo(ch chan<- prometheus.Metric, npu *HuaWeiNPUCard, chip *HuaWeiAIChip, cardLabelsVal []string) {
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescTxPBW, float64(chip.PcieBwInfo.PcieTxPBw.PcieAvgBw), avgPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescTxNpBW, float64(chip.PcieBwInfo.PcieTxNPBw.PcieAvgBw), avgPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescTxCplBW, float64(chip.PcieBwInfo.PcieTxCPLBw.PcieAvgBw), avgPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescRxPBW, float64(chip.PcieBwInfo.PcieRxPBw.PcieAvgBw), avgPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescRxNpBW, float64(chip.PcieBwInfo.PcieRxNPBw.PcieAvgBw), avgPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescRxCplBW, float64(chip.PcieBwInfo.PcieRxCPLBw.PcieAvgBw), avgPcieBw))

	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescTxPBW, float64(chip.PcieBwInfo.PcieTxPBw.PcieMinBw), minPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescTxNpBW, float64(chip.PcieBwInfo.PcieTxNPBw.PcieMinBw), minPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescTxCplBW, float64(chip.PcieBwInfo.PcieTxCPLBw.PcieMinBw), minPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescRxPBW, float64(chip.PcieBwInfo.PcieRxPBw.PcieMinBw), minPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescRxNpBW, float64(chip.PcieBwInfo.PcieRxNPBw.PcieMinBw), minPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescRxCplBW, float64(chip.PcieBwInfo.PcieRxCPLBw.PcieMinBw), minPcieBw))

	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescTxPBW, float64(chip.PcieBwInfo.PcieTxPBw.PcieMaxBw), maxPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescTxNpBW, float64(chip.PcieBwInfo.PcieTxNPBw.PcieMaxBw), maxPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescTxCplBW, float64(chip.PcieBwInfo.PcieTxCPLBw.PcieMaxBw), maxPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescRxPBW, float64(chip.PcieBwInfo.PcieRxPBw.PcieMaxBw), maxPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescRxNpBW, float64(chip.PcieBwInfo.PcieRxNPBw.PcieMaxBw), maxPcieBw))
	ch <- prometheus.NewMetricWithTimestamp(npu.Timestamp,
		metricWithPcieBw(cardLabelsVal, npuChipInfoDescRxCplBW, float64(chip.PcieBwInfo.PcieRxCPLBw.PcieMaxBw), maxPcieBw))
}
