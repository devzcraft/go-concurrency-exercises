//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync/atomic"
	"time"
)

const timeQuota = 10

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

func (u *User) addSecond() {
	u.TimeUsed += 1
	atomic.AddInt64(&u.TimeUsed, 1)
}

func (u *User) exceededQuota() bool {
	return u.TimeUsed >= timeQuota
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	t := time.Tick(time.Second)

	finished := make(chan struct{})

	go func() {
		process()
		finished <- struct{}{}
	}()

	for {
		select {
		case <-t:
			u.addSecond()
			if u.exceededQuota() {
				return false
			}
		case <-finished:
			return true

		}
	}

}

func main() {
	RunMockServer()
}
