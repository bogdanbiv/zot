package userdynamo_test

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	guuid "github.com/gofrs/uuid"
	"github.com/rs/zerolog"
	. "github.com/smartystreets/goconvey/convey"

	"zotregistry.io/zot/pkg/log"

	zlog "zotregistry.io/zot/pkg/log"
	metaParams "zotregistry.io/zot/pkg/meta/params"
	"zotregistry.io/zot/pkg/meta/repodb/dynamodb-wrapper/iterator"
	"zotregistry.io/zot/pkg/meta/userdb"
	userdynamo "zotregistry.io/zot/pkg/meta/userdb/dynamodb-wrapper"
)

func TestIterator(t *testing.T) {
	const (
		endpoint = "http://localhost:4566"
		region   = "us-east-2"
	)

	uuid, err := guuid.NewV4()
	if err != nil {
		panic(err)
	}

	userMetaTablename := "UserMetadataTable" + uuid.String()
	log := zlog.NewLogger("INFO", "stdout")

	Convey("TestIterator", t, func() {
		dynamoWrapper, err := userdynamo.NewDynamoDBWrapper(metaParams.DBDriverParameters{
			Endpoint:          endpoint,
			Region:            region,
			UserMetaTablename: userMetaTablename,
		}, log)
		So(err, ShouldBeNil)

		So(dynamoWrapper.ResetUserMetaTable(), ShouldBeNil)
		userMeta, err := dynamoWrapper.GetUserMeta("user1")
		So(userMeta, ShouldBeBlank)
		So(err, ShouldBeNil)

		userMeta2 := userdb.UserMetadata{
			StarredRepos:    []string{"Abc", "Xyz", "Mno"},
			BookmarkedRepos: []string{"repo1", "repo2"},
		}
		err = setUserMeta(dynamoWrapper.Client, dynamoWrapper.UserMetaTablename, "user1", userMeta2)
		So(err, ShouldBeNil)

		So(err, ShouldBeNil)

		repoMetaAttributeIterator := iterator.NewBaseDynamoAttributesIterator(
			dynamoWrapper.Client,
			userMetaTablename,
			"UserMetadata",
			1,
			log,
		)

		attribute, err := repoMetaAttributeIterator.First(context.Background())
		So(err, ShouldBeNil)
		So(attribute, ShouldNotBeNil)

		attribute, err = repoMetaAttributeIterator.Next(context.Background())
		So(err, ShouldBeNil)
		So(attribute, ShouldNotBeNil)

	})
}

func TestIteratorErrors(t *testing.T) {
	Convey("errors", t, func() {
		customResolver := aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           "endpoint",
					SigningRegion: region,
				}, nil
			})

		cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("region"),
			config.WithEndpointResolverWithOptions(customResolver))
		So(err, ShouldBeNil)

		repoMetaAttributeIterator := iterator.NewBaseDynamoAttributesIterator(
			dynamodb.NewFromConfig(cfg),
			"UserMetadataTable",
			"UserMetadata",
			1,
			log.Logger{Logger: zerolog.New(os.Stdout)},
		)

		_, err = repoMetaAttributeIterator.First(context.Background())
		So(err, ShouldNotBeNil)
	})
}

func setBadUserMeta(client *dynamodb.Client, userMetadataTableName string,
	userName string, userMeta userdb.UserMetadata) error {
	userAttributeValue, err := attributevalue.Marshal("string")
	if err != nil {
		return err
	}

	_, err = client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]string{
			"#RM": "UserMetadata",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":UserMetadata": userAttributeValue,
		},
		Key: map[string]types.AttributeValue{
			"RepoName": &types.AttributeValueMemberS{
				Value: userName,
			},
		},
		TableName:        aws.String(userMetadataTableName),
		UpdateExpression: aws.String("SET #RM = :UserMetadata"),
	})

	return err
}

func setUserMeta(client *dynamodb.Client, userMetadataTableName string,
	userid string, userMeta userdb.UserMetadata) error {

	repoAttributeValue, err := attributevalue.Marshal(userMeta)
	if err != nil {
		return err
	}

	_, err = client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]string{
			"#RM": "UserMetadata",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":UserMetadata": repoAttributeValue,
		},
		Key: map[string]types.AttributeValue{
			"RepoName": &types.AttributeValueMemberS{
				Value: userid,
			},
		},
		TableName:        aws.String(userMetadataTableName),
		UpdateExpression: aws.String("SET #RM = :UserMetadata"),
	})

	return err
}
