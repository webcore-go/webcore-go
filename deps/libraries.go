package deps

import (
	"github.com/semanggilab/webcore-go/app/core"
	"github.com/semanggilab/webcore-go/lib/auth/apikey"
	"github.com/semanggilab/webcore-go/lib/authstore/yaml"
	"github.com/semanggilab/webcore-go/lib/mongo"
	"github.com/semanggilab/webcore-go/lib/pubsub"
)

var APP_LIBRARIES = map[string]core.LibraryLoader{
	// "database:postgres":     &postgres.PostgresLoader{},
	"database:mongodb": &mongo.MongoLoader{},
	// "redis":           &redis.RedisLoader{},
	"pubsub":                &pubsub.PubSubLoader{},
	"authstorage:yaml":      &yaml.YamlLoader{},
	"authentication:apikey": &apikey.ApiKeyLoader{},

	// Add your library here
}
