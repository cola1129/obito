package main

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/sideshow/apns2"
)

func startConsume(cert tls.Certificate) {
	var client *apns2.Client
	if *prod { // prod env
		client = apns2.NewClient(cert).Production()
	} else { // test env
		client = apns2.NewClient(cert).Development()
	}

	log.Printf("start to wait for notifications...")
	for i := range daemon {
		log.Printf("got new notification to push: %v", i)
		res, err := client.Push(i.apns)

		if err != nil {
			log.Printf("push(%v) failed(apns: %s, reason: %s), retry after 3 seconds...", i, res.ApnsID, res.Reason)
			time.Sleep(time.Second * time.Duration(*retryAfter))
			retryQueue <- i
		} else {
			log.Printf("%s success with reason: %s", res.ApnsID, res.Reason)
		}
	}
}
