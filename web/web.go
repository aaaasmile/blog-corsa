package web

import (
	"context"
	"corsa-blog/conf"
	"corsa-blog/idl"
	"corsa-blog/mail"
	"corsa-blog/telegram"
	"corsa-blog/util"
	"corsa-blog/web/app"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/crewjam/saml/samlsp"
)

type ObserveComments struct {
	simulation bool
	debug      bool
	newCmt     chan idl.CmtItem
}

func RunService(configfile string, simulate bool) error {

	if _, err := conf.ReadConfig(configfile); err != nil {
		return err
	}
	serverurl := conf.Current.ServiceURL
	serverurl = strings.Replace(serverurl, "0.0.0.0", "localhost", 1)
	serverurl = strings.Replace(serverurl, "127.0.0.1", "localhost", 1)
	dashboardServURL := fmt.Sprintf("http://%s%s", serverurl, conf.Current.RootURLPattern)
	log.Println("Server started with URL ", serverurl)
	log.Println("Try this url for Dashboard: ", dashboardServURL)

	// Dashboard app
	//http.HandleFunc(conf.Current.RootURLPattern, app.APiHandler)
	if err := prepareSaml(dashboardServURL); err != nil {
		return err
	}

	staticPathHnd := conf.Current.RootURLPattern + "static/"
	staticDirSrv := http.Dir(util.GetFullPath("static"))
	log.Println("static handler for dashboard on ", staticPathHnd, staticDirSrv)
	http.Handle(staticPathHnd, http.StripPrefix(staticPathHnd, http.FileServer(staticDirSrv)))
	// blog site should be /
	log.Println("Try this url for Blog: ", fmt.Sprintf("http://%s", serverurl))
	staticBlogDirSrv := http.Dir(util.GetFullPath(fmt.Sprintf("static/%s", conf.Current.StaticBlogDir)))
	log.Println("static blog dir", staticBlogDirSrv)
	http.Handle("/", http.StripPrefix("/", http.FileServer(staticBlogDirSrv)))

	chShutdown := make(chan struct{}, 1)
	go func(chs chan struct{}) {
		sch := ObserveComments{simulation: (conf.Current.SimulateAlarm || simulate),
			debug: conf.Current.Debug,
		}
		if err := sch.doObserving(); err != nil {
			log.Println("Server is not observing anymore because: ", err)
			chs <- struct{}{}
		}
	}(chShutdown)

	srv := &http.Server{
		Addr: serverurl,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      nil,
	}
	go func() {
		log.Println("start listening web with http")
		if err := srv.ListenAndServe(); err != nil {
			log.Println("Server is not listening anymore: ", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	log.Println("Enter in server blocking loop")

loop:
	for {
		select {
		case <-sig:
			log.Println("stop because interrupt")
			break loop
		case <-chShutdown:
			log.Println("stop because service shutdown on observing")
			log.Fatal("Force with an error to restart the service")
		}
	}

	log.Println("Bye, service")
	return nil
}

func prepareSaml(dashboardUri string) error {
	log.Println("prepare  saml on ", dashboardUri)
	keyPair, err := tls.LoadX509KeyPair("./cert/igorrun.cert", "./cert/igorrun.key")
	if err != nil {
		log.Println("ERROR 509 key pair")
		return err
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		log.Println("ERROR certificate")
		return err
	}
	myReader := strings.NewReader(`<EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" validUntil="2015-12-03T01:57:09Z" entityID="http://localhost:5572/blog-admin/saml/metadata"><SPSSODescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" validUntil="0001-01-01T00:00:00Z" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol" AuthnRequestsSigned="false" WantAssertionsSigned="true"><KeyDescriptor use="signing"><KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#"><X509Data xmlns="http://www.w3.org/2000/09/xmldsig#"><X509Certificate xmlns="http://www.w3.org/2000/09/xmldsig#">MIIB7zCCAVgCCQDFzbKIp7b3MTANBgkqhkiG9w0BAQUFADA8MQswCQYDVQQGEwJVUzELMAkGA1UECAwCR0ExDDAKBgNVBAoMA2ZvbzESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTEzMTAwMjAwMDg1MVoXDTE0MTAwMjAwMDg1MVowPDELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkdBMQwwCgYDVQQKDANmb28xEjAQBgNVBAMMCWxvY2FsaG9zdDCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEA1PMHYmhZj308kWLhZVT4vOulqx/9ibm5B86fPWwUKKQ2i12MYtz07tzukPymisTDhQaqyJ8Kqb/6JjhmeMnEOdTvSPmHO8m1ZVveJU6NoKRn/mP/BD7FW52WhbrUXLSeHVSKfWkNk6S4hk9MV9TswTvyRIKvRsw0X/gfnqkroJcCAwEAATANBgkqhkiG9w0BAQUFAAOBgQCMMlIO+GNcGekevKgkakpMdAqJfs24maGb90DvTLbRZRD7Xvn1MnVBBS9hzlXiFLYOInXACMW5gcoRFfeTQLSouMM8o57h0uKjfTmuoWHLQLi6hnF+cvCsEFiJZ4AbF+DgmO6TarJ8O05t8zvnOwJlNCASPZRH/JmF8tX0hoHuAQ==</X509Certificate></X509Data></KeyInfo></KeyDescriptor><KeyDescriptor use="encryption"><KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#"><X509Data xmlns="http://www.w3.org/2000/09/xmldsig#"><X509Certificate xmlns="http://www.w3.org/2000/09/xmldsig#">MIIB7zCCAVgCCQDFzbKIp7b3MTANBgkqhkiG9w0BAQUFADA8MQswCQYDVQQGEwJVUzELMAkGA1UECAwCR0ExDDAKBgNVBAoMA2ZvbzESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTEzMTAwMjAwMDg1MVoXDTE0MTAwMjAwMDg1MVowPDELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkdBMQwwCgYDVQQKDANmb28xEjAQBgNVBAMMCWxvY2FsaG9zdDCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEA1PMHYmhZj308kWLhZVT4vOulqx/9ibm5B86fPWwUKKQ2i12MYtz07tzukPymisTDhQaqyJ8Kqb/6JjhmeMnEOdTvSPmHO8m1ZVveJU6NoKRn/mP/BD7FW52WhbrUXLSeHVSKfWkNk6S4hk9MV9TswTvyRIKvRsw0X/gfnqkroJcCAwEAATANBgkqhkiG9w0BAQUFAAOBgQCMMlIO+GNcGekevKgkakpMdAqJfs24maGb90DvTLbRZRD7Xvn1MnVBBS9hzlXiFLYOInXACMW5gcoRFfeTQLSouMM8o57h0uKjfTmuoWHLQLi6hnF+cvCsEFiJZ4AbF+DgmO6TarJ8O05t8zvnOwJlNCASPZRH/JmF8tX0hoHuAQ==</X509Certificate></X509Data></KeyInfo><EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#aes128-cbc"></EncryptionMethod><EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#aes192-cbc"></EncryptionMethod><EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#aes256-cbc"></EncryptionMethod><EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p"></EncryptionMethod></KeyDescriptor><AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="http://localhost:5572/blog-admin/sam/acs" index="1"></AssertionConsumerService></SPSSODescriptor></EntityDescriptor>`)
	req, err := http.NewRequest("PUT", "http://localhost:8000/services/sp",
		myReader)
	if err != nil {
		log.Println("ERROR on PUT service")
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("ERROR on PUT service")
		return err
	} else {
		fmt.Println("*** res", res)
	}

	idpMetadataURL, err := url.Parse("http://localhost:8000/metadata")
	if err != nil {
		log.Println("ERROR url parse  samlidp")
		return err
	}
	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient,
		*idpMetadataURL)
	if err != nil {
		log.Println("ERROR fetch metadata")
		return err
	}

	rootURL, err := url.Parse(dashboardUri)
	if err != nil {
		log.Println("ERROR url parse dashboard")
		return err
	}

	samlSP, _ := samlsp.New(samlsp.Options{
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMetadata,
	})
	log.Println("saml setup ok")
	appApi := http.HandlerFunc(app.APiHandler)
	http.Handle(conf.Current.RootURLPattern, samlSP.RequireAccount(appApi))
	http.Handle("/saml/", samlSP)

	log.Println("saml handler setup completed")
	return nil
}

func (sch *ObserveComments) doObserving() error {
	log.Println("starting observe loop")
	var err error
loop:
	for {
		select {
		case newCmt := <-sch.newCmt:
			log.Println("new comment recognized", newCmt.Id)
			if err = sch.sendNewCommentNtfy(&newCmt); err != nil {
				break loop
			}
		}
	}
	close(sch.newCmt)
	return err
}

func (sch *ObserveComments) sendNewCommentNtfy(cmt *idl.CmtItem) error {
	templ := "templates/cmt/comment-mail.html"
	if err := sendEmail(templ, sch.simulation, cmt); err != nil {
		return err
	}
	if err := sendTelegram(templ, sch.simulation, cmt, sch.debug); err != nil {
		return err
	}
	return nil
}

func sendEmail(templFileName string, simulation bool, cmtItem *idl.CmtItem) error {
	mail := mail.MailSender{}
	mail.FillConf(simulation)
	if err := mail.BuildEmailMsg(templFileName, cmtItem); err != nil {
		return err
	}
	if err := mail.SendEmailViaRelay(); err != nil {
		return err
	}
	return nil
}

func sendTelegram(templFileName string, simulation bool, cmtItem *idl.CmtItem, debug bool) error {
	ts := telegram.TelegramSender{}
	ts.FillConf(simulation, debug)

	if err := ts.BuildMsg(templFileName, cmtItem); err != nil {
		return err
	}
	if err := ts.Send(); err != nil {
		return err
	}
	return nil
}
