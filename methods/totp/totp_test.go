package totp

import (
	"testing"
	"time"
)

func TestTOTP_Generate(t *testing.T) {
	RFCSecret := []byte("12345678901234567890")

	RFCSha1OTP := TOTP{
		Secret: RFCSecret,
		Digits: 8,
		Period: 30,
	}

	type args struct {
		ts time.Time
	}
	tests := []struct {
		name string
		tr   *TOTP
		args args
		want string
	}{
		{
			name: "RFCSha1 Generate 0",
			tr:   &RFCSha1OTP,
			args: args{
				ts: time.Date(1970, 1, 1, 0, 0, 59, 0, time.UTC),
			},
			want: "94287082",
		},
		{
			name: "RFCSha1 Generate 1",
			tr:   &RFCSha1OTP,
			args: args{
				ts: time.Date(2005, 3, 18, 1, 58, 29, 0, time.UTC),
			},
			want: "07081804",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Generate(tt.args.ts); got != tt.want {
				t.Errorf("TOTP.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
