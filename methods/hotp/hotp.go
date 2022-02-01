package hotp

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/sha3"
)

type HOTP struct {
	Secret  []byte
	Digits  int
	Counter int

	Algorithm string
}

func pow(base, exp uint64) uint64 {
	val := uint64(1)
	for i := uint64(0); i < exp; i++ {
		val *= base
	}
	return val
}

func (h *HOTP) Generate() string {
	alg := sha1.New
	switch h.Algorithm {
	case "sha1":
		alg = sha1.New
	case "sha256":
		alg = sha256.New
	case "sha384":
		alg = sha512.New384
	case "sha512":
		alg = sha512.New
	case "sha3-256":
		alg = sha3.New256
	case "sha3-384":
		alg = sha3.New384
	case "sha3-512":
		alg = sha3.New512
	}

	hm := hmac.New(alg, h.Secret)
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(h.Counter))
	hm.Write(b[:])
	data := hm.Sum(nil)
	//offset := data[19] & 0xf
	offset := uint64(data[len(data)-1] & 0xf)
	bincode := (uint64(data[offset]) & 0x7f) << 24
	bincode |= (uint64(data[offset+1]) & 0xff) << 16
	bincode |= (uint64(data[offset+2]) & 0xff) << 8
	bincode |= (uint64(data[offset+3]) & 0xff) << 0
	otp := bincode % pow(10, uint64(h.Digits))
	return fmt.Sprintf("%0*d", h.Digits, otp)
}
