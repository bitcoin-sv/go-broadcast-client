package broadcast

// TxStatus is the status of the transaction
type TxStatus string

// List of statuses available here: https://github.com/bitcoin-sv/arc
const (
	// Unknown status means that transaction has been sent to metamorph, but no processing has taken place. This should never be the case, unless something goes wrong.
	Unknown TxStatus = "UNKNOWN" // 0
	// Queued status means that transaction has been queued for processing.
	Queued TxStatus = "QUEUED" // 1
	// Received status means that transaction has been properly received by the metamorph processor.
	Received TxStatus = "RECEIVED" // 2
	// Stored status means that transaction has been stored in the metamorph store. This should ensure the transaction will be processed and retried if not picked up immediately by a mining node.
	Stored TxStatus = "STORED" // 3
	// AnnouncedToNetwork status means that transaction has been announced (INV message) to the Bitcoin network.
	AnnouncedToNetwork TxStatus = "ANNOUNCED_TO_NETWORK" // 4
	// RequestedByNetwork status means that transaction has been requested from metamorph by a Bitcoin node.
	RequestedByNetwork TxStatus = "REQUESTED_BY_NETWORK" // 5
	// SentToNetwork status means that transaction has been sent to at least 1 Bitcoin node.
	SentToNetwork TxStatus = "SENT_TO_NETWORK" // 6
	// AcceptedByNetwork status means that transaction has been accepted by a connected Bitcoin node on the ZMQ interface. If metamorph is not connected to ZQM, this status will never by set.
	AcceptedByNetwork TxStatus = "ACCEPTED_BY_NETWORK" // 7
	// SeenOnNetwork status means that transaction has been seen on the Bitcoin network and propagated to other nodes. This status is set when metamorph receives an INV message for the transaction from another node than it was sent to.
	SeenOnNetwork TxStatus = "SEEN_ON_NETWORK" // 8
	// Mined status means that transaction has been mined into a block by a mining node.
	Mined TxStatus = "MINED" // 9
	// SeenInOrphanMempool means that transaction has been sent to at least 1 Bitcoin node but parent transaction was not found.
	SeenInOrphanMempool TxStatus = "SEEN_IN_ORPHAN_MEMPOOL" // 10
	// Confirmed status means that transaction is marked as confirmed when it is in a block with 100 blocks built on top of that block.
	Confirmed TxStatus = "CONFIRMED" // 108
	// Rejected status means that transaction has been rejected by the Bitcoin network.
	Rejected TxStatus = "REJECTED" // 109
)

// String returns the string representation of the TxStatus
func (s TxStatus) String() string {
	return string(s)
}

// MapTxStatusToInt maps the TxStatus to an int value
func MapTxStatusToInt(status TxStatus) (int, bool) {
	var value int
	var ok bool = true

	switch status {
	case Unknown:
		value = 0
	case Queued:
		value = 1
	case Received:
		value = 2
	case Stored:
		value = 3
	case AnnouncedToNetwork:
		value = 4
	case RequestedByNetwork:
		value = 5
	case SentToNetwork:
		value = 6
	case AcceptedByNetwork:
		value = 7
	case SeenOnNetwork:
		value = 8
	case Mined:
		value = 9
	case SeenInOrphanMempool:
		value = 10
	case Confirmed:
		value = 108
	case Rejected:
		value = 109
	default:
		ok = false
	}

	return value, ok
}
