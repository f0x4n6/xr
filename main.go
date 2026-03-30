// https://github.com/libyal/libevtx/blob/main/documentation/Windows%20XML%20Event%20Log%20(EVTX).asciidoc
// https://www.researchgate.net/publication/222426407_Introducing_the_Microsoft_Vista_event_log_file_format
// https://blog.fox-it.com/2019/06/04/export-corrupts-windows-event-log-files/
// https://blog.fox-it.com/2017/12/08/detection-and-recovery-of-nsas-covered-up-tracks/
// https://ernw.de/download/EventManipulation.pdf
// https://parsiya.net/blog/2018-11-01-windows-filetime-timestamps-and-byte-wrangling-with-go/
package main

import (
	"bufio"
	"fmt"
	"os"

	"go.foxforensics.dev/tri/internal"
)

func main() {
	var count uint64

	reader := bufio.NewReader(os.Stdin)

	for internal.ReadUntil(reader, internal.Signature) {
		record := internal.NewRecord(reader)

		if record.IsSkipped() {
			fmt.Printf("[-] skip record #%d\n\n", record.Id)
			continue
		}

		if record.IsValid() {
			fmt.Printf("[+] found record #%d\n\n", record.Id)
		} else {
			fmt.Printf("[!] found record #%d (corrupt)\n\n", record.Id)
		}

		count++

		fmt.Println(record.String())
		//fmt.Printf("%s %04d\n", internal.FileTime(record.Time).UTC().Format(time.RFC3339), record.EventId)
	}

	fmt.Printf("[+] found %d records\n", count)
}
