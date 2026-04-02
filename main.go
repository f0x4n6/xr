package main

import (
	"bufio"
	"maps"
	"os"

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
			utils.Debug("[+] found record #%d\n%s\n", record.Id, record.String())
		} else {
			utils.Debug("[-] skip record #%d\n", record.Id)
		}

		//if record.Fragment != nil {
		//	fmt.Printf("%s\t%5d\t%s\n",
		//		internal.FileTime(record.Time).UTC().Format(time.RFC3339),
		//		record.Fragment.EventId,
		//		record.Fragment.Computer,
		//	)
		//}
	}

	for k, v := range maps.All(evtx.Computers) {
		utils.Debug("[+] Computer [%08x] %s\n", k, v)
	}

	utils.Debug("\n[*] found %d records\n", count)
}
