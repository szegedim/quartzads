package cards

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"showmycard.com/metadata"
	"strings"
)

//Licensed under Creative Commons CC0.
//
//To the extent possible under law, the author(s) have dedicated all copyright and related and
//neighboring rights to this software to the public domain worldwide.
//This software is distributed without any warranty.
//You should have received a copy of the CC0 Public Domain Dedication along with this software.
//If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

func Setup() {
	http.HandleFunc("/png", func(writer http.ResponseWriter, request *http.Request) {
		apiKey := request.URL.Query().Get("apikey")
		_, _ = io.Copy(writer, bytes.NewBuffer(GetPicture(metadata.SGUID(apiKey))))
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
			BookImpressionOrCLick(metadata.SGUID(split[0]), string(message))
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
			privateKey := metadata.GenerateSGUID()
			Set(string(privateKey), []byte(apiKey))
			http.Redirect(writer, request, "/rp?apikey="+string(privateKey), http.StatusTemporaryRedirect)
		}
	})
	http.HandleFunc("/ad", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			form, _ := os.Open("res/testupload.html")
			defer form.Close()

			billing := os.Getenv("PAYMENTURL")
			if billing != "" {
				buf := bytes.Buffer{}
				_, _ = io.Copy(&buf, form)
				ret := strings.ReplaceAll(buf.String(), "https://buy.stripe.com/test_5kA4gMaZYaiY3JK8ww", billing)
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

		apiKey := request.URL.Query().Get("apikey")
		SetPicture(metadata.SGUID(apiKey), png.Bytes())

		message := request.FormValue("message")
		var url string
		_, _ = fmt.Sscanf(message, "Point the ad card to %s and inject ads.", &url)
		if len(url) > 4 {
			SetTarget(metadata.SGUID(apiKey), url)
			// TODO cleanup after 24 hours
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
		target := GetTarget(metadata.SGUID(apiKey))
		s0 := strings.ReplaceAll(buf.String(), "https://showmycard.com/7d5f7d3c5a885cead3102434c9becd8a", target)
		s0 = strings.ReplaceAll(s0, "f7e77596ddf4b12e38e469421d78f3cc.png", "/png?apikey="+apiKey)
		begin := strings.Index(s0, "<!-- Data begins here -->")
		end := strings.Index(s0, "<!-- Data ends here -->")
		if begin != -1 && end != -1 {
			clicks, impressions := GetStatistics(metadata.SGUID(apiKey))
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
		s0 = strings.ReplaceAll(s0, "<!-- Location goes here -->", "<a href=\""+GetLocation(metadata.SGUID(apiKey))+"\">Paid media</a>")
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
		name := fmt.Sprintf("card%04d", i)
		cardId := NewCard(req.URL.Path + "#" + name)
		target := GetTarget(cardId)
		card := fmt.Sprintf(`
		<div class="showmycard" aria-label="Description of the image">
			<img class="showmycardimg" id='%s' src="%s" alt="Descriptive text" style="width: 3in;height: auto;" onclick="clicked(event.target, '%s')">
		</div>
		`, cardId, "/png?apikey="+cardId, target)
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
