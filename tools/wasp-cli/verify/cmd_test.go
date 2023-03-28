package verify_test

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotaledger/wasp/tools/wasp-cli/util"
	"github.com/iotaledger/wasp/tools/wasp-cli/verify"
)

var (
	cwd string
)

func TestMain(m *testing.M) {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestValidateMapping(t *testing.T) {
	testMappings := map[string]struct {
		Error  error
		Source verify.ContractSource
	}{
		path.Join(cwd, "..", "..", "..", "packages", "vm", "core", "evm", "iscmagic"):   {nil, verify.FilesystemContractSource},
		path.Join("..", "..", "..", "packages", "vm", "core", "evm", "iscmagic"):        {nil, verify.FilesystemContractSource},
		path.Join("packages", "vm", "core", "evm", "iscmagic"):                          {fs.ErrNotExist, verify.InvalidContractSource},
		"https://github.com/iotaledger/wasp/tree/develop/packages/vm/core/evm/iscmagic": {nil, verify.HTTPContractSource},
		"git://github.com/iotaledger/wasp/tree/develop/packages/vm/core/evm/iscmagic":   {verify.ErrInvalidURLScheme{}, verify.InvalidContractSource},
		"http:///iotaledger/wasp/tree/develop/packages/vm/core/evm/iscmagic":            {verify.ErrInvalidURLHostname{}, verify.InvalidContractSource},
	}
	for mapping, expectedErr := range testMappings {
		if source, err := verify.ValidateSource(mapping); err != nil && !errors.Is(err, expectedErr.Error) {
			t.Logf("error validating %s", mapping)
			t.Error(err)
		} else if err != nil && errors.Is(err, expectedErr.Error) {
			t.Logf("successfully errored with %s", err)
		} else if source == expectedErr.Source {
			t.Logf("valid mapping found: %s", mapping)
		} else {
			t.Errorf("correct contract source not identified for %s", mapping)
		}
	}
}

func TestParseContract(t *testing.T) {
	testMappings := map[string]error{
		path.Join("@iscmagic", "ISCAccounts.sol"):                                                        nil,
		path.Join(cwd, "..", "..", "..", "packages", "vm", "core", "evm", "iscmagic", "ISCAccounts.sol"): nil,
		path.Join("..", "..", "..", "packages", "vm", "core", "evm", "iscmagic", "ISCAccounts.sol"):      nil,
		path.Join("packages", "vm", "core", "evm", "iscmagic", "ISCAccounts.sol"):                        fs.ErrNotExist,
		// "https://github.com/iotaledger/wasp/tree/develop/packages/vm/core/evm/iscmagic/ISCAccounts.sol":  nil,
		// "git://github.com/iotaledger/wasp/tree/develop/packages/vm/core/evm/iscmagic/ISCAccounts.sol":    verify.ErrInvalidURLScheme{},
		// "http:///iotaledger/wasp/tree/develop/packages/vm/core/evm/iscmagic/ISCAccounts.sol":             verify.ErrInvalidURLHostname{},
	}
	for mapping, expectedErr := range testMappings {
		c := verify.NewContract(
			new(common.Hash),
			"test",
			mapping,
			util.SolcVersion,
			"",
			"",
			false,
			false,
			0,
			map[string]string{
				"@iscmagic": path.Join(cwd, "..", "..", "..", "packages", "vm", "core", "evm", "iscmagic"),
			},
		)
		sources := make(map[string]verify.BlockscoutContractSource)
		if sources, err := verify.ParseContract(c.ContractSourceCodeFilePath, sources, c.ImportRemappings); err != nil && !errors.Is(err, expectedErr) {
			t.Logf("error validating %s", mapping)
			t.Error(err)
		} else if err != nil && errors.Is(err, expectedErr) {
			t.Logf("successfully errored with %s", err)
		} else {
			t.Logf("valid sources found: %v", sources)
		}
	}
}
