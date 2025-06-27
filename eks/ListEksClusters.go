package eks

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// listEKSClusters retrieves all EKS clusters across all AWS regions.

func ListEKSClusters(sess *session.Session, clusters *[]string, ctx context.Context) error {

	logging := log.FromContext(ctx)

	// Get all available regions
	ec2Svc := ec2.New(sess)
	regionsOutput, err := ec2Svc.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		logging.Error(err, "failed to describe regions")
		return err
	}

	// Iterate over all regions and list EKS clusters
	for _, region := range regionsOutput.Regions {
		regionName := aws.StringValue(region.RegionName)
		eksSvc := eks.New(sess, &aws.Config{Region: aws.String(regionName)})

		listClustersInput := &eks.ListClustersInput{}
		for {
			listClustersOutput, err := eksSvc.ListClusters(listClustersInput)
			if err != nil {
				logging.Error(err, "failed to list EKS clusters", "region", regionName)
				break
			}

			*clusters = append(*clusters, aws.StringValueSlice(listClustersOutput.Clusters)...)

			if listClustersOutput.NextToken == nil {
				break
			}
			listClustersInput.NextToken = listClustersOutput.NextToken
		}
	}

	return nil
}
