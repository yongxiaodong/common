package golib

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

type CertBaseInfo struct {
	CertPath   string
	KeyPath    string
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func (base CertBaseInfo) verify_validity() {

}

func (base CertBaseInfo) verify_encrypt() {
	// 使用公钥加密
	originalData := "data to be encrypted"
	data := []byte(originalData)
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, base.PublicKey, data, nil)
	if err != nil {
		log.Fatalf("Failed to encrypt data: %v", err)
	}
	// 使用私钥解密
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, base.PrivateKey, ciphertext, nil)
	if err != nil {
		log.Fatalf("Failed to decrypt data: %v", err)
	}
	if data[0] == plaintext[0] {
		fmt.Println("encrypted pass")

	}
	//fmt.Printf("Original data: %s\n", data)
	//fmt.Printf("Decrypted data: %s\n", plaintext)
}

func (base *CertBaseInfo) Verify_tls() {

	// 读取证书文件
	certBytes, err := os.ReadFile(base.CertPath)
	if err != nil {
		log.Fatalf("Failed to read certificate file: %v", err)
	}

	// 解码 PEM 格式的证书
	block, _ := pem.Decode(certBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		log.Fatal("Failed to decode certificate PEM")
	}

	// 解析证书
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse certificate: %v", err)
	}
	base.PublicKey = cert.PublicKey.(*rsa.PublicKey)

	// 打印证书有效期
	fmt.Printf("Certificate Validity:\n")
	fmt.Printf("  - Not Before: %v\n", cert.NotBefore)
	fmt.Printf("  - Not After: %v\n", cert.NotAfter)

	// 读取私钥文件
	keyBytes, err := os.ReadFile(base.KeyPath)
	if err != nil {
		log.Fatalf("Failed to read private key file: %v", err)
	}

	// 解码 PEM 格式的私钥
	block, _ = pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		log.Fatal("Failed to decode private key PEM")
	}
	base.PrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)

	// 尝试解析私钥（不会验证其有效性，仅检查格式）
	_, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	fmt.Println("Private key is valid.")
	base.verify_encrypt()

}
