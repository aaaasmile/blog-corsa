package web

import (
	"corsa-blog/conf"
	"corsa-blog/idl"
	"corsa-blog/mail"
	"corsa-blog/telegram"
	"corsa-blog/util"
	"corsa-blog/web/app"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
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

	staticPathHnd := conf.Current.RootURLPattern + "static/"
	staticDirSrv := http.Dir(util.GetFullPath("static"))
	log.Println("static handler for dashboard on ", staticPathHnd, staticDirSrv)
	http.Handle(staticPathHnd, http.StripPrefix(staticPathHnd, http.FileServer(staticDirSrv)))
	// blog site should be /
	log.Println("Try this url for Blog: ", fmt.Sprintf("http://%s", serverurl))
	staticBlogDirSrv := http.Dir(util.GetFullPath(fmt.Sprintf("static/%s", conf.Current.StaticBlogDir)))
	log.Println("static blog dir", staticBlogDirSrv)
	http.Handle("/", http.StripPrefix("/", http.FileServer(staticBlogDirSrv)))

	// Dashboard app
	http.HandleFunc(conf.Current.RootURLPattern, app.APiHandler)

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
	templ := "templates/comment-mail.html"
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
