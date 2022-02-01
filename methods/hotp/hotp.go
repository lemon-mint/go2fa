package hotp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
)

type HOTP struct {
	Secret  []byte
	Digits  int
	Counter int
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
