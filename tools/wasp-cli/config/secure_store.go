package config

import (
	"encoding/base64"
	"errors"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/99designs/keyring"
	"github.com/awnumar/memguard"
	stronghold_go "github.com/iotaledger/stronghold-bindings/go"
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
	if IsStrongholdScheme() {
		if log.VerboseFlag {
			stronghold_go.SetLogLevel(stronghold_go.LogLevelTrace)
		}

		if log.DebugFlag {
			stronghold_go.SetLogLevel(stronghold_go.LogLevelInfo)
		}
	}

	return &SecureStore{}
}

func (s *SecureStore) createNewStrongholdEnvironment(strongholdPtr *stronghold_go.StrongholdNative, vaultPath string, addressIndex uint32) error {
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

func (s *SecureStore) InitializeNewStronghold() error {
	if _, err := os.Stat(s.StrongholdVaultPath()); err == nil {
		err = os.Remove(s.StrongholdVaultPath())

		if err != nil {
			return err
		}
	}

	seed := cryptolib.NewUntypedSeed()
	encodedSeed := make([]byte, base64.StdEncoding.EncodedLen(len(seed)))
	base64.StdEncoding.Encode(encodedSeed, seed)

	defer func() {
		zeroKeyBuffer(&seed)
		zeroKeyBuffer(&encodedSeed)
	}()

	item := keyring.Item{
		Key:         strongholdKey,
		Data:        encodedSeed,
		Label:       "Stronghold",
		Description: "Key to be used to unlock stronghold",
	}

	err := s.store.Set(item)
	if err != nil {
		return err
	}

	vaultPath := s.StrongholdVaultPath()
	strongholdPtr := stronghold_go.NewStronghold(encodedSeed)
	defer strongholdPtr.Close()

	if _, err := os.Stat(s.StrongholdVaultPath()); errors.Is(err, os.ErrNotExist) {
		err = s.createNewStrongholdEnvironment(strongholdPtr, vaultPath, 0)

		return err
	}

	return nil
}

func (s *SecureStore) StrongholdKey() (*memguard.Enclave, error) {
	item, err := s.store.Get(strongholdKey)
	if err != nil {
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

func (s *SecureStore) OpenStronghold(addressIndex uint32) (*stronghold_go.StrongholdNative, error) {
	keyEnclave, err := s.StrongholdKey()
	if err != nil {
		return nil, err
	}

	vaultPath := s.StrongholdVaultPath()

	strongholdPtr := stronghold_go.NewStrongholdWithEnclave(keyEnclave)

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
