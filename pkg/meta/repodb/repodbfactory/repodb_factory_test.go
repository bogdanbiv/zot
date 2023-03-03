package repodbfactory_test

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	zlog "zotregistry.io/zot/pkg/log"
	bolt "zotregistry.io/zot/pkg/meta/repodb/boltdb-wrapper"
	dynamoParams "zotregistry.io/zot/pkg/meta/repodb/dynamodb-wrapper/params"
	"zotregistry.io/zot/pkg/meta/repodb/repodbfactory"
)

func TestCreateDynamo(t *testing.T) {
	skipDynamo(t)

	Convey("Create", t, func() {
		dynamoDBDriverParams := dynamoParams.DBDriverParameters{
			Endpoint:              os.Getenv("DYNAMODBMOCK_ENDPOINT"),
			RepoMetaTablename:     "RepoMetadataTable",
			ManifestDataTablename: "ManifestDataTable",
			IndexDataTablename:    "IndexDataTable",
			VersionTablename:      "Version",
			Region:                "us-east-2",
		}

		log := zlog.NewLogger("INFO", "stdout")
		repoDB, err := repodbfactory.Create("dynamodb", dynamoDBDriverParams, log)
		So(repoDB, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})

	Convey("Fails", t, func() {
		log := zlog.NewLogger("INFO", "stdout")
		So(func() {
			_, _ = repodbfactory.Create("dynamodb",
				bolt.DBParameters{RootDir: "root"}, log)
		}, ShouldPanic)

		repoDB, err := repodbfactory.Create("random", bolt.DBParameters{RootDir: "root"}, log)
		So(repoDB, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

func TestCreateBoltDB(t *testing.T) {
	Convey("Create", t, func() {
		rootDir := t.TempDir()
		log := zlog.NewLogger("INFO", "stdout")

		repoDB, err := repodbfactory.Create("boltdb", bolt.DBParameters{
			RootDir: rootDir,
		}, log)
		So(repoDB, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})

	Convey("fails", t, func() {
		log := zlog.NewLogger("INFO", "stdout")
		So(func() {
			_, _ = repodbfactory.Create("boltdb", dynamoParams.DBDriverParameters{}, log)
		}, ShouldPanic)
	})
}

func skipDynamo(t *testing.T) {
	t.Helper()

	if os.Getenv("DYNAMODBMOCK_ENDPOINT") == "" {
		t.Skip("Skipping testing without AWS DynamoDB mock server")
	}
}
