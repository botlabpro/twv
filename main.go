package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/botlabpro/twv/common"
	"github.com/botlabpro/twv/node"
	"github.com/botlabpro/twv/stats"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"time"
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

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"ORGANIZATION_NAME"},
			Country:       []string{"COUNTRY_CODE"},
			Province:      []string{"PROVINCE"},
			Locality:      []string{"CITY"},
			StreetAddress: []string{"ADDRESS"},
			PostalCode:    []string{"POSTAL_CODE"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey

	// Sign the certificate
	certificate, _ := x509.CreateCertificate(rand.Reader, cert, cert, pub, priv)

	certBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate})
	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	// Generate a key pair from your pem-encoded cert and key ([]byte).
	x509Cert, _ := tls.X509KeyPair(certBytes, keyBytes)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{x509Cert}}
	server := http.Server{Addr: ":3000", Handler: router, TLSConfig: tlsConfig}

	if err := server.ListenAndServeTLS("", ""); err != nil {
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
