package socket

import (
	"fmt"
	tp "github.com/henrylee2cn/teleport"
	"log"
	"time"
)

//"github.com/google/gopacket"
//" net配合goroutine+channel 就可以了"""
func TcpMain() error {
	// graceful

	go tp.GraceSignal()

	// server peer
	srv := tp.NewPeer(tp.PeerConfig{
		CountTime:   true,
		ListenPort:  9090,
		PrintDetail: true,
	})

	// router
	srv.RouteCall(new(Math))
	srv.RouteCall(new(User))
	// broadcast per 5s
	go func() {
		for {
			time.Sleep(time.Second * 5)
			srv.RangeSession(func(sess tp.Session) bool {
				sess.Push(
					"/push/status",
					fmt.Sprintf("this is a broadcast, server time: %v\n", time.Now()),
				)
				return true
			})
		}
	}()

	// listen and serve
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

// Math handler
type Math struct {
	tp.CallCtx
}

// Math handler
type User struct {
	tp.CallCtx
}

func (u *User) Login(id *string) (bool, *tp.Rerror) {
	fmt.Printf("proto login id : %v\n", *id)
	return true, nil
}

// Add handles addition request
func (m *Math) Add(arg *[]int) (int, *tp.Rerror) {
	// test query parameter
	//tp.Infof("author: %s", m.Query().Get("author"))
	fmt.Printf("Add() calling\n")
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}

/////////////////////////////client
func client() {
	// proto level
	tp.SetLoggerLevel("ERROR")

	cli := tp.NewPeer(tp.PeerConfig{})
	defer func() {
		err := cli.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	cli.RoutePush(new(Push))

	sess, err := cli.Dial(":9090")
	if err != nil {
		tp.Fatalf("%v", err)
	}

	var result int
	rerr := sess.Call("/math/add",
		[]int{1, 2, 3, 4, 5},
		&result,
	).Rerror()
	if rerr != nil {
		tp.Fatalf("%v", rerr)
	}
	tp.Printf("result: %d", result)

	var result2 bool
	rerr2 := sess.Call("/proto/login",
		"fewfwe",
		&result2,
	).Rerror()
	if rerr2 != nil {
		tp.Fatalf("%v", rerr2)
	}

	tp.Printf("result2: %v", result2)

	tp.Printf("wait for 10s...")
	time.Sleep(time.Second * 10)
}

// Push push handler
type Push struct {
	tp.PushCtx
}

// Push handles '/push/status' message
func (p *Push) Status(arg *string) *tp.Rerror {
	tp.Printf("%s", *arg)
	return nil
}
