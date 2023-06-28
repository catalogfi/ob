// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.18;

/**
 * @dev Interface of the ERC20 standard as defined in the EIP.
 */
interface IERC20 {
    /**
     * @dev Emitted when `value` tokens are moved from one account (`from`) to
     * another (`to`).
     *
     * Note that `value` may be zero.
     */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /**
     * @dev Emitted when the allowance of a `spender` for an `owner` is set by
     * a call to {approve}. `value` is the new allowance.
     */
    event Approval(
        address indexed owner,
        address indexed spender,
        uint256 value
    );

    /**
     * @dev Returns the amount of tokens in existence.
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev Returns the amount of tokens owned by `account`.
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev Moves `amount` tokens from the caller's account to `to`.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transfer(address to, uint256 amount) external returns (bool);

    /**
     * @dev Returns the remaining number of tokens that `spender` will be
     * allowed to spend on behalf of `owner` through {transferFrom}. This is
     * zero by default.
     *
     * This value changes when {approve} or {transferFrom} are called.
     */
    function allowance(
        address owner,
        address spender
    ) external view returns (uint256);

    /**
     * @dev Sets `amount` as the allowance of `spender` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * IMPORTANT: Beware that changing an allowance with this method brings the risk
     * that someone may use both the old and the new allowance by unfortunate
     * transaction ordering. One possible solution to mitigate this race
     * condition is to first reduce the spender's allowance to 0 and set the
     * desired value afterwards:
     * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
     *
     * Emits an {Approval} event.
     */
    function approve(address spender, uint256 amount) external returns (bool);

    /**
     * @dev Moves `amount` tokens from `from` to `to` using the
     * allowance mechanism. `amount` is then deducted from the caller's
     * allowance.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) external returns (bool);
}

contract AtomicSwapV1 {
    IERC20 immutable token;

    struct Order {
        address redeemer;
        address initiator;
        uint256 expiry;
        uint256 amount;
        bool isFullfilled;
    }
    mapping(bytes32 => Order) AtomicSwapOrders;

    event Redeemed(bytes32 indexed secrectHash, bytes _secret);
    event Initiated(bytes32 indexed secrectHash, uint256 amount);
    event Refunded(bytes32 indexed secrectHash);

    modifier checkSafe(
        address redeemer,
        address intiator,
        uint256 expiry
    ) {
        _;
        require(
            redeemer != address(0),
            "AtomicSwap: redeemer cannot be null address"
        );
        require(
            intiator != redeemer,
            "AtomicSwap: initiator cannot be equal to redeemer"
        );
        require(
            expiry > block.number,
            "AtomicSwap: expiry cannot be lower than current block"
        );
    }

    constructor(address _token) {
        token = IERC20(_token);
    }

    function initiate(
        address _redeemer,
        uint256 _expiry,
        uint256 _amount,
        bytes32 _secretHash
    ) external checkSafe(_redeemer, msg.sender, _expiry) {
        Order memory order = AtomicSwapOrders[_secretHash];
        require(!order.isFullfilled, "AtomicSwap: cannot reuse secret");
        require(
            order.redeemer == address(0x0),
            "AtomicSwap: order already fullfilled"
        );
        Order memory newOrder = Order({
            redeemer: _redeemer,
            initiator: msg.sender,
            expiry: _expiry,
            amount: _amount,
            isFullfilled: false
        });
        token.transferFrom(msg.sender, address(this), newOrder.amount);
        AtomicSwapOrders[_secretHash] = newOrder;
        emit Initiated(_secretHash, newOrder.amount);
    }

    function redeem(bytes calldata _secret) external {
        bytes32 secretHash = sha256(_secret);
        Order storage order = AtomicSwapOrders[secretHash];
        require(
            order.redeemer != address(0x0),
            "AtomicSwap: invalid secret or order not initiated"
        );
        require(!order.isFullfilled, "AtomicSwap: order already fullfilled");
        order.isFullfilled = true;
        token.transfer(order.redeemer, order.amount);
        emit Redeemed(secretHash, _secret);
    }

    function refund(bytes32 _secretHash) external {
        Order storage order = AtomicSwapOrders[_secretHash];
        require(
            order.redeemer != address(0x0),
            "AtomicSwap: order not initated"
        );
        require(!order.isFullfilled, "AtomicSwap: order already fullfilled");
        require(block.number > order.expiry, "AtomicSwap: lock not expired");
        order.isFullfilled = true;
        token.transfer(order.initiator, order.amount);
        emit Refunded(_secretHash);
    }
}
