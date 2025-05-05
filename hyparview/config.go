package hyparview

type HyParViewConfig struct {
	Fanout,
	PartialViewSize,
	ARWL,
	PRWL,
	ShuffleInterval,
	Ka,
	Kp int
}

type Config struct {
	NodeID,
	NodeAddress,
	ContactNodeAddress string
	HyParViewConfig
}
