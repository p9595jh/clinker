// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.17;

contract Sample {
    event SampleEvent(address indexed a);

    function hello(address a) public returns (address) {
        emit SampleEvent(a);
        return msg.sender;
    }

    function hi() public view returns (uint256) {
        // emit SampleEvent(address(0));
        return block.number;
    }

    function what(address a) public returns (uint256) {
        emit SampleEvent(a);
        return block.number;
    }
}
