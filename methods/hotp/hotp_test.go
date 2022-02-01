package hotp

import (
	"encoding/hex"
	"testing"
)

func MustBytes(v []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return v
}

func TestHOTP_Generate(t *testing.T) {
	RFCSecret := MustBytes(hex.DecodeString("3132333435363738393031323334353637383930"))

	tests := []struct {
		name string
		h    HOTP
		want string
	}{
		{
			"Test HOTP_Generate 0",
			HOTP{
				Secret:  RFCSecret,
				Digits:  6,
				Counter: 0,
			},
			"755224",
		},
		{
			"Test HOTP_Generate 1",
			HOTP{
				Secret:  RFCSecret,
				Digits:  6,
				Counter: 1,
			},
			"287082",
		},
		{
			"Test HOTP_Generate 2",
			HOTP{
				Secret:  RFCSecret,
				Digits:  6,
				Counter: 2,
			},
			"359152",
		},
		{
			"Test HOTP_Generate 3",
			HOTP{
				Secret:  RFCSecret,
				Digits:  6,
				Counter: 3,
			},
			"969429",
		},
		{
			"Test HOTP_Generate 4",
			HOTP{
				Secret:  RFCSecret,
				Digits:  6,
				Counter: 4,
			},
			"338314",
		},
		{
			"Test HOTP_Generate 5",
			HOTP{
				Secret:  RFCSecret,
				Digits:  6,
				Counter: 5,
			},
			"254676",
		},
		{
			"Test HOTP_Generate 6",
			HOTP{
				Secret:  RFCSecret,
				Digits:  6,
				Counter: 6,
			},
			"287922",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Generate(); got != tt.want {
				t.Errorf("HOTP.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
