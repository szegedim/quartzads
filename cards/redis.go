package cards

import (
	"bytes"
	"sync"
)

//Licensed under Creative Commons CC0.
//
//To the extent possible under law, the author(s) have dedicated all copyright and related and
//neighboring rights to this software to the public domain worldwide.
//This software is distributed without any warranty.
//You should have received a copy of the CC0 Public Domain Dedication along with this software.
//If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

// TODO Add actual redis stub
// TODO Isolated security locking SGUIDs

var lock sync.Mutex

func singleton() {
	//go func() {
	//	url := os.Getenv("REDISURL")
	//	if url != "" {
	//		opt, err := redis0.ParseURL(url)
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		client := redis0.NewClient(opt)
	//
	//		ctx := context.Background()
	//		val, _ := client.Get(ctx, "showmycard.com").Result()
	//		fmt.Println(val)
	//	}
	//}()
}

func Set(key string, value []byte) {
	singleton()
	lock.Lock()
	defer lock.Unlock()

	redis[key] = value
}

func Get(key string) (value []byte) {
	singleton()
	lock.Lock()
	defer lock.Unlock()

	value, ok := redis[key]
	if !ok {
		value = make([]byte, 0)
	}
	return
}

func List() (keys string) {
	singleton()
	lock.Lock()
	defer lock.Unlock()

	b := bytes.Buffer{}
	for k, _ := range redis {
		b.WriteString(k + "\n")
	}
	keys = b.String()
	return
}
