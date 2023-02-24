package params

type DBDriverParameters struct {
	Endpoint, Region, RepoMetaTablename, ManifestDataTablename, VersionTablename string
	StarredReposTablename, BookedReposTablename, UserMetaTablename               string
	RootDir                                                                      string
}
