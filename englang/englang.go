package englang

import "strings"

//Licensed under Creative Commons CC0.
//
//To the extent possible under law, the author(s) have dedicated all copyright and related and
//neighboring rights to this software to the public domain worldwide.
//This software is distributed without any warranty.
//You should have received a copy of the CC0 Public Domain Dedication along with this software.
//If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

func Englang(str string, pattern string) (items []string) {
	//func TestA(t *testing.T) {
	//	t.Log(Englang("The color is red today and green tomorrow.", "The color is %s today and %s tomorrow."))
	//	t.Log(Englang("The color is red today and tomorrow. blue", "The color is %s today and tomorrow. %s"))
	//	t.Log(Englang("yellow The color is red today and tomorrow.", "%s The color is %s today and tomorrow."))
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
