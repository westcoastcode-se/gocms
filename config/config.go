package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

type Config struct {
	ListenAddr        string
	Author            bool
	PublicKeyPath     string
	PrivateKeyPath    string
	UserDatabasePath  string
	ACLDatabasePath   string
	CacheDatabasePath string
	ContentDirectory  string
	StaticURIPrefix   string
}

func GetConfig() *Config {
	var configPath string
	flag.StringVar(&configPath, "config-path", "", "Path to the a file where configuration can be found")
	config := loadConfigFromPath(configPath)

	flag.StringVar(&config.ListenAddr, "listen-addr", config.ListenAddr, "Address where we listen for incoming requests")
	flag.BoolVar(&config.Author, "author", config.Author, "If this instance allows for authoring of content")
	flag.StringVar(&config.PublicKeyPath, "public-key-path", config.PublicKeyPath, "Path to the public key")
	flag.StringVar(&config.PrivateKeyPath, "private-key-path", config.PrivateKeyPath, "Path to the private key")
	flag.StringVar(&config.UserDatabasePath, "user-db-path", config.UserDatabasePath, "Path to a database containing user information")
	flag.StringVar(&config.ACLDatabasePath, "acl-db-path", config.ACLDatabasePath, "Path to a database containing the access control list")
	flag.StringVar(&config.CacheDatabasePath, "cache-db-path", config.CacheDatabasePath, "Path to a database containing the access control list")
	flag.StringVar(&config.ContentDirectory, "content-path", config.ContentDirectory, "Path to where content can be found")
	flag.StringVar(&config.StaticURIPrefix, "static-uri-prefix", config.StaticURIPrefix, "URI prefix for")
	flag.Parse()

	return config
}

func loadConfigFromPath(path string) *Config {
	config := &Config{
		ListenAddr:        "127.0.0.1:8080",
		Author:            true,
		PublicKeyPath:     "config/key.pub",
		PrivateKeyPath:    "config/key.pem",
		UserDatabasePath:  "content/config/users.json",
		ACLDatabasePath:   "content/config/acl.json",
		CacheDatabasePath: "content/config/cache.json",
		ContentDirectory:  "content",
		StaticURIPrefix:   "/assets",
	}

	if len(path) > 0 {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bytes, config)
		if err != nil {
			panic(err)
		}
	}

	return config
}
