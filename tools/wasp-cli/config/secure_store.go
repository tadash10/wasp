package config

import (
	"errors"
	"github.com/99designs/keyring"
	"github.com/lmoe/stronghold.rs/bindings/native/go"
	"golang.org/x/xerrors"
	"os"
	"path"
)

const strongholdKey = "stronghold.key"
const jwtTokenKey = "auth.jwt"

type SecureStore struct {
	store      keyring.Keyring
	stronghold *stronghold.StrongholdNative
}

func NewSecureStore() *SecureStore {
	return &SecureStore{}
}

func (s *SecureStore) Open() error {
	var err error

	s.store, err = keyring.Open(keyring.Config{
		KeychainName: "IOTA_Foundation/wasp-cli",
		//FileDir:     "wasp-cli.secure.json",
	})

	if err != nil {
		return err
	}

	key, err := s.StrongholdKey()

	if err != nil {
		return err
	}

	s.stronghold = stronghold.NewStronghold(key)

	cwd, _ := os.Getwd()
	vault := path.Join(cwd, "wasp-cli.vault")
	success, err := s.stronghold.Create(vault)

	if err != nil {
		return err
	}

	if !success {
		return xerrors.New("failed to open vault with an unknown error")
	}

	return nil
}

func (s *SecureStore) getString(key string) (string, error) {
	token, err := s.store.Get(key)

	if err != nil {
		return "", err
	}

	return string(token.Data), nil
}

func (s *SecureStore) Token() (string, error) {
	return s.getString(jwtTokenKey)
}

func (s *SecureStore) GenerateNewStrongholdKey() error {
	key := "SuperSecretPassword"
	return s.store.Set(keyring.Item{Key: strongholdKey, Data: []byte(key)})
}

func (s *SecureStore) StrongholdKey() (string, error) {
	item, err := s.store.Get(strongholdKey)

	if errors.Is(err, keyring.ErrKeyNotFound) {
		err = s.GenerateNewStrongholdKey()

		if err != nil {
			return "", err
		}
	}

	item, err = s.store.Get(strongholdKey)

	if err != nil {
		return "", err
	}

	return string(item.Data), nil
}

func (s *SecureStore) Stronghold() *stronghold.StrongholdNative {
	return s.stronghold
}
