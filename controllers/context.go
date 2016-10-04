package controllers

// Context contains info on the current request
type Context struct {
	// Signer is the public key that signed the transaction (may be nil)
	signer []byte
}

func (c *Context) IsAnon() bool {
	return c == nil || c.signer == nil || len(c.signer) == 0
}

func (c *Context) Signer() []byte {
	return c.signer
}

// TODO: store db
