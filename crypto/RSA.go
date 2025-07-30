package crypto

import (
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path"
)

func NewRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privkey, err := rsa.GenerateKey(cryptorand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}
	return privkey, &privkey.PublicKey, nil
}

func WriteRSAKeyPair(pkiPath, name string, privkey *rsa.PrivateKey, pubkey *rsa.PublicKey) error {
	err := WriteRSAPublicKey(pkiPath, name, pubkey)
	if err != nil {
		return err
	}
	err = WriteRSAPrivateKey(pkiPath, name, privkey)
	if err != nil {
		return err
	}
	return nil
}

func WriteRSAPublicKey(pkiPath, name string, pubkey *rsa.PublicKey) error {
	publicFile, err := os.Create(path.Join(pkiPath, name+"RSAPublic.pem"))
	if err != nil {
		return err
	}
	defer publicFile.Close()
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return err
	}
	pubkeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubkeyBytes,
	}
	err = pem.Encode(publicFile, pubkeyBlock)
	if err != nil {
		return err
	}
	return nil
}

func WriteRSAPrivateKey(pkiPath, name string, privkey *rsa.PrivateKey) error {
	privateFile, err := os.Create(path.Join(pkiPath, name+"RSAPrivate.pem"))
	if err != nil {
		return err
	}
	defer privateFile.Close()
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privkeyBytes,
	}
	err = pem.Encode(privateFile, privkeyBlock)
	if err != nil {
		return err
	}
	return nil
}

func TryLoadRSAKeyPairFromDisk(pkiPath, name string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	rsapubkey, err := TryLoadRSAPublicKeyFromDisk(pkiPath, name)
	if err != nil {
		return nil, nil, err
	}
	rsaprivkey, err := TryLoadRSAPrivateKeyFromDisk(pkiPath, name)
	if err != nil {
		return nil, nil, err
	}

	return rsaprivkey, rsapubkey, nil
}

func TryLoadRSAPublicKeyFromDisk(pkiPath, name string) (*rsa.PublicKey, error) {
	contentBytes, err := os.ReadFile(path.Join(pkiPath, name+"RSAPublic.pem"))
	if err != nil {
		return nil, err
	}
	return TryLoadRSAPublicKeyFromContent(contentBytes)
}

func TryLoadRSAPrivateKeyFromDisk(pkiPath, name string) (*rsa.PrivateKey, error) {
	contentBytes, err := os.ReadFile(path.Join(pkiPath, name+"RSAPrivate.pem"))
	if err != nil {
		return nil, err
	}
	data, rest := pem.Decode(contentBytes)
	if len(rest) > 0 {
		return nil, errors.New("remainder of content found")
	}
	return x509.ParsePKCS1PrivateKey(data.Bytes)
}

func TryLoadRSAPublicKeyFromContent(contentBytes []byte) (*rsa.PublicKey, error) {
	data, rest := pem.Decode([]byte(contentBytes))
	if len(rest) > 0 {
		return nil, errors.New("remainder of content found")
	}
	pubkey, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		return nil, err
	}
	switch rsapubkey := pubkey.(type) {
	case *rsa.PublicKey:
		return rsapubkey, nil
	default:
		return nil, errors.New("public key type is not rsa")
	}
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	return string(privkeyPEM)
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)
	return string(pubkeyPEM), nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	pubkey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch rsapubkey := pubkey.(type) {
	case *rsa.PublicKey:
		return rsapubkey, nil
	default:
		return nil, errors.New("key type is not rsa")
	}
}

func RSAEncrypt(text []byte, pubkey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), cryptorand.Reader, pubkey, text, nil)
}

func RSADecrypt(ctext []byte, privkey *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), cryptorand.Reader, privkey, ctext, nil)
}
