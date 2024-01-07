package englang

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"gitlab.com/eper.io/quartzads/metadata"
	"io"
	"os"
)

func Backup(in map[string][]byte) {
	b := bytes.Buffer{}
	keys := make([]string, 0)
	for k, _ := range in {
		keys = append(keys, k)
	}
	for _, k := range keys {
		b.WriteString(fmt.Sprintf("\nSegment %s has %d bytes that follow.\n", k, len(in[k])))
		b.Write(in[k])
	}

	hash := fmt.Sprintf("/tmp/snapshot%x", sha256.Sum256(b.Bytes()))
	f, _ := os.Create(hash)
	defer f.Close()
	_, _ = io.Copy(f, &b)
	metadata.LatestSnapshotFile = hash
	fmt.Printf("Backed up %d records to %s file.\n", len(keys), metadata.LatestSnapshotFile)
}

func Restore(out *map[string][]byte) {
	if metadata.LatestSnapshotFile == "" {
		fmt.Println("no snapshot to restore")
		return
	}
	f, _ := os.Open(metadata.LatestSnapshotFile)
	scanner := bufio.NewScanner(f)
	segment := ""
	length := 0
	n := 0
	for scanner.Scan() {
		if segment == "" {
			_, _ = fmt.Sscanf(scanner.Text(), "Segment %s has %d bytes that follow.", &segment, &length)
		} else {
			if length <= len(scanner.Bytes()) {
				(*out)[segment] = scanner.Bytes()[0:length]
				n++
			}
			segment = ""
			length = 0
		}
	}
	fmt.Printf("Restored %d records from %s file.\n", n, metadata.LatestSnapshotFile)
}
