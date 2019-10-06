package aws

import "testing"

func Test_createCdn(t *testing.T) {
	type args struct {
		bucketName string
		domainName string
		acmArn     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test create basic cdn with valid certificate",
			args:    args{bucketName: "golive-1-test", domainName: "golive-test-3.apps.goclip.co", acmArn: "arn:aws:acm:us-east-1:258743747205:certificate/5b878a88-4c5e-499f-bc4d-6a08f49c0634"},
			wantErr: false,
		},
		{
			name:    "test create basic cdn with invalid certificate",
			args:    args{bucketName: "golive-1-test", domainName: "golive-test-4.apps.goclip.co", acmArn: "a"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createCdn(tt.args.bucketName, tt.args.domainName, tt.args.acmArn, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("createCdn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("createCdn() got = %v, want %v", got, tt.want)
			}
		})
	}
}
