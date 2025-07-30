package crypto_test

import (
	"crypto/x509"
	"fmt"
	"net"
	"testing"

	"github.com/wsva/lib_go/crypto"
)

func TestCA(T *testing.T) {
	rootconfig := crypto.CertConfig{
		CertConfigBase: crypto.CertConfigBase{
			CommonName:   "OW",
			Organization: []string{"OW"},
		},
		PublicKeyAlgorithm: x509.ECDSA,
	}
	rootcrt, rootkey, err := crypto.NewCertificateAuthority(&rootconfig)
	if err != nil {
		fmt.Println("root", err)
		return
	}
	err = crypto.WriteCertAndKey("certs", "GMKarRoot", rootcrt, rootkey)
	if err != nil {
		fmt.Println("root write", err)
		return
	}
	ccconfig := crypto.CertConfig{
		CertConfigBase: crypto.CertConfigBase{
			CommonName: "10.0.0.1",
			AltNames: crypto.AltNames{
				IPs: []net.IP{net.ParseIP("10.0.0.1")},
			},
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		PublicKeyAlgorithm: x509.ECDSA,
	}
	cccrt, cckey, err := crypto.NewCertAndKey(rootcrt, rootkey, &ccconfig)
	if err != nil {
		fmt.Println("control center", err)
		return
	}
	err = crypto.WriteCertAndKey("./certs", "GMKarControlCenter", cccrt, cckey)
	if err != nil {
		fmt.Println("control center write", err)
		return
	}
	clientconfig := crypto.CertConfig{
		CertConfigBase: crypto.CertConfigBase{
			CommonName: "ow",
			AltNames: crypto.AltNames{
				IPs: []net.IP{net.ParseIP("10.0.0.1")},
			},
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		PublicKeyAlgorithm: x509.ECDSA,
	}
	clientcrt, clientkey, err := crypto.NewCertAndKey(rootcrt, rootkey, &clientconfig)
	if err != nil {
		fmt.Println("control center", err)
		return
	}
	err = crypto.WriteCertAndKey("./certs", "Client10.0.0.1", clientcrt, clientkey)
	if err != nil {
		fmt.Println("control center write", err)
		return
	}
	fmt.Println("success")
}
