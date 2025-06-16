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

    function verifyProof(uint256 batchTime, bytes32 leaf, bytes32[] calldata proof) external view returns (bool) {
        bytes32 computedHash = leaf;

        for (uint256 i = 0; i < proof.length; i++) {
            bytes32 proofElement = proof[i];
            
            if (computedHash < proofElement) {
                computedHash = sha256(abi.encodePacked(computedHash, proofElement));
            } else {
                computedHash = sha256(abi.encodePacked(proofElement, computedHash));
            }
        }

        return computedHash == roots[batchTime];
    }
}