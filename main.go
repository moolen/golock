package main

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func main() {
	X, err := xgb.NewConn()
	if err != nil {
		fmt.Println(err)
		return
	}
	screen := xproto.Setup(X).DefaultScreen(X)

	<-time.After(time.Millisecond * 100)
	grabc := xproto.GrabKeyboard(X, false, screen.Root, xproto.TimeCurrentTime,
		xproto.GrabModeAsync, xproto.GrabModeAsync,
	)
	repk, err := grabc.Reply()
	if err != nil {
		fmt.Printf("error grabbing Keyboard: %#v", err)
		os.Exit(1)
	}
	if repk.Status != xproto.GrabStatusSuccess {
		fmt.Printf("could not grab keyboard: %v", repk.Status)
		os.Exit(1)
	}
	grabp := xproto.GrabPointer(X, false, screen.Root, (xproto.EventMaskKeyPress|xproto.EventMaskKeyRelease)&0,
		xproto.GrabModeAsync, xproto.GrabModeAsync, xproto.WindowNone, xproto.CursorNone, xproto.TimeCurrentTime)
	repp, err := grabp.Reply()
	if err != nil {
		fmt.Printf("error grabbing pointer: %#v", err)
		os.Exit(1)
	}
	if repp.Status != xproto.GrabStatusSuccess {
		fmt.Printf("could not grab pointer: %v", repp.Status)
		os.Exit(1)
	}

	// todo: check for password
	for {
		ev, xerr := X.WaitForEvent()
		if ev == nil && xerr == nil {
			fmt.Println("Both event and error are nil. Exiting...")
			return
		}
		if ev != nil {
			fmt.Printf("Event: %v\n", ev)
			keyEvent, ok := ev.(xproto.KeyReleaseEvent)
			if ok {
				if keyEvent.Detail == 9 {
					fmt.Printf("ESC Key")
					return
				}
			}
		}
		if xerr != nil {
			fmt.Printf("Error: %s\n", xerr)
		}
	}
}
