// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
pragma abicoder v2;

contract ExternalCall {
    struct Call {
        address target;
        uint256 gasLimit;
        bytes callData;
    }

    function callExternal(Call memory call) public returns (bool) {
        (bool success, bytes memory ret) = call.target.call{gas: call.gasLimit}(call.callData);
        return success;
    }
}