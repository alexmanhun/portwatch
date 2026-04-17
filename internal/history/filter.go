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
