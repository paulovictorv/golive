package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"time"
)

var acmClient = acm.New(sess)

func findCertificate(domainName string) (string, error) {
	output, e := acmClient.ListCertificates(&acm.ListCertificatesInput{
		CertificateStatuses: aws.StringSlice([]string{acm.CertificateStatusIssued}),
		MaxItems:            aws.Int64(30),
	})

	if e != nil {
		return "", errors.New(parseAwsError(e))
	}

	for _, summary := range output.CertificateSummaryList {
		if *summary.DomainName == domainName {
			return *summary.CertificateArn, nil
		}
	}

	return "", errors.New(fmt.Sprintf("certificate not found for domain %s", domainName))
}

type DNSValidation struct {
	DomainName     *string
	CertificateArn *string
	ResourceRecord *acm.ResourceRecord
}

func requestCertificate(domainName string) (*DNSValidation, error) {
	output, e := acmClient.RequestCertificate(&acm.RequestCertificateInput{
		DomainName:       aws.String(domainName),
		ValidationMethod: aws.String(acm.ValidationMethodDns),
	})

	if e != nil {
		return nil, errors.New(parseAwsError(e))
	}

	//todo: pushback request to wait for DNS fields to be available
	time.Sleep(8 * time.Second)

	certificateOutput, e := acmClient.DescribeCertificate(&acm.DescribeCertificateInput{
		CertificateArn: output.CertificateArn,
	})

	if e != nil {
		return nil, errors.New(parseAwsError(e))
	}

	validations := certificateOutput.Certificate.DomainValidationOptions
	if len(validations) > 0 {
		return &DNSValidation{
			DomainName:     aws.String(domainName),
			CertificateArn: output.CertificateArn,
			ResourceRecord: validations[0].ResourceRecord,
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("failure while getting domain validation options while "+
		"requesting certificate for domain %s", domainName))
}

func deleteCertificate(certificateArn string) error {
	_, e := acmClient.DeleteCertificate(&acm.DeleteCertificateInput{CertificateArn: aws.String(certificateArn)})

	if e != nil {
		return errors.New(parseAwsError(e))
	}

	return nil
}
