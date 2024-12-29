package watch

import (
	"corsa-blog/conf"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func RunWatcher(configfile string, targetFile string) error {
	if targetFile == "" {
		return fmt.Errorf("target file is empty")
	}
	log.Println("watching ", targetFile)
	if _, err := os.Stat(targetFile); err != nil {
		return err
	}
	if _, err := conf.ReadConfig(configfile); err != nil {
		return err
	}

	chShutdown := make(chan struct{}, 1)
	go func(chs chan struct{}) {
		// sch := Scheduler{datafileName: conf.Current.DataFileName,
		// 	simulation: (conf.Current.SimulateAlarm || simulate),
		// 	debug:      conf.Current.Debug,
		// }
		// if err := sch.doSchedule(); err != nil {
		// 	log.Println("Server is not scheduling anymore: ", err)
		// 	chs <- struct{}{}
		// }
	}(chShutdown)

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
			log.Println("stop because service shutdown on scheduling")
			log.Fatal("Force with an error to restart the service")
		}
	}

	log.Println("Bye, service")
	return nil
}
