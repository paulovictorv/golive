package aws

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()

	run := m.Run()

	os.Exit(run)
}

func Test_findCertificate(t *testing.T) {
	type args struct {
		domainName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "find certificate for financamcz.com.br",
			args:    args{domainName: "financamcz.com.br"},
			want:    "arn:aws:acm:us-east-1:258743747205:certificate/a9d70ef4-6636-4287-9626-a61b1b089b5e",
			wantErr: false,
		},
		{
			name:    "find certificate for goclip.com.br",
			args:    args{domainName: "goclip.com.br"},
			want:    "arn:aws:acm:us-east-1:258743747205:certificate/6faa9590-021f-4656-9cec-2676d39da483",
			wantErr: false,
		},
		{
			name:    "find certificate for tasking.io",
			args:    args{domainName: "tasking.io"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findCertificate(tt.args.domainName)
			if (err != nil) != tt.wantErr {
				t.Errorf("findCertificate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("findCertificate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_requestCertificate(t *testing.T) {
	type args struct {
		domainName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "requesting certificate for tasking.io",
			args:    args{domainName: "tasking.io"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := requestCertificate(tt.args.domainName)
			if (err != nil) != tt.wantErr {
				t.Errorf("requestCertificate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if err := deleteCertificate(*got.CertificateArn); err != nil {
					t.Errorf("error while cleaning up")
				}
			}
		})
	}
}
