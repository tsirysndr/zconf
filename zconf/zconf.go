package zconf

import (
	"log"

	"github.com/grandcat/zeroconf"
	"github.com/tsirysndr/mdns"
	"crypto/md5"
	"encoding/hex"

	"github.com/lithammer/shortuuid"
)

func RegisterService(name, protocol, id string, port int) {
	if id == "" {
		hasher := md5.New()
		hasher.Write([]byte(shortuuid.New()))
		id = hex.EncodeToString(hasher.Sum(nil)) 
	}
	log.Println("id=" + id)
	meta := []string{
		"txtv=0",
		"lo=1",
		"la=2",
		"id=" + id,
		"fn=" + name,
	}
	_, err := zeroconf.Register(
		name,
		protocol,
		"local.",
		port,
		meta,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	for {
	}

}

func ListServices(protocol string, limit int) []*mdns.ServiceEntry {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, limit)
	go func() {
		mdns.Query(&mdns.QueryParam{
			Service: protocol,
			Domain:  "local",
			// Timeout: time.Second * 1,
			Entries: entriesCh,
		})
		close(entriesCh)
	}()

	entries := []*mdns.ServiceEntry{}
	for entry := range entriesCh {
		entries = append(entries, entry)
	}

	return entries
}
