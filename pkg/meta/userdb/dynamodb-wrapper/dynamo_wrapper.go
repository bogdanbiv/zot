package userdynamo

import (
	"context"
	// "encoding/json"

	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	// godigest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"

	zerr "zotregistry.io/zot/errors"
	mcommon "zotregistry.io/zot/pkg/common"
	zlog "zotregistry.io/zot/pkg/log"

	// "zotregistry.io/zot/pkg/meta/repodb" //nolint:go-staticcheck
	// "zotregistry.io/zot/pkg/meta/repodb/common"
	// "zotregistry.io/zot/pkg/meta/repodb/dynamodb-wrapper/iterator"
	metaParams "zotregistry.io/zot/pkg/meta/params"
	"zotregistry.io/zot/pkg/meta/repodb/version"
	"zotregistry.io/zot/pkg/meta/userdb"

	// localCtx "zotregistry.io/zot/pkg/requestcontext"

	// zlog "zotregistry.io/zot/pkg/log"
	msConfig "zotregistry.io/zot/pkg/meta/config"
	// merrors "zotregistry.io/zot/pkg/meta/errors"
)

type DBWrapper struct {
	Client *dynamodb.Client
	// StarsMetaTablename    string
	// BooksMetaTablename    string
	UserMetaTablename string
	// RepoMetaTablename string
	// // ManifestDataTablename string
	VersionTablename string
	Patches          []func(client *dynamodb.Client, tableNames map[string]string) error
	Log              zlog.Logger
}

func NewDynamoDBWrapper(params metaParams.DBDriverParameters, log zlog.Logger) (*DBWrapper, error) {
	// custom endpoint resolver to point to localhost
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           params.Endpoint,
				SigningRegion: region,
			}, nil
		})

	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(params.Region),
		config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		return nil, err
	}

	dynamoWrapper := DBWrapper{
		Client: dynamodb.NewFromConfig(cfg),
		// RepoMetaTablename: params.RepoMetaTablename,
		UserMetaTablename: params.UserMetaTablename,
		// ManifestDataTablename: params.ManifestDataTablename,
		VersionTablename: params.VersionTablename,
		Patches:          version.GetDynamoDBPatches(),
		Log:              log,
	}

	err = dynamoWrapper.createVersionTable()
	if err != nil {
		return nil, err
	}

	err = dynamoWrapper.createUserMetaTable()
	if err != nil {
		return nil, err
	}

	// Using the Config value, create the DynamoDB client
	return &dynamoWrapper, nil
}

func (dwr *DBWrapper) ToggleBookmarkRepo(userid, reponame string) ( //nolint:dupl
	msConfig.UserState, error,
) {
	var res msConfig.UserState = msConfig.NotChanged
	userMeta, err := dwr.GetUserMeta(reponame)
	if err != nil {
		return res, err
	}

	if !mcommon.Contains(userMeta.BookmarkedRepos, reponame) {
		res = msConfig.Added
		userMeta.BookmarkedRepos = append(userMeta.BookmarkedRepos, reponame)
	} else {
		userMeta.BookmarkedRepos = mcommon.RemoveFrom(userMeta.BookmarkedRepos, reponame)
		res = msConfig.Removed
	}

	if res != msConfig.NotChanged {
		err = dwr.setUserMeta(userid, userMeta)
	}

	if err != nil {
		res = msConfig.NotChanged
		return res, err
	}

	return res, nil
}

func (dwr DBWrapper) ResetUserMetaTable() error {
	err := dwr.deleteUserMetaTable()
	if err != nil {
		return err
	}

	return dwr.createUserMetaTable()
}

func (dwr *DBWrapper) GetBookmarkedRepos(userid string) ([]string, error) {
	userMeta, err := dwr.GetUserMeta(userid)

	return userMeta.BookmarkedRepos, err
}

func (dwr *DBWrapper) ToggleStarRepo(userid, reponame string) ( //nolint:dupl
	msConfig.UserState, error,
) {
	var res msConfig.UserState = msConfig.NotChanged
	userMeta, err := dwr.GetUserMeta(reponame)
	if err != nil {
		return res, err
	}

	if !mcommon.Contains(userMeta.StarredRepos, reponame) {
		res = msConfig.Added
		userMeta.StarredRepos = append(userMeta.StarredRepos, reponame)
	} else {
		userMeta.StarredRepos = mcommon.RemoveFrom(userMeta.StarredRepos, reponame)
		res = msConfig.Removed
	}

	if res != msConfig.NotChanged {
		err = dwr.setUserMeta(userid, userMeta)
	}

	if err != nil {
		res = msConfig.NotChanged
		return res, err
	}

	return res, nil
}

