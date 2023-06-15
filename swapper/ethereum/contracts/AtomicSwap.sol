// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.17;

contract AtomicSwap {
    address immutable redeemer;
    address immutable refunder;
    bytes32 immutable secretHash;
    uint256 immutable expiry;

    event redeemed(string _secret);

    constructor(
        address _redeemer,
        address _refunder,
        bytes32 _secretHash,
        uint256 _expiry
    ) {
        redeemer = _redeemer;
        refunder = _refunder;
        secretHash = _secretHash;
        expiry = _expiry;
    }

    function execute(address _token, bytes memory _secret) external {
        if (keccak256(_secret) == secretHash) {
            emit redeemed(string(_secret));
            _transferBalance(_token, redeemer);
        } else if (block.number > expiry) {
            _transferBalance(_token, refunder);
        } else {
            revert("invalid secret or expiry");
        }
    }

    function _transferBalance(address _token, address _to) internal {
        // Transfer ERC20 balance
        (bool _ok1, bytes memory balanceData) = _token.call(
            abi.encodeWithSelector(0x70a08231, address(this))
        );
        require(_ok1, "AtomicSwap: ERC20 balanceOf did not succeed");
        uint256 _amount = abi.decode(balanceData, (uint256));
        if (_amount > 0) {
            (bool _ok2, bytes memory transferData) = _token.call(
                abi.encodeWithSelector(0xa9059cbb, _to, _amount)
            );
            require(_ok2, "AtomicSwap: ERC20 transfer did not succeed");
            if (transferData.length > 0) {
                // Return data is optional
                require(
                    abi.decode(transferData, (bool)),
                    "AtomicSwap: ERC20 transfer did not succeed"
                );
            }
        }
        selfdestruct(payable(_to));
    }
}
