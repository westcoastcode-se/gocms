package config

import "flag"

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
	config := &Config{}

	flag.StringVar(&config.ListenAddr, "listen-addr", "127.0.0.1:8080", "Address where we listen for incoming requests")
	flag.BoolVar(&config.Author, "author", true, "If this instance allows for authoring of content")
	flag.StringVar(&config.PublicKeyPath, "public-key-path", "config/key.pub", "Path to the public key")
	flag.StringVar(&config.PrivateKeyPath, "private-key-path", "config/key.pem", "Path to the private key")
	flag.StringVar(&config.UserDatabasePath, "user-db-path", "content/config/users.json", "Path to a database containing user information")
	flag.StringVar(&config.ACLDatabasePath, "acl-db-path", "content/config/acl.json", "Path to a database containing the access control list")
	flag.StringVar(&config.CacheDatabasePath, "cache-db-path", "content/config/cache.json", "Path to a database containing the access control list")
	flag.StringVar(&config.ContentDirectory, "content-path", "content", "Path to where content can be found")
	flag.StringVar(&config.StaticURIPrefix, "static-uri-prefix", "/assets", "URI prefix for")
	flag.Parse()

	return config
}
