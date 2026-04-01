// https://github.com/libyal/libevtx/blob/main/documentation/Windows%20XML%20Event%20Log%20(EVTX).asciidoc
// https://www.researchgate.net/publication/222426407_Introducing_the_Microsoft_Vista_event_log_file_format
// https://blog.fox-it.com/2019/06/04/export-corrupts-windows-event-log-files/
// https://blog.fox-it.com/2017/12/08/detection-and-recovery-of-nsas-covered-up-tracks/
// https://ernw.de/download/EventManipulation.pdf
// https://parsiya.net/blog/2018-11-01-windows-filetime-timestamps-and-byte-wrangling-with-go/
package main

import (
	"bufio"
	"maps"
	"os"

	"go.foxforensics.dev/tri/internal"
)

func main() {
	var count uint64

	reader := bufio.NewReader(os.Stdin)

	for internal.ReadUntil(reader, internal.Signature) {
		record := internal.NewRecord(reader)

		if record.IsSizeValid() && record.IsTimeValid() {
			count++
			internal.Debug("[+] found record #%d\n%s\n", record.Id, record.String())
		} else {
			internal.Debug("[-] skip record #%d\n", record.Id)
		}

		//if record.Fragment != nil {
		//	fmt.Printf("%s %s %d\n",
		//		internal.FileTime(record.Time).UTC().Format(time.RFC3339),
		//		record.Fragment.Computer,
		//		record.Fragment.EventId,
		//	)
		//}
	}

	for k, v := range maps.All(internal.Computers) {
		internal.Debug("[+] Computer [%08x] %s\n", k, v)
	}

	internal.Debug("\n[*] found %d records\n", count)
}
