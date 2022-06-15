package config

import (
	"errors"
	"github.com/99designs/keyring"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	"github.com/mr-tron/base58"
	"golang.org/x/term"
	"golang.org/x/xerrors"
	"path/filepath"
	"syscall"
)

const strongholdKey = "wasp-cli.stronghold.key"
const jwtTokenKey = "wasp-cli.auth.jwt"
const seedKey = "wasp-cli.seed"

type SecureStore struct {
	store keyring.Keyring
}

func NewSecureStore() *SecureStore {
	return &SecureStore{}
}

func passwordCallback(m string) (string, error) {
	if len(FileStorePassword) > 0 {
		return FileStorePassword, nil
	}

	log.Printf("No password set (args[--file-password/-p], env[%v])'\n", FilePasswordEnvKey)
	log.Printf("Enter password manually: ")

	passwordBytes, err := term.ReadPassword(int(syscall.Stdin)) //nolint:unconvert // int cast is needed for windows
	log.Printf("\n")

	return string(passwordBytes), err
}

func (s *SecureStore) Open() error {
	var err error

	keyring.Debug = true

	dir, _ := filepath.Abs(ConfigPath)
	s.store, err = keyring.Open(keyring.Config{
		KeychainName:         "IOTA Foundation",
		ServiceName:          "IOTA Foundation",
		FileDir:              filepath.Dir(dir),
		KeychainPasswordFunc: passwordCallback,
		FilePasswordFunc:     passwordCallback,
	})

	return err
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

func (s *SecureStore) CanAccessSeed() bool {
	return WalletScheme() == WalletSchemePlain
}

func (s *SecureStore) Seed() (string, error) {
	if !s.CanAccessSeed() {
		return "", xerrors.New("Can not access seed")
	}

	return s.getString(seedKey)
}

func (s *SecureStore) StoreNewSeed() error {
	seed := cryptolib.NewSeed()
	seedBase58 := base58.Encode(seed[:])

	return s.SetSeed(seedBase58)
}

func (s *SecureStore) SetSeed(seedBase58 string) error {
	if s.CanAccessSeed() {
		item := keyring.Item{
			Key:         seedKey,
			Data:        []byte(seedBase58),
			Label:       "Seed",
			Description: "The private seed of the wasp cli",
		}

		return s.store.Set(item)
	}

	return nil
}

func (s *SecureStore) Token() (string, error) {
	return s.getString(jwtTokenKey)
}

func (s *SecureStore) SetToken(token string) error {
	item := keyring.Item{
		Key:         jwtTokenKey,
		Data:        []byte(token),
		Label:       "JWT",
		Description: "A token to authenticate against a wasp node",
	}

	return s.store.Set(item)
}

func (s *SecureStore) GenerateNewStrongholdKey() error {
	key := cryptolib.NewSeed()
	item := keyring.Item{
		Key:         strongholdKey,
		Data:        key[:],
		Label:       "Stronghold",
		Description: "Key to be used to unlock stronghold",
	}

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
