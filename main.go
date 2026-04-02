// https://github.com/libyal/libevtx/blob/main/documentation/Windows%20XML%20Event%20Log%20(EVTX).asciidoc
// https://github.com/libyal/libfwnt/blob/main/documentation/Security%20Descriptor.asciidoc
// https://blog.fox-it.com/2019/06/04/export-corrupts-windows-event-log-files/
// https://blog.fox-it.com/2017/12/08/detection-and-recovery-of-nsas-covered-up-tracks/
// https://www.researchgate.net/publication/222426407_Introducing_the_Microsoft_Vista_event_log_file_format
// https://parsiya.net/blog/2018-11-01-windows-filetime-timestamps-and-byte-wrangling-with-go/
// https://ernw.de/download/EventManipulation.pdf
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
	var count uint64

	reader := bufio.NewReader(os.Stdin)

	for utils.ReadUntil(reader, evtx.Signature) {
		record := evtx.NewRecord(reader)

		if record.IsSizeValid() && record.IsTimeValid() {
			count++
			utils.Debug("[+] Found record #%d\n%s\n", record.Id, record.String())
		} else {
			utils.Debug("[-] Skip record #%d\n", record.Id)
			continue
		}

		if record.Fragment != nil {
			fmt.Printf("%s %s %d\n",
				utils.FileTime(record.Time).UTC().Format(time.RFC3339),
				record.Fragment.Computer,
				record.Fragment.EventId,
			)
		}
	}

	for k, v := range maps.All(evtx.Computers) {
		utils.Debug("[+] Computer [%08x] %s\n", k, v)
	}

	utils.Debug("\n[*] Found %d records\n", count)
}
