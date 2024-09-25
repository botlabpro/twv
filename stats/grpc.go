package stats

import (
	"context"
	"github.com/botlabpro/twv/node"
	api "github.com/sagernet/sing-box/experimental/v2rayapi"
	"google.golang.org/grpc"
	"time"
)

type Service struct {
	Node   *node.Node
	client api.StatsServiceClient
}

func Init(nodes []*node.Node) {
	for _, n := range nodes {
		s := &Service{
			Node: n,
		}
		conn, err := grpc.NewClient(n.Ip+":8080", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		s.client = api.NewStatsServiceClient(conn)
		go s.run()
	}

}

func (s *Service) run() {
	for {
		inboundOut := s.getStats("inbound>>>vppl_proxy>>>traffic>>>uplink")
		inboundIn := s.getStats("inbound>>>vppl_proxy>>>traffic>>>downlink")
		outboundOut := s.getStats("inbound>>>vppl>>>traffic>>>uplink") + s.getStats("inbound>>>vless>>>traffic>>>uplink")
		outboundIn := s.getStats("inbound>>>vppl>>>traffic>>>downlink") + s.getStats("inbound>>>vless>>>traffic>>>downlink")
		oldTraffic := s.Node.InboundIn + s.Node.InboundOut + s.Node.OutboundIn + s.Node.OutboundOut
		newTraffic := inboundIn + inboundOut + outboundIn + outboundOut
		if newTraffic > oldTraffic {
			s.Node.AddTraffic(newTraffic - oldTraffic)
		} else {
			s.Node.AddTraffic(0)
		}
		s.Node.InboundIn = inboundIn
		s.Node.InboundOut = inboundOut
		s.Node.OutboundIn = outboundIn
		s.Node.OutboundOut = outboundOut
		sysStats := s.getSysStats()
		if sysStats != nil {
			s.Node.Online = true
			s.Node.StartAt = time.Now().Add(time.Duration(sysStats.Uptime) * -1 * time.Second)
		} else {
			s.Node.Online = false
		}
		time.Sleep(time.Second)
	}
}

func (s *Service) getStats(name string) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	request := &api.GetStatsRequest{
		Name: name,
	}
	stats, err := s.client.GetStats(ctx, request)
	if err != nil {
		//fmt.Println(err.Error())
		return 0
	}
	return stats.Stat.Value
}

func (s *Service) getSysStats() *api.SysStatsResponse {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	request := &api.SysStatsRequest{}
	stats, err := s.client.GetSysStats(ctx, request)
	if err != nil {
		//fmt.Println(err.Error())
		return nil
	}
	return stats
}
