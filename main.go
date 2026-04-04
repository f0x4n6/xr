package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"time"

	"go.foxforensics.dev/tri/pkg/evtx"
	"go.foxforensics.dev/tri/pkg/utils"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--debug" {
		evtx.Debug = true // activate debug mode
	}

	var count uint64

	reader := bufio.NewReader(os.Stdin)

	for utils.ReadUntil(reader, evtx.Signature) {
		record := evtx.NewRecord(reader)

		if evtx.Debug {
			if record.IsSizeValid() && record.IsTimeValid() {
				count++
				utils.Debug("[+] Found record #%d\n%s\n", record.Id, record.String())
			} else {
				utils.Debug("[-] Skip record #%d\n", record.Id)
				continue
			}
		} else {
			if record != nil && record.Fragment != nil {
				fmt.Printf("%s %s %d\n",
					utils.FileTime(record.Time).UTC().Format(time.RFC3339),
					record.Fragment.Computer,
					record.Fragment.EventId,
				)
			}
		}
	}

	if evtx.Debug {
		utils.Debug("[+] Found %d computers\n", len(evtx.Computers))

		for k, v := range maps.All(evtx.Computers) {
			utils.Debug("[+]  [%08x] %s\n", k, v)
		}

		utils.Debug("[+] Found %d records\n", count)
	}
}
