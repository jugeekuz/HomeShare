package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func GetARecord(hostedZoneID string, recordName string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", err
	}

	client := route53.NewFromConfig(cfg)

	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
		StartRecordName: aws.String(recordName),
		StartRecordType: types.RRTypeA,
	}

	result, err := client.ListResourceRecordSets(context.TODO(), input)
	if err != nil {
		return "", err
	}

	for _, record := range result.ResourceRecordSets {
		if *record.Name == recordName && record.Type == types.RRTypeA {
            if len(record.ResourceRecords) > 0 {
                return *record.ResourceRecords[0].Value, nil
            }
        }
	}

	return "", fmt.Errorf("a record not found %s", hostedZoneID)
}

func UpdateARecord(hostedZoneId string, recordName string, newIp string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := route53.NewFromConfig(cfg)

	changeBatch := &types.ChangeBatch{
        Changes: []types.Change{{
            Action: types.ChangeActionUpsert,
            ResourceRecordSet: &types.ResourceRecordSet{
                Name: aws.String(recordName),
                Type: types.RRTypeA,
                TTL:  aws.Int64(300),
                ResourceRecords: []types.ResourceRecord{{
                    Value: aws.String(newIp),
                }},
            },
        }},
    }

	input := &route53.ChangeResourceRecordSetsInput{
        HostedZoneId: aws.String(hostedZoneId),
        ChangeBatch:  changeBatch,
    }

	result, err := client.ChangeResourceRecordSets(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to update A record: %v", err)
	}

	fmt.Printf("A record updated successfully. Change ID: %s\n", *result.ChangeInfo.Id)
    return nil
}