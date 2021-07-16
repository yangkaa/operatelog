package config

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

// C represents a global configuration.
var C *Config

// Config is the global configurations.
type Config struct {
	DB              *DB
	ListenPort      int
	ListenHost      string
	EnableAccessLog bool
	LogLevel        string
	AccessToken     string
}

// DB holds configurations for database.
type DB struct {
	Addr        string
	User        string
	Pass        string
	Name        string
	TablePrefix string
}

//GetDefaultConfig get default config
func GetDefaultConfig() *Config {
	return &Config{
		DB: &DB{
			Addr: "127.0.0.1:3306",
			User: "root",
			Pass: "goodrain",
			Name: "operatelog",
		},
		ListenHost:  "0.0.0.0",
		ListenPort:  8080,
		LogLevel:    "info",
		AccessToken: "test@token",
	}
}

// Parse -
func Parse() {
	c := GetDefaultConfig()
	kingpin.Flag("db-addr", "The addr of database").Default(c.DB.Addr).Envar("MYSQL_ADDR").StringVar(&c.DB.Addr)
	kingpin.Flag("db-user", "The user name of database").Default(c.DB.User).Envar("MYSQL_USER").StringVar(&c.DB.User)
	kingpin.Flag("db-pass", "The password of database").Default(c.DB.Pass).Envar("MYSQL_PASS").StringVar(&c.DB.Pass)
	kingpin.Flag("db-name", "The database name of database").Default(c.DB.Name).Envar("MYSQL_DB").StringVar(&c.DB.Name)
	kingpin.Flag("db-table-prefix", "The prefix of table").Default(c.DB.TablePrefix).Envar("MYSQL_TABLE_PREFIX").StringVar(&c.DB.TablePrefix)
	kingpin.Flag("listen-host", "The listen host of server").Default(c.ListenHost).Envar("LISTEN_HOST").StringVar(&c.ListenHost)
	kingpin.Flag("listen-port", "The listen port of server").Default(fmt.Sprintf("%d", c.ListenPort)).Envar("PORT").IntVar(&c.ListenPort)
	kingpin.Flag("access-token", "token use create log").Default(c.AccessToken).Envar("ACCESS_TOKEN").StringVar(&c.AccessToken)

	kingpin.CommandLine.GetFlag("help").Short('h')
	kingpin.Parse()
	C = c
}
