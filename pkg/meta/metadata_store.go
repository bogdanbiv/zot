package meta

import (
	zlog "zotregistry.io/zot/pkg/log"
	msConfig "zotregistry.io/zot/pkg/meta/config"
)

type MetadataStore struct {
	userDB UserStore
}

type UserStore interface {
	GetStarredRepos(userid string) ([]string, error)
	GetBookmarkedRepos(userid string) ([]string, error)
	ToggleStarRepo(userid, reponame string) error
	ToggleBookmarkRepo(userid, reponame string) error
}

type RepoDB interface {
	RepoDB()
}

func NewBaseMetaDB(msc msConfig.MetadataStoreConfig, log zlog.Logger) (MetadataStore, error) {
	var (
		userdata UserStore
	)

	if msc.UserMetadataStoreConfig != nil {
		userdata = NewUserMetadataLocalStore(
			msc.UserMetadataStoreConfig, msConfig.UserMetadataLocalFile, log)
	}

	return MetadataStore{
		// config: msc,
		userDB: userdata,
	}, nil
}

func (m MetadataStore) GetBookmarkedRepos(userid string) ([]string, error) {
	return m.userDB.GetBookmarkedRepos(userid)
	// return []string{}, nil
}

func (m MetadataStore) ToggleStarRepo(userid, reponame string) error {
	return m.userDB.ToggleStarRepo(userid, reponame)
	// return nil
}

func (m MetadataStore) GetStarredRepos(userid string) ([]string, error) {
	return m.userDB.GetStarredRepos(userid)
	// return []string{}, nil
}

func (m MetadataStore) ToggleBookmarkRepo(userid, reponame string) error {
	return m.userDB.ToggleBookmarkRepo(userid, reponame)
	// return nil
}
