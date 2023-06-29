// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.18;
uint256 constant FEE = 30;

contract AtomicSwap {
    address immutable redeemer;
    address immutable initiator;
    address immutable token;
    address immutable feeCollecter;
    bytes32 immutable secretHash;
    uint256 immutable expiry;
    uint256 immutable amount;

    event Redeemed(bytes _secret);

    constructor(
        address _redeemer,
        address _initiator,
        address _token,
        address _feeCollector,
        bytes32 _secretHash,
        uint256 _expiry,
        uint256 _amount
    ) {
        redeemer = _redeemer;
        initiator = _initiator;
        token = _token;
        secretHash = _secretHash;
        expiry = _expiry;
        feeCollecter = _feeCollector;
        amount = _amount;
    }

    function redeem(bytes calldata secret) external {
        require(sha256(secret) == secretHash, "AtomicSwap: secret invalid");
        uint256 balance = _getBalance();
        require(balance >= amount, "AtomicSwap: contract not initiated");
        uint256 fee = (balance * FEE) / 10000;
        _transfer(feeCollecter, fee);
        _transfer(redeemer, amount - fee);
        emit Redeemed(secret);
    }

    function refund() external {
        require(block.number > expiry, "AtomicSwap: lock not expired");
        uint256 balance = _getBalance();
        _transfer(initiator, balance);
    }

    function _getBalance() internal returns (uint256 balance) {
        // Transfer ERC20 balance
        (bool _ok, bytes memory balanceData) = token.call(
            abi.encodeWithSelector(0x70a08231, address(this))
        );
        require(_ok, "AtomicSwap: ERC20 balanceOf did not succeed");
        require(
            balanceData.length > 0,
            "AtomicSwap: ERC20 balanceOf did not return data"
        );
        balance = abi.decode(balanceData, (uint256));
    }

    function _transfer(address _to, uint256 _amount) internal {
        (bool _ok, bytes memory transferData) = token.call(
            abi.encodeWithSelector(0xa9059cbb, _to, _amount)
        );
        require(_ok, "AtomicSwap: ERC20 transfer did not succeed (bool)");
        require(
            transferData.length > 0,
            "AtomicSwap: ERC20 transfer did not return data"
        );
        require(
            abi.decode(transferData, (bool)),
            "AtomicSwap: ERC20 transfer failed"
        );
    }
}
