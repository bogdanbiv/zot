package config

type UserMetadataStoreConfig struct {
	RootDir string
	Driver string
}

type MetadataStoreConfig struct {
	*UserMetadataStoreConfig
}

const (
	UserMetadataLocalDriver = "local"
	UserMetadataLocalFile = "storageName"
)

