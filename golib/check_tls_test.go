package golib

import "testing"

func TestCertBaseInfo_Verify_tls(t *testing.T) {
	base := CertBaseInfo{
		CertPath:   "../cert_storage/buydance.vip.crt",
		KeyPath:    "../cert_storage/buydance.vip.key",
		PublicKey:  nil,
		PrivateKey: nil,
	}
	base.Verify_tls()

}
