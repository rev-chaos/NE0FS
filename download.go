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
				for {
					err := func() error {
						channel := make(chan payload, ntry)
						for i := 0; i < ntry; i++ {
							go func() {
								defer func() {
									recover()
								}()
								p := create()
								p.get(data)
								channel <- p
							}()
						}

						defer func() {
							close(channel)
						}()

						for i := 0; i < ntry; i++ {
							p := <-channel
							bin := p.getData()
							if len(bin) == 0 {
								continue
							}

							t.sendPayloads(p)
							nl := p.getNL()
							nr := p.getNR()
							if nl != hashzero {
								t.sendHash(nl)
							}
							if nr != hashzero {
								t.sendHash(nr)
							}
							return nil
						}

						return errRPC{}
					}()
					if err == nil {
						break
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
