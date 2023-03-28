package verify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/slices"
)

var (
	urlRegex     = regexp.MustCompile(`(?m)^\w+:\/\/`)
	importRegex  = regexp.MustCompile(`(?m)^import\s+"(?P<import>[@\w/\.]*)";$`)
	validSchemes = []string{"https", "http"}
)

const (
	module = "contract"
	action = "verify"
)

type OptimizationUsed string

const (
	NoOptimizationWasUsed OptimizationUsed = "0"
	OptimizationWasUsed   OptimizationUsed = "1"
)

type ContractSource byte

const (
	InvalidContractSource    ContractSource = iota
	HTTPContractSource       ContractSource = iota
	FilesystemContractSource ContractSource = iota
)

type BlockscoutContractSource struct {
	Content string `json:"content"`
}

type BlockscoutContract struct {
	Language string                              `json:"language"`
	Sources  map[string]BlockscoutContractSource `json:"sources"`
}

type Contract struct {
	// A string with the name of the module to be invoked.
	//
	// Must be set to: contract
	module string

	// A string with the name of the action to be invoked.
	//
	// Must be set to: verify
	action string

	// The address of the contract.
	AddressHash *common.Hash `json:"addressHash"`

	// The name of the contract.
	Name string `json:"name"`

	// The compiler version for the contract.
	CompilerVersion string `json:"compilerVersion"`

	// Whether or not compiler optimizations were enabled.
	Optimization bool `json:"optimization"`

	// The path to the source code of the contract.
	//
	// ContractSourceCodeFilePath is mutually exclusive with ContractSourceCode, setting ContractSourceCode will cause ContractSourceCodeFilePath to be ignored.
	ContractSourceCodeFilePath string `json:"-"`

	// The source code of the contract.
	//
	// ContractSourceCode is mutually exclusive with ContractSourceCodeFilePath, setting ContractSourceCode will cause ContractSourceCodeFilePath to be ignored.
	ContractSourceCode *BlockscoutContract `json:"contractSourceCode,omitempty"`

	// The constructor argument data provided.
	ConstructorArguments string `json:"constructorArguments,omitempty"`

	// Whether or not automatically detect constructor argument.
	AutoDetectConstructorArguments bool `json:"autodetectConstructorArguments,omitempty"`

	// The EVM version for the contract.
	EvmVersion string `json:"evmVersion,omitempty"`

	// The number of optimization runs used during compilation
	OptimizationRuns int `json:"optimizationRuns,omitempty"`

	// Mapping replacements
	ImportRemappings map[string]string `json:"-"`

	// The name of the first library used.
	Library1Name string `json:"library1Name,omitempty"`

	// The address of the first library used.
	Library1Address common.Address `json:"library1Address,omitempty"`

	// The name of the second library used.
	Library2Name string `json:"library2Name,omitempty"`

	// The address of the second library used.
	Library2Address common.Address `json:"library2Address,omitempty"`

	// The name of the third library used.
	Library3Name string `json:"library3Name,omitempty"`

	// The address of the third library used.
	Library3Address common.Address `json:"library3Address,omitempty"`

	// The name of the fourth library used.
	Library4Name string `json:"library4Name,omitempty"`

	// The address of the fourth library used.
	Library4Address common.Address `json:"library4Address,omitempty"`

	// The name of the fifth library used.
	Library5Name string `json:"library5Name,omitempty"`

	// The address of the fifth library used.
	Library5Address common.Address `json:"library5Address,omitempty"`
}

type BlockScoutVerifyContractModel struct {
	ABI              string           `json:"ABI"`
	Address          common.Address   `json:"Address"`
	ContractName     string           `json:"ContractName"`
	OptimizationUsed OptimizationUsed `json:"OptimizationUsed"`
}

type BlockScoutContractVerifyResponse struct {
	Message string                         `json:"message"`
	Result  *BlockScoutVerifyContractModel `json:"result"`
	Status  string                         `json:"status"`
}

