package main

import (
	"log"
)

type download struct {
	payloads []byte
}

func (t *download) run(data []byte) []byte {
	log.Printf("start downloading ...")
	var head [32]byte
	copy(head[:], data)
	t.iter(head)
	log.Printf("file downloaded: %d byts", len(t.payloads))
	return t.payloads
}

func (t *download) iter(data [32]byte) {
	p := create()
	log.Printf("downloading: %d bytes recieved", len(t.payloads)+1)
	p.get(data)
	t.payloads = append(t.payloads, p.getData()...)
	nl := p.getNL()
	nr := p.getNR()
	if nl != hashzero {
		t.iter(nl)
	}
	if nr != hashzero {
		t.iter(nr)
	}
}
