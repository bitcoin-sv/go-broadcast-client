package broadcast

// TxStatus is the status of the transaction
type TxStatus string

// List of statuses available here: https://github.com/bitcoin-sv/arc
const (
	// Unknown contains value for unknown status
	Unknown TxStatus = "UNKNOWN" // 0
	// Queued contains value for queued status
	Queued TxStatus = "QUEUED" // 1
	// Received contains value for received status
	Received TxStatus = "RECEIVED" // 2
	// Stored contains value for stored status
	Stored TxStatus = "STORED" // 3
	// AnnouncedToNetwork contains value for announced to network status
	AnnouncedToNetwork TxStatus = "ANNOUNCED_TO_NETWORK" // 4
	// RequestedByNetwork contains value for requested by network status
	RequestedByNetwork TxStatus = "REQUESTED_BY_NETWORK" // 5
	// SentToNetwork contains value for sent to network status
	SentToNetwork TxStatus = "SENT_TO_NETWORK" // 6
	// AcceptedByNetwork contains value for accepted by network status
	AcceptedByNetwork TxStatus = "ACCEPTED_BY_NETWORK" // 7
	// SeenOnNetwork contains value for seen on network status
	SeenOnNetwork TxStatus = "SEEN_ON_NETWORK" // 8
	// Mined contains value for mined status
	Mined TxStatus = "MINED" // 9
	// Confirmed contains value for confirmed status
	Confirmed TxStatus = "CONFIRMED" // 108
	// Rejected contains value for rejected status
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
	case Confirmed:
		value = 108
	case Rejected:
		value = 109
	default:
		ok = false
	}

	return value, ok
}
