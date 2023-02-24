package userbolt

import (
	"encoding/json"
	"path"
	"time"

	"go.etcd.io/bbolt"

	"zotregistry.io/zot/pkg/common"
	zlog "zotregistry.io/zot/pkg/log"
	msConfig "zotregistry.io/zot/pkg/meta/config"
	merrors "zotregistry.io/zot/pkg/meta/errors"
	metaParams "zotregistry.io/zot/pkg/meta/params"
	"zotregistry.io/zot/pkg/meta/repodb"
	"zotregistry.io/zot/pkg/meta/repodb/version"
)

const (
	MetadataDBName          = "metadata"
	UserMetadataName        = "Users"
	DBExtensionName         = ".db"
	dbCacheLockCheckTimeout = 10 * time.Second
)

const (
	UserMetadataBucket = "UserMeta"
	VersionBucket      = "Version"
)

type DBParameters struct {
	RootDir string
}

type DBWrapper struct {
	DB      *bbolt.DB
	Patches []func(DB *bbolt.DB) error
	Log     zlog.Logger
}

// type UserMetadataLocalStore struct {
// 	// xTODOx: not yet logging.
// 	userMetaConfig msConfig.UserMetadataStoreConfig
// 	db             *bbolt.DB
// 	log            zlog.Logger
// }

//nolint:dupl
func (d *DBWrapper) ToggleStarRepo(userid, reponame string) (
	msConfig.UserState, error,
) { //nolint:dupl
	var res msConfig.UserState

	if err := d.DB.Update(func(tx *bbolt.Tx) error {
		userdb := tx.Bucket([]byte(UserMetadataName))
		userBucket, err := userdb.CreateBucketIfNotExists([]byte(userid))
		if err != nil {
			// this is a serious failure
			return merrors.ErrUnableToCreateUserBucket
		}

		mdata := userBucket.Get([]byte(msConfig.StarredReposKey))
		unpacked := []string{}
		if mdata != nil {
			if err = json.Unmarshal(mdata, &unpacked); err != nil {
				return merrors.ErrInvalidOldUserStarredRepos
			}
		}

		if unpacked == nil {
			// should we panic now?
			return merrors.ErrUnmarshalledRepoListIsNil
		}

		if !common.Contains(unpacked, reponame) {
			res = msConfig.Added
			unpacked = append(unpacked, reponame)
		} else {
			unpacked = common.RemoveFrom(unpacked, reponame)
			res = msConfig.Removed
		}

		var repacked []byte
		if repacked, err = json.Marshal(unpacked); err != nil {
			return merrors.ErrCouldNotMarshalStarredRepos
		}

		err = userBucket.Put([]byte(msConfig.StarredReposKey), repacked)
		if err != nil {
			return merrors.ErrCouldNotPersistData
		}

		return nil
	}); err != nil {
		return msConfig.NotChanged, err
	}

	return res, nil
}

func (d *DBWrapper) GetStarredRepos(userid string) ([]string, error) {
	starredRepos := make([]string, 0)

	err := d.DB.View(func(tx *bbolt.Tx) error { //nolint:dupl
		if userid == "" {
			return nil
		}

		userdb := tx.Bucket([]byte(UserMetadataName))
		userBucket := userdb.Bucket([]byte(userid))

		if userBucket == nil {
			return nil
		}

		mdata := userBucket.Get([]byte(msConfig.StarredReposKey))
		if mdata == nil {
			return nil
		}

		if err := json.Unmarshal(mdata, &starredRepos); err != nil {
			d.Log.Info().Str("user", userid).Err(err).Msg("unmarshal error")

			return merrors.ErrInvalidOldUserStarredRepos
		}

		if starredRepos == nil {
			starredRepos = make([]string, 0)
		}

		return nil
	})

	return starredRepos, err
}

