package cards

import (
	"fmt"
	"gitlab.com/eper.io/quartzads/englang"
	"gitlab.com/eper.io/quartzads/metadata"
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

func NewCard(path string) englang.SGUID {
	id := string(Get(path + ".path"))
	if id == "" {
		id = "public" + string(englang.GenerateSGUID())
		SetTarget(englang.SGUID(id), "/ad?apikey="+id)
		Set(path+".path", []byte(id))
		png, _ := os.ReadFile("res/img.png")
		SetPicture(englang.SGUID(id), png)
		SetLocation(englang.SGUID(id), path)
	}
	return englang.SGUID(id)
}

func DeleteCard(path string) {
	Set(path+".path", []byte(""))
}

func LockAd(apiKey string) {
	if strings.HasPrefix(GetTarget(englang.SGUID(apiKey)), "/ad") {
		expiry := time.Now().Add(metadata.DefaultPurchaseTime)
		png, _ := os.ReadFile("res/beingedited.png")
		SetPicture(englang.SGUID(apiKey), png)
		SetTarget(englang.SGUID(apiKey), ".")
		AddActivity(englang.SGUID(apiKey), fmt.Sprintf("Card will expire on %s and revert to ad.", expiry.Format(time.RFC822Z)))
	}
}

func GetCard(path string, i int) (englang.SGUID, string) {
	name := fmt.Sprintf("card%04d", i)
	cardId := NewCard(path + "#" + name)
	expiries := FindActivity(cardId, "Card will expire on %s and revert to ad.")
	expiryLog := ""
	for _, item := range expiries {
		expiry := englang.SplitEnglang(item, "Card will expire on %s and revert to ad.")
		if len(expiry) > 0 {
			t, err := time.Parse(time.RFC822Z, expiry[0])
			if err == nil {
				if t.After(time.Now()) {
					expiryLog = fmt.Sprintf("<!--Expiry of %s is in %s seconds.-->\n", cardId, t.Sub(time.Now()))
				} else {
					expiryLog = fmt.Sprintf("<!--Expiry of %s was in %s seconds.-->\n", cardId, time.Now().Sub(t))
					AddActivity(cardId, string(cardId)+": Card expired.")
					DeleteCard(path + "#" + name)
					cardId = NewCard(path + "#" + name)
				}
			}
		}
	}
	return cardId, expiryLog
}

func SetPicture(id englang.SGUID, png []byte) {
	picId := string(id) + ".png"
	if len(png) < 10000000 {
		Set(picId, png)
	}
	// TODO replace with default
}

func GetPicture(id englang.SGUID) (png []byte) {
	picId := string(id) + ".png"
	png = Get(picId)
	return
}

func SetTarget(id englang.SGUID, url string) {
	Set(string(id), []byte(fmt.Sprintf("Card %s points to %s url.", id, url)))
}

func GetTarget(id englang.SGUID) string {
	var target, text, id1 string
	text = string(Get(string(id)))
	n, _ := fmt.Sscanf(text, "Card %s points to %s url.", &id1, &target)
	if n < 2 {
		target = ""
	}
	return target
}

func SetLocation(id englang.SGUID, url string) {
	Set(string(id)+".location", []byte(fmt.Sprintf("Card %s is at %s url.", id, url)))
}

func GetLocation(id englang.SGUID) string {
	var location, text, id1 string
	text = string(Get(string(id) + ".location"))
	n, _ := fmt.Sscanf(text, "Card %s is at %s url.", &id1, &location)
	if n < 2 {
		location = ""
	}
	return location
}

func AddActivity(id englang.SGUID, log string) {
	logId := string(id) + ".activity"
	Add(logId, []byte(log))
}

func GetActivities(id englang.SGUID) (current string) {
	logId := string(id) + ".activity"
	current = string(Get(logId))
	return
}

func SetActivities(id englang.SGUID, fullLog string) {
	logId := string(id) + ".activity"
	// TODO lock
	Set(logId, []byte(fullLog))
}

func FindActivity(id englang.SGUID, pattern string) (activity []string) {
	activities := GetActivities(id)
	patterns := strings.Split(pattern, "%s")
	index := 0
	activity = []string{}
	for index != -1 {
		index = strings.LastIndex(activities[index:], patterns[0])
		if index == -1 {
			break
		}
		end := strings.Index(activities[index+1:], patterns[len(patterns)-1])
		activity = append(activity, activities[index:index+1+end+len(patterns[len(patterns)-1])])
		index = index + 1
	}
	return
}

func InvalidateStatistics(id englang.SGUID) {
	current := GetActivities(id)
	current = strings.ReplaceAll(current, "Element clicked", "Purchased Element clicked")
	current = strings.ReplaceAll(current, "Element became visible", "Purchased Element clicked")
	SetActivities(id, current)
}

func GetStatistics(id englang.SGUID) (clicks, impressions int) {
	current := GetActivities(id)
	clicks = strings.Count(current, string(id)+": Element clicked")
	impressions = strings.Count(current, string(id)+": Element became visible")
	return
}
