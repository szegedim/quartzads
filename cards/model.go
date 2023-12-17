package cards

import (
	"fmt"
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

func NewCard(path string) metadata.SGUID {
	id := string(Get(path + ".path"))
	if id == "" {
		id = "showmycardcom" + string(metadata.GenerateSGUID())
		SetTarget(metadata.SGUID(id), "/ad?apikey="+id)
		Set(path+".path", []byte(id))
		png, _ := os.ReadFile("res/img.png")
		SetPicture(metadata.SGUID(id), png)
		SetLocation(metadata.SGUID(id), path)
	}
	return metadata.SGUID(id)
}

func DeleteCard(path string) {
	Set(path+".path", []byte(""))
}

func SetPicture(id metadata.SGUID, png []byte) {
	picId := string(id) + ".png"
	if len(png) < 10000000 {
		Set(picId, png)
	}
	// TODO replace with default
}

func GetPicture(id metadata.SGUID) (png []byte) {
	picId := string(id) + ".png"
	png = Get(picId)
	return
}

func SetTarget(id metadata.SGUID, url string) {
	Set(string(id), []byte(fmt.Sprintf("Card %s points to %s url.", id, url)))
}

func GetTarget(id metadata.SGUID) string {
	var target, text, id1 string
	text = string(Get(string(id)))
	n, _ := fmt.Sscanf(text, "Card %s points to %s url.", &id1, &target)
	if n < 2 {
		target = ""
	}
	return target
}

func SetLocation(id metadata.SGUID, url string) {
	Set(string(id)+".location", []byte(fmt.Sprintf("Card %s is at %s url.", id, url)))
}

func GetLocation(id metadata.SGUID) string {
	var location, text, id1 string
	text = string(Get(string(id) + ".location"))
	n, _ := fmt.Sscanf(text, "Card %s is at %s url.", &id1, &location)
	if n < 2 {
		location = ""
	}
	return location
}

func AddActivity(id metadata.SGUID, log string) {
	logId := string(id) + ".activity"
	Add(logId, []byte(log))
}

func GetActivities(id metadata.SGUID) (current string) {
	logId := string(id) + ".activity"
	current = string(Get(logId))
	return
}

func FindActivity(id metadata.SGUID, pattern string) (activity []string) {
	activities := GetActivities(id)
	patterns := strings.Split(pattern, "%s")
	index := 0
	activity = []string{}
	for index != -1 {
		index = strings.Index(activities[index:], patterns[0])
		if index == -1 {
			break
		}
		end := strings.Index(activities[index+1:], patterns[len(patterns)-1])
		activity = append(activity, activities[index:index+1+end+len(patterns[len(patterns)-1])])
		index = index + 1
	}
	return
}

func Parse(str string, pattern string) (items []string) {
	//func TestA(t *testing.T) {
	//	t.Log(Parse("The color is red today and green tomorrow.", "The color is %s today and %s tomorrow."))
	//	t.Log(Parse("The color is red today and tomorrow. blue", "The color is %s today and tomorrow. %s"))
	//	t.Log(Parse("yellow The color is red today and tomorrow.", "%s The color is %s today and tomorrow."))
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

func GetStatistics(id metadata.SGUID) (clicks, impressions int) {
	current := GetActivities(id)
	clicks = strings.Count(current, "Element clicked")
	impressions = strings.Count(current, "Element became visible")
	return
}
