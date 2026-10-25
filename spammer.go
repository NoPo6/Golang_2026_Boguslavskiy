package main

import (
	"fmt"
	"sort"
	"sync"
)

func RunPipeline(cmds ...cmd) {
	chansNum := len(cmds) + 1
	chans := make([]chan interface{}, chansNum)
	for i := range chans {
		chans[i] = make(chan interface{})
	}

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
	users := make(chan User)
	var wg sync.WaitGroup
	for email := range in {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			users <- GetUser(email)
		}(email.(string))
	}

	go func() {
		wg.Wait()
		close(users)
	}()

	idMapa := make(map[uint64]struct{})
	for user := range users {
		if _, ok := idMapa[user.ID]; !ok {
			idMapa[user.ID] = struct{}{}
			out <- user
		}
	}
}

func SelectMessages(in, out chan interface{}) {
	// 	in - User
	// 	out - MsgID
	msgIDs := make(chan MsgID)
	var wg sync.WaitGroup

	runBatch := func(batch []User) {
		defer wg.Done()
		ids, _ := GetMessages(batch...)
		for _, id := range ids {
			msgIDs <- id
		}
	}

	userSlice := []User{}
	for user := range in {
		userSlice = append(userSlice, user.(User))
		if len(userSlice) == 2 {
			wg.Add(1)
			batch := userSlice
			go runBatch(batch)
			userSlice = []User{}
		}
	}
	if len(userSlice) > 0 {
		wg.Add(1)
		go runBatch(userSlice)
	}

	go func() {
		wg.Wait()
		close(msgIDs)
	}()

	for id := range msgIDs {
		out <- id
	}
}

func CheckSpam(in, out chan interface{}) {
	// in - MsgID
	// out - MsgData
	msgData := make(chan MsgData)
	sem := make(chan struct{}, HasSpamMaxAsyncRequests)
	var wg sync.WaitGroup

	for msgID := range in {
		wg.Add(1)
		go func(id MsgID) {
			defer wg.Done()
			sem <- struct{}{}
			spam, _ := HasSpam(id)
			<-sem
			msgData <- MsgData{ID: id, HasSpam: spam}
		}(msgID.(MsgID))
	}

	go func() {
		wg.Wait()
		close(msgData)
	}()

	for md := range msgData {
		out <- md
	}
}

func CombineResults(in, out chan interface{}) {
	// in - MsgData
	// out - string
	msgDataSlice := make([]MsgData, 0, 64)
	for msgData := range in {
		msgDataSlice = append(msgDataSlice, msgData.(MsgData))
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
