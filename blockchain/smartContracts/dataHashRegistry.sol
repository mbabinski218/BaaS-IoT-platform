// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract DataHashRegistry {

    struct DataRecord { 
        bytes32 dataHash;
        bytes16 iotId;
        uint256 timestamp;
        address sender;
    } 

    mapping(bytes16 => DataRecord) public records; 

    function storeHash(bytes16 id, bytes32 dataHash, bytes16 iotId) external {
        require(records[id].timestamp == 0, "Hash already exists for this Id");

        records[id] = DataRecord(dataHash, iotId, block.timestamp, msg.sender);
    }

    function verifyHash(bytes16 id, bytes32 providedHash) external view returns (bool) {
        return records[id].dataHash == providedHash;
    }
}