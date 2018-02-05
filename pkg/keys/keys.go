package keys

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/lytics/base62"
)

// GenerateKeyPair generates a new ecdsa keypair.
func GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
}

// GetAddress returns the base62 encoded public address of a keypair.
func GetAddress(keyPair *ecdsa.PrivateKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&keyPair.PublicKey)
	if err != nil {
		return "", err
	}

	return base62.StdEncoding.EncodeToString(publicKeyBytes), nil
}

// DecodeAddress decodes a base62 address and return a public key.
func DecodeAddress(address string) (*ecdsa.PublicKey, error) {
	bytes, err := base62.StdEncoding.DecodeString(address)
	if err != nil {
		return nil, err
	}

	pub, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		return nil, err
	}

	return pub.(*ecdsa.PublicKey), nil
}

// WriteKeys writes a keypair to a pem file.
func WriteKeys(keyPair *ecdsa.PrivateKey, pemFile string) error {
	f, err := os.Create(pemFile)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	err = write(keyPair, w)
	if err != nil {
		return err
	}

	err = w.Flush()
	return err
}

// LoadKeys loads a keypair from a pem file.
func LoadKeys(filename string) (*ecdsa.PrivateKey, error) {
	keyPairBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return decode(keyPairBytes)
}

func write(keyPair *ecdsa.PrivateKey, writer io.Writer) error {
	x509Encoded, err := x509.MarshalECPrivateKey(keyPair)
	if err != nil {
		return err
	}

	err = pem.Encode(writer, &pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	if err != nil {
		return err
	}

	x509EncodedPub, err := x509.MarshalPKIXPublicKey(&keyPair.PublicKey)
	if err != nil {
		return err
	}

	return pem.Encode(writer, &pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
}

func decode(pemBytes []byte) (*ecdsa.PrivateKey, error) {
	keyPair := &ecdsa.PrivateKey{}
	var (
		err   error
		block *pem.Block
	)
	for {
		block, pemBytes = pem.Decode(pemBytes)
		if block == nil {
			break
		}
		if block.Type == "PRIVATE KEY" {
			keyPair, err = x509.ParseECPrivateKey(block.Bytes)
		}

		if block.Type == "PUBLIC KEY" {
			var pub interface{}
			pub, err = x509.ParsePKIXPublicKey(block.Bytes)
			keyPair.PublicKey = *(pub.(*ecdsa.PublicKey))

		}
	}
	return keyPair, err
}

// GetDefaultKeyPath returns the default path to a user's keys.
func GetDefaultKeyPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	os.MkdirAll(fmt.Sprintf("%s/%s", usr.HomeDir, ".lolachain/"), os.ModePerm)

	return fmt.Sprintf("%s/%s", usr.HomeDir, ".lolachain/key.pem"), nil

}

//LoadOrGenerateKeys loads the keys at the provided path, or generates them if needed.
func LoadOrGenerateKeys(filename string) (*ecdsa.PrivateKey, error) {
	var keyPair *ecdsa.PrivateKey
	if _, err := os.Stat(filename); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		keyPair, err = GenerateKeyPair()
		if err != nil {
			return nil, err
		}

		err = WriteKeys(keyPair, filename)
		if err != nil {
			return nil, err
		}
		return keyPair, nil
	}

	return LoadKeys(filename)
}
