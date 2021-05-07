package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDsnParse_Ok(t *testing.T) {
	cases := []*DSN{
		{
			Dsn:      "mongodb://user:password@db1:27017,db2:27018,db3:27019/test?readPreference=primary&ssl=true&w=majority",
			Protocol: "mongodb",
			Auth: &Auth{
				User:     "user",
				Password: "password",
			},
			Hosts: []*Host{
				{Host: "db1", Port: "27017"},
				{Host: "db2", Port: "27018"},
				{Host: "db3", Port: "27019"},
			},
			Database: "test",
			Options:  map[string]string{"readPreference": "primary", "ssl": "true", "w": "majority"},
		},
		{
			Dsn:      "postgres://user:password@db1:5432/test?sslmode=disable",
			Protocol: "postgres",
			Auth: &Auth{
				User:     "user",
				Password: "password",
			},
			Hosts: []*Host{
				{Host: "db1", Port: "5432"},
			},
			Database: "test",
			Options:  map[string]string{"sslmode": "disable"},
		},
		{
			Dsn:      "mysql://user:@db1:3306/test",
			Protocol: "mysql",
			Auth: &Auth{
				User:     "user",
				Password: "",
			},
			Hosts: []*Host{
				{Host: "db1", Port: "3306"},
			},
			Database: "test",
			Options:  map[string]string{},
		},
		{
			Dsn:      "mongodb://db1:27017/test?readPreference=primary&ssl=true",
			Protocol: "mongodb",
			Auth: &Auth{
				User:     "",
				Password: "",
			},
			Hosts: []*Host{
				{Host: "db1", Port: "27017"},
			},
			Database: "test",
			Options:  map[string]string{"readPreference": "primary", "ssl": "true"},
		},
		{
			Dsn:      "mysql://user@db1:3306/test",
			Protocol: "mysql",
			Auth: &Auth{
				User:     "user",
				Password: "",
			},
			Hosts: []*Host{
				{Host: "db1", Port: "3306"},
			},
			Database: "test",
			Options:  map[string]string{},
		},
		{
			Dsn:      "mongodb://db1:27017/?readPreference=primary&ssl=true",
			Protocol: "mongodb",
			Auth: &Auth{
				User:     "",
				Password: "",
			},
			Hosts: []*Host{
				{Host: "db1", Port: "27017"},
			},
			Database: "",
			Options:  map[string]string{"readPreference": "primary", "ssl": "true"},
		},
		{
			Dsn:      "mongodb://db1:/?readPreference=primary&ssl=true",
			Protocol: "mongodb",
			Auth: &Auth{
				User:     "",
				Password: "",
			},
			Hosts: []*Host{
				{Host: "db1", Port: ""},
			},
			Database: "",
			Options:  map[string]string{"readPreference": "primary", "ssl": "true"},
		},
		{
			Dsn:      "mongodb://db1/?readPreference=primary&ssl=true",
			Protocol: "mongodb",
			Auth: &Auth{
				User:     "",
				Password: "",
			},
			Hosts: []*Host{
				{Host: "db1", Port: ""},
			},
			Database: "",
			Options:  map[string]string{"readPreference": "primary", "ssl": "true"},
		},
		{
			Dsn:      "mongodb://db1:27017?readPreference&ssl=true",
			Protocol: "mongodb",
			Auth: &Auth{
				User:     "",
				Password: "",
			},
			Hosts: []*Host{
				{Host: "db1", Port: "27017"},
			},
			Database: "",
			Options:  map[string]string{"readPreference": "", "ssl": "true"},
		},
	}

	for _, val := range cases {
		res, err := New(val.Dsn)
		assert.NoError(t, err)
		assert.Equal(t, val.Dsn, res.Dsn)
		assert.Equal(t, val.Protocol, res.Protocol)
		assert.NotNil(t, res.Auth)
		assert.Equal(t, val.Auth.User, res.Auth.User)
		assert.Equal(t, val.Auth.Password, res.Auth.Password)

		assert.Len(t, res.Hosts, len(val.Hosts))

		for key, host := range val.Hosts {
			assert.Equal(t, res.Hosts[key].Host, host.Host)
			assert.Equal(t, res.Hosts[key].Port, host.Port)
		}

		assert.Equal(t, val.Database, res.Database)
		assert.Equal(t, val.Options, res.Options)
	}
}

func TestDsnParse_ProtocolNotExists_Error(t *testing.T) {
	cases := []string{
		"://user:password@db1:5432/test?sslmode=disable",
		"//user:password@db1:5432/test?sslmode=disable",
		"/user:password@db1:5432/test?sslmode=disable",
		"user:password@db1:5432/test?sslmode=disable",
	}

	for _, val := range cases {
		_, err := New(val)
		assert.Error(t, err)
		assert.Equal(t, err, ErrorProtocolNotFound)
	}
}

func TestDsnParse_BadHost_Error(t *testing.T) {
	cases := []struct {
		Val      string
		Expected error
	}{
		{Val: "postgres://user:password@:5432/test?sslmode=disable", Expected: ErrorHostNameCanNotBeEmpty},
		{Val: "postgres://user:password@/test?sslmode=disable", Expected: ErrorHostsNotFound},
		{Val: "postgres:///test?sslmode=disable", Expected: ErrorHostsNotFound},
	}

	for _, val := range cases {
		_, err := New(val.Val)
		assert.Error(t, err)
		assert.Equal(t, err, val.Expected)
	}
}
