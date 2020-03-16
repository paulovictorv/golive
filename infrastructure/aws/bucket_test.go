package aws

import "testing"

func Test_createBucket(t *testing.T) {
	type args struct {
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test create bucket with correct website configuration",
			args:    args{bucketName: ""},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createBucket(tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("createBucket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createBucket() got = %v, want %v", got, tt.want)
			}
		})
	}
}
