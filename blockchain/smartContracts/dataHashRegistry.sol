// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

contract DataHashRegistry {

    struct DataRecord { 
        bytes32 dataHash;
        bytes16 iotId;
        address sender;
    } 

    mapping(bytes16 => DataRecord) public records;

    event HashStored(bytes16 indexed id, bytes32 dataHash);

    function storeHash(bytes16 id, bytes32 dataHash, bytes16 iotId) external {
        require(id != bytes16(0), "Zero id not allowed");
        require(dataHash != bytes32(0), "Zero hash not allowed");
        require(iotId != bytes16(0), "Zero iotId not allowed");

        DataRecord storage record = records[id];
        require(record.dataHash == bytes32(0), "Hash already exists for this Id");
        
        records[id] = DataRecord(dataHash, iotId, msg.sender);

        emit HashStored(id, dataHash);
    }

    function verifyHash(bytes16 id, bytes32 providedHash) external view returns (bool) {
        require(id == bytes16(0), "Zero id not allowed");
        require(providedHash == bytes32(0), "Zero providedHash not allowed");

        DataRecord storage record = records[id];
        return record.dataHash == providedHash;
    }
}