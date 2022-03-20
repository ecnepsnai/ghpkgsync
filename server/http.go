package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"

	"github.com/ecnepsnai/web/router"
)

func startHTTPS(server *router.Server) error {
	var pKey crypto.PrivateKey
	var err error
	pKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return fmt.Errorf("crypto error: %s", err.Error())
	}

	pub := pKey.(crypto.Signer).Public()
	serial := big.NewInt(1)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return fmt.Errorf("crypto error: %s", err.Error())
	}
	h := sha1.Sum(publicKeyBytes)

	tpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: ""},
		NotBefore:             time.Unix(0, 0),
		NotAfter:              time.Now().UTC().AddDate(100, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment,
		BasicConstraintsValid: true,
		SubjectKeyId:          h[:],
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, tpl, tpl, pub, pKey)
	if err != nil {
		return fmt.Errorf("crypto error: %s", err.Error())
	}

	l, err := tls.Listen("tcp", "0.0.0.0:443", &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{certBytes},
				PrivateKey:  pKey,
			},
		},
	})
	log.PDebug("Listen", map[string]interface{}{
		"address": "0.0.0.0:443",
	})
	if err != nil {
		return fmt.Errorf("listen error: %s", err.Error())
	}
	return server.Serve(l)
}

func startHTTP(server *router.Server) error {
	return server.ListenAndServe("0.0.0.0:80")
}
