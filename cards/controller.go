package cards

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"showmycard.com/englang"
	"showmycard.com/metadata"
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
	billing := os.Getenv("PAYMENTURL")
	if billing == "" {
		// Localhost test behavior running in
		metadata.DefaultAdTime = metadata.TestAdTime
		metadata.DefaultPurchaseTime = metadata.TestPurchaseTime
	}
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
		if handleReport(writer, request) {
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
			items := englang.Englang(string(buf), "%s<!-- Repeat this for all cards -->%s<!-- ... -->%s")
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
			form, _ := os.Open("./res/testsetup.html")
			defer form.Close()
			_, _ = io.Copy(writer, form)
			return
		}

		message := request.FormValue("message")
		if message != "" {
			proxy := ""
			_, _ = fmt.Sscanf(message, "Add advertisement to %s page that you have rights to.", &proxy)
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
			privateKey := englang.GenerateSGUID()
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
			expiry := time.Now().Add(metadata.DefaultPurchaseTime)
			png, _ := os.ReadFile("res/waitingwithmessage.png")
			SetPicture(englang.SGUID(apiKey), png)
			SetTarget(englang.SGUID(apiKey), ".")
			AddActivity(englang.SGUID(apiKey), fmt.Sprintf("Card will expire on %s and revert to ad.", expiry.Format(time.RFC822Z)))

			form, _ := os.Open("res/testupload.html")
			defer form.Close()

			billing := os.Getenv("PAYMENTURL")
			if billing != "" {
				buf := bytes.Buffer{}
				_, _ = io.Copy(&buf, form)
				ret := strings.ReplaceAll(buf.String(), "https://buy.stripe.com/test_00gfZueca62I1BC9AB", billing)
				_, _ = io.WriteString(writer, ret)
			} else {
				_, _ = io.Copy(writer, form)
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
	http.HandleFunc("/", proxyCore)
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
		s0 := strings.ReplaceAll(buf.String(), "https://showmycard.com/7d5f7d3c5a885cead3102434c9becd8a", target)
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

func proxyCore(res http.ResponseWriter, req *http.Request) {
	if metadata.ProxySite == "" {
		http.Redirect(res, req, "/up", http.StatusTemporaryRedirect)
		return
	}
	url := metadata.ProxySite + req.URL.Path

	request, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	for key, value := range req.Header {
		request.Header.Set(key, value[0])
	}
	request.Header.Set("Accept-Encoding", "identity")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		//http.Error(res, err.Error(), http.StatusBadGateway)
		return
	}
	defer response.Body.Close()

	for key, value := range response.Header {
		res.Header().Set(key, value[0])
	}

	res.WriteHeader(response.StatusCode)
	res.Header().Set("Cache-Control", "no-store")

	content, _ := io.ReadAll(response.Body)
	contentWithCards := string(content)

	placeholders := strings.Count(contentWithCards, metadata.Placeholder)
	adBlocker := fmt.Sprintf("<div style=\"text-align: center\"><p>Block ads for your convenience. <a href=\"%s\">🍁(hop)</a><!--%s--> </p></div>\n", metadata.ProxySite, metadata.ProxySite)
	if placeholders == 0 {
		contentWithCards = strings.ReplaceAll(contentWithCards, "<body", "<body><br><br><br><br>"+adBlocker+metadata.Placeholder+"<div")
		contentWithCards = strings.ReplaceAll(contentWithCards, "</body>", "</div>"+metadata.Placeholder+"</body>")
		placeholders = 2
	} else {
		contentWithCards = strings.ReplaceAll(contentWithCards, "<body", "<body>"+adBlocker+"<div")
		contentWithCards = strings.ReplaceAll(contentWithCards, "</body>", "</div>"+"</body>")
	}

	for i := 0; i < placeholders; i++ {
		cardId, expiryLog := GetCard(req.URL.Path, i)

		target := GetTarget(cardId)
		report := fmt.Sprintf("<button style=\"position:absolute;top:90%%;right:1%%;opacity:70%%;font-size: xx-small;text-align: right;color: darkorchid;font-family: system-ui\" onclick=\"fetch('%s');location.reload();\">Report</button>", "/pg18?apikey="+cardId)
		if strings.HasPrefix(target, "/") {
			report = ""
		}
		card := fmt.Sprintf(`
		<div class="showmycard" aria-label="Description of the image">
			%s
			<img class="showmycardimg" id='%s' src="%s" alt="Descriptive text" style="width: 3in;height: auto;" onclick="clicked(event.target, '%s')">
			%s
		</div>
		`, expiryLog, cardId, "/png?apikey="+cardId, target, report)
		contentWithCards = strings.Replace(contentWithCards, metadata.Placeholder, card, 1)
	}
	base, _ := os.ReadFile("./res/testcard.html")
	baseString := string(base)
	beginHeader := strings.Index(baseString, "<!-- Style starts here -->")
	endHeader := strings.Index(baseString, "<!-- Style ends here -->")
	contentWithCards = strings.ReplaceAll(contentWithCards, "</head>", baseString[beginHeader:endHeader+len("<!-- Style ends here -->")]+"</head>")

	beginScriptHeader := strings.Index(baseString, "<!-- Script starts here -->")
	endScriptHeader := strings.Index(baseString, "<!-- Script ends here -->")
	contentWithCards = strings.ReplaceAll(contentWithCards, "</body>", baseString[beginScriptHeader:endScriptHeader+len("<!-- Script ends here -->")]+"</body>")

	contentWithCards = strings.ReplaceAll(contentWithCards, metadata.ProxySite+req.URL.Path, req.URL.Path)
	contentWithCards = strings.ReplaceAll(contentWithCards, "utm_content=sitename", "utm_content="+metadata.SiteName)

	//if strings.HasPrefix(response.Header.Get("Content-Type"), "text/html") {
	//	fmt.Println(contentWithCards)
	//}
	b2 := bytes.NewBufferString(contentWithCards)
	_, _ = io.Copy(res, b2)
}
