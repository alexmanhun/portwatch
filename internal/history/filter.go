package history

import "time"

// FilterOptions defines criteria for filtering history records.
type FilterOptions struct {
	Port    int
	Event   string
	Since   time.Time
	Until   time.Time
	Limit   int
}

// Filter returns records matching all non-zero criteria in opts.
// Records are filtered by port, event type, and time range before
// the limit is applied, keeping the most recent matching records.
func Filter(records []Record, opts FilterOptions) []Record {
	var out []Record
	for _, r := range records {
		if opts.Port != 0 && r.Port != opts.Port {
			continue
		}
		if opts.Event != "" && r.Event != opts.Event {
			continue
		}
		if !opts.Since.IsZero() && r.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && r.Timestamp.After(opts.Until) {
			continue
		}
		out = append(out, r)
	}
	if opts.Limit > 0 && len(out) > opts.Limit {
		out = out[len(out)-opts.Limit:]
	}
	return out
}

// FilterByPort is a convenience wrapper that returns all records for the given port.
func FilterByPort(records []Record, port int) []Record {
	return Filter(records, FilterOptions{Port: port})
}
