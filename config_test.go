package go_config

import (
	"os"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestNewDBConfig(t *testing.T) {
	thost := "TEST_DB_HOST"
	tport := "TEST_DB_PORT"
	tschema := "TEST_DB_SCHEMA"
	tuser := "TEST_DB_USER"
	tpassword := "TEST_DB_PASSWORD"
	ttls := "TEST_DB_TLS"
	ttype := "TEST_DB_TYPE"
	tlevel := "TEST_LOG_LEVEL"
	toutput := "TEST_LOG_OUTPUT"
	tdomain := "TEST_AUTH_DOMAIN"
	tclientid := "TEST_AUTH_CLIENT_ID"
	tsecret := "TEST_AUTH_SECRET"
	test := "testing"
	envKeys := []string{
		thost, tport, tschema, tuser, tpassword, toutput, tdomain, tclientid, tsecret,
	}
	t.Run("test defaults", func(t *testing.T) {
		os.Clearenv()
		c, _ := New()
		assert.Equal(t, c.LogLevel, 0)
		assert.Equal(t, c.DBConnection, "root:root@tcp(localhost:3306)/test?parseTime=true&tls=preferred&charset=utf8")
		assert.Equal(t, c.AuthConfig, AuthConfig{
			Domain:   "test.com",
			ClientID: "test",
			Secret:   "test",
		})
	})
	t.Run("test defaults, postgres", func(t *testing.T) {
		os.Clearenv()
		os.Setenv(ttype, "postgres")
		c, _ := New()
		assert.Equal(t, c.LogLevel, 0)
		assert.Equal(t, c.DBConnection, "host=localhost port=5432 user=root password=root dbname=test sslmode=preferred")
		assert.Equal(t, c.AuthConfig, AuthConfig{
			Domain:   "test.com",
			ClientID: "test",
			Secret:   "test",
		})
	})
	t.Run("test changes w/ env variables", func(t *testing.T) {
		os.Clearenv()
		for _, key := range envKeys {
			os.Setenv(key, test)
		}
		os.Setenv(tlevel, "1")
		os.Setenv(ttls, "true")
		c, _ := New()
		assert.Equal(t, c.LogLevel, 1)
		assert.Equal(t, c.DBConnection, "testing:testing@tcp(testing:testing)/testing?parseTime=true&tls=true&charset=utf8")
		assert.Equal(t, c.AuthConfig, AuthConfig{
			Domain:   "testing",
			ClientID: "testing",
			Secret:   "testing",
		})
	})
	t.Run("test changes w/ env variables, postgres", func(t *testing.T) {
		os.Clearenv()
		for _, key := range envKeys {
			os.Setenv(key, test)
		}
		os.Setenv(tlevel, "1")
		os.Setenv(ttls, "true")
		os.Setenv(ttype, "postgres")
		c, _ := New()
		assert.Equal(t, c.LogLevel, 1)
		assert.Equal(t, c.DBConnection, "host=testing port=testing user=testing password=testing dbname=testing sslmode=true")
		assert.Equal(t, c.AuthConfig, AuthConfig{
			Domain:   "testing",
			ClientID: "testing",
			Secret:   "testing",
		})
	})
	t.Run("test changes w/flags, mysql", func(t *testing.T) {
		os.Args = []string{"test", "-H", test, "-u", test, "-s", test, "-p", test, "-P", test, "-t", "true", "-d", "mysql", "-l", "1", "-D", test, "-S", test, "-c", test}
		c, _ := New()
		assert.Equal(t, c.LogLevel, 1)
		assert.Equal(t, c.DBConnection, "testing:testing@tcp(testing:testing)/testing?parseTime=true&tls=true&charset=utf8")
		assert.Equal(t, c.AuthConfig, AuthConfig{
			Domain:   "testing",
			ClientID: "testing",
			Secret:   "testing",
		})
	})
	t.Run("test changes w/flags, postgres", func(t *testing.T) {
		os.Args = []string{"test", "-H", test, "-u", test, "-s", test, "-p", test, "-P", test, "-t", "true", "-d", "postgres", "-l", "1", "-D", test, "-S", test, "-c", test}
		c, _ := New()
		assert.Equal(t, c.LogLevel, 1)
		assert.Equal(t, c.DBConnection, "host=testing port=testing user=testing password=testing dbname=testing sslmode=true")
		assert.Equal(t, c.AuthConfig, AuthConfig{
			Domain:   "testing",
			ClientID: "testing",
			Secret:   "testing",
		})
	})
}
