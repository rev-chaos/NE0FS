package main

import (
	"log"
	"sync"
)

type download struct {
	buf       map[[32]byte]payload
	hashes    chan [32]byte
	payloads  chan payload
	wgHash    sync.WaitGroup
	wgPayload sync.WaitGroup
	ret       []byte
}

func (t *download) run(data []byte) []byte {
	var head [32]byte
	copy(head[:], data)
	t.hashes = make(chan [32]byte)
	t.payloads = make(chan payload)
	t.buf = make(map[[32]byte]payload)
	t.sendHash(head)

	log.Printf("start downloading ...")
	t.start()

	t.wgHash.Wait()
	close(t.hashes)
	t.wgPayload.Wait()
	close(t.payloads)

	t.iter(head)
	log.Printf("file downloaded: %d bytes", len(t.ret))
	return t.ret
}

func (t *download) sendHash(data [32]byte) {
	t.wgHash.Add(1)
	go func() {
		t.hashes <- data
	}()
}

func (t *download) sendPayloads(data payload) {
	t.wgPayload.Add(1)
	go func() {
		t.payloads <- data
	}()
}

func (t *download) start() {
	for i := 0; i < threads; i++ {
		go func() {
			for data := range t.hashes {
				p := create()
				err := p.get(data)
				if err != nil {
					t.sendHash(data)
				} else {
					t.sendPayloads(p)
					nl := p.getNL()
					nr := p.getNR()
					if nl != hashzero {
						t.sendHash(nl)
					}
					if nr != hashzero {
						t.sendHash(nr)
					}
				}
				t.wgHash.Done()
			}
		}()
	}

	go func() {
		for data := range t.payloads {
			t.buf[data.hash()] = data
			log.Printf("recieved %d packets", len(t.buf))
			t.wgPayload.Done()
		}
	}()
}

func (t *download) iter(data [32]byte) {
	p := t.buf[data]
	t.ret = append(t.ret, p.getData()...)
	nl := p.getNL()
	nr := p.getNR()
	if nl != hashzero {
		t.iter(nl)
	}
	if nr != hashzero {
		t.iter(nr)
	}
}
