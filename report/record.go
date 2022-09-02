package report

import (
	"sync"

	"github.com/pterm/pterm"
)

var m = &sync.Mutex{}

var records = make(map[string]Entry)

func Record(e Entry) {
	defer m.Unlock()
	m.Lock()
	records[e.Identifier] = e
}

// RecordBatch accepts a BatchEntry that contains a slice of identifiers, loops through them and converts each identifier to
// a standard Entry. This is useful for supporting batch delete workflows in cloud-nuke (such as cloudwatch_dashboards)
func RecordBatch(e BatchEntry) {
	for _, identifier := range e.Identifiers {
		entry := Entry{
			Identifier:   identifier,
			ResourceType: e.ResourceType,
			Error:        e.Error,
		}
		Record(entry)
	}
}

func Print() {
	renderSection("Nuking complete:")
	data := make([][]string, len(records))
	entriesToDisplay := []Entry{}
	for _, entry := range records {
		entriesToDisplay = append(entriesToDisplay, entry)
	}

	for idx, entry := range entriesToDisplay {
		var errSymbol string
		if entry.Error != nil {
			errSymbol = "  ❌  "
		} else {
			errSymbol = "  ✅  "
		}
		data[idx] = []string{entry.Identifier, entry.ResourceType, errSymbol}
	}

	renderTableWithHeader([]string{"Identifier", "Resource Type", "Deleted Successfully"}, data)
}

func renderSection(sectionTitle string) {
	pterm.DefaultSection.Style = pterm.NewStyle(pterm.FgLightCyan)
	pterm.DefaultSection.WithLevel(0).Println(sectionTitle)
}

func renderTableWithHeader(headers []string, data [][]string) {
	tableData := pterm.TableData{
		headers,
	}
	for idx := range data {
		tableData = append(tableData, data[idx])
	}
	pterm.DefaultTable.
		WithHasHeader().
		WithBoxed(true).
		WithRowSeparator("-").
		WithData(tableData).
		Render()
}

// Custom types

type Entry struct {
	Identifier   string
	ResourceType string
	Error        error
}

type BatchEntry struct {
	Identifiers  []string
	ResourceType string
	Error        error
}
