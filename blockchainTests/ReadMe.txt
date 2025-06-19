Tests:
Powershell
- cd .\blockchainTests\
- forge test
- forge test --json > test-report.json     
- forge test --gas-report > gas-report.txt 


Coverage:
Powershell
- cd .\blockchainTests\
- forge coverage --report lcov
- generate html:
    - MSYS2
    - cd /d/Studia_mgr/Praca_magisterska/Code/BaaS-IoT-platform/blockchainTests
    - pacman -S lcov sed coreutils
    - genhtml lcov.info -o coverage-html


Security:
Powershell
cd .\blockchain\smartContracts\
slither . --json slither-report.json