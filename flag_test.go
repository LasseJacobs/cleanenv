package cleanenv

import (
	"os"
	"reflect"
	"testing"
)

func TestFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"-test1", "10"}

	type Combined struct {
		Empty   int
		Default int `env:"TEST0" env-default:"1"`
		Global  int `env:"TEST1" env-default:"1"`
		local   int `env:"TEST2" env-default:"1"`
	}

	tests := []struct {
		name    string
		env     map[string]string
		cfg     interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "combined",
			env: map[string]string{
				"TEST1": "2",
				"TEST2": "3",
			},
			cfg: &Combined{},
			want: &Combined{
				Empty:   0,
				Default: 1,
				Global:  10,
				local:   0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for env, val := range tt.env {
				os.Setenv(env, val)
			}
			defer os.Clearenv()

			if err := ReadEnv(tt.cfg, ""); (err != nil) != tt.wantErr {
				t.Errorf("wrong error behavior %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.cfg, tt.want) {
				t.Errorf("wrong data %v, want %v", tt.cfg, tt.want)
			}
		})
	}
}
