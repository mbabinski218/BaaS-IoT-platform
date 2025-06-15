// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BatchRegistry {
    mapping(uint256 => bytes32) public roots;

    function storeRoot(uint256 batchTime, bytes32 merkleRoot) external {
        require(roots[batchTime] == 0, "Already stored for this batch time");

        roots[batchTime] = merkleRoot;
    }

    function verifyRoot(uint256 batchTime, bytes32 providedMerkleRoot) external view returns (bool) {
        return roots[batchTime] == providedMerkleRoot;
    }
}