func NewContract(
	addressHash *common.Hash,
	name,
	sourceCodeFilePath,
	compilerVersion,
	constructorArguments,
	evmVersion string,
	optimization,
	autoDetectConstructorArguments bool,
	optimizationRuns int,
	remap map[string]string,
) *Contract {
	return &Contract{
		module:                         module,
		action:                         action,
		AddressHash:                    addressHash,
		Name:                           name,
		CompilerVersion:                compilerVersion,
		Optimization:                   optimization,
		ContractSourceCodeFilePath:     sourceCodeFilePath,
		ConstructorArguments:           constructorArguments,
		AutoDetectConstructorArguments: autoDetectConstructorArguments,
		EvmVersion:                     evmVersion,
		OptimizationRuns:               optimizationRuns,
		ImportRemappings:               remap,
	}
}

func ValidateSource(input string) (ContractSource, error) {
	if len(urlRegex.FindString(input)) > 0 {
		if u, uErr := url.ParseRequestURI(input); uErr != nil {
			return InvalidContractSource, fmt.Errorf("error validating as a url: %w", uErr)
		} else if !slices.Contains(validSchemes, u.Scheme) {
			return InvalidContractSource, NewErrInvalidURLScheme(u.Scheme)
		} else if len(u.Host) == 0 {
			return InvalidContractSource, NewErrInvalidURLHostname(u.Host)
		} else {
			return HTTPContractSource, nil
		}
	} else {
		if _, fErr := os.Stat(input); fErr != nil && !(errors.Is(fErr, os.ErrNotExist) || errors.Is(fErr, os.ErrPermission)) {
			return InvalidContractSource, fmt.Errorf("error checking local file system: %w", fErr)
		} else if fErr != nil {
			return InvalidContractSource, fErr
		} else {
			return FilesystemContractSource, nil
		}
	}
}

func ParseContract(path string, knownSources map[string]BlockscoutContractSource, importRemapping map[string]string) (map[string]BlockscoutContractSource, error) {
	importPath := path
	for k, v := range importRemapping {
		if strings.Contains(path, k) {
			importPath = strings.Replace(path, k, v, 1)
			break
		}
	}

	source, err := ValidateSource(importPath)
	if err != nil {
		return nil, err
	}

	var b []byte

	switch source {
	case FilesystemContractSource:
		b, err = os.ReadFile(importPath)
		if err != nil {
			return nil, err
		}
	case HTTPContractSource:
		resp, err := http.Get(importPath)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	case InvalidContractSource:
		return nil, fmt.Errorf("cannot retrieve contract from %s", importPath)
	}

	var imports []string
	for _, match := range importRegex.FindAllSubmatch(b, -1) {
		for i, name := range importRegex.SubexpNames() {
			if i != 0 && name != "" {
				imports = append(imports, string(match[i]))
			}
		}
	}

	knownSources[path] = BlockscoutContractSource{
		Content: string(b),
	}

	for _, _import := range imports {
		if _, ok := knownSources[importPath]; !ok {
			importedSources, err := ParseContract(_import, knownSources, importRemapping)
			if err != nil {
				return nil, err
			}

			for k, v := range importedSources {
				knownSources[k] = v
			}
		}
	}
	fmt.Printf("imports found in %s, %v\n", path, imports)
	return knownSources, nil
}

func VerifyContract(explorerAPI string, contractToVerify *Contract) error {
	if contractToVerify.module != module {
		return fmt.Errorf("at the time of writing this client, Blockscout only supports a module value of '%s', received '%s'", module, contractToVerify.module)
	}

	if contractToVerify.action != action {
		return fmt.Errorf("at the time of writing this client, Blockscout only supports an action value of '%s', received '%s'", action, contractToVerify.action)
	}

	if contractToVerify.ContractSourceCode == nil {
		sources, err := ParseContract(contractToVerify.ContractSourceCodeFilePath, make(map[string]BlockscoutContractSource), contractToVerify.ImportRemappings)
		if err != nil {
			return err
		}

		contractToVerify.ContractSourceCode = &BlockscoutContract{
			Language: "Solidity",
			Sources:  sources,
		}
	}

	jsonData, err := json.Marshal(contractToVerify)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf(
			"%s?%s=%s&%s=%s",
			explorerAPI,
			module,
			url.QueryEscape(contractToVerify.module),
			action,
			url.QueryEscape(contractToVerify.action),
		),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r BlockScoutContractVerifyResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return err
	}

	return nil
}
