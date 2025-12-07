// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title CommodityDTARequestManager
 * @notice Manages subscription and redemption requests for commodity-backed tokens in the Digital Transfer Agent flow.
 * @dev This contract tracks the lifecycle of DTA requests from creation through settlement.
 * 
 * Request lifecycle statuses:
 *  0 = Pending      - Request created, awaiting pricing and compliance checks
 *  1 = Priced       - Price calculated, ready for compliance review
 *  2 = Approved     - Passed compliance, ready for settlement
 *  3 = Rejected     - Failed compliance or other checks
 *  4 = Settled      - Settlement complete (tokens minted/burned, collateral transferred)
 *  5 = Cancelled    - Request cancelled by requester or timeout
 */
contract CommodityDTARequestManager {
    // ============ Structs ============

    struct Request {
        uint256 id;
        address requester;
        uint8 requestType;      // 0 = Subscription, 1 = Redemption
        uint256 quantityTons;   // Commodity amount in tons (scaled appropriately)
        address paymentToken;   // Payment token address (e.g., USDC)
        uint256 paymentAmount;  // Payment token amount (for subscriptions)
        uint8 status;           // Request lifecycle status (0-5)
        uint256 createdAt;
        uint256 updatedAt;
    }

    // ============ State Variables ============

    uint256 private s_nextRequestId;
    mapping(uint256 => Request) private s_requests;

    // TODO: Add role-based access control for updateRequestStatus
    // For MVP, any address can update status; production should restrict to authorized orchestrator

    // ============ Events ============

    event SubscriptionRequested(
        uint256 indexed id,
        address indexed requester,
        uint256 quantityTons,
        address paymentToken
    );

    event RedemptionRequested(
        uint256 indexed id,
        address indexed requester,
        uint256 quantityTons
    );

    event RequestStatusUpdated(
        uint256 indexed id,
        uint8 status
    );

    // ============ Constructor ============

    constructor() {
        s_nextRequestId = 1; // Start IDs at 1 for easier validation
    }

    // ============ External Functions ============

    /**
     * @notice Create a new subscription request (buy commodity-backed tokens)
     * @param quantityTons Amount of commodity in tons requested
     * @param paymentToken ERC20 token address for payment (e.g., USDC)
     * @return requestId The ID of the newly created request
     */
    function createSubscriptionRequest(
        uint256 quantityTons,
        address paymentToken
    ) external returns (uint256) {
        require(quantityTons > 0, "Quantity must be > 0");
        require(paymentToken != address(0), "Invalid payment token");

        uint256 requestId = s_nextRequestId++;
        
        s_requests[requestId] = Request({
            id: requestId,
            requester: msg.sender,
            requestType: 0, // Subscription
            quantityTons: quantityTons,
            paymentToken: paymentToken,
            paymentAmount: 0, // Will be set after pricing
            status: 0, // Pending
            createdAt: block.timestamp,
            updatedAt: block.timestamp
        });

        emit SubscriptionRequested(requestId, msg.sender, quantityTons, paymentToken);
        
        return requestId;
    }

    /**
     * @notice Create a new redemption request (redeem tokens for physical commodity or cash)
     * @param quantityTons Amount of commodity in tons to redeem
     * @return requestId The ID of the newly created request
     */
    function createRedemptionRequest(
        uint256 quantityTons
    ) external returns (uint256) {
        require(quantityTons > 0, "Quantity must be > 0");

        uint256 requestId = s_nextRequestId++;
        
        s_requests[requestId] = Request({
            id: requestId,
            requester: msg.sender,
            requestType: 1, // Redemption
            quantityTons: quantityTons,
            paymentToken: address(0), // Not applicable for redemptions
            paymentAmount: 0,
            status: 0, // Pending
            createdAt: block.timestamp,
            updatedAt: block.timestamp
        });

        emit RedemptionRequested(requestId, msg.sender, quantityTons);
        
        return requestId;
    }

    /**
     * @notice Update the status of a request
     * @dev TODO: Add access control - only authorized orchestrator should call this
     * @param id Request ID
     * @param newStatus New status value (0-5)
     */
    function updateRequestStatus(uint256 id, uint8 newStatus) external {
        require(id > 0 && id < s_nextRequestId, "Invalid request ID");
        require(newStatus <= 5, "Invalid status");
        
        // TODO: Add role check - require(hasRole(ORCHESTRATOR_ROLE, msg.sender), "Not authorized");
        
        Request storage request = s_requests[id];
        request.status = newStatus;
        request.updatedAt = block.timestamp;

        emit RequestStatusUpdated(id, newStatus);
    }

    /**
     * @notice Get request details by ID
     * @param id Request ID
     * @return Request struct
     */
    function getRequest(uint256 id) external view returns (Request memory) {
        require(id > 0 && id < s_nextRequestId, "Invalid request ID");
        return s_requests[id];
    }

    /**
     * @notice Get the next request ID (for external tracking)
     * @return Next request ID that will be assigned
     */
    function getNextRequestId() external view returns (uint256) {
        return s_nextRequestId;
    }
}

