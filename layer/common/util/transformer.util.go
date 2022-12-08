package util

import (
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/p9595jh/transform"
)

var transformer = transform.New(
	transform.I{
		Name: "str2big",
		F: transform.F2(func(s1, s2 string) *big.Int {
			if s1 == "" {
				return nil
			}
			i := new(big.Int)
			i.SetString(s1, 10)
			return i
		}),
	}, transform.I{
		Name: "big2str",
		F: transform.F2(func(i *big.Int, s string) string {
			if i == nil {
				return ""
			}
			return i.String()
		}),
	}, transform.I{
		Name: "trimPre",
		F: transform.F2(func(s1, s2 string) string {
			s1 = strings.TrimPrefix(s1, s2)
			s1 = strings.ToLower(s1)
			return s1
		}),
	}, transform.I{
		Name: "time2str",
		F: transform.F2(func(t *time.Time, s string) string {
			if t == nil {
				return ""
			} else {
				return t.Format(DateFormat)
			}
		}),
	}, transform.I{
		Name: "str2time",
		F: transform.F2(func(s1, s2 string) *time.Time {
			if t, err := time.Parse(DateFormat, s1); err != nil {
				return &time.Time{}
			} else {
				return &t
			}
		}),
	}, transform.I{
		Name: "unix2time",
		F: transform.F2(func(i int64, s string) *time.Time {
			t := time.Unix(i, 0)
			return &t
		}),
	}, transform.I{
		Name: "time2unix",
		F: transform.F2(func(t *time.Time, s string) int64 {
			return t.Unix()
		}),
	}, transform.I{
		Name: "addr2str",
		F: transform.F2(func(a common.Address, _ string) string {
			s := a.String()
			s = strings.TrimPrefix(s, "0x")
			s = strings.ToLower(s)
			return s
		}),
	}, transform.I{
		Name: "str2addr",
		F: transform.F2(func(s, _ string) common.Address {
			return common.HexToAddress(s)
		}),
	},
)

func Transformer() transform.Transformer {
	return transformer
}
