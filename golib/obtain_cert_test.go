package golib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func TestMyUser_Obtain_cert(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	user := MyUser{
		Email:   "your@mail.com",
		Key:     privateKey,
		Domains: []string{"buydance.vip", "*.buydance.vip"},
	}
	key := AliKey{
		ALICLOUD_ACCESS_KEY: "yourkey",
		ALICLOUD_SECRET_KEY: "yoursecret",
	}
	user.Obtain_cert(key)

}
