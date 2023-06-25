package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

func generateCertAndKey() (bytes.Buffer, bytes.Buffer, error) {
	var certPEM, privateKeyPEM bytes.Buffer

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return certPEM, privateKeyPEM, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return certPEM, privateKeyPEM, err
	}

	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return certPEM, privateKeyPEM, nil
}

func CreateCertAndKeyFiles() (string, string, error) {
	certPath := "localhost.pem"
	keyPath := "localhost.key"

	crtBytes, keyBytes, err := generateCertAndKey()
	if err != nil {
		return certPath, keyPath, err
	}

	err = os.WriteFile(certPath, crtBytes.Bytes(), 0644)
	if err != nil {
		return certPath, keyPath, err
	}
	err = os.WriteFile(keyPath, keyBytes.Bytes(), 0644)

	if err != nil {
		return certPath, keyPath, err
	}

	return certPath, keyPath, nil
}