func (dwr *DBWrapper) GetStarredRepos(userid string) ([]string, error) {
	userMeta, err := dwr.GetUserMeta(userid)

	return userMeta.StarredRepos, err
}

func (dwr DBWrapper) GetUserMeta(userid string) (userdb.UserMetadata, error) {
	resp, err := dwr.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(dwr.UserMetaTablename),
		Key: map[string]types.AttributeValue{
			"UserName": &types.AttributeValueMemberS{Value: userid},
		},
	})
	if err != nil {
		return userdb.UserMetadata{}, err
	}

	if resp.Item == nil {
		return userdb.UserMetadata{}, zerr.ErrRepoMetaNotFound
	}

	var userMeta userdb.UserMetadata

	err = attributevalue.Unmarshal(resp.Item["UserMetadata"], &userMeta)
	if err != nil {
		return userdb.UserMetadata{}, err
	}

	return userMeta, nil
}

func (dwr DBWrapper) createUserMetaTable() error {
	_, err := dwr.Client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: aws.String(dwr.UserMetaTablename),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("UserMeta"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("UserMeta"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})

	if err != nil && !strings.Contains(err.Error(), "Table already exists") {
		return err
	}

	return dwr.waitTableToBeCreated(dwr.UserMetaTablename)
}

func (dwr DBWrapper) deleteUserMetaTable() error {
	_, err := dwr.Client.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
		TableName: aws.String(dwr.UserMetaTablename),
	})

	if temp := new(types.ResourceNotFoundException); errors.As(err, &temp) {
		return nil
	}

	return dwr.waitTableToBeDeleted(dwr.UserMetaTablename)
}

func (dwr DBWrapper) waitTableToBeCreated(tableName string) error {
	const maxWaitTime = 20 * time.Second

	waiter := dynamodb.NewTableExistsWaiter(dwr.Client)

	return waiter.Wait(context.Background(), &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}, maxWaitTime)
}

func (dwr DBWrapper) waitTableToBeDeleted(tableName string) error {
	const maxWaitTime = 20 * time.Second

	waiter := dynamodb.NewTableNotExistsWaiter(dwr.Client)

	return waiter.Wait(context.Background(), &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}, maxWaitTime)
}

func (dwr *DBWrapper) createVersionTable() error {
	_, err := dwr.Client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: aws.String(dwr.VersionTablename),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("VersionKey"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("VersionKey"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		if strings.Contains(err.Error(), "Table already exists") {
			return nil
		}

		return err
	}

	err = dwr.waitTableToBeCreated(dwr.VersionTablename)
	if err != nil {
		return err
	}

	if err == nil {
		mdAttributeValue, err := attributevalue.Marshal(version.CurrentVersion)
		if err != nil {
			return err
		}

		_, err = dwr.Client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]string{
				"#V": "Version",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":Version": mdAttributeValue,
			},
			Key: map[string]types.AttributeValue{
				"VersionKey": &types.AttributeValueMemberS{
					Value: version.DBVersionKey,
				},
			},
			TableName:        aws.String(dwr.VersionTablename),
			UpdateExpression: aws.String("SET #V = :Version"),
		})

		if err != nil {
			return err
		}
	}

	return nil
}

/*
// func (dwr DBWrapper) createManifestDataTable() error {
// 	_, err := dwr.Client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
// 		TableName: aws.String(dwr.ManifestDataTablename),
// 		AttributeDefinitions: []types.AttributeDefinition{
// 			{
// 				AttributeName: aws.String("Digest"),
// 				AttributeType: types.ScalarAttributeTypeS,
// 			},
// 		},
// 		KeySchema: []types.KeySchemaElement{
// 			{
// 				AttributeName: aws.String("Digest"),
// 				KeyType:       types.KeyTypeHash,
// 			},
// 		},
// 		BillingMode: types.BillingModePayPerRequest,
// 	})

// 	if err != nil && !strings.Contains(err.Error(), "Table already exists") {
// 		return err
// 	}

// 	return dwr.waitTableToBeCreated(dwr.ManifestDataTablename)
// }
*/

func (dwr DBWrapper) setUserMeta(userid string, userMeta userdb.UserMetadata) error {
	repoAttributeValue, err := attributevalue.Marshal(userMeta)
	if err != nil {
		return err
	}

	_, err = dwr.Client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]string{
			"#UM": "UserMetadata",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":UserMetadata": repoAttributeValue,
		},
		Key: map[string]types.AttributeValue{
			"UserName": &types.AttributeValueMemberS{
				Value: userid,
			},
		},
		TableName:        aws.String(dwr.UserMetaTablename),
		UpdateExpression: aws.String("SET #UM = :UserMetadata"),
	})

	return err
}