func (d *DBWrapper) ToggleBookmarkRepo(userid, reponame string) ( //nolint:dupl
	msConfig.UserState, error,
) {
	var res msConfig.UserState

	if err := d.DB.Update(func(tx *bbolt.Tx) error { //nolint:dupl
		userdb := tx.Bucket([]byte(UserMetadataName))
		userBucket, err := userdb.CreateBucketIfNotExists([]byte(userid))
		if err != nil {
			// this is a serious failure
			return merrors.ErrUnableToCreateUserBucket
		}

		mdata := userBucket.Get([]byte(msConfig.BookmarkedReposKey))
		unpacked := []string{}
		if mdata != nil {
			if err = json.Unmarshal(mdata, &unpacked); err != nil {
				return merrors.ErrInvalidOldUserBookmarkedRepos
			}
		}

		if unpacked == nil {
			return merrors.ErrUnmarshalledRepoListIsNil
			// should we panic now?
		}

		if !common.Contains(unpacked, reponame) {
			res = msConfig.Added
			unpacked = append(unpacked, reponame)
		} else {
			unpacked = common.RemoveFrom(unpacked, reponame)
			res = msConfig.Removed
		}

		var repacked []byte
		if repacked, err = json.Marshal(unpacked); err != nil {
			return merrors.ErrCouldNotMarshalBookmarkedRepos
		}

		err = userBucket.Put([]byte(msConfig.BookmarkedReposKey), repacked)
		if err != nil {
			return merrors.ErrUnableToCreateUserBucket
		}

		return nil
	}); err != nil {
		return msConfig.NotChanged, err
	}

	return res, nil
}

func (d *DBWrapper) GetBookmarkedRepos(userid string) ([]string, error) {
	bookmarkedRepos := []string{}

	err := d.DB.View(func(tx *bbolt.Tx) error { //nolint:dupl
		if userid == "" {
			return nil
		}

		userdb := tx.Bucket([]byte(UserMetadataName))
		userBucket := userdb.Bucket([]byte(userid))

		if userBucket == nil {
			return nil
		}

		mdata := userBucket.Get([]byte(msConfig.BookmarkedReposKey))
		if mdata == nil {
			return nil
		}

		if err := json.Unmarshal(mdata, &bookmarkedRepos); err != nil {
			d.Log.Info().Str("user", userid).Err(err).Msg("unmarshal error")

			return merrors.ErrInvalidOldUserBookmarkedRepos
		}

		if bookmarkedRepos == nil {
			bookmarkedRepos = make([]string, 0)
		}

		return nil
	})

	return bookmarkedRepos, err
}

// Constructor for bbolt based drivers that implement UserMetadata.
func NewBoltDBWrapper(params metaParams.DBDriverParameters,
	log zlog.Logger) (msConfig.UserStore, error) {
	// umsc *msConfig.UserMetadataStoreConfig,
	// storageName string,
	var (
		metadataDB *bbolt.DB
		err        error
	)

	const perms = 0o600
	dbPath := path.Join(params.RootDir, "user.db")
	dbOpts := &bbolt.Options{
		Timeout:      time.Second * 10,
		FreelistType: bbolt.FreelistArrayType,
	}

	metadataDB, err = bbolt.Open(dbPath, perms, dbOpts) //nolint:gomnd
	if err != nil {
		log.Error().Err(err).Str("dbPath", dbPath).Msg("unable to create user db")

		return nil, err
	}

	// log outside the clojure
	err = metadataDB.Update(func(trax *bbolt.Tx) error {
		versionBuck, err := trax.CreateBucketIfNotExists([]byte(repodb.VersionBucket))
		if err != nil {
			return err
		}

		err = versionBuck.Put([]byte(version.DBVersionKey), []byte(version.CurrentVersion))
		if err != nil {
			return err
		}

		_, err = trax.CreateBucketIfNotExists([]byte(msConfig.UserMetadataName))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		// something went wrong
		log.Error().Err(err).Msg("unable to create a cache")

		return nil, err
	}

	return &DBWrapper{
		DB:      metadataDB,
		Patches: version.GetBoltDBPatches(),
		Log:     log,
	}, nil
}
