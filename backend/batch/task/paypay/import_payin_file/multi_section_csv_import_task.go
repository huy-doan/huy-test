package task

import (
	"context"
	"log"
	"strings"
)

// MultiSectionCSVImportTask handles CSV files with summary and detail sections
// The first section (lines before the first blank line) is summary, the second section is detail
// summaryHeaders: header for summary section (first non-empty line)
// detailHeaders: header for detail section (first non-empty line after blank line)
type MultiSectionCSVImportTask struct {
	SummaryInsert func(ctx context.Context, payinFileID int, records []map[string]string) error
	DetailInsert  func(ctx context.Context, payinFileID int, records []map[string]string) error
}

func NewMultiSectionCSVImportTask(
	summaryInsert func(ctx context.Context, payinFileID int, records []map[string]string) error,
	detailInsert func(ctx context.Context, payinFileID int, records []map[string]string) error,
) *MultiSectionCSVImportTask {
	return &MultiSectionCSVImportTask{
		SummaryInsert: summaryInsert,
		DetailInsert:  detailInsert,
	}
}

// Do: parse lines, split into summary and detail, insert to DB
func (t *MultiSectionCSVImportTask) Do(ctx context.Context, payinFileID int, lines []string) error {
	// Only get the second non-empty line as summary record, ignore all other lines for summary
	var summaryHeaders []string
	var summaryRecords []map[string]string
	var detailHeaders []string
	var detailRecords []map[string]string
	var foundHeader bool
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		fields := parseCSVLine(trimmed)
		if !foundHeader {
			summaryHeaders = fields
			foundHeader = true
			continue
		}
		if len(summaryRecords) == 0 {
			record := make(map[string]string)
			for j, v := range fields {
				if j < len(summaryHeaders) {
					record[summaryHeaders[j]] = v
				}
			}
			summaryRecords = append(summaryRecords, record)
			continue
		}
		// All further lines are ignored for summary, but could be detail in fallback logic
	}

	// Fallback logic for detail section (if needed)
	if summaryHeaders != nil && len(summaryRecords) == 0 && len(lines) > 2 {
		parsedHeader := summaryHeaders
		var summaryLine, detailLines []string
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			fields := parseCSVLine(trimmed)
			if i == 0 {
				continue // header
			} else if i == 1 {
				summaryLine = fields
			} else {
				detailLines = append(detailLines, trimmed)
			}
		}
		if summaryLine != nil {
			record := make(map[string]string)
			for j, v := range summaryLine {
				if j < len(parsedHeader) {
					record[parsedHeader[j]] = v
				}
			}
			summaryRecords = append(summaryRecords, record)
		}
		if len(detailLines) > 0 {
			detailHeaders = parsedHeader
			for _, line := range detailLines {
				fields := parseCSVLine(line)
				record := make(map[string]string)
				for j, v := range fields {
					if j < len(detailHeaders) {
						record[detailHeaders[j]] = v
					}
				}
				detailRecords = append(detailRecords, record)
			}
		}
	} else if len(lines) >= 4 {
		parsedHeader := parseCSVLine(strings.TrimSpace(lines[0]))
		detailHeaderIdx := 3
		if len(lines) > detailHeaderIdx {
			detailHeaders = parseCSVLine(strings.TrimSpace(lines[detailHeaderIdx]))
			for i := 1; i < detailHeaderIdx; i++ {
				fields := parseCSVLine(strings.TrimSpace(lines[i]))
				record := make(map[string]string)
				for j, v := range fields {
					if j < len(parsedHeader) {
						record[parsedHeader[j]] = v
					}
				}
				summaryRecords = append(summaryRecords, record)
			}
			for i := detailHeaderIdx + 1; i < len(lines); i++ {
				fields := parseCSVLine(strings.TrimSpace(lines[i]))
				if len(fields) == 0 || (len(fields) == 1 && fields[0] == "") {
					continue
				}
				record := make(map[string]string)
				for j, v := range fields {
					if j < len(detailHeaders) {
						record[detailHeaders[j]] = v
					}
				}
				detailRecords = append(detailRecords, record)
			}
		}
	}

	// --- Normalize and validate summary section ---
	if summaryHeaders != nil && len(summaryRecords) > 0 {
		mappingTask := NewCSVHeaderMappingTask()
		normSummaryHeaders, normSummaryRecords := mappingTask.Do(summaryHeaders, recordsToSlice(summaryRecords, summaryHeaders))
		validateTask := NewValidateCSVFieldsTask()
		isValid, _ := validateTask.CheckHeaders(normSummaryHeaders, 0)
		if !isValid {
			log.Printf("[MultiSectionCSVImportTask] Summary section missing required headers. Records will not be imported.")
		} else {
			log.Printf("[MultiSectionCSVImportTask] First normalized summary record: %+v", normSummaryRecords[0])
			err := t.SummaryInsert(ctx, payinFileID, normSummaryRecords)
			if err != nil {
				return err
			}
		}
	}
	// --- Normalize and validate detail section ---
	if detailHeaders != nil && len(detailRecords) > 0 {
		log.Printf("[MultiSectionCSVImportTask] Validating detail section using RequiredPayinDetailHeaders: %v", RequiredPayinDetailHeaders)
		mappingTask := NewCSVHeaderMappingTask()
		normDetailHeaders, normDetailRecords := mappingTask.Do(detailHeaders, recordsToSlice(detailRecords, detailHeaders))
		validateTask := NewValidateCSVFieldsTask()
		isValid, required := validateTask.CheckHeaders(normDetailHeaders, 1)
		log.Printf("[MultiSectionCSVImportTask] Normalized detail headers: %v", normDetailHeaders)
		log.Printf("[MultiSectionCSVImportTask] Required detail headers: %v", required)
		if !isValid {
			log.Printf("[MultiSectionCSVImportTask] Detail section missing required headers. Records will not be imported.")
		} else {
			log.Printf("[MultiSectionCSVImportTask] Preparing to insert %d detail records for payinFileID=%d", len(normDetailRecords), payinFileID)
			log.Printf("[MultiSectionCSVImportTask] First normalized detail record: %+v", normDetailRecords[0])
			err := t.DetailInsert(ctx, payinFileID, normDetailRecords)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// parseCSVLine is a simple CSV parser for a single line (no quote support)
func parseCSVLine(line string) []string {
	return splitCSV(line)
}

// splitCSV splits a CSV line by comma (no quote support)
func splitCSV(line string) []string {
	return strings.Split(line, ",")
}

// Helper: convert []map[string]string to [][]string for mapping
func recordsToSlice(records []map[string]string, headers []string) [][]string {
	var out [][]string
	for _, rec := range records {
		row := make([]string, len(headers))
		for i, h := range headers {
			row[i] = rec[h]
		}
		out = append(out, row)
	}
	return out
}
