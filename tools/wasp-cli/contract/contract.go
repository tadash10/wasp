package contract

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/slices"
)

var (
	urlRegex           = regexp.MustCompile(`(?m)^\w+:\/\/`)
	importRegex        = regexp.MustCompile(`(?m)^import\s+"(?P<import>[@\w/\.]*)";$`)
	validSchemes       = []string{"https", "http"}
	blockscoutSettings = map[string]any{
		"optimizer": map[string]any{
			"enabled": false,
			"runs":    200,
		},
		"outputSelection": map[string]any{
			"*": map[string]any{
				"*": []string{
					"abi",
					"evm.bytecode",
					"evm.deployedBytecode",
					"evm.methodIdentifiers",
				},
				"": []string{
					"ast",
				},
			},
		},
	}
)

const (
	module                          string = "contract"
	action                          string = "verifysourcecode"
	codeFormat                      string = "solidity-standard-json-input"
	alreadyVerifiedContractResponse string = "Smart-contract already verified."
)

type VerificationStatusResult string

const (
	verifyAction string                   = "checkverifystatus"
	pending      VerificationStatusResult = "Pending in queue"
	pass         VerificationStatusResult = "Pass - Verified"
	fail         VerificationStatusResult = "Fail - Unable to verify"
	unknown      VerificationStatusResult = "Unknown UID"
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
	Settings any                                 `json:"settings"`
}

type Contract struct {
	// A string with the name of the module to be invoked.
	//
	// Must be set to: contract
	module string

	// A string with the name of the action to be invoked.
	//
	// Must be set to: verifysourcecode
	action string

	// Format of sourceCode
	//
	// must be "solidity-standard-json-input"
	CodeFormat string `json:"codeformat"`

	// The address of the contract.
	ContractAddress common.Address `json:"contractaddress"`

	// The name of the contract.
	ContractName string `json:"contractname"`

	// The compiler version for the contract.
	CompilerVersion string `json:"compilerversion"`

	// The path to the source code of the contract.
	//
	// ContractSourceCodeFilePath is mutually exclusive with ContractSourceCode, setting ContractSourceCode will cause ContractSourceCodeFilePath to be ignored.
	ContractSourceCodeFilePath string `json:"-"`

	// The source code of the contract.
	//
	// SourceCode is mutually exclusive with ContractSourceCodeFilePath, setting SourceCode will cause ContractSourceCodeFilePath to be ignored.
	SourceCode string `json:"sourceCode,omitempty"`

	// The constructor argument data provided.
	ConstructorArguments string `json:"constructorArguments,omitempty"`

	// Whether or not automatically detect constructor argument.
	AutoDetectConstructorArguments bool `json:"autodetectConstructorArguments,omitempty"`

	// Mapping replacements
	ImportRemappings map[string]string `json:"-"`
}

type BlockScoutVerifyContractModel struct {
	ABI              string           `json:"ABI"`
	Address          common.Address   `json:"Address"`
	ContractName     string           `json:"ContractName"`
	OptimizationUsed OptimizationUsed `json:"OptimizationUsed"`
}

type BlockScoutContractVerifyResponse struct {
	Message string `json:"message"`
	Result  string `json:"result"`
	Status  string `json:"status"`
}

