package internal

import (
	"bufio"
	"fmt"
	"sort"

	"encoding/csv"

	"github.com/sha1n/benchy/api"
)

// textReportWriter a simple human readable test report writer
type csvReportWriter struct {
	writer *csv.Writer
}

// NewCsvReportWriter returns a CSV report write handler.
func NewCsvReportWriter(writer *bufio.Writer) api.WriteReportFn {
	w := &csvReportWriter{
		writer: csv.NewWriter(writer),
	}

	return w.Write
}

func (rw *csvReportWriter) Write(summary api.Summary, config *api.BenchmarkSpec) (err error) {
	rw.writer.Write([]string{"Timestamp", "Scenario", "Min", "Max", "Mean", "Median", "Percentile 90", "StdDev", "Errors"})

	timeStr := summary.Time().Format("2006-01-02T15:04:05Z07:00")
	sortedIds := getSortedScenarioIds(summary)

	for i := range sortedIds {
		id := sortedIds[i]
		stats := summary.Get(id)

		rw.writer.Write([]string{
			timeStr,
			id,
			rw.writeStatNano(stats.Min),
			rw.writeStatNano(stats.Max),
			rw.writeStatNano(stats.Mean),
			rw.writeStatNano(stats.Median),
			rw.writeStatNano(func() (float64, error) { return stats.Percentile(90) }),
			rw.writeStatNano(stats.StdDev),
			rw.writeErrorRateStat(stats.ErrorRate),
		})

		rw.writer.Flush()
	}

	return nil
}

func getSortedScenarioIds(summary api.Summary) []api.ID {
	statsMap := summary.All()
	sortedIds := make([]api.ID, 0, len(statsMap))
	for k := range statsMap {
		sortedIds = append(sortedIds, k)
	}
	sort.Strings(sortedIds)

	return sortedIds
}

func (rw *csvReportWriter) writeStatNano(f func() (float64, error)) string {
	value, err := f()
	if err == nil {
		return fmt.Sprintf("%.3f", value)
	}

	return "ERROR"
}

func (rw *csvReportWriter) writeErrorRateStat(f func() float64) string {
	value := f()
	errorRate := int(value * 100)

	return fmt.Sprintf("%d", errorRate)
}