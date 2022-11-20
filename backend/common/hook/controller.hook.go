package hook

type Controller interface {
	// allow access without jwt
	Accessible()

	// need jwt to access
	Restricted()
}
