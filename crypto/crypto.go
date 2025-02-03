package crypto

import (
	"corsa-blog/util"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/argon2"
)

type UserCred struct {
	UserName     string
	PasswordHash string
	Salt         string
	CredFile     string
	PemFile      string
	PubPemfile   string
	SecretPem    string
	RsaLen       int
	MySigningKey *rsa.PrivateKey
	MyPubKey     *rsa.PublicKey
}

type Token struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	TokenType     string `json:"token_type"`
	Expire        string `json:"expiry"`
	RefreshExpire string `json:"refresh_expiry"`
}

func NewUserCred() *UserCred {
	res := UserCred{
		CredFile:   "./cert/cred.json",
		PemFile:    "./cert/key.pem",
		PubPemfile: "./cert/pubkey.pem",
		RsaLen:     1024,
	}
	return &res
}

func (uc *UserCred) CreateAdminCredentials() error {
	err := uc.CredFromFile()
	if err != nil {
		log.Println("Missed od malformed credential: ", err)
		if err := uc.credFromPrompt(); err != nil {
			return err
		}
		return uc.saveCredential()
	}
	return fmt.Errorf("Unable to create a new crediantial. File %s already exist. Please delete it if you want a new crendential", uc.CredFile)
}

func (uc *UserCred) saveCredential() error {
	path := uc.CredFile
	log.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %v", err)
	}
	defer f.Close()
	cred := struct {
		UserName     string
		PasswordHash string
		Salt         string
		SecretPem    string
	}{
		UserName:     uc.UserName,
		PasswordHash: uc.PasswordHash,
		Salt:         uc.Salt,
		SecretPem:    uc.SecretPem,
	}

	return json.NewEncoder(f).Encode(cred)
}

func (uc *UserCred) credFromPrompt() error {
	var user, pwd, pwdcfm, pwdpem string
	fmt.Println("Please enter the username for administrator")
	fmt.Scanln(&user)
	//fmt.Println("*** user: ", user)

	fmt.Println("Please enter the password")
	fmt.Scanln(&pwd)
	if len(pwd) < 6 {
		return fmt.Errorf("Password must be al least 6 character lenght")
	}

	fmt.Println("Please confirm the password")
	fmt.Scanln(&pwdcfm)
	if pwd != pwdcfm {
		return fmt.Errorf("Password confirm is different")
	}

	genKeyFile := uc.PemFile
	if _, err := os.Stat(genKeyFile); err != nil {
		fmt.Println("Generate also the: ", genKeyFile)
		fmt.Println("Please enter the pwd for the ", genKeyFile)
		fmt.Scanln(&pwdpem)
		if len(pwdpem) < 6 {
			return fmt.Errorf("The pem key is too short")
		}
		priv, _ := rsa.GenerateKey(rand.Reader, uc.RsaLen)
		err := savePrivateKeyInFile(genKeyFile, priv, pwdpem)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("CAUTION: Please enter the right pwd for the ", genKeyFile)
		fmt.Println("A wrong secret for the pem wil make the service unusable. If you recreate the key, please remember that old encrypted files are not available anymore.")
		fmt.Scanln(&pwdpem)
		if len(pwdpem) < 6 {
			return fmt.Errorf("The pem key is too short")
		}
	}

	salt := make([]byte, 32)
	io.ReadFull(rand.Reader, salt)

	uc.UserName = user
	//fmt.Println("*** pwd: ", pwd)
	uc.PasswordHash = hashPassword(pwd, salt)
	uc.Salt = base64.StdEncoding.EncodeToString(salt)
	uc.SecretPem = pwdpem

	return nil
}

func (uc *UserCred) CredFromFile() error {
	file := uc.CredFile
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	cred := struct {
		UserName     string
		PasswordHash string
		Salt         string
		SecretPem    string
	}{}

	err = json.NewDecoder(f).Decode(&cred)
	uc.UserName = cred.UserName
	uc.PasswordHash = cred.PasswordHash
	uc.Salt = cred.Salt
	uc.SecretPem = cred.SecretPem
	return err
}

