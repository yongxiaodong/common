package golib

import (
	"crypto"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
	"github.com/go-acme/lego/v4/registration"
	"log"
	"os"
	"time"
)

var (
	storage_path = "../cert_storage/"
)

type MyUser struct {
	Email        string
	Registration *registration.Resource
	Key          crypto.PrivateKey
	Domains      []string
}

type AliKey struct {
	ALICLOUD_ACCESS_KEY string
	ALICLOUD_SECRET_KEY string
}

func (u *MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}

func (myUser MyUser) Obtain_cert(aliKey AliKey) {
	//privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//myUser = MyUser{
	//	Email: "you@yours.com",
	//	key:   privateKey,
	//}

	config := lego.NewConfig(&myUser)

	// This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	config.CADirURL = "https://acme-v02.api.letsencrypt.org/directory"
	config.Certificate.KeyType = certcrypto.RSA2048

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Setenv("ALICLOUD_ACCESS_KEY", aliKey.ALICLOUD_ACCESS_KEY)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("ALICLOUD_SECRET_KEY", aliKey.ALICLOUD_SECRET_KEY)
	if err != nil {
		panic(err)
	}

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Fatal(err)
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: myUser.Domains,
		Bundle:  true,
	}

	dnsprovider, err := alidns.NewDNSProvider()
	if err != nil {
		log.Fatalf("Error creating AliDNS provider: %v", err)
	}
	err = client.Challenge.SetDNS01Provider(dnsprovider, dns01.AddDNSTimeout(6*time.Minute))
	if err != nil {
		panic(err)
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Fatal(err)
	}
	//filename := strings.Split(myUser.Domains[0], ".")[-1]

	err = os.WriteFile(storage_path+myUser.Domains[0]+".crt", certificates.Certificate, 0755)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(storage_path+myUser.Domains[0]+".key", certificates.PrivateKey, 0755)
	if err != nil {
		panic(err)
	}

}
