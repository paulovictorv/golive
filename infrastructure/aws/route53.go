package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

func findHostedZone(domainName *string) (*route53.HostedZone, error) {
	output, e := route53Client.ListHostedZonesByName(&route53.ListHostedZonesByNameInput{
		DNSName: domainName,
	})

	if e != nil {
		return nil, errors.New(parseAwsError(e))
	}

	if len(output.HostedZones) > 0 {
		return output.HostedZones[0], nil
	}

	return nil, errors.New(fmt.Sprintf("no hostedzone found for domain %s", *domainName))

}

func addCertificateDetails(details *DNSValidation) error {
	hostedZone, err := findHostedZone(details.DomainName)

	if err != nil {
		return errors.New(parseAwsError(err))
	}

	_, e := route53Client.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: details.ResourceRecord.Name,
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(fmt.Sprintf("\"%s\"", *details.ResourceRecord.Value)),
							},
						},
						TTL:  aws.Int64(300),
						Type: details.ResourceRecord.Type,
					},
				},
			},
			Comment: aws.String("Automatically created DNS entries for ACM validation by GoLive"),
		},
		HostedZoneId: hostedZone.Id,
	})

	if e != nil {
		return errors.New(parseAwsError(e))
	}

	return nil
}

func removeCertificateDetails(details *DNSValidation) error {
	hostedZone, findError := findHostedZone(details.DomainName)

	if findError != nil {
		return errors.New(parseAwsError(findError))
	}

	_, deleteError := route53Client.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("DELETE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: details.ResourceRecord.Name,
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(fmt.Sprintf("\"%s\"", *details.ResourceRecord.Value)),
							},
						},
						TTL:  aws.Int64(300),
						Type: details.ResourceRecord.Type,
					},
				},
			},
			Comment: aws.String("Tearing down environment provisioned by GoLive"),
		},
		HostedZoneId: hostedZone.Id,
	})

	if deleteError != nil {
		return errors.New(parseAwsError(deleteError))
	}

	return nil
}