func (uc *UserCred) fetchPrivateRsaKey() error {
	if uc.MySigningKey == nil {
		mySigningKey, err := privateKeyFromPemFile(util.GetFullPath(uc.PemFile), uc.SecretPem)
		if err != nil {
			return err
		}
		uc.MySigningKey = mySigningKey
	}

	return nil
}

func (uc *UserCred) GetJWTToken(user string, expInSec int, resTk *Token) error {
	log.Println("Using key: ", uc.PemFile)

	err := uc.fetchPrivateRsaKey()
	if err != nil {
		return err
	}
	//log.Printf("Signing key %q \n", mySigningKey)

	exp := time.Now()
	strForSec := "300s" // min validity is 5 minutes. There isn't time sync between portal and service,
	// so avoid little time sync issue like one ore two minutes async.
	if expInSec > 300 {
		strForSec = fmt.Sprintf("%ds", expInSec)
	}
	log.Printf("JWT will Expire in %s seconds\n", strForSec)
	duration, _ := time.ParseDuration(strForSec)
	exp = exp.Add(duration)
	claims := jwt.MapClaims{
		"sub": user,
		"exp": exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	//log.Println("Signing key: ", mySigningKey)
	//tk, err := token.SigningString()
	tk, err := token.SignedString(uc.MySigningKey)
	if err != nil {
		return err
	}
	resTk.AccessToken = tk
	resTk.TokenType = "token_type"
	resTk.Expire = exp.Format(time.RFC3339)
	//RefreshToken
	duration, err = time.ParseDuration(fmt.Sprintf("%ds", 3600*24*7))
	if err != nil {
		return err
	}
	exp = exp.Add(duration)

	claims2 := jwt.MapClaims{
		"sub": user,
		"exp": exp.Unix(),
		"aud": "auth",
	}
	token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims2)
	tk, err = token.SignedString(uc.MySigningKey)
	if err != nil {
		return err
	}
	resTk.RefreshToken = tk
	resTk.RefreshExpire = exp.Format(time.RFC3339)

	return nil

}

func (uc *UserCred) ParseJwtToken(tokenString string) (string, error) {
	if uc.MyPubKey == nil {
		rawdata, err := os.ReadFile(util.GetFullPath(uc.PubPemfile))
		if err != nil {
			return "", err
		}
		pubkey, err := jwt.ParseRSAPublicKeyFromPEM(rawdata)
		if err != nil {
			return "", err
		}
		uc.MyPubKey = pubkey
	}
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return uc.MyPubKey, nil
	})

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			err = fmt.Errorf("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			err = fmt.Errorf("Timing is everything, token is expired")
		} else {
			err = fmt.Errorf("Couldn't handle this token: %v", err)
		}
	}
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Println("Refresh token claims: ", claims)
		if claims["aud"] == "auth" {
			if user, ok := claims["sub"].(string); ok {
				return user, nil
			}
		}
		err = fmt.Errorf("Token is not valid for refresh")
	}
	return "", err
}

func (uc *UserCred) String() string {
	return fmt.Sprintf("Username: %s, Hash: %s, Salt %s", uc.UserName, uc.PasswordHash, uc.Salt)
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

func privateKeyFromPemFile(file string, pwd string) (*rsa.PrivateKey, error) {
	der, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// https://gist.github.com/Northern-Lights/8685a823e5c5503511e89068d855994c
	pemBlock, _ := pem.Decode(der)
	// der, err = x509.DecryptPEMBlock(pemBlock, []byte(pwd)) // deprecated!
	// if err != nil {
	// 	return nil, err
	// }
	// priv, err := x509.ParsePKCS1PrivateKey(der)
	priv, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)

	// Public key can be obtained through priv.PublicKey
	return priv, err
}

func savePrivateKeyInFile(file string, priv *rsa.PrivateKey, pwd string) error {
	//der := x509.MarshalPKCS1PrivateKey(priv)
	// pp := []byte(pwd)
	// block, err := x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", der, pp, x509.PEMCipherAES256)
	// if err != nil {
	// 	return err
	// }
	//return os.WriteFile(file, pem.EncodeToMemory(block), 0644)
	block := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)
	//pem.Encode(der, block)
	log.Println("Save the key in ", file)
	return os.WriteFile(file, block, 0644)
}
