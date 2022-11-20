// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.17;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

contract ClinkerERC721 is ERC721URIStorage, Ownable {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIds;
    
    enum Status {
        PENDING,
        USING,
        DISUSED
    }

    struct User {
        address requester;  // which wanted to be a new address
        Status status;
    }

    mapping(address => User) users;
    mapping(address => address) linker;

    event AddressAvailable(address indexed available);

    constructor() ERC721("Clinker", "CLK") {}

    function mint(address user, string memory tokenURI) public onlyOwner returns (uint256) {
        uint256 newItemId = _tokenIds.current();
        _mint(user, newItemId);
        _setTokenURI(newItemId, tokenURI);

        _tokenIds.increment();
        return newItemId;
    }

    modifier isNew() {
        require(msg.sender != address(0), "sender is 0");
        require(users[msg.sender].status == Status.PENDING, "already used address");
        _;
    }

    modifier isPrevious() {
        require(msg.sender != address(0), "sender is 0");
        require(users[msg.sender].status == Status.USING, "address is not available");
        require(users[msg.sender].requester != address(0), "request not exists");
        _;
    }

    function _setNewAddress(address newAddr, address prevAddr) internal {
        users[newAddr] = User(address(0), Status.USING);
        linker[newAddr] = prevAddr;
        emit AddressAvailable(newAddr);
    }

    function proveNew(address previous) external isNew {
        if (previous == address(0)) {
            _setNewAddress(msg.sender, previous);
        } else {
            users[previous].requester = msg.sender;
        }
    }

    function approveModify() external isPrevious {
        users[msg.sender].status = Status.DISUSED;
        _setNewAddress(users[msg.sender].requester, msg.sender);
    }

    function getLinkedAddress(address a) external view returns (address) {
        return linker[a];
    }

    function getStatus(address a) external view returns (Status) {
        return users[a].status;
    }
}
