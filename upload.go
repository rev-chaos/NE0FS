package main

import (
	"log"
)

type upload struct {
	payloads []payload
}

func (t *upload) run(data []byte) []byte {
	log.Printf("start uploading ...")
	log.Printf("data size: %d bytes", len(data))
	ret := t.iter(data)
	log.Printf("packets built: %d packets", len(t.payloads))
	for i, v := range t.payloads {
		for {
			log.Printf("sending packet %d/%d ...", i+1, len(t.payloads))
			err := v.put()
			if err == nil {
				break
			}
			log.Printf("sending packet %d/%d error: %s", i+1, len(t.payloads), err.Error())
		}
	}
	return ret[:]
}

func (t *upload) iter(data []byte) [32]byte {
	p := create()
	p.setData(data[:940])
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
	t.payloads = append(t.payloads, p)
	return p.hash()
}
