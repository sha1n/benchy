package report

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/sha1n/benchy/api"
	"github.com/sha1n/benchy/test"
	"github.com/stretchr/testify/assert"
)

var randomLabels = test.RandomLabels()

// GetRawDataHandler a provider for a RawDataHandler instance
type GetRawDataHandler = func(*bufio.Writer, *api.ReportContext) RawDataHandler

// ParseRecords parses report records from a reader
type ParseRecords = func(io.Reader) ([][]string, error)

func assertRawTraceRecord(t *testing.T, trace api.Trace, actualRecord []string) {
	expectedLabels := strings.Join(randomLabels, ",")

	_, err := time.Parse(TabularReportDateFormat, actualRecord[0])
	assert.NoError(t, err)
	assert.Equal(t, trace.ID(), actualRecord[1])
	assert.Equal(t, expectedLabels, actualRecord[2])
	assert.Equal(t, fmt.Sprint(trace.Elapsed().Nanoseconds()), actualRecord[3])
	assert.Equal(t, fmt.Sprint(trace.Error() != nil), actualRecord[4])
}

func writeRawReport(t *testing.T, getHandler GetRawDataHandler, parseRecords ParseRecords, includeHeaders bool, traces ...api.Trace) [][]string {
	buf := new(bytes.Buffer)
	ctx := &api.ReportContext{
		Labels:         randomLabels,
		IncludeHeaders: includeHeaders,
	}

	handler := getHandler(bufio.NewWriter(buf), ctx)

	for _, trace := range traces {
		assert.NoError(t, handler.Handle(trace))
	}

	allRecords, err := parseRecords(buf)

	assert.NoError(t, err)

	return allRecords
}
