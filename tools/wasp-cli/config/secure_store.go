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
	stronghold *stronghold.StrongholdNative //nolint:typecheck
}

func NewSecureStore() *SecureStore {
	return &SecureStore{}
}

func (s *SecureStore) Open() error {
	var err error

	s.store, err = keyring.Open(keyring.Config{
		KeychainName: "IOTA_Foundation/wasp-cli",
		// TODO: Make configurable
		FileDir: "wasp-cli.secure.json",
	})

	if err != nil {
		return err
	}

	err = s.initializeStronghold()

	return err
}

func (s *SecureStore) zeroKeyBuffer(data *[]byte) {
	for i := 0; i < len(*data); i++ {
		(*data)[i] = 0
	}
}

func (s *SecureStore) initializeStronghold() error {
	key, err := s.StrongholdKey()

	if err != nil {
		return err
	}

	s.stronghold = stronghold.NewStronghold(key)
	s.zeroKeyBuffer(&key)

	// TODO: Make configurable
	cwd, _ := os.Getwd()
	vaultPath := path.Join(cwd, "wasp-cli.vault")
	//

	success, err := s.stronghold.OpenOrCreate(vaultPath)

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

func (s *SecureStore) Reset() error {
	keys, err := s.store.Keys()

	if err != nil {
		return nil
	}

	for _, key := range keys {
		s.store.Remove(key)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SecureStore) Token() (string, error) {
	return s.getString(jwtTokenKey)
}

func (s *SecureStore) SetToken(token string) error {
	item := keyring.Item{Key: jwtTokenKey, Data: []byte(token)}

	return s.store.Set(item)
}

func (s *SecureStore) GenerateNewStrongholdKey() error {
	key := "U)%&Usdfj95ÄÖsdfinsd" // Super secret password randomly generated, promised!
	item := keyring.Item{Key: strongholdKey, Data: []byte(key)}
	return s.store.Set(item)
}

func (s *SecureStore) StrongholdKey() ([]byte, error) {
	item, err := s.store.Get(strongholdKey)

	if errors.Is(err, keyring.ErrKeyNotFound) {
		err = s.GenerateNewStrongholdKey()

		if err != nil {
			return []byte{}, err
		}
	}

	item, err = s.store.Get(strongholdKey)

	if err != nil {
		return []byte{}, err
	}

	return item.Data, nil
}

func (s *SecureStore) Stronghold() *stronghold.StrongholdNative {
	return s.stronghold
}
