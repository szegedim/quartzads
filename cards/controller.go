package cards

import (
	"bufio"
	"bytes"
	"fmt"
	"gitlab.com/eper.io/quartzads/englang"
	"gitlab.com/eper.io/quartzads/metadata"
	"io"
	"net/http"
	"os"
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

func Setup() {
	implementation := os.Getenv("IMPLEMENTATION")
	englang.EnglangRemoteImplementation(implementation)

	http.HandleFunc("/1223c99f-70fe-40be-abe3-bf1c6ba1bdb6.txt", func(writer http.ResponseWriter, request *http.Request) {
		ret := metadata.GetDefaultImplementation()
		_, _ = io.WriteString(writer, ret)
	})
	http.HandleFunc("/png", func(writer http.ResponseWriter, request *http.Request) {
		apiKey := request.URL.Query().Get("apikey")
		_, _ = io.Copy(writer, bytes.NewBuffer(GetPicture(englang.SGUID(apiKey))))
	})
	http.HandleFunc("/logo.png", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "image/png")
		form, _ := os.Open("./res/img.png")
		defer form.Close()
		_, _ = io.Copy(writer, form)
		return
	})
	http.HandleFunc("/rp", func(writer http.ResponseWriter, request *http.Request) {
		if ReportOnCard(writer, request) {
			return
		}
		return
	})
	http.HandleFunc("/bids", func(writer http.ResponseWriter, request *http.Request) {
		privateKey := request.URL.Query().Get("apikey")
		apiKey := string(Get(privateKey))
		if apiKey != "" {
			// TODO list uptime, impressions and clicks
			// TODO Security allows to hide/enable ads here.
			buf, _ := os.ReadFile("./res/testbids.html")
			ret := Customize(string(buf))
			items := englang.SplitEnglang(ret, "%s<!-- Repeat this for all cards -->%s<!-- ... -->%s")
			if len(items) == 3 {
				_, _ = io.WriteString(writer, items[0])
				for k, _ := range redis {
					png := GetPicture(englang.SGUID(k))
					if len(png) > 0 {
						_, _ = io.WriteString(writer, strings.ReplaceAll(items[1], "img.png", "/png?apikey="+k))
					}
				}
				_, _ = io.WriteString(writer, items[2])
			}
		}
	})
	http.HandleFunc("/up", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			form, _ := os.ReadFile("./res/testsetup.html")
			s0 := Customize(string(form))
			_, _ = io.WriteString(writer, s0)
			return
		}

		message := request.FormValue("message")
		if message != "" {
			proxy := ""
			_, _ = fmt.Sscanf(message, "Add advertisement to %s page. Make sure that you have copyright.", &proxy)
			if proxy != "" {
				metadata.ProxySite = proxy
			}
			return
		}
	})
	http.HandleFunc("/log", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "PUT" {
			message, _ := io.ReadAll(request.Body)
			split := strings.SplitN(string(message), ":", 2)
			AddActivity(englang.SGUID(split[0]), string(message))
		}
	})
	http.HandleFunc("/englang", func(writer http.ResponseWriter, request *http.Request) {
		keysText := List()

		buf := bufio.NewWriter(writer)
		_, _ = buf.WriteString(keysText)
		_ = buf.Flush()

		keys := strings.Split(keysText, "\n")
		for _, k := range keys {
			if strings.HasSuffix(k, ".activity") {
				_, _ = buf.Write(Get(k))
				_, _ = buf.WriteString("\n")
				_ = buf.Flush()
			}
		}
	})
	http.HandleFunc("/paid", func(writer http.ResponseWriter, request *http.Request) {
		apiKey := request.URL.Query().Get("utm_source")
		redirectSite := request.URL.Query().Get("utm_content")
		if apiKey != "" {
			http.Redirect(writer, request, redirectSite+"/paid?apikey="+apiKey, http.StatusTemporaryRedirect)
		}
		apiKey = request.URL.Query().Get("apikey")
		if apiKey != "" {
			privateKey := "private" + englang.GenerateSGUID()
			Set(string(privateKey), []byte(apiKey))
			expiry := time.Now().Add(metadata.DefaultAdTime)
			AddActivity(englang.SGUID(apiKey), fmt.Sprintf("Card will expire on %s and revert to ad.", expiry.Format(time.RFC822Z)))
			http.Redirect(writer, request, "/rp?apikey="+string(privateKey), http.StatusTemporaryRedirect)
		}
	})
	http.HandleFunc("/pg18", func(writer http.ResponseWriter, request *http.Request) {
		apiKey := request.URL.Query().Get("apikey")
		if apiKey != "" {
			png, _ := os.ReadFile("res/noncompliant.png")
			SetPicture(englang.SGUID(apiKey), png)
		}
	})
	http.HandleFunc("/ad", func(writer http.ResponseWriter, request *http.Request) {
		apiKey := request.URL.Query().Get("apikey")
		if request.Method == "GET" {
			LockAd(apiKey)

			form, _ := os.ReadFile("res/testupload.html")
			cardUpload := Customize(string(form))

			contact := GetContactInfo()
			cardUpload = strings.Replace(cardUpload, "</body>", contact+"</body>", 1)

			if metadata.PaymentUrl != "" {
				ret := strings.ReplaceAll(cardUpload, "https://buy.stripe.com/test_00gfZueca62I1BC9AB", metadata.PaymentUrl)
				ret = Customize(ret)
				_, _ = io.WriteString(writer, ret)
			} else {
				_, _ = io.WriteString(writer, cardUpload)
			}
			return
		}

		file, _, err := request.FormFile("ad")
		if err != nil {
			http.Error(writer, "Error retrieving the file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		png := bytes.Buffer{}
		_, _ = io.Copy(&png, file)

		SetPicture(englang.SGUID(apiKey), png.Bytes())

		message := request.FormValue("message")
		var url string
		_, _ = fmt.Sscanf(message, "Point the ad card to %s and inject ads.", &url)
		if len(url) > 4 {
			SetTarget(englang.SGUID(apiKey), url)
		} else {
			SetTarget(englang.SGUID(apiKey), ".")
		}
	})
	http.HandleFunc("/", CoreProxy)
}

func Customize(ret string) string {
	if metadata.SiteTitle != "" {
		ret = strings.ReplaceAll(ret, "Show My Card℠", metadata.SiteTitle)
	}
	if metadata.ProxySite != "" {
		ret = strings.ReplaceAll(ret, "https://hq45a13f0d8b0f0.wordpress.com", metadata.ProxySite)
	}
	if metadata.DefaultAdTime.Hours() > 4*24 && int(metadata.DefaultAdTime.Hours())%24 == 0 {
		ret = strings.ReplaceAll(ret, "The card is displayed for 24 hours.", fmt.Sprintf("The card is displayed for %00.1f days.", metadata.DefaultAdTime.Hours()/24))
	} else if metadata.DefaultAdTime.Hours() > 1 {
		ret = strings.ReplaceAll(ret, "The card is displayed for 24 hours.", fmt.Sprintf("The card is displayed for %00.1f hours.", metadata.DefaultAdTime.Hours()))
	} else if metadata.DefaultAdTime.Minutes() > 1 {
		ret = strings.ReplaceAll(ret, "The card is displayed for 24 hours.", fmt.Sprintf("The card is displayed for %00.1f minutes.", metadata.DefaultAdTime.Minutes()))
	} else {
		ret = strings.ReplaceAll(ret, "The card is displayed for 24 hours.", "The card is displayed for some time.")
	}
	return ret
}

func handleReport(writer http.ResponseWriter, request *http.Request) bool {
	privateKey := request.URL.Query().Get("apikey")
	apiKey := string(Get(privateKey))
	if request.Method == "GET" {
		form, _ := os.Open("./res/testreport.html")
		defer form.Close()
		buf := bytes.Buffer{}
		_, _ = io.Copy(&buf, form)
		target := GetTarget(englang.SGUID(apiKey))
		s0 := Customize(buf.String())
		s0 = strings.ReplaceAll(s0, "https://showmycard.com/7d5f7d3c5a885cead3102434c9becd8a", target)
		s0 = strings.ReplaceAll(s0, "f7e77596ddf4b12e38e469421d78f3cc.png", "/png?apikey="+apiKey)
		begin := strings.Index(s0, "<!-- Data begins here -->")
		end := strings.Index(s0, "<!-- Data ends here -->")
		if begin != -1 && end != -1 {
			clicks, impressions := GetStatistics(englang.SGUID(apiKey))
			clickbar := clicks
			if clickbar > 60 {
				clickbar = 60
			}
			impressionbar := impressions
			if impressionbar > 60 {
				impressionbar = 60
			}
			pattern := `
				<rect class="impression" x="20" y="%d" width="10" height="%d" rx="2.5"></rect>
				<rect class="click" x="60" y="%d" width="10" height="%d" rx="2.5"></rect>
				<text class="impression" x="20" y="%d">%d</text>
				<text class="click" x="60" y="%d">%d</text>
				`
			s0 = strings.ReplaceAll(s0, s0[begin:end+len("<!-- Data ends here -->")], fmt.Sprintf(pattern, 85-impressionbar, impressionbar, 85-clickbar, clickbar, 96, impressions, 96, clicks))
		}
		s0 = strings.ReplaceAll(s0, "<!-- Location goes here -->", "<a href=\""+GetLocation(englang.SGUID(apiKey))+"\">Paid media</a>")
		buffered := bufio.NewWriter(writer)
		_, _ = buffered.WriteString(s0)
		_ = buffered.Flush()
		return true
	}
	return false
}
