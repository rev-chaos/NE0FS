package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type detectnode struct {
	inNodes  []string
	outNodes []string
	tested   map[string]interface{}
}

func (t *detectnode) run(data []byte) []byte {
	t.inNodes = []string{
		"seed1.ngd.network:10332",
		"seed2.ngd.network:10332",
		"seed3.ngd.network:10332",
		"seed4.ngd.network:10332",
		"seed5.ngd.network:10332",
		"seed6.ngd.network:10332",
		"seed7.ngd.network:10332",
		"seed8.ngd.network:10332",
		"seed9.ngd.network:10332",
		"seed10.ngd.network:10332",
	}
	t.tested = make(map[string]interface{})
	for _, v := range t.inNodes {
		t.tested[v] = true
	}
	http.DefaultClient.Timeout = time.Duration(latency) * time.Millisecond

	for {
		if len(t.inNodes) == 0 {
			break
		}

		log.Printf("current state: innodes = %d, outnotes = %d", len(t.inNodes), len(t.outNodes))

		address := t.inNodes[0]
		t.inNodes = t.inNodes[1:]

		values := url.Values{
			"jsonrpc": []string{"2.0"},
			"method":  []string{"getpeers"},
			"params":  []string{"[]"},
			"id":      []string{"1"},
		}
		addr := url.URL{
			Scheme:   "http",
			Host:     address,
			RawQuery: values.Encode(),
		}

		resp, err := http.Get(addr.String())
		if err != nil {
			continue
		}

		defer resp.Body.Close()

		var ret struct {
			Result struct {
				Connected []struct {
					Address string `json:"address"`
				} `json:"connected"`
				Unconnected []struct {
					Address string `json:"address"`
				} `json:"unconnected"`
			} `json:"result"`
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&ret)
		if err != nil {
			continue
		}

		t.outNodes = append(t.outNodes, address)
		log.Printf("detect node alive: %s", address)

		for _, v := range ret.Result.Connected {
			address := v.Address + ":10332"
			if t.tested[address] == true {
				continue
			}
			t.inNodes = append(t.inNodes, address)
			t.tested[address] = true
		}

		for _, v := range ret.Result.Unconnected {
			address := v.Address + ":10332"
			if t.tested[address] == true {
				continue
			}
			t.inNodes = append(t.inNodes, address)
			t.tested[address] = true
		}
	}
	ret, err := json.Marshal(t.outNodes)
	if err != nil {
		return nil
	}
	return ret
}
