package monero_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cirocosta/go-monero/pkg/monero"
)

func TestSeed(t *testing.T) {
	for _, tc := range []struct {
		name string

		pk             []byte
		mnemonic       []string
		primaryAddress string
	}{
		{name: "full 0-s",

			pk: []byte{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			mnemonic: []string{
				"abbey", "abbey", "abbey", "abbey",
				"abbey", "abbey", "abbey", "abbey",
				"abbey", "abbey", "abbey", "abbey",
				"abbey", "abbey", "abbey", "abbey",
				"abbey", "abbey", "abbey", "abbey",
				"abbey", "abbey", "abbey", "abbey",

				"abbey",
			},
			primaryAddress: "41fJjQDhryD11111111111111111111111111111111112N1GuTZeagfRbbKcALdcZev4QXGGuoLh2x36LhaxLSxCc2YDhi",
		},
		{name: "first 8, and last 8",

			pk: []byte{
				1, 2, 3, 4, 5, 6, 7, 8,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				8, 7, 6, 5, 4, 3, 2, 1,
			},
			mnemonic: []string{
				"object", "anxiety", "asked", "stockpile",
				"saucepan", "skew", "abbey", "abbey",
				"abbey", "abbey", "abbey", "abbey",
				"abbey", "abbey", "abbey", "abbey",
				"abbey", "abbey", "huddle", "excess",
				"fever", "dagger", "nibs", "nineteen",

				"abbey",
			},
			primaryAddress: "4953Se8CDGeZHr8sWmL61WNhKJatXZRSv6eJHB4hbBXF2RFmmKcxHQDR9i8nDkk94uHddmZohDAaPMcoqWgj4oMa748VFsf",
		},
		{name: "full 1-s",

			pk: []byte{
				1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1,
			},
			mnemonic: []string{
				"myriad", "vane", "vector", "myriad",
				"vane", "vector", "myriad", "vane",
				"vector", "myriad", "vane", "vector",
				"myriad", "vane", "vector", "myriad",
				"vane", "vector", "myriad", "vane",
				"vector", "myriad", "vane", "vector",

				"vector",
			},
			primaryAddress: "42Lxp5b63YJ8mVZTzcioVnCk9WQCPAMk4RH7e7ygPTkzEexugYvy5PvQS1MoWJE5ugbSCx9jYHgdbWhd8tQNvDbwFXdWDQC",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := monero.NewSeed(tc.pk)

			assert.Equal(t, tc.mnemonic, s.Mnemonic())
			assert.Equal(t, tc.primaryAddress, s.PrimaryAddress())
		})
	}
}
