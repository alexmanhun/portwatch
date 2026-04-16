// Package history provides persistent storage and querying of port change
// events recorded by the portwatch monitor.
//
// Events are appended to a JSON file on disk so that history survives
// process restarts. Use New to open (or create) a history file, Record to
// append events, and Query or Last to retrieve them.
package history