func NewContract(
	addressHash common.Address,
	name,
	sourceCodeFilePath,
	compilerVersion,
	constructorArguments string,
	autoDetectConstructorArguments bool,
	remap map[string]string,
) *Contract {
	return &Contract{
		module:                         module,
		action:                         action,
		CodeFormat:                     codeFormat,
		ContractAddress:                addressHash,
		ContractName:                   fmt.Sprintf("%s:%s", sourceCodeFilePath, name),
		CompilerVersion:                compilerVersion,
		ContractSourceCodeFilePath:     sourceCodeFilePath,
		ConstructorArguments:           constructorArguments,
		AutoDetectConstructorArguments: autoDetectConstructorArguments,
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

func MapImport(importPath, currentSource string, importRemapping map[string]string) (string, error) {
	// Realize full path to imported file
	// i.e. Convert ./IERC20.sol to @openzeppelin/contracts/token/ERC20/IERC20.sol
	separator := string(filepath.Separator)
	if strings.HasPrefix(importPath, ".") && currentSource != "" {
		if strings.HasPrefix(currentSource, "http") {
			separator = "/"
			importPath = fmt.Sprintf("%s/%s", currentSource, importPath)
		} else {
			importPath = path.Join(currentSource, importPath)
		}
	}

	// Remove any /./ from the path, they are useless
	importPath = strings.Replace(importPath, fmt.Sprintf("%s.%s", separator, separator), separator, -1)

	// loop to remove /../ from the path
	for {
		// split the path on the first /../
		before, after, ok := strings.Cut(importPath, fmt.Sprintf("%s..%s", separator, separator))
		if !ok {
			// break the loop when we run out of /../ sequences
			break
		}
		if before == ".." {
			// break the loop if the path begins with ../
			break
		}
		// get the parts before the first /../
		parts := strings.Split(before, separator)
		// make sure we can move up a dir
		if len(parts) < 1 {
			return "", fmt.Errorf("cannot move to a parent directory above %s", before)
		}
		// create a new path without the dir before the first /../
		importPath = strings.Join(append(parts[:len(parts)-1], after), separator)
		// repeat as needed
	}
	return importPath, nil
}

func parseContract(_path, statedPath, currentSource string, knownSources map[string]BlockscoutContractSource, importRemapping map[string]string) (map[string]BlockscoutContractSource, error) {
	// local var to parse the full import path
	importPath := _path

	// Realize full path to imported file
	// Convert an imported file's relative imports (i.e. ./IERC20.sol => @openzeppelin/contracts/token/ERC20/IERC20.sol)
	separator := string(filepath.Separator)
	if strings.HasPrefix(importPath, ".") && currentSource != "" {
		if strings.HasPrefix(currentSource, "http") {
			separator = "/"
			importPath = fmt.Sprintf("%s/%s", currentSource, importPath)
		} else {
			importPath = path.Join(currentSource, importPath)
		}
	}
	fmt.Printf("importing %s\n", importPath)

	// Remove any /./ from the path, they are useless
	importPath = strings.Replace(importPath, fmt.Sprintf("%s.%s", separator, separator), separator, -1)

	// loop to remove /../ from the path
	for {
		// split the path on the first /../
		before, after, ok := strings.Cut(importPath, fmt.Sprintf("%s..%s", separator, separator))
		if !ok {
			break
		}
		// get the parts before the first /../
		parts := strings.Split(before, separator)
		// make sure we can move up a dir
		if len(parts) < 1 {
			return nil, fmt.Errorf("cannot move to a parent directory above %s", before)
		}
		// create a new path without the dir before the first /../
		importPath = strings.Join(append(parts[:len(parts)-1], after), separator)
		// repeat as needed
	}
	// set the current source dir to the directory containing the current file
	currentSource = filepath.Dir(importPath)

	// map the import mappings so we can pull the file
	filePath := importPath
	for k, v := range importRemapping {
		if strings.Contains(importPath, k) {
			filePath = strings.Replace(importPath, k, v, 1)
			break
		}
	}

	fmt.Printf("\trealized import path %s\n", filePath)

	// determine if mapped import is a local or remote file
	source, err := ValidateSource(filePath)
	if err != nil {
		return nil, err
	}

	var b []byte
	switch source {
	case FilesystemContractSource:
		// read the file from the filesystem
		b, err = os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
	case HTTPContractSource:
		// read the file from the remote source
		resp, err := http.Get(filePath)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	case InvalidContractSource:
		// we don't know what it is...
		return nil, fmt.Errorf("cannot retrieve contract from %s", filePath)
	}

	// Save the source code for the current file
	knownSources[importPath] = BlockscoutContractSource{
		Content: string(b),
	}

	// Parse all imports from the file
	// var imports []string
	for _, match := range importRegex.FindAllSubmatch(b, -1) {
		for i, name := range importRegex.SubexpNames() {
			if i != 0 && name != "" {
				importPath, err := MapImport(string(match[i]), currentSource, importRemapping)
				if err != nil {
					return nil, fmt.Errorf("error mapping import path: %w", err)
				}

				// imports = append(imports, importPath)
				if _, ok := knownSources[importPath]; !ok {
					// only parse an import if we haven't seen it before, no need to duplicate imported files
					// recursion!
					importedSources, err := parseContract(importPath, string(match[i]), currentSource, knownSources, importRemapping)
					if err != nil {
						return nil, err
					}

					// save the imported imports to our known imports from the recursive imports
					for k, v := range importedSources {
						knownSources[k] = v
					}
				}
			}
		}
	}

	return knownSources, nil
}

// Public function the begins the recursive process of parsing a contract.
// It will follow imports in a given contract, pulling from your local file system or HTTP sources
func ParseContract(_path string, knownSources map[string]BlockscoutContractSource, importRemapping map[string]string) (map[string]BlockscoutContractSource, error) {
	return parseContract(_path, _path, "", knownSources, importRemapping)
}

// Verifies a contract with the verifysourcecode method on Etherscan based APIs
func VerifyContract(explorerAPI string, contractToVerify *Contract) error {
	if contractToVerify.module != module {
		return fmt.Errorf("at the time of writing this client, Blockscout only supports a module value of '%s', received '%s'", module, contractToVerify.module)
	}

	if contractToVerify.action != action {
		return fmt.Errorf("at the time of writing this client, Blockscout only supports an action value of '%s', received '%s'", action, contractToVerify.action)
	}

	fmt.Printf("uploading with compiler version %s\n", contractToVerify.CompilerVersion)

	if contractToVerify.SourceCode == "" {
		sources, err := ParseContract(contractToVerify.ContractSourceCodeFilePath, make(map[string]BlockscoutContractSource), contractToVerify.ImportRemappings)
		if err != nil {
			return fmt.Errorf("error parsing contract: %w", err)
		}

		keys := make([]string, 0, len(sources))
		for k := range sources {
			keys = append(keys, k)
		}
		fmt.Printf("found sources: %v\n", keys)

		jsonData, err := json.Marshal(&BlockscoutContract{
			Language: "Solidity",
			Sources:  sources,
			Settings: blockscoutSettings,
		})
		if err != nil {
			return fmt.Errorf("error marshalling contracts as json: %w", err)
		}
		contractToVerify.SourceCode = string(jsonData)
	}

	jsonData, err := json.Marshal(contractToVerify)
	if err != nil {
		return fmt.Errorf("error marshalling contract to json: %w", err)
	}

	uri := fmt.Sprintf(
		"%s?module=%s&action=%s",
		explorerAPI,
		url.QueryEscape(contractToVerify.module),
		url.QueryEscape(contractToVerify.action),
	)
	fmt.Printf("verifying contract at %s\n", uri)

	resp, err := http.Post(
		uri,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("error sending contract to blockscout: %w", err)
	}
	defer resp.Body.Close()

	var r BlockScoutContractVerifyResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return fmt.Errorf("error reading blockscout response: %w", err)
	}
	switch r.Result {
	case alreadyVerifiedContractResponse:
		fmt.Printf("cannot re-verify this contract: %s\n", r.Result)
	default:
		fmt.Printf("check the status of your verification with this guid: %s\n", r.Result)
		fmt.Printf("\twasp-cli contract check-status %s %s\n", explorerAPI, r.Result)
	}

	return nil
}

func CheckVerificationStatus(explorerAPI, guid string) error {
	uri := fmt.Sprintf(
		"%s?module=%s&action=%s&guid=%s",
		explorerAPI,
		url.QueryEscape(module),
		url.QueryEscape(verifyAction),
		url.QueryEscape(guid),
	)

	fmt.Printf("checking contract verification status at %s\n", uri)

	resp, err := http.Post(
		uri,
		"application/json",
		nil,
	)
	if err != nil {
		return fmt.Errorf("error sending contract to blockscout: %w", err)
	}
	defer resp.Body.Close()

	var r BlockScoutContractVerifyResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return fmt.Errorf("error reading blockscout response: %w", err)
	}
	switch VerificationStatusResult(r.Result) {
	case pending:
		fmt.Printf("contract verification is still pending: %s\n", r.Result)
	case pass:
		fmt.Printf("contract successfully verified: %s\n", r.Result)
	case fail:
		fmt.Printf("failed to verify the contract: %s\n", r.Result)
		fmt.Println("If your contract was recently uploaded, it may not have been indexed yet. Try again in a few minutes.")
	case unknown:
		fallthrough
	default:
		fmt.Printf("something went wrong: %s\n", r.Result)
	}

	return nil
}
