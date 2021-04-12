package settings

import "testing"

func TestGather(t *testing.T) {
	type args struct {
		opts ReadOptions
		out  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test 1",
			args{
				opts: Options(),
				out:  nil,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Gather(tt.args.opts, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("Gather() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
