// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import * as wasmtypes from '../wasmtypes';

export const ScName        = 'governance';
export const ScDescription = 'Governance contract';
export const HScName       = new wasmtypes.ScHname(0x17cf909f);

export const ParamAccessAPI     = 'ia';
export const ParamAccessOnly    = 'i';
export const ParamActions       = 'n';
export const ParamAddress       = 'S';
export const ParamCertificate   = 'ic';
export const ParamChainOwner    = 'o';
export const ParamFeePolicy     = 'g';
export const ParamGasLimits     = 'l';
export const ParamGasRatio      = 'e';
export const ParamMetadata      = 'md';
export const ParamPayoutAddress = 's';
export const ParamPubKey        = 'ip';
export const ParamPublicURL     = 'x';
export const ParamSetMinSD      = 'ms';

export const ResultAccessNodeCandidates = 'an';
export const ResultAccessNodes          = 'ac';
export const ResultChainID              = 'c';
export const ResultChainOwner           = 'o';
export const ResultChainOwnerID         = 'o';
export const ResultControllers          = 'a';
export const ResultFeePolicy            = 'g';
export const ResultGasLimits            = 'l';
export const ResultGasRatio             = 'e';
export const ResultGetMinSD             = 'ms';
export const ResultMetadata             = 'md';
export const ResultPayoutAddress        = 's';
export const ResultPublicURL            = 'x';
export const ResultStatus               = 'm';

export const FuncAddAllowedStateControllerAddress    = 'addAllowedStateControllerAddress';
export const FuncAddCandidateNode                    = 'addCandidateNode';
export const FuncChangeAccessNodes                   = 'changeAccessNodes';
export const FuncClaimChainOwnership                 = 'claimChainOwnership';
export const FuncDelegateChainOwnership              = 'delegateChainOwnership';
export const FuncRemoveAllowedStateControllerAddress = 'removeAllowedStateControllerAddress';
export const FuncRevokeAccessNode                    = 'revokeAccessNode';
export const FuncRotateStateController               = 'rotateStateController';
export const FuncSetEVMGasRatio                      = 'setEVMGasRatio';
export const FuncSetFeePolicy                        = 'setFeePolicy';
export const FuncSetGasLimits                        = 'setGasLimits';
export const FuncSetMetadata                         = 'setMetadata';
export const FuncSetMinSD                            = 'setMinSD';
export const FuncSetPayoutAddress                    = 'setPayoutAddress';
export const FuncStartMaintenance                    = 'startMaintenance';
export const FuncStopMaintenance                     = 'stopMaintenance';
export const ViewGetAllowedStateControllerAddresses  = 'getAllowedStateControllerAddresses';
export const ViewGetChainInfo                        = 'getChainInfo';
export const ViewGetChainNodes                       = 'getChainNodes';
export const ViewGetChainOwner                       = 'getChainOwner';
export const ViewGetEVMGasRatio                      = 'getEVMGasRatio';
export const ViewGetFeePolicy                        = 'getFeePolicy';
export const ViewGetGasLimits                        = 'getGasLimits';
export const ViewGetMaintenanceStatus                = 'getMaintenanceStatus';
export const ViewGetMetadata                         = 'getMetadata';
export const ViewGetMinSD                            = 'getMinSD';
export const ViewGetPayoutAddress                    = 'getPayoutAddress';

export const HFuncAddAllowedStateControllerAddress    = new wasmtypes.ScHname(0x9469d567);
export const HFuncAddCandidateNode                    = new wasmtypes.ScHname(0xb745b382);
export const HFuncChangeAccessNodes                   = new wasmtypes.ScHname(0x7bca3700);
export const HFuncClaimChainOwnership                 = new wasmtypes.ScHname(0x03ff0fc0);
export const HFuncDelegateChainOwnership              = new wasmtypes.ScHname(0x93ecb6ad);
export const HFuncRemoveAllowedStateControllerAddress = new wasmtypes.ScHname(0x31f69447);
export const HFuncRevokeAccessNode                    = new wasmtypes.ScHname(0x5459512d);
export const HFuncRotateStateController               = new wasmtypes.ScHname(0x244d1038);
export const HFuncSetEVMGasRatio                      = new wasmtypes.ScHname(0xaae22338);
export const HFuncSetFeePolicy                        = new wasmtypes.ScHname(0x5b791c9f);
export const HFuncSetGasLimits                        = new wasmtypes.ScHname(0xd72fb355);
export const HFuncSetMetadata                         = new wasmtypes.ScHname(0x0eb3a798);
export const HFuncSetMinSD                            = new wasmtypes.ScHname(0x9cad5084);
export const HFuncSetPayoutAddress                    = new wasmtypes.ScHname(0x65e7c531);
export const HFuncStartMaintenance                    = new wasmtypes.ScHname(0x742f0521);
export const HFuncStopMaintenance                     = new wasmtypes.ScHname(0x4e017b6a);
export const HViewGetAllowedStateControllerAddresses  = new wasmtypes.ScHname(0xf3505183);
export const HViewGetChainInfo                        = new wasmtypes.ScHname(0x434477e2);
export const HViewGetChainNodes                       = new wasmtypes.ScHname(0xe1832289);
export const HViewGetChainOwner                       = new wasmtypes.ScHname(0x9b2ef0ac);
export const HViewGetEVMGasRatio                      = new wasmtypes.ScHname(0xb81c8c34);
export const HViewGetFeePolicy                        = new wasmtypes.ScHname(0xf8c89790);
export const HViewGetGasLimits                        = new wasmtypes.ScHname(0x3a493455);
export const HViewGetMaintenanceStatus                = new wasmtypes.ScHname(0x61fe5443);
export const HViewGetMetadata                         = new wasmtypes.ScHname(0x79ad1ac6);
export const HViewGetMinSD                            = new wasmtypes.ScHname(0x37f53a59);
export const HViewGetPayoutAddress                    = new wasmtypes.ScHname(0x2af7a8c3);
