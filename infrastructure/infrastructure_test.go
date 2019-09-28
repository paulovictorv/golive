package infrastructure

import (
	"reflect"
	"testing"
)

func TestCreateInfra(t *testing.T) {
	type args struct {
		provider Provider
	}
	tests := []struct {
		name string
		args args
		want Infrastructure
	}{
		{
			"test get aws provider",
			args{provider: AWS},
			AmazonInfrastructure{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateInfra(tt.args.provider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateInfra() = %v, want %v", got, tt.want)
			}
		})
	}
}
