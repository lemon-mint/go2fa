package totp

import (
	"time"

	"github.com/lemon-mint/go2fa/methods/hotp"
)

type TOTP struct {
	Secret []byte
	Digits int
	Period int

	Algorithm string
}

func (t *TOTP) Generate(ts time.Time) string {
	ctr := ts.Unix() / int64(t.Period)
	h := hotp.HOTP{
		Secret:  t.Secret,
		Digits:  t.Digits,
		Counter: int(ctr),

		Algorithm: t.Algorithm,
	}
	return h.Generate()
}
