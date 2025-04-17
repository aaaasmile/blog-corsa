package admin

import (
	"corsa-blog/conf"
	"corsa-blog/crypto"
	"encoding/json"
	"fmt"
	"log"
)

type LoginParam struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type LoginReq struct {
	Params LoginParam
}

func (ah *AdminHandler) doLogin() error {
	//fmt.Println("*** raw", string(ah.rawbody))
	loginReq := LoginReq{}
	if err := json.Unmarshal(ah.rawbody, &loginReq); err != nil {
		return err
	}
	credReq := loginReq.Params
	if credReq.User == "" || credReq.Password == "" {
		return fmt.Errorf("user or password wrong")
	}
	refCred := conf.Current.AdminCred
	if credReq.User == refCred.UserName {
		log.Println("Check password for user ", credReq.User)
		//fmt.Println("*** refcred", refCred)
		hash := crypto.GetHashOfSecret(credReq.Password, refCred.Salt)
		//log.Println("Hash is: ", hash)
		//fmt.Println("*** hash is ", hash)
		if hash == refCred.PasswordHash {
			return tokenResult(200, credReq.User, ah._w)
		}
	}

	return tokenResult(403, credReq.User, ah._w)
}
