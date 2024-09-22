package stats

import (
	"context"
	"github.com/botlab/twv/node"
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
		s.Node.InboundOut = s.getStats("inbound>>>vppl_proxy>>>traffic>>>uplink")
		s.Node.InboundIn = s.getStats("inbound>>>vppl_proxy>>>traffic>>>downlink")
		s.Node.OutboundOut = s.getStats("inbound>>>vppl>>>traffic>>>uplink") + s.getStats("inbound>>>vless>>>traffic>>>uplink")
		s.Node.OutboundIn = s.getStats("inbound>>>vppl>>>traffic>>>downlink") + s.getStats("inbound>>>vless>>>traffic>>>downlink")
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
