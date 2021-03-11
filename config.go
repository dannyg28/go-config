package go_config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	appName        = os.Args[0]
	dbType         = "db-type"
	host           = "db-host"
	password       = "db-password"
	user           = "db-user"
	schema         = "db-schema"
	port           = "db-port"
	tls            = "db-tls"
	logLevel       = "log-level"
	logOutput      = "log-output"
	dbTypeMysql    = "mysql"
	dbTypePostgres = "postgres"
	mysqlPort      = "3306"
	postgresPort   = "5432"
	authDomain     = "auth-domain"
	authClientID   = "auth-client-id"
	authSecret     = "auth-secret"
)

type Config struct {
	DBConfig   DBConfig
	AuthConfig AuthConfig
	LogLevel   int
}

type AuthConfig struct {
	Domain   string
	ClientID string
	Secret   string
}

type DBConfig struct {
	DBConnection string
	DBType       string
}

type dbConfig struct {
	DBType     string
	DBHost     string
	DBPassword string
	DBUser     string
	DBSchema   string
	DBPort     string
	DBTls      string
}

type loggerConfig struct {
	LogOutput string
	LogLevel  int
}

var rootCmd = &cobra.Command{
	Use: appName,
	Run: func(c *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.Flags().StringP(dbType, "d", "mysql", "type of database to use for db connection string, mysql and postgres acceptable")
	rootCmd.Flags().StringP(host, "H", "localhost", "host for database")
	rootCmd.Flags().StringP(password, "p", "root", "password for database")
	rootCmd.Flags().StringP(port, "P", "", "port for database, uses default appropriate for db type if not set")
	rootCmd.Flags().StringP(tls, "t", "preferred", "configuration to use for db tls connection.Options are true, false and preferred")
	rootCmd.Flags().StringP(user, "u", "root", "user for database")
	rootCmd.Flags().StringP(schema, "s", "test", "schema for database")
	rootCmd.Flags().StringP(logOutput, "o", "", "file to write logs to.No value writes to stdout")
	rootCmd.Flags().StringP(authDomain, "D", "test.com", "domain to use for authentication backend")
	rootCmd.Flags().StringP(authClientID, "c", "test", "client-id to use for authentication backend")
	rootCmd.Flags().StringP(authSecret, "S", "test", "secret to use for authentication backend")

	rootCmd.Flags().IntP(logLevel, "l", 0, "log level, 0-4 4 being debug, 4 being fatal")
}

func initialize() {
	viper.BindPFlag(dbType, rootCmd.Flags().Lookup(dbType))
	viper.BindPFlag(host, rootCmd.Flags().Lookup(host))
	viper.BindPFlag(schema, rootCmd.Flags().Lookup(schema))
	viper.BindPFlag(user, rootCmd.Flags().Lookup(user))
	viper.BindPFlag(password, rootCmd.Flags().Lookup(password))
	viper.BindPFlag(port, rootCmd.Flags().Lookup(port))
	viper.BindPFlag(tls, rootCmd.Flags().Lookup(tls))
	viper.BindPFlag(logLevel, rootCmd.Flags().Lookup(logLevel))
	viper.BindPFlag(logOutput, rootCmd.Flags().Lookup(logOutput))
	viper.BindPFlag(authDomain, rootCmd.Flags().Lookup(authDomain))
	viper.BindPFlag(authClientID, rootCmd.Flags().Lookup(authClientID))
	viper.BindPFlag(authSecret, rootCmd.Flags().Lookup(authSecret))

	if len(appName) > 10 {
		viper.SetEnvPrefix("test")
	} else {
		viper.SetEnvPrefix(appName)
	}

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

}

func New() (*Config, error) {
	initialize()
	if err := rootCmd.Execute(); err != nil {
		return nil, fmt.Errorf("unable to execute config %w", err)
	}
	l := newLoggerConfig()
	l.SetLogging()
	dbc := newDBConfig()
	c := &Config{
		DBConfig: DBConfig{
			DBConnection: dbc.newDbConnectionString(),
			DBType:       dbc.DBType,
		},
		LogLevel:   l.LogLevel,
		AuthConfig: newAuthConfig(),
	}
	return c, nil
}

func (c *loggerConfig) SetLogging() {
	var f *os.File
	var err error
	if c.LogOutput != "" {
		f, err = os.OpenFile(c.LogOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic("unable to open  log file!")
		}

	} else {
		f = os.Stdout
	}
	if c.LogLevel == 0 {
		output := zerolog.ConsoleWriter{Out: f, TimeFormat: time.RFC3339, NoColor: false}
		log.Logger = log.Output(output).
			Level(zerolog.Level(c.LogLevel)).
			With().
			Timestamp().
			Caller().
			Stack().
			Logger()
	} else {
		log.Logger = zerolog.New(f).
			Level(zerolog.Level(c.LogLevel)).
			With().
			Timestamp().
			Caller().
			Stack().
			Logger()
	}
}

func newLoggerConfig() *loggerConfig {
	return &loggerConfig{
		LogLevel:  viper.GetInt(logLevel),
		LogOutput: viper.GetString(logOutput),
	}
}

func newDBConfig() *dbConfig {
	dbPort := viper.GetString(port)
	DBtype := viper.GetString(dbType)
	if DBtype != dbTypeMysql && DBtype != dbTypePostgres {
		panic("db-type not set to valid value!")
	}
	if dbPort == "" {
		if DBtype == dbTypeMysql {
			dbPort = mysqlPort
		} else if DBtype == dbTypePostgres {
			dbPort = postgresPort
		}
	}
	dbTLS := viper.GetString(tls)
	if dbTLS != "true" && dbTLS != "false" && dbTLS != "preferred" {
		panic("invalid tls configuration!")
	}
	return &dbConfig{
		DBType:     DBtype,
		DBHost:     viper.GetString(host),
		DBPassword: viper.GetString(password),
		DBUser:     viper.GetString(user),
		DBSchema:   viper.GetString(schema),
		DBPort:     dbPort,
		DBTls:      dbTLS,
	}
}

func (c *dbConfig) newDbConnectionString() string {
	if c.DBType == dbTypeMysql {
		tc, err := mysql.ParseDSN("root:root@(localhost:3306)/test?charset=utf8&parseTime=True")
		if err != nil {
			panic(err)
		}
		mc := tc.Clone()
		mc.User = c.DBUser
		mc.Passwd = c.DBPassword
		mc.Addr = fmt.Sprintf("%s:%s", c.DBHost, c.DBPort)
		mc.DBName = c.DBSchema
		mc.TLSConfig = c.DBTls
		return mc.FormatDSN()

	}

	if c.DBType == dbTypePostgres {
		return c.createPostgresConnectionString()
	}
	panic("db-type not set to valid value!")
}

func (c *dbConfig) createPostgresConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBSchema, c.DBTls)
}

func newAuthConfig() AuthConfig {
	return AuthConfig{
		Secret:   viper.GetString(authSecret),
		Domain:   viper.GetString(authDomain),
		ClientID: viper.GetString(authClientID),
	}
}
