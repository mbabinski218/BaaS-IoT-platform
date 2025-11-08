// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "forge-std/Test.sol";
import "smartContracts/DataHashRegistry.sol";

contract DataHashRegistryTest is Test {
    DataHashRegistry public registry;

    function setUp() public {
        registry = new DataHashRegistry();
    }

    function testStoreAndVerifyHashCorrect() public {
        bytes16 id = 0x11112222333344445555666677778888;
        bytes32 hash = sha256("example-data");
        bytes16 iotId = 0x99990000111122223333444455556666;

        registry.storeHash(id, hash, iotId);

        (bytes32 storedHash, bytes16 storedIotId, address sender) = registry.records(id);
        assertEq(storedHash, hash);
        assertEq(sender, address(this));
        assertTrue(registry.verifyHash(id, hash));
    }

    function testStoreAndVerifyHashIncorrect() public {
        bytes16 id = 0x11112222333344445555666677778888;
        bytes32 hash = sha256("example-data");
        bytes32 incorrectHash = sha256("incorrect");
        bytes16 iotId = 0x99990000111122223333444455556666;

        registry.storeHash(id, hash, iotId);

        (bytes32 storedHash, bytes16 storedIotId, address sender) = registry.records(id);
        assertEq(storedHash, hash);
        assertEq(sender, address(this));
        assertFalse(registry.verifyHash(id, incorrectHash));
    }

    function testRevertOnDuplicateStore() public {
        bytes16 id = 0xABCDEFABCDEFABCDEFABCDEFABCDEFAB;
        bytes32 hash = sha256("some-data");
        bytes16 iotId = 0x11112222333344445555666677778888;

        registry.storeHash(id, hash, iotId);

        vm.expectRevert("Hash already exists for this Id");
        registry.storeHash(id, hash, iotId);
    }
}