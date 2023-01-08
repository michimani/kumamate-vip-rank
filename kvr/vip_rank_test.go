package kvr_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/michimani/kumamate-vip-rank/kvr"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/cache
var cache []byte

func Test_VIPRrankList_GetTitle(t *testing.T) {
	var vrl kvr.VIPRankList
	err := json.Unmarshal(cache, &vrl)
	assert.NoError(t, err)

	cases := []struct {
		name   string
		power  uint64
		expect string
	}{
		{
			name:   "zero",
			power:  0,
			expect: "未VIP発射台",
		},
		{
			name:   "sakurai",
			power:  99999999,
			expect: "桜井",
		},
		{
			name:   "just vip",
			power:  11452732,
			expect: "VIP到達！",
		},
		{
			name:   "lower vip",
			power:  11452731,
			expect: "VIPまであと2-3勝！",
		},
		{
			name:   "upper vip",
			power:  11578712,
			expect: "VIP入りたて",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			title := vrl.GetTitle(c.power)
			asst.Equal(c.expect, title)
		})
	}
}
