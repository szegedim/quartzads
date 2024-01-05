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
)

//Licensed under Creative Commons CC0.
//
//To the extent possible under law, the author(s) have dedicated all copyright and related and
//neighboring rights to this software to the public domain worldwide.
//This software is distributed without any warranty.
//You should have received a copy of the CC0 Public Domain Dedication along with this software.
//If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

func CoreProxy(res http.ResponseWriter, req *http.Request) {
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
	adBlocker := fmt.Sprintf(metadata.AdBlocker, metadata.ProxySite, metadata.ProxySite)

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
		<div class="quartzads" aria-label="Description of the image">
			%s
			<img class="quartzadsimg" id='%s' src="%s" alt="Descriptive text" style="width: 3in;height: auto;" onclick="clicked(event.target, '%s')">
			%s
		</div>
		`, expiryLog, cardId, "/png?apikey="+cardId, target, report)
		contentWithCards = strings.Replace(contentWithCards, metadata.Placeholder, card, 1)
	}
	if metadata.SiteTitle != "" {
		contentWithCards = strings.Replace(contentWithCards, "<title>", "<!--<title>", 1)
		contentWithCards = strings.Replace(contentWithCards, "</title>", "</title>-->", 1)
		contentWithCards = strings.Replace(contentWithCards, "<!--<title>", fmt.Sprintf("<title>%s</title><!--<title>", metadata.SiteTitle), 1)
	}

	base, _ := os.ReadFile("./res/testcard.html")
	baseString := string(base)
	baseString = Customize(baseString)
	beginHeader := strings.Index(baseString, "<!-- Style starts here -->")
	endHeader := strings.Index(baseString, "<!-- Style ends here -->")
	contentWithCards = strings.ReplaceAll(contentWithCards, "</head>", baseString[beginHeader:endHeader+len("<!-- Style ends here -->")]+"</head>")

	beginScriptHeader := strings.Index(baseString, "<!-- Script starts here -->")
	endScriptHeader := strings.Index(baseString, "<!-- Script ends here -->")
	contentWithCards = strings.ReplaceAll(contentWithCards, "</body>", baseString[beginScriptHeader:endScriptHeader+len("<!-- Script ends here -->")]+"</body>")

	contentWithCards = strings.ReplaceAll(contentWithCards, metadata.ProxySite+req.URL.Path, req.URL.Path)
	contentWithCards = strings.ReplaceAll(contentWithCards, "utm_content=sitename", "utm_content="+metadata.SiteName)

	contact := GetContactInfo()
	contentWithCards = strings.Replace(contentWithCards, "</body>", contact+"</body>", 1)

	_, _ = io.WriteString(res, contentWithCards)
}

func GetContactInfo() string {
	contact := fmt.Sprintf(metadata.AdDescription+" <a href=\"%s\">(hop)</a> ", "https://www.showmycard.com")
	if metadata.Contact != "" {
		contact = contact + fmt.Sprintf(" Contact & Refunds <a href=\"%s\">(hop)</a> ", metadata.Contact)
	}
	if metadata.Terms != "" {
		contact = contact + fmt.Sprintf(" Privacy & Terms <a href=\"%s\">(hop)</a>", metadata.Terms)
	}
	contact = "<div style=\"text-align: center\">" + contact + "</div>"
	return contact
}

func ReportOnCard(writer http.ResponseWriter, request *http.Request) bool {
	privateKey := request.URL.Query().Get("apikey")
	apiKey := string(Get(privateKey))
	if request.Method == "GET" {
		form, _ := os.Open("./res/testreport.html")
		defer form.Close()
		buf := bytes.Buffer{}
		_, _ = io.Copy(&buf, form)
		target := GetTarget(englang.SGUID(apiKey))
		formString := Customize(buf.String())
		formString = strings.ReplaceAll(formString, "https://showmycard.com/7d5f7d3c5a885cead3102434c9becd8a", target)
		formString = strings.ReplaceAll(formString, "f7e77596ddf4b12e38e469421d78f3cc.png", "/png?apikey="+apiKey)
		begin := strings.Index(formString, "<!-- Data begins here -->")
		end := strings.Index(formString, "<!-- Data ends here -->")
		if begin != -1 && end != -1 {
			clicks, impressions := GetStatistics(englang.SGUID(apiKey))
			clickBar := clicks
			if clickBar > 60 {
				clickBar = 60
			}
			impressionBar := impressions
			if impressionBar > 60 {
				impressionBar = 60
			}
			pattern := `
				<rect class="impression" x="20" y="%d" width="10" height="%d" rx="2.5"></rect>
				<rect class="click" x="60" y="%d" width="10" height="%d" rx="2.5"></rect>
				<text class="impression" x="20" y="%d">%d</text>
				<text class="click" x="60" y="%d">%d</text>
				`
			formString = strings.ReplaceAll(formString, formString[begin:end+len("<!-- Data ends here -->")], fmt.Sprintf(pattern, 85-impressionBar, impressionBar, 85-clickBar, clickBar, 96, impressions, 96, clicks))
		}
		formString = strings.ReplaceAll(formString, "<!-- Location goes here -->", "<a href=\""+GetLocation(englang.SGUID(apiKey))+"\">Paid media</a>")
		buffered := bufio.NewWriter(writer)
		_, _ = buffered.WriteString(formString)
		_ = buffered.Flush()
		return true
	}
	return false
}
