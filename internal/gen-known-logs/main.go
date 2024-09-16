package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	out = flag.String("out", "", "Output file")
)

// https://github.com/google/certificate-transparency-community-site/blob/master/docs/google/known-logs.md
const knownLogsAddr = "https://www.gstatic.com/ct/log_list/v2/all_logs_list.json"

type ctLog struct {
	operator string
	url      string
}

func main() {
	flag.Parse()

	resp, err := http.Get(knownLogsAddr)
	if err != nil {
		log.Fatalf("Failed to fetch the list of known CT logs: %v", err)
	}
	defer resp.Body.Close()

	var logs struct {
		Operators []struct {
			Name string `json:"name"`
			Logs []struct {
				ID  string `json:"log_id"`
				URL string `json:"url"`
			} `json:"logs"`
		} `json:"operators"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		log.Fatalf("Failed to parse the list of known CT logs: %v", err)
	}

	f, err := os.Create(*out)
	if err != nil {
		log.Fatalf("Failed to open file %s", *out)
	}
	defer f.Close()

	fmt.Fprintf(f, `// Code generated by github.com/square/certigo/internal/gen-known-logs; DO NOT EDIT
// Generated at %s
package %s

type ctLog struct {
	operator string
	url      string
}

var knownLogs = map[string]*ctLog{
`, time.Now().Format(time.RFC3339), os.Getenv("GOPACKAGE"))
	knownLogs := make(map[string]*ctLog)
	for _, op := range logs.Operators {
		for _, l := range op.Logs {
			fmt.Fprintf(f, "\t%q: {operator: %q, url: %q},\n", l.ID, op.Name, l.URL)
			knownLogs[l.ID] = &ctLog{
				operator: op.Name,
				url:      l.URL,
			}
		}
	}
	fmt.Fprintln(f, "}")
}