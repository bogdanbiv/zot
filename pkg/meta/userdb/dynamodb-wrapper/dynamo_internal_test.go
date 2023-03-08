package userdynamo

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	guuid "github.com/gofrs/uuid"
	"github.com/rs/zerolog"
	. "github.com/smartystreets/goconvey/convey"

	"zotregistry.io/zot/pkg/log" //nolint:go-staticcheck
	"zotregistry.io/zot/pkg/meta/repodb/version"
	// userdynamo "zotregistry.io/zot/pkg/meta/userdb/dynamodb-wrapper"
)

func TestWrapperErrors(t *testing.T) {
	const (
		endpoint = "http://localhost:4566"
		region   = "us-east-2"
	)

	uuid, err := guuid.NewV4()
	if err != nil {
		panic(err)
	}

	userMetaTablename := "UserMetadataTable" + uuid.String()
	versionTablename := "Version" + uuid.String()

	Convey("Create table errors", t, func() {
		badEndpoint := endpoint + "1"

		customResolver := aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           badEndpoint,
					SigningRegion: region,
				}, nil
			},
		)

		cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region),
			config.WithEndpointResolverWithOptions(customResolver))
		So(err, ShouldBeNil)

		dynamoWrapper := DBWrapper{
			Client:            dynamodb.NewFromConfig(cfg),
			UserMetaTablename: userMetaTablename,
			VersionTablename:  versionTablename,
			Patches:           version.GetDynamoDBPatches(),
			Log:               log.Logger{Logger: zerolog.New(os.Stdout)},
		}

		// The table creation should fail as the endpoint is not configured correctly
		err = dynamoWrapper.createUserMetaTable()
		So(err, ShouldNotBeNil)

		err = dynamoWrapper.createVersionTable()
		So(err, ShouldNotBeNil)
	})

	Convey("Delete table errors", t, func() {
		customResolver := aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			},
		)

		cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region),
			config.WithEndpointResolverWithOptions(customResolver))
		So(err, ShouldBeNil)

		dynamoWrapper := DBWrapper{
			Client:            dynamodb.NewFromConfig(cfg),
			UserMetaTablename: userMetaTablename,
			VersionTablename:  versionTablename,
			Patches:           version.GetDynamoDBPatches(),
			Log:               log.Logger{Logger: zerolog.New(os.Stdout)},
		}

		// The tables were not created so delete calls fail, but dynamoWrapper should not error
		err = dynamoWrapper.deleteUserMetaTable()
		So(err, ShouldBeNil)
	})
}
