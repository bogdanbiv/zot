package config

type UserMetadataStoreConfig struct {
	RootDir string `json:"rootDir"`
	Driver  string `json:"driver"`
	Enabled *bool  `json:"enabled,omitempty"`
}

type MetadataStoreConfig struct {
	User *UserMetadataStoreConfig `mapstructure:"user,omitempty"`
}

const (
	UserMetadataLocalDriver  = "local"
	UserMetadataDynamoDriver = "dynamo"
	UserMetadataLocalFile    = "metadata_user"
	StarredReposKey          = "starredReposKey"
	BookmarkedReposKey       = "bookmarkedReposKey"
	UserMetadataName         = "UserMetadataKey"
)

type UserState int

const (
	NotChanged UserState = iota
	Added
	Removed
)

type UserStore interface {
	GetStarredRepos(userid string) ([]string, error)
	GetBookmarkedRepos(userid string) ([]string, error)
	ToggleStarRepo(userid, reponame string) (UserState, error)
	ToggleBookmarkRepo(userid, reponame string) (UserState, error)
}

type UserMetadata struct {
	// data for each user.
	StarredRepos    []string
	BookmarkedRepos []string
}
