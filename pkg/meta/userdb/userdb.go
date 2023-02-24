package userdb

import (
	msConfig "zotregistry.io/zot/pkg/meta/config"
)

type UserStore interface {
	GetStarredRepos(userid string) ([]string, error)
	GetBookmarkedRepos(userid string) ([]string, error)
	ToggleStarRepo(userid, reponame string) (msConfig.UserState, error)
	ToggleBookmarkRepo(userid, reponame string) (msConfig.UserState, error)
}

type UserMetadata struct {
	// data for each user.
	StarredRepos    []string
	BookmarkedRepos []string
}

type UserMetadataEmptyStore struct{}

func (umes *UserMetadataEmptyStore) ToggleBookmarkRepo(userid, reponame string) (msConfig.UserState, error) {
	return msConfig.NotChanged, nil
}

func (umes *UserMetadataEmptyStore) GetBookmarkedRepos(userid string) ([]string, error) {
	return []string{}, nil
}

func (umes *UserMetadataEmptyStore) ToggleStarRepo(userid, reponame string) (msConfig.UserState, error) {
	return msConfig.NotChanged, nil
}

func (umes *UserMetadataEmptyStore) GetStarredRepos(userid string) ([]string, error) {
	return []string{}, nil
}
