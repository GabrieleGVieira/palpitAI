package importer

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/data/models"
)

const (
	sourceInternationalResults = "international-results"
	sourceFifaRanking          = "fifa-ranking-historical"
)

type ImportResult struct {
	models.ImportSummary
}

func headerIndex(header []string) map[string]int {
	index := make(map[string]int, len(header))
	for i, name := range header {
		index[normalizeHeader(name)] = i
	}

	return index
}

func normalizeHeader(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func field(record []string, index map[string]int, names ...string) string {
	for _, name := range names {
		if i, ok := index[normalizeHeader(name)]; ok && i < len(record) {
			return strings.TrimSpace(record[i])
		}
	}

	return ""
}

func requireFields(values map[string]string) error {
	for name, value := range values {
		if strings.TrimSpace(value) == "" {
			return errors.New("missing required field: " + name)
		}
	}

	return nil
}

func parseDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	for _, layout := range []string{"2006-01-02", "1/2/2006", "01/02/2006", "02/01/2006", "2006/01/02"} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, errors.New("invalid date: " + value)
}

func parseInt(value string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(value))
}

func parseOptionalInt(value string) (*int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func parseOptionalFloat(value string) (*float64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func parseBool(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "t", "1", "yes", "y":
		return true
	default:
		return false
	}
}

func readHeader(reader *csv.Reader) (map[string]int, error) {
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	return headerIndex(header), nil
}

func nextRecord(reader *csv.Reader) ([]string, error) {
	record, err := reader.Read()
	if err != nil {
		return nil, err
	}

	return record, nil
}

func isEOF(err error) bool {
	return errors.Is(err, io.EOF)
}
