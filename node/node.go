package node

import "time"

type BaseNode struct {
	Ip          string `json:"ip"`
	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
	CityName    string `json:"cityName"`
}

type Stats struct {
	InboundIn    int64     `json:"inboundIn"`
	InboundOut   int64     `json:"inboundOut"`
	OutboundIn   int64     `json:"outboundIn"`
	OutboundOut  int64     `json:"outboundOut"`
	StartAt      time.Time `json:"startAt"`
	Online       bool      `json:"online"`
	TrafficBySec []int64   `json:"trafficBySec"`
}

type Node struct {
	BaseNode
	Stats
}

func (n *Node) AddTraffic(b int64) {
	if len(n.TrafficBySec) >= 60 {
		n.TrafficBySec = n.TrafficBySec[1:]
	}
	n.TrafficBySec = append(n.TrafficBySec, b)
}
