package params

type DBDriverParameters struct {
	Endpoint, Region, RepoMetaTablename, ManifestDataTablename, VersionTablename string
	UserMetaTablename                                                            string
	RootDir, Driver                                                              string
}

type DynamoDriverParameters struct {
	Endpoint, Region, RepoMetaTablename, ManifestDataTablename, VersionTablename string
	UserMetaTablename                                                            string
}

type RepoDynamoDriverParameters struct {
	Endpoint, Region, RepoMetaTablename, ManifestDataTablename, VersionTablename string
}

type UserDynamoDriverParameters struct {
	UserMetaTablename string
}

type BoltDriverParameters struct {
	RootDir string
}
