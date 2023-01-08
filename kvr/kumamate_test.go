package kvr_test

import (
	"net/http"
	"testing"

	"github.com/michimani/kumamate-vip-rank/kvr"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	cases := []struct {
		name      string
		hc        *http.Client
		cachePath string
		wantErr   bool
	}{
		{
			name:      "ok",
			hc:        http.DefaultClient,
			cachePath: "tmp-path",
			wantErr:   false,
		},
		{
			name:      "error: http.Client is nil",
			hc:        nil,
			cachePath: "tmp-path",
			wantErr:   true,
		},
		{
			name:      "error: cachePath is empty",
			hc:        http.DefaultClient,
			cachePath: "",
			wantErr:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			if c.cachePath != "" {
				tt.Setenv(kvr.Exported_cacheFilePathKey, c.cachePath)
			}

			client, err := kvr.New(c.hc)
			if c.wantErr {
				asst.Error(err)
				asst.Nil(client)
				return
			}

			asst.NoError(err)
			asst.NotNil(client)
		})
	}
}
