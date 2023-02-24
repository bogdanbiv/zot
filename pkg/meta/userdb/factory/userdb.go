package userdbfactory

import (
	"zotregistry.io/zot/errors"
	zlog "zotregistry.io/zot/pkg/log"
	metaParams "zotregistry.io/zot/pkg/meta/params"
	userdb "zotregistry.io/zot/pkg/meta/userdb"
	userbolt "zotregistry.io/zot/pkg/meta/userdb/boltdb-wrapper"
	userdynamo "zotregistry.io/zot/pkg/meta/userdb/dynamodb-wrapper"
)

func Create(
	dbtype string, parameters interface{}, log zlog.Logger,
) (userdb.UserStore, error) {
	switch dbtype {
	case "boltdb":
		{
			properParameters, ok := parameters.(metaParams.DBDriverParameters)
			if !ok {
				panic("failed type assertion")
			}

			return userbolt.NewBoltDBWrapper(properParameters, log)
		}
	case "dynamodb":
		{
			properParameters, ok := parameters.(metaParams.DBDriverParameters)
			if !ok {
				panic("failed type assertion")
			}

			return userdynamo.NewDynamoDBWrapper(properParameters, log)
		}
	default:
		{
			return nil, errors.ErrBadConfig
		}
	}
}
