package internal

import (
	"encoding/json"
	"github.com/jedib0t/go-pretty/table"
	"io"
)

type WalShowOutputType int

const (
	TableOutput WalShowOutputType = iota + 1
	JsonOutput
)

// WalShowOutputWriter writes the output of wal-show command execution result
type WalShowOutputWriter interface {
	Write(timelineInfos []*TimelineInfo) error
}

// WalShowJsonOutputWriter writes the detailed JSON output
type WalShowJsonOutputWriter struct {
	output io.Writer
}

func (writer *WalShowJsonOutputWriter) Write(timelineInfos []*TimelineInfo) error {
	bytes, err := json.Marshal(timelineInfos)
	if err != nil {
		return err
	}
	_, err = writer.output.Write(bytes)
	return err
}

// WalShowTableOutputWriter writes the output in compact pretty table
type WalShowTableOutputWriter struct {
	output io.Writer
	includeBackups bool
}

func (writer *WalShowTableOutputWriter) Write(timelineInfos []*TimelineInfo) error {
	tableWriter := table.NewWriter()
	tableWriter.SetOutputMirror(writer.output)
	defer tableWriter.Render()

	header := table.Row{"TLI", "Parent TLI", "Switchpoint LSN", "Start segment",
		"End segment", "Segment range", "Segments count", "Status"}
	if writer.includeBackups {
		header = append(header, "Backups count")
	}
	tableWriter.AppendHeader(header)

	for _, tl := range timelineInfos {
		row := table.Row{tl.Id, tl.ParentId, tl.SwitchPointLsn, tl.StartSegment,
			tl.EndSegment, tl.SegmentRangeSize, tl.SegmentsCount, tl.Status}
		if writer.includeBackups {
			row = append(row, len(tl.Backups))
		}
		tableWriter.AppendRow(row)
	}

	return nil
}

func NewWalShowOutputWriter(outputType WalShowOutputType, output io.Writer, includeBackups bool) WalShowOutputWriter {
	switch outputType {
	case TableOutput:
		return &WalShowTableOutputWriter{output: output, includeBackups: includeBackups}
	case JsonOutput:
		return &WalShowJsonOutputWriter{output: output}
	default:
		return &WalShowTableOutputWriter{output: output, includeBackups: includeBackups}
	}
}
