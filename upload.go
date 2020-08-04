package main

import (
	"log"
	"sync"
)

type upload struct {
	payloads chan payload
}

func (t *upload) run(data []byte) []byte {
	log.Printf("start uploading ...")
	log.Printf("data size: %d bytes", len(data))

	t.payloads = make(chan payload, len(data)/940+0xf)
	ret := t.iter(data)
	close(t.payloads)

	log.Printf("packets built: %d packets", len(t.payloads))

	var wg sync.WaitGroup
	wg.Add(threads)
	defer wg.Wait()

	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Done()
			for v := range t.payloads {
				for {
					log.Printf("packet left %d ...", len(t.payloads))
					ret := func() error {
						channel := make(chan error, ntry)
						for i := 0; i < ntry; i++ {
							go func() {
								defer func() {
									recover()
								}()
								channel <- v.put()
							}()
						}

						defer func() {
							close(channel)
						}()

						for i := 0; i < ntry; i++ {
							err := <-channel
							if err == nil {
								return nil
							}
						}
						return errRPC{}
					}()
					if ret == nil {
						break
					}
				}
			}
		}()
	}

	return ret[:]
}

func (t *upload) iter(data []byte) [32]byte {
	p := create()
	size := len(data)
	if size > 940 {
		size = 940
	}
	p.setData(data[:size])
	if len(data) > 940 {
		extra := data[940:]

		size := len(extra)
		size = size / 2
		size = size + 939
		size /= 940
		size *= 940

		dl := extra[:size]
		nl := t.iter(dl)
		p.setNL(nl)

		if size < len(extra) {
			dr := extra[size:]
			nr := t.iter(dr)
			p.setNR(nr)
		}
	}
	t.payloads <- p
	return p.hash()
}
