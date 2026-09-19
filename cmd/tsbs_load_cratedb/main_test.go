package main

import (
	"testing"

	"github.com/blagojts/viper"
	"github.com/spf13/pflag"
	"github.com/timescale/tsbs/pkg/targets/crate"
)

// TestCrateDBStorageSettings is a regression test for a bug where
// --replicas/--shards passed on the command line were silently ignored.
func TestCrateDBStorageSettings(t *testing.T) {
	cases := []struct {
		desc         string
		args         []string
		wantReplicas int
		wantShards   int
	}{
		{
			desc:         "CLI values are honored",
			args:         []string{"--replicas=1", "--shards=16"},
			wantReplicas: 1,
			wantShards:   16,
		},
		{
			desc:         "defaults apply when unset",
			args:         []string{},
			wantReplicas: 0,
			wantShards:   5,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			crate.NewTarget().TargetSpecificFlags("", fs)
			if err := fs.Parse(c.args); err != nil {
				t.Fatalf("%s: failed to parse flags: %v", c.desc, err)
			}

			v := viper.New()
			if err := v.BindPFlags(fs); err != nil {
				t.Fatalf("%s: failed to bind flags: %v", c.desc, err)
			}

			gotReplicas, gotShards := crateDBStorageSettings(v)
			if gotReplicas != c.wantReplicas {
				t.Errorf("%s: replicas = %d, want %d", c.desc, gotReplicas, c.wantReplicas)
			}
			if gotShards != c.wantShards {
				t.Errorf("%s: shards = %d, want %d", c.desc, gotShards, c.wantShards)
			}
		})
	}
}
