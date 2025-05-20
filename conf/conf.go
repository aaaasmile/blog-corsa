package conf

import (
	"corsa-blog/crypto"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ServiceURL     string
	RootURLPattern string
	ServerName     string
	StaticBlogDir  string
	PostSubDir     string
	PageSubDir     string
	VuetifyLibName string
	VueLibName     string
	Relay          *Relay
	Telegram       *Telegram
	Database       *Database
	SimulateAlarm  bool
	Debug          bool
	AdminCred      *crypto.UserCred
	Comment        *Comment
}

type Comment struct {
	AllowEmptyMail   bool
	HasDateInCmtForm bool
	ModerateCmt      bool
}

type Database struct {
	DbFileName string
	SQLDebug   bool
}

type Telegram struct {
	SendTelegram bool
	ChatID       int64
	APIString    string
}

type Relay struct {
	SendMail    bool
	MailFrom    string
	Secret      string
	Host        string
	User        string
	EmailTarget string
}

var Current = &Config{}

func ReadConfig(configfile, baseDirCert string) (*Config, error) {
	_, err := os.Stat(configfile)
	if err != nil {
		return nil, err
	}
	if _, err := toml.DecodeFile(configfile, &Current); err != nil {
		return nil, err
	}

	log.Println("Configuration: ", Current.Relay.MailFrom, Current.Relay.Host, Current.Relay.MailFrom, Current.Telegram.SendTelegram)
	return Current, nil
}
