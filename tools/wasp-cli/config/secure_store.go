package config

import (
	"encoding/base64"
	"errors"
	"github.com/99designs/keyring"
	"github.com/awnumar/memguard"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	stronghold "github.com/lmoe/stronghold.rs/bindings/native/go"
	"github.com/mr-tron/base58"
	"golang.org/x/term"
	"os"
	"path"
	"path/filepath"
	"syscall"
)

const strongholdKey = "wasp-cli.stronghold.key"
const jwtTokenKey = "wasp-cli.auth.jwt"
const seedKey = "wasp-cli.seed"

type SecureStore struct {
	store keyring.Keyring
}

func zeroKeyBuffer(data *[]byte) {
	for i := 0; i < len(*data); i++ {
		(*data)[i] = 0
	}
}

func passwordCallback(m string) (string, error) {
	storePasswordBuffer, err := StorePassword.Open()

	if err != nil {
		return "", err
	}

	defer storePasswordBuffer.Destroy()

	if len(storePasswordBuffer.String()) > 0 {
		return storePasswordBuffer.String(), nil
	}

	log.Printf("No password set (args[--file-password/-p], env[%v])'\n", StorePasswordEnvKey)
	log.Printf("Enter password manually: ")

	passwordBytes, err := term.ReadPassword(int(syscall.Stdin)) //nolint:unconvert // int cast is needed for windows
	log.Printf("\n")

	return string(passwordBytes), err
}

func NewSecureStore() *SecureStore {
	return &SecureStore{}
}

func (s *SecureStore) getString(key string) (string, error) {
	token, err := s.store.Get(key)

	if err != nil {
		return "", err
	}

	return string(token.Data), nil
}

func (s *SecureStore) createNewStrongholdEnvironment(strongholdPtr *stronghold.StrongholdNative, vaultPath string, addressIndex uint32) error {
	_, err := strongholdPtr.Create(vaultPath)

	if err != nil {
		return err
	}

	_, err = strongholdPtr.GenerateSeed()

	if err != nil {
		return err
	}

	_, err = strongholdPtr.DeriveSeed(addressIndex)

	if err != nil {
		return err
	}

	return nil
}

func (s *SecureStore) Open() error {
	var err error

	dir, _ := filepath.Abs(ConfigPath)
	s.store, err = keyring.Open(keyring.Config{
		ServiceName:          "IOTA_Foundation",
		FileDir:              filepath.Dir(dir),
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

func (s *SecureStore) StoreNewPlainSeed() error {
	seed := cryptolib.NewSeed()
	return s.SetSeed(base58.Encode(seed[:]))
}

func (s *SecureStore) Seed() (*memguard.Enclave, error) {
	seed, err := s.store.Get(seedKey)

	if err != nil {
		return nil, err
	}

	return memguard.NewEnclave(seed.Data), nil
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

func (s *SecureStore) InitializeNewStronghold() error {
	if _, err := os.Stat(s.StrongholdVaultPath()); err == nil {
		err = os.Remove(s.StrongholdVaultPath())

		if err != nil {
			return err
		}
	}

	seed := cryptolib.NewSeed()
	key := base64.StdEncoding.EncodeToString(seed[:])
	item := keyring.Item{
		Key:         strongholdKey,
		Data:        []byte(key),
		Label:       "Stronghold",
		Description: "Key to be used to unlock stronghold",
	}

	err := s.store.Set(item)

	if err != nil {
		return err
	}

	vaultPath := s.StrongholdVaultPath()
	strongholdPtr := stronghold.NewStronghold([]byte(key))
	defer strongholdPtr.Close()

	if _, err := os.Stat(s.StrongholdVaultPath()); errors.Is(err, os.ErrNotExist) {
		err = s.createNewStrongholdEnvironment(strongholdPtr, vaultPath, 0)

		return err
	}

	return nil
}

func (s *SecureStore) StrongholdKey() (*memguard.Enclave, error) {
	item, err := s.store.Get(strongholdKey)

	if errors.Is(err, keyring.ErrKeyNotFound) {
		return nil, err
	}

	enclave := memguard.NewEnclave(item.Data)

	zeroKeyBuffer(&item.Data)

	return enclave, nil
}

func (s *SecureStore) StrongholdVaultPath() string {
	cwd, _ := os.Getwd()
	return path.Join(cwd, "wasp-cli.vault")
}

func (s *SecureStore) OpenStronghold(addressIndex uint32) (*stronghold.StrongholdNative, error) {
	keyEnclave, err := s.StrongholdKey()

	if err != nil {
		return nil, err
	}

	key, err := keyEnclave.Open()
	defer key.Destroy()

	if err != nil {
		return nil, err
	}

	vaultPath := s.StrongholdVaultPath()
	strongholdPtr := stronghold.NewStronghold(key.Bytes())

	_, err = strongholdPtr.Open(vaultPath)

	if err != nil {
		return nil, err
	}

	_, err = strongholdPtr.DeriveSeed(addressIndex)

	if err != nil {
		return nil, err
	}

	return strongholdPtr, err
}
