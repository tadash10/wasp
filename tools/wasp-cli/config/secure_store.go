package config

import (
	"path/filepath"
	"syscall"

	"github.com/99designs/keyring"
	"github.com/awnumar/memguard"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	"github.com/mr-tron/base58"
	"golang.org/x/term"
)

const (
	strongholdKey = "wasp-cli.stronghold.key"
	jwtTokenKey   = "wasp-cli.auth.jwt"
	seedKey       = "wasp-cli.seed"
)

type SecureStore struct {
	store keyring.Keyring
}

func zeroKeyBuffer(data *[]byte) {
	for i := 0; i < len(*data); i++ {
		(*data)[i] = 0
	}
}

func passwordCallback(m string) (string, error) {
	if len(StorePassword) > 0 {
		return StorePassword, nil
	}

	log.Printf("No password set (args[--file-password/-s], env[%v])'\n", StorePasswordEnvKey)
	log.Printf("Enter password manually: ")

	passwordBytes, err := term.ReadPassword(int(syscall.Stdin)) //nolint:unconvert // int cast is needed for windows
	log.Printf("\n")

	return string(passwordBytes), err
}

func NewSecureStore() *SecureStore {
	return &SecureStore{}
}

func (s *SecureStore) Open() error {
	dir, err := filepath.Abs(ConfigPath)
	if err != nil {
		return err
	}

	s.store, err = keyring.Open(keyring.Config{
		ServiceName:          "IOTA_Foundation",
		FileDir:              filepath.Join(filepath.Dir(dir), "secret_store"),
		KeychainPasswordFunc: passwordCallback,
		FilePasswordFunc:     passwordCallback,
	})

	return err
}

func (s *SecureStore) Reset() error {
	keys, err := s.store.Keys()
	if err != nil {
		return nil
	}

	for _, key := range keys {
		err := s.store.Remove(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SecureStore) Dump() (map[string]string, error) {
	vaultDump := make(map[string]string)
	keys, err := s.store.Keys()
	if err != nil {
		return vaultDump, err
	}

	for _, key := range keys {
		if key == strongholdKey {
			vaultDump[key] = "**REDACTED**"
			continue
		}

		item, err := s.store.Get(key)
		if err != nil {
			return vaultDump, err
		}

		vaultDump[key] = string(item.Data)
	}

	return vaultDump, nil
}

func (s *SecureStore) GenerateAndStorePlainSeed() error {
	seed := cryptolib.NewUntypedSeed()
	err := s.SetSeed(base58.Encode(seed))
	if err != nil {
		return err
	}

	zeroKeyBuffer(&seed)

	return nil
}

func (s *SecureStore) Seed() (*memguard.Enclave, error) {
	seed, err := s.store.Get(seedKey)
	if err != nil {
		return nil, err
	}

	seedEnclave := memguard.NewEnclave(seed.Data)
	zeroKeyBuffer(&seed.Data)

	return seedEnclave, nil
}

func (s *SecureStore) SetSeed(seed string) error {
	item := keyring.Item{
		Key:         seedKey,
		Data:        []byte(seed),
		Label:       "Seed",
		Description: "The private seed of the wasp cli",
	}

	return s.store.Set(item)
}

func (s *SecureStore) Token() (string, error) {
	token, err := s.store.Get(jwtTokenKey)
	if err != nil {
		return "", err
	}

	return string(token.Data), nil
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
