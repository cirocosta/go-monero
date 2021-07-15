package constant

const (
	// AtomicUnit refers to the smallest fraction of a monero.
	//
	AtomicUnit = 1

	// A monero corresponds to 1e12 atomic units.
	//
	XMR      = 1_000_000_000_000 * AtomicUnit
	MilliXMR = 1_000_000_000 * AtomicUnit
	MicroXMR = 1_000_000 * AtomicUnit
)
