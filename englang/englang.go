package englang

import (
	"fmt"
	"gitlab.com/eper.io/quartzads/metadata"
	"io"
	"net/http"
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

func RunEnglang(instructions string) {
	s := metadata.GetDefaultImplementation()
	response, err := http.Get(instructions)
	if err == nil && response != nil && response.Body != nil {
		defer response.Body.Close()
		buf, _ := io.ReadAll(response.Body)
		s = string(buf)
	}
	fmt.Println(s)
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		tokens := SplitEnglang(strings.TrimSpace(line), "Set the payment url to %s address.")
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
	}
}
