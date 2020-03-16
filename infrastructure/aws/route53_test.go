package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"testing"
)

func Test_addCertificateDetails(t *testing.T) {
	type args struct {
		details *DNSValidation
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test adding entries for existing domain",
			args: args{details: &DNSValidation{
				DomainName:     aws.String("goclip.co"),
				CertificateArn: nil,
				ResourceRecord: &acm.ResourceRecord{
					Name:  aws.String("golive-test.goclip.co"),
					Type:  aws.String("TXT"),
					Value: aws.String("golive-test"),
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := addCertificateDetails(tt.args.details); (err != nil) != tt.wantErr {
				t.Errorf("addCertificateDetails() error = %v, wantErr %v", err, tt.wantErr)
			}

			if removeCertificateDetails(tt.args.details) != nil {
				t.Errorf("error tearing down")
			}

		})
	}
}
