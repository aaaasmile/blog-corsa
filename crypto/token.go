package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"golang.org/x/crypto/argon2"
)

type Token struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	TokenType     string `json:"token_type"`
	Expire        string `json:"expiry"`
	RefreshExpire string `json:"refresh_expiry"`
}

func GetHashOfSecret(pwd, salt string) string {
	ss, _ := base64.StdEncoding.DecodeString(salt)
	return hashPassword(pwd, ss)
}

func hashPassword(pwd string, salt []byte) string {
	ram := 512 * 1024
	t0 := time.Now()
	//fmt.Printf("*** salt is %x\n", salt)
	//fmt.Printf("*** password is %s\n", pwd)
	key := argon2.IDKey([]byte(pwd), salt, 1, uint32(ram), uint8(runtime.NumCPU()<<1), 32)
	log.Printf("hash time: %v, key: %x, salt: %x\n", time.Since(t0), key, salt)
	return fmt.Sprintf("%x", key)
}

func privateKeyFromPemFile(file string) (*rsa.PrivateKey, error) {
	der, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// https://gist.github.com/Northern-Lights/8685a823e5c5503511e89068d855994c
	pemBlock, _ := pem.Decode(der)
	priv, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	// Public key can be obtained through priv.PublicKey
	return priv, err
}

func savePrivateKeyInFile(file string, priv *rsa.PrivateKey) error {
	block := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)
	log.Println("Save the key in ", file)
	return os.WriteFile(file, block, 0644)
}
