package hyparview

type HyParViewConfig struct {
	NodeID,
	NodeAddress,
	ContactNodeAddress string
	Fanout,
	PartialViewSize,
	ARWL,
	PRWL,
	ShuffleInterval,
	Ka,
	Kp uint
}
