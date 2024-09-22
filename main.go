package main

import (
	"github.com/botlabpro/twv/common"
	"github.com/botlabpro/twv/node"
	"github.com/botlabpro/twv/stats"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config := common.LoadConfig("./config.json")
	var allNodes []*node.Node
	for _, n := range config.Nodes {
		allNodes = append(allNodes, &node.Node{BaseNode: *n})
	}

	stats.Init(allNodes)

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/nodes", func(context *gin.Context) {
		context.JSON(200, allNodes)
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
	/*for {
		var info []string
		for _, n := range allNodes {
			info = append(info, fmt.Sprintf("Node %s InIn:%d InOut:%d OutIn:%d OutOut:%d", n.Ip, n.InboundIn, n.InboundOut, n.OutboundIn, n.OutboundOut))
		}
		fmt.Println(strings.Join(info, " "))
		time.Sleep(time.Second)
	}*/
}
