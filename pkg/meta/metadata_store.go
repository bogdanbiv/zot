package meta

import (
	"zotregistry.io/zot/errors"
	zlog "zotregistry.io/zot/pkg/log"
	metaErrors "zotregistry.io/zot/pkg/meta/errors"
	metaParams "zotregistry.io/zot/pkg/meta/params"
	"zotregistry.io/zot/pkg/meta/repodb"
	"zotregistry.io/zot/pkg/meta/repodb/repodbfactory"
	"zotregistry.io/zot/pkg/meta/userdb"
	userdbfactory "zotregistry.io/zot/pkg/meta/userdb/factory"
)

type MetadataStore struct {
	userdb.UserStore
	repodb.RepoDB
}

func Create(dbtype string, dbparam interface{}, log zlog.Logger) (MetadataStore, error) {
	var (
		userdata           userdb.UserStore
		reposdata          repodb.RepoDB
		err                error
		emptyMetadataStore = MetadataStore{
			UserStore: nil,
			RepoDB:    nil,
		}
	)

	switch dbtype {
	case "boltdb":
		{
			params, ok := dbparam.(metaParams.BoltDriverParameters)
			if !ok {
				return emptyMetadataStore, metaErrors.ErrInvalidConfiguration
			}
			userdata, err = userdbfactory.Create("boltdb", params, log)
		}
	case "dynamodb":
		{
			params, ok := dbparam.(metaParams.BoltDriverParameters)
			if !ok {
				return emptyMetadataStore, metaErrors.ErrInvalidConfiguration
			}
			reposdata, err = repodbfactory.Create("dynamodb", params, log)
		}
	default:
		{
			return emptyMetadataStore, errors.ErrBadConfig
		}
	}

	return MetadataStore{
		// config: msc,
		UserStore: userdata,
		RepoDB:    reposdata,
	}, err
}

// func (m MetadataStore) GetBookmarkedRepos(userid string) ([]string, error) {
// 	return m.UserStore.GetBookmarkedRepos(userid)
// }

// func (m MetadataStore) ToggleStarRepo(userid, reponame string) (msConfig.UserState, error) {
// 	return m.UserStore.ToggleStarRepo(userid, reponame)
// }

// func (m MetadataStore) GetStarredRepos(userid string) ([]string, error) {
// 	return m.UserStore.GetStarredRepos(userid)
// }

// func (m MetadataStore) ToggleBookmarkRepo(userid, reponame string) (msConfig.UserState, error) {
// 	return m.UserStore.ToggleBookmarkRepo(userid, reponame)
// }
