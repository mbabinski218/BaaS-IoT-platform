// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "forge-std/Test.sol";
import "smartContracts/BatchRegistry.sol";

contract BatchRegistryTest is Test {
    BatchRegistry public registry;

    bytes32 h0;
    bytes32 h1;
    bytes32 h2;
    bytes32 h3;
    bytes32 hash0_1;
    bytes32 hash2_3;
    bytes32 root;
    uint256 batchTime;

    function setUp() public {
        registry = new BatchRegistry();
    }

    function testStoreAndVerifyRootCorrect() public {
        batchTime = 1234567890;
        root = sha256("root-node");

        registry.storeRoot(batchTime, root);

        assertEq(registry.roots(batchTime), root);
        assertTrue(registry.verifyRoot(batchTime, root));
    }

    function testStoreAndVerifyRootIncorrectRoot() public {
        batchTime = 1234567890;
        root = sha256("root-node");
        bytes32 incorrectRoot = sha256("incorrect");

        registry.storeRoot(batchTime, root);

        assertEq(registry.roots(batchTime), root);
        assertFalse(registry.verifyRoot(batchTime, incorrectRoot));
    }

    function testStoreAndVerifyRootIncorrectBatchTime() public {
        batchTime = 1234567890;
        root = sha256("root-node");
        uint256 incorrectBatchTime = 99999999;

        registry.storeRoot(batchTime, root);

        assertEq(registry.roots(batchTime), root);
        assertFalse(registry.verifyRoot(incorrectBatchTime, root));
    }

    function testRevertOnDuplicateRoot() public {
        batchTime = 1111;
        root = sha256("root-hash");

        registry.storeRoot(batchTime, root);

        vm.expectRevert("Already stored for this batch time");
        registry.storeRoot(batchTime, root);
    }

    function setUpForMerkleProof() private {
        batchTime = 123;

        h0 = sha256(abi.encodePacked("a"));
        h1 = sha256(abi.encodePacked("b"));
        h2 = sha256(abi.encodePacked("c"));
        h3 = sha256(abi.encodePacked("d"));

        hash0_1 = h0 < h1
            ? sha256(abi.encodePacked(h0, h1))
            : sha256(abi.encodePacked(h1, h0));

        hash2_3 = h2 < h3
            ? sha256(abi.encodePacked(h2, h3))
            : sha256(abi.encodePacked(h3, h2));

        root = hash0_1 < hash2_3
            ? sha256(abi.encodePacked(hash0_1, hash2_3))
            : sha256(abi.encodePacked(hash2_3, hash0_1));

        registry.storeRoot(batchTime, root);
    }

    function testVerifyProofSuccess() public {
        setUpForMerkleProof();

        bytes32[] memory proof = new bytes32[](2);
        proof[0] = h1;
        proof[1] = hash2_3;

        bool result = registry.verifyProof(batchTime, h0, proof);
        assertTrue(result);
    }

    function testVerifyProofFail() public {
        setUpForMerkleProof();

        bytes32 fakeLeaf = sha256(abi.encodePacked("fake"));

        bytes32[] memory proof = new bytes32[](2);
        proof[0] = h1;
        proof[1] = hash2_3;

        bool result = registry.verifyProof(batchTime, fakeLeaf, proof);
        assertFalse(result);
    }
}