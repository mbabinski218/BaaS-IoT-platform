// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract IoTDataContract {

    struct IoTData {
        string deviceId;
        string encryptedData;
        uint256 timestamp;
    }

    mapping(string => mapping(string => IoTData[])) public deviceData;

    event DataStored(string deviceId, string deviceType, string region, string encryptedData, uint256 timestamp);

    function storeData(string memory deviceId, string memory deviceType, string memory region, string memory encryptedData) public {
        require(bytes(deviceId).length > 0, "Device id is required");
        require(bytes(deviceType).length > 0, "Device type is required");
        require(bytes(region).length > 0, "Region is required");
        require(bytes(encryptedData).length > 0, "Encrypted data is required");

        IoTData memory newData = IoTData({
            deviceId: deviceId,
            encryptedData: encryptedData,
            timestamp: block.timestamp
        });

        deviceData[deviceType][region].push(newData);
        emit DataStored(deviceId, deviceType, region, encryptedData, block.timestamp);
    }

    function getData(string memory deviceType, string memory region) public view returns (IoTData[] memory) {
        return deviceData[deviceType][region];
    }
}