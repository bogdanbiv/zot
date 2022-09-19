package metadata

import (
	zlog "zotregistry.io/zot/pkg/log"
	msConfig "zotregistry.io/zot/pkg/metadata/config"
)

type Store struct {
	UserStore
	RepoDB
}

type StoreLocal struct {
	config msConfig.MetadataStoreConfig
	Store
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

func NewBaseMetaDB(msc msConfig.MetadataStoreConfig, log zlog.Logger) (StoreLocal, error) {
	return StoreLocal{
		config: msc,
		Store:Store{
			NewUserMetadataLocalStore(msc.RootDir, "storageName", log),
			BaseRepoDB{},
		},
		
	}, nil
}

func doSmth() {
	storeLocal, _ := NewBaseMetaDB(msConfig.MetadataStoreConfig{},zlog.Logger{})
	storeLocal.GetStarredRepos("")
}

// func (m StoreLocal) ToggleStarRepo(userid, reponame string) error {
// 	return m.userDB.ToggleStarRepo(userid, reponame)
// }

// func (m StoreLocal) GetStarredRepos(userid string) ([]string, error) {
// 	return m.userDB.GetStarredRepos(userid)
// }

// func (m StoreLocal) ToggleBookmarkRepo(userid, reponame string) error {
// 	return m.userDB.ToggleBookmarkRepo(userid, reponame)
// }

// func (m StoreLocal) GetBookmarkedRepos(userid string) ([]string, error) {
// 	return m.userDB.GetBookmarkedRepos(userid)
// }

// func (m StoreLocal) RepoDB() {
// 	m.repoDB.RepoDB()
// }

type BaseRepoDB struct{}

func (bo BaseRepoDB) RepoDB() {}
