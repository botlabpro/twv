package node

import "time"

type BaseNode struct {
	Ip          string `json:"ip"`
	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
	CityName    string `json:"cityName"`
}

type Stats struct {
	InboundIn   int64     `json:"inboundIn"`
	InboundOut  int64     `json:"inboundOut"`
	OutboundIn  int64     `json:"outboundIn"`
	OutboundOut int64     `json:"outboundOut"`
	StartAt     time.Time `json:"startAt"`
	Online      bool      `json:"online"`
}

type Node struct {
	BaseNode
	Stats
}
