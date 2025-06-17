// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "forge-std/Test.sol";
import "smartContracts/BatchRegistry.sol";

contract BatchRegistryTest is Test {
    BatchRegistry public registry;

    function setUp() public {
        registry = new BatchRegistry();
    }

    function testStoreAndVerifyRootCorrect() public {
        uint256 batchTime = 1234567890;
        bytes32 root = sha256("root-node");

        registry.storeRoot(batchTime, root);

        assertEq(registry.roots(batchTime), root);
        assertTrue(registry.verifyRoot(batchTime, root));
    }

    function testStoreAndVerifyRootIncorrectRoot() public {
        uint256 batchTime = 1234567890;
        bytes32 root = sha256("root-node");
        bytes32 incorrectRoot = sha256("incorrect");

        registry.storeRoot(batchTime, root);

        assertEq(registry.roots(batchTime), root);
        assertFalse(registry.verifyRoot(batchTime, incorrectRoot));
    }

    function testStoreAndVerifyRootIncorrectBatchTime() public {
        uint256 batchTime = 1234567890;
        uint256 incorrectBatchTime = 99999999;
        bytes32 root = sha256("root-node");

        registry.storeRoot(batchTime, root);

        assertEq(registry.roots(batchTime), root);
        assertFalse(registry.verifyRoot(incorrectBatchTime, root));
    }

    function testRevertOnDuplicateRoot() public {
        uint256 batchTime = 1111;
        bytes32 root = sha256("root-hash");

        registry.storeRoot(batchTime, root);

        vm.expectRevert("Already stored for this batch time");
        registry.storeRoot(batchTime, root);
    }

    // function testVerifyMerkleProof() public {
    //     bytes32 leaf = sha256("leaf");
    //     bytes32 sibling = sha256("sibling");
    //     bytes32 root;

    //     // simulate Merkle: sha256(min||max)
    //     if (leaf < sibling) {
    //         root = sha256(abi.encodePacked(leaf, sibling));
    //     } else {
    //         root = sha256(abi.encodePacked(sibling, leaf));
    //     }

    //     uint256 batchTime = 999;
    //     registry.storeRoot(batchTime, root);

    //     bytes32 ;
    //     proof[0] = sibling;

    //     assertTrue(registry.verifyProof(batchTime, leaf, proof));
    // }
}