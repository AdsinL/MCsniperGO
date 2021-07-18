package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gookit/color"
	"github.com/kqzz/mcgo"
)

var accounts []*mcgo.MCaccount

func init() {
	color.Printf(genHeader())
	accStrs, err := readLines("accounts.txt")
	if err != nil {
		logFatal(err.Error())
	}

	accounts = loadAccSlice(accStrs)
}

func main() {
	if len(accounts) < 1 {
		logFatal("Please put one account in the accounts.txt file!")
	}

	if len(accounts) > 1 {
		logWarn("Using more than 1 account is not recommended")
	}

	targetName := userInput("target username")
	offsetStr := userInput("offset")
	offset, err := strconv.ParseFloat(offsetStr, 64)
	if err != nil {
		logFatal(fmt.Sprintf("%v is not a valid integer", offsetStr))
	}

	droptime, err := coolkidmachoDroptime(targetName)
	if err != nil {
		logFatal(err.Error())
	}

	logInfo(fmt.Sprintf("Sniping %v at %v", targetName, droptime.Format("2006/01/02 15:04:05")))

	time.Sleep(time.Until(droptime.Add(-time.Hour * 8))) // sleep until 8 hours before droptime

	for _, acc := range accounts {
		authErr := acc.MojangAuthenticate()
		if authErr != nil {
			logErr(fmt.Sprintf("Failed to authenticate %v, %v", acc.Email, err.Error()))
		} else {
			logSuccess(fmt.Sprintf("successfully authenticated %v", acc.Email))
		}
	}

	changeTime := droptime.Add(time.Millisecond * time.Duration(0-offset))

	var wg sync.WaitGroup

	var resps []mcgo.NameChangeReturn

	for _, acc := range accounts {
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				resp, err := acc.ChangeName(targetName, changeTime, false)
				if err != nil {
					logErr(fmt.Sprintf("encountered err on nc for %v: %v", acc.Email, err.Error()))
				} else {
					resps = append(resps, resp)
				}
			}()
		}
	}

	wg.Wait()

	for _, resp := range resps {
		logInfo(fmt.Sprintf("sent @ %v", resp.SendTime))
	}

	for _, resp := range resps {
		logInfo(fmt.Sprintf("[%v] recv @ %v", resp.StatusCode, resp.ReceiveTime))
	}

}