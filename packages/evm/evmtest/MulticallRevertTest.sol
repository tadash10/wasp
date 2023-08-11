// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MulticallRevertTest {
    uint32 private count = 0;

    function testRevert() public returns (uint32) {
        count = count + 1;
        require(false, "rip");
        return count;
    }

    function callRevertTestByItself() public returns (uint32) {
        count = count + 1;
        return this.testRevert();
    }

    function increaseCount() public {
        count = count + 1;
    }

    function getCount() public view returns (uint32) {
        return count;
    }
}