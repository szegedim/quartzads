package englang

import (
	"bytes"
	"fmt"
	"gitlab.com/eper.io/quartzads/metadata"
	"gitlab.com/eper.io/quartzads/storage"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

//Licensed under Creative Commons CC0.
//
//To the extent possible under law, the author(s) have dedicated all copyright and related and
//neighboring rights to this software to the public domain worldwide.
//This software is distributed without any warranty.
//You should have received a copy of the CC0 Public Domain Dedication along with this software.
//If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

func SplitEnglang(str string, pattern string) (items []string) {
	//func TestA(t *testing.T) {
	//	t.Log(SplitEnglang("The color is red today and green tomorrow.", "The color is %s today and %s tomorrow."))
	//	t.Log(SplitEnglang("The color is red today and tomorrow. blue", "The color is %s today and tomorrow. %s"))
	//	t.Log(SplitEnglang("yellow The color is red today and tomorrow.", "%s The color is %s today and tomorrow."))
	//}

	str = "^" + str + "$"
	patterns := strings.Split("^"+pattern+"$", "%s")
	index := 0
	items = []string{}

	for len(patterns) > 1 {
		next := strings.Index(str[index:], patterns[0])
		if next == -1 {
			break
		}
		end := strings.Index(str[index+next:], patterns[1])
		items = append(items, str[index+next+len(patterns[0]):index+next+end])
		index = index + next + end
		patterns = patterns[1:]
	}
	return
}

func EnglangRemoteImplementation(instructions string) {
	implementationFile := metadata.GetDefaultImplementation()
	response, err := http.Get(instructions)
	if err == nil && response != nil && response.Body != nil {
		//goland:noinspection ALL
		defer response.Body.Close()
		buf, _ := io.ReadAll(response.Body)
		implementationFile = string(buf)
	}
	EnglangImplementation(implementationFile, err)
}

func EnglangImplementation(implementationFile string, err error) {
	fmt.Println(implementationFile)
	lines := strings.Split(implementationFile, "\n")
	for _, line := range lines {
		tokens := SplitEnglang(strings.TrimSpace(line), "Replace all references of %s in the resource project to %s in the current implementation.")
		if len(tokens) == 2 {
			list, _ := os.ReadDir("./res")
			for _, item := range list {
				if !item.IsDir() && strings.HasSuffix(item.Name(), ".html") {
					name := path.Join("./res", item.Name())
					in, _ := os.ReadFile(name)
					out := bytes.Replace(in, []byte(tokens[0]), []byte(tokens[1]), 1)
					_ = os.WriteFile(name, out, 0600)
				}
			}
			metadata.AdBlocker = strings.ReplaceAll(metadata.AdBlocker, tokens[0], tokens[1])
			metadata.AdDescription = strings.ReplaceAll(metadata.AdDescription, tokens[0], tokens[1])
		}
		tokens = SplitEnglang(strings.TrimSpace(line), "Set the payment url to %s address.")
		if len(tokens) > 0 {
			metadata.PaymentUrl = tokens[0]
		}
		tokens = SplitEnglang(strings.TrimSpace(line), "Set the title to %s text.")
		if len(tokens) > 0 {
			metadata.SiteTitle = tokens[0]
		}
		tokens = SplitEnglang(strings.TrimSpace(line), "Proxy the %s site.")
		if len(tokens) > 0 {
			metadata.ProxySite = tokens[0]
		}
		tokens = SplitEnglang(strings.TrimSpace(line), "Point contact to %s site.")
		if len(tokens) > 0 {
			metadata.Contact = tokens[0]
		}
		tokens = SplitEnglang(strings.TrimSpace(line), "Point terms to %s site.")
		if len(tokens) > 0 {
			metadata.Terms = tokens[0]
		}
		tokens = SplitEnglang(strings.TrimSpace(line), "Set ad time to %s hours.")
		if len(tokens) > 0 {
			hours, _ := strconv.ParseInt(tokens[0], 10, 32)
			if hours > 0 && hours < 10000 {
				metadata.DefaultAdTime = time.Duration(hours) * time.Hour
			}
		}
		tokens = SplitEnglang(strings.TrimSpace(line), "Upload snapshot every %s seconds to %s site.")
		if len(tokens) == 2 {
			AutoBackup(tokens, line)
		}
		tokens = SplitEnglang(strings.TrimSpace(line), "Download %s to %s file.")
		if len(tokens) == 2 {
			DownloadFile(tokens[0], tokens[1])
		}
		tokens = SplitEnglang(strings.TrimSpace(line), "Restore from %s file.")
		if len(tokens) > 0 {
			Restore(&storage.Redis, tokens[0])
		}
	}
}

func DownloadFile(url string, path string) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	} else {
		out, _ := os.Create(path)
		defer out.Close()
		_, _ = io.Copy(out, res.Body)
	}
}

func AutoBackup(tokens []string, line string) {
	seconds, _ := strconv.ParseInt(tokens[0], 10, 32)
	if seconds > 0 && seconds < 10000 {
		if strings.HasPrefix(tokens[1], "https://") || strings.HasPrefix(tokens[1], "file://") {
			fmt.Println(strings.TrimSpace(line))
			go func(timer time.Duration, site string) {
				for {
					time.Sleep(timer)
					Backup(storage.Redis)
					_, err := os.Stat(metadata.LatestSnapshotFile)
					if err != nil {
						continue
					}
					if strings.HasPrefix(site, "file://") {
						copyFile(site)
					}
					if strings.HasPrefix(site, "https://") {
						uploadFile(site)
					}
					fmt.Printf("Backup finished to %s location.\n", site)
				}
			}(time.Duration(seconds)*time.Second, tokens[1])
		}
	}
}

func copyFile(site string) {
	// TODO do some sec checks
	site = site[len("file://"):]
	out, _ := os.Create(site)
	in, _ := os.Open(metadata.LatestSnapshotFile)
	_, _ = io.Copy(out, in)
	defer out.Close()
	defer in.Close()
}

func uploadFile(site string) {
	if !strings.HasPrefix(site, "https://") {
		return
	}
	// TODO do some sec checks
	in, _ := os.Open(metadata.LatestSnapshotFile)
	defer in.Close()
	resp, _ := http.Post(site, "application/octet-stream", in)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			fmt.Printf("Upload error %d code.\n", resp.StatusCode)
		}
	}
}
