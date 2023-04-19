package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {

	accessKey := os.Getenv("aws_access_key_id")
	secretKey := os.Getenv("aws_secret_access_key")
	token := os.Getenv("aws_session_token")
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("us-east-1"),
			// Credentials: credentials.NewStaticCredentials("ASIAZ3CUIHQKJWMN3L5I", "b+kizi4X6Zj9u/ymLJQTVTsfjBX3XCh1Z9f9D59v", "IQoJb3JpZ2luX2VjEBgaCXVzLXdlc3QtMiJHMEUCIQC7w9Ns9V05v/1u3c/vQzNROvaS0sOodAjF0sn5g0EqOgIgYap5Vy6hMP28HvOO+0bIFamGRS7TBNprmpfCOmTPEfUqkgMIERABGgw2NzY2MzQwNDk1NTYiDDyNC5A446QtToCWtirvAhcVqlSU8sWgpGp9vsZ53yzHKX4EZT+YuOSF05tRETnGqgIwVfegXAbQZzhwvZ/kQDq5jtf0deOiNBHrbKc2MY3rYgC1a+2HADXf/tpDUUoI3/7BkZR3GG3MYZpd1BKbWIAqzezspqYnryr6Z8Or5wB+MLLmEImLJ7IYhRpL1/M+jma9sRUKMn2J6atjS0D6GCLtaMlg0SM4AOdla86eO1TnRQjDd1ki6x7ahV/AlDOkGm/JhS5s5n8iE7jHscg5/SImqSXmjyTu/4cuNxve/dmrHdIjDGG0m2Ix1MiXDs95Zyt08Z3V8hqV2RALWMEJHi0H7Yrrk+13ehOryXDYvcXpyVWCgXx1fFZGDshLQxCJqozbvm8T2nQ4GCktw30pkmvCTNdbHJiC2aJC9MS4sCzd1TGOC/8j2cI13IlW2RKICq/eGlfC+EfMbvGL6bR8DZ+Q+Z9hB4aGF2qeC7nnXEz4Nw7HXp7+2HUrQDY/F2Uw9bDpoQY6pgETXyl7FtNou7YrzaAG2py6/PoRy7bE/X8+ENmVP73n/ZajUDnzoCL57JsTl0uVInh32qh4O20RELla8QtsGGyCRVIbXFWYjxh/3GnYCYXf7WbL/f7H1R1hFeHCGuhObqa8VwocgVBOJ3VqhPA/w9H+GHrMeEo8SWnGm6mhK8Sw0IDjnNx1LlYw1nCf6D9tYumeDQaqBmye6Bx6yCRaclbprmKurG4v"),
			Credentials: credentials.NewStaticCredentials(accessKey, secretKey, token),
		},
	})

	if err != nil {
		fmt.Printf("Failed to initialize new session: %v", err)
		return
	}
	// Create a new session using the default AWS configuration
	// sess, err := session.NewSession()
	// if err != nil {
	// 	panic(err)
	// }

	// Create a new EC2 client using the session
	ec2Svc := ec2.New(sess)

	// Create a new EC2 instance
	resp, err := ec2Svc.RunInstances(&ec2.RunInstancesInput{

		ImageId:      aws.String("ami-06e46074ae430fba6"), // Replace with the AMI ID of your choice
		InstanceType: aws.String("t2.micro"),              // Replace with the instance type of your choice
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf(err.Error())

	// Get the instance ID of the new instance

	instanceID := *resp.Instances[0].InstanceId
	fmt.Printf("Created instance with ID: %s\n", instanceID)

	input := &ec2.CreateTagsInput{
		Resources: []*string{aws.String(instanceID)},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("Raghav"),
			},
			{
				Key:   aws.String("Purpose"),
				Value: aws.String("practice"),
			},
		},
	}

	// Call the CreateTags API to create the tags for the instance.
	_, err = ec2Svc.CreateTags(input)
	if err != nil {
		fmt.Println("Error creating tags:", err)
		os.Exit(1)
	}
	// input2 := &ec2.StopInstancesInput{
	// 	InstanceIds: []*string{
	// 		aws.String(instanceID),
	// 	},
	// }

	// // Use the EC2 client to stop the instance
	// _, err = ec2Svc.StopInstances(input2)
	// if err != nil {
	// 	fmt.Println("Error stopping instance:", err)
	// 	os.Exit(1)
	// }
	//----------------------------------------------------------------------------------------------------------------------------
	input3 := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			},
		},
	}

	// Use the EC2 client to describe the instances
	result, err := ec2Svc.DescribeInstances(input3)
	if err != nil {
		fmt.Println("Error describing instances:", err)
		os.Exit(1)
	}

	// Extract the instance IDs from the result
	instanceIDs := make([]string, 0)
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, *instance.InstanceId)
		}
	}

	fmt.Println("Running instances:", instanceIDs)

}
