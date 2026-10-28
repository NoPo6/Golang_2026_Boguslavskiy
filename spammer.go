package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
)

func RunPipeline(cmds ...cmd) {
	chansNum := len(cmds) + 1
	chans := make([]chan interface{}, chansNum)
	for i := range chans {
		chans[i] = make(chan interface{})
	}
	close(chans[0])

	var wg sync.WaitGroup
	for i := 0; i < len(cmds); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cmds[i](chans[i], chans[i+1])
			close(chans[i+1])
		}(i)
	}
	wg.Wait()
}

func SelectUsers(in, out chan interface{}) {
	// 	in - string
	// 	out - User
	idMapa := make(map[uint64]struct{})
	var mu sync.Mutex
	var wg sync.WaitGroup

	for v := range in {
		email, ok := v.(string)
		if !ok {
			log.Printf("bad email type %T", v)
			continue
		}
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			user := GetUser(email)
			mu.Lock()
			_, exists := idMapa[user.ID]
			idMapa[user.ID] = struct{}{}
			mu.Unlock()
			if !exists {
				out <- user
			}
		}(email)
	}
	wg.Wait()
}

func SelectMessages(in, out chan interface{}) {
	// 	in - User
	// 	out - MsgID
	var wg sync.WaitGroup

	runBatch := func(batch []User) {
		defer wg.Done()
		ids, err := GetMessages(batch...)
		if err != nil {
			log.Printf("GetMessages: %v", err)
			return
		}
		for _, id := range ids {
			out <- id
		}
	}

	userSlice := make([]User, 0, GetMessagesMaxUsersBatch)
	for v := range in {
		user, ok := v.(User)
		if !ok {
			log.Printf("bad user type %T", v)
			continue
		}
		userSlice = append(userSlice, user)
		if len(userSlice) == GetMessagesMaxUsersBatch {
			wg.Add(1)
			go runBatch(userSlice)
			userSlice = make([]User, 0, GetMessagesMaxUsersBatch)
		}
	}
	if len(userSlice) > 0 {
		wg.Add(1)
		go runBatch(userSlice)
	}

	wg.Wait()
}

func CheckSpam(in, out chan interface{}) {
	// in - MsgID
	// out - MsgData
	sem := make(chan struct{}, HasSpamMaxAsyncRequests)
	var wg sync.WaitGroup

	for v := range in {
		id, ok := v.(MsgID)
		if !ok {
			log.Printf("bad msgid type %T", v)
			continue
		}
		wg.Add(1)
		go func(id MsgID) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			spam, err := HasSpam(id)
			if err != nil {
				log.Printf("HasSpam(%d): %v", id, err)
				return
			}
			out <- MsgData{ID: id, HasSpam: spam}
		}(id)
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	// in - MsgData
	// out - string
	msgDataSlice := make([]MsgData, 0, 64)
	for v := range in {
		md, ok := v.(MsgData)
		if !ok {
			log.Printf("bad msgdata type %T", v)
			continue
		}
		msgDataSlice = append(msgDataSlice, md)
	}
	sort.Slice(msgDataSlice, func(i, j int) bool {
		if msgDataSlice[i].HasSpam != msgDataSlice[j].HasSpam {
			return msgDataSlice[i].HasSpam
		}
		return msgDataSlice[i].ID < msgDataSlice[j].ID
	})
	for _, msgData := range msgDataSlice {
		out <- fmt.Sprintf("%v %v", msgData.HasSpam, msgData.ID)
	}
}
