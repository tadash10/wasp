package verify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

var (
	regex = regexp.MustCompile(`(?m)^import\s+"(?P<import>[@\w/\.]*)";$`)
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
	AddressHash common.Hash `json:"addressHash"`

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
	AutodetectConstructorArguments bool `json:"autodetectConstructorArguments,omitempty"`

	// The EVM version for the contract.
	EvmVersion string `json:"evmVersion,omitempty"`

	// The number of optimization runs used during compilation
	OptimizationRuns int `json:"optimizationRuns,omitempty"`

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

func VerifyContract(explorerAPI string, contractToVerify Contract, importRemapping map[string]string) error {
	if contractToVerify.module != module {
		return fmt.Errorf("at the time of writing this client, Blockscout only supports a module value of '%s', received '%s'", module, contractToVerify.module)
	}

	if contractToVerify.action != action {
		return fmt.Errorf("at the time of writing this client, Blockscout only supports an action value of '%s', received '%s'", action, contractToVerify.action)
	}

	if contractToVerify.ContractSourceCode == nil {
		b, err := os.ReadFile(contractToVerify.ContractSourceCodeFilePath)
		if err != nil {
			return err
		}

		sources := map[string]BlockscoutContractSource{
			contractToVerify.ContractSourceCodeFilePath: {
				Content: string(b),
			},
		}

		var imports []string
		for _, match := range regex.FindAllSubmatch(b, -1) {
			for i, name := range regex.SubexpNames() {
				if i != 0 && name != "" {
					imports = append(imports, string(match[i]))
				}
			}
		}

		
		for _, i := range imports {
			// allow the user to specify solidity import remappings
			// replace `@someImport` part of import paths with the import remapping
			// TODO: Assume same directory and then current directory if not supplied
			f := i
			for k, v := range importRemapping {
				if strings.Contains(i, k) {
					f = strings.Replace(i, k, v, 1)
				}
			}

			b, err := os.ReadFile(f)
			if err != nil {
				return err
			}
	
			sources[i] = BlockscoutContractSource{
				Content: string(b),
			}
		}

		contractToVerify.ContractSourceCode = &BlockscoutContract{
			Language: "Solidity",
			Sources: sources,
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
