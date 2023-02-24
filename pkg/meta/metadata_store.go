package meta

import (
	"zotregistry.io/pkg/meta/repodb/repodbfactory"
	"zotregistry.io/zot/pkg/api/config"
	zlog "zotregistry.io/zot/pkg/log"
	msConfig "zotregistry.io/zot/pkg/meta/config"
	"zotregistry.io/zot/pkg/meta/repodb"
	bolt "zotregistry.io/zot/pkg/meta/repodb/boltdb-wrapper"
	dynamoParams "zotregistry.io/zot/pkg/meta/repodb/dynamodb-wrapper/params"
	"zotregistry.io/zot/pkg/meta/userdb"
	userdbfactory "zotregistry.io/zot/pkg/meta/userdb/factory"
)

type MetadataStore struct {
	userdb.UserStore
	repodb.RepoDB
}

func Create(storageConfig config.StorageConfig, log zlog.Logger) (MetadataStore, error) {
	dynamoParams := getDynamoParams(storageConfig.CacheDriver, log)
	params := bolt.DBParameters{}
	params.RootDir = storageConfig.RootDirectory

	userdata, err := userdbfactory.Create("boltdb", params, log)
	// repodbfactory.Create("dynamodb", dynamoParams)

	return MetadataStore{
		// config: msc,
		UserStore: userdata,
		RepoDB:    repodbfactory.Create("boltdb", params),
	}, err
}

func (m MetadataStore) GetBookmarkedRepos(userid string) ([]string, error) {
	return m.UserStore.GetBookmarkedRepos(userid)
}

func (m MetadataStore) ToggleStarRepo(userid, reponame string) (msConfig.UserState, error) {
	return m.UserStore.ToggleStarRepo(userid, reponame)
}

func (m MetadataStore) GetStarredRepos(userid string) ([]string, error) {
	return m.UserStore.GetStarredRepos(userid)
}

func (m MetadataStore) ToggleBookmarkRepo(userid, reponame string) (msConfig.UserState, error) {
	return m.UserStore.ToggleBookmarkRepo(userid, reponame)
}

func getDynamoParams(cacheDriverConfig map[string]interface{}, log zlog.Logger) dynamoParams.DBDriverParameters {
	allParametersOk := true

	endpoint, ok := toStringIfOk(cacheDriverConfig, "endpoint", log)
	allParametersOk = allParametersOk && ok

	region, ok := toStringIfOk(cacheDriverConfig, "region", log)
	allParametersOk = allParametersOk && ok

	repoMetaTablename, ok := toStringIfOk(cacheDriverConfig, "repometatablename", log)
	allParametersOk = allParametersOk && ok

	manifestDataTablename, ok := toStringIfOk(cacheDriverConfig, "manifestdatatablename", log)
	allParametersOk = allParametersOk && ok

	versionTablename, ok := toStringIfOk(cacheDriverConfig, "versiontablename", log)
	allParametersOk = allParametersOk && ok

	if !allParametersOk {
		panic("dynamo parameters are not specified correctly, can't proceede")
	}

	return dynamoParams.DBDriverParameters{
		Endpoint:              endpoint,
		Region:                region,
		RepoMetaTablename:     repoMetaTablename,
		ManifestDataTablename: manifestDataTablename,
		VersionTablename:      versionTablename,
	}
}

func toStringIfOk(cacheDriverConfig map[string]interface{}, param string, log zlog.Logger) (string, bool) {
	val, ok := cacheDriverConfig[param]

	if !ok {
		log.Error().Msgf("parsing CacheDriver config failed, field '%s' is not present", param)

		return "", false
	}

	str, ok := val.(string)

	if !ok {
		log.Error().Msgf("parsing CacheDriver config failed, parameter '%s' isn't a string", param)

		return "", false
	}

	if str == "" {
		log.Error().Msgf("parsing CacheDriver config failed, field '%s' is is empty", param)

		return "", false
	}

	return str, ok
}
