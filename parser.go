package parser

import (
	"errors"
	"strings"
)

var (
	ErrorProtocolNotFound      = errors.New("protocol not found")
	ErrorHostsNotFound         = errors.New("hosts not found")
	ErrorHostNameCanNotBeEmpty = errors.New("host name can't be empty")
)

type DSN struct {
	Dsn      string
	Protocol string
	*Auth
	Hosts    []*Host
	Database string
	Options  map[string]string
}

type Auth struct {
	User     string
	Password string
}

type Host struct {
	Host string
	Port string
}

func New(dsn string) (*DSN, error) {
	obj := &DSN{
		Dsn:     dsn,
		Options: make(map[string]string),
		Auth:    &Auth{},
	}
	dsn, err := obj.parseProtocol(dsn)

	if err != nil {
		return nil, err
	}

	index := strings.Index(dsn, "@")

	if index > -1 {
		tmp := dsn[0:index]
		data := strings.Split(tmp, ":")

		if len(data) == 2 {
			obj.Auth.User = data[0]
			obj.Auth.Password = data[1]
		} else {
			obj.Auth.User = data[0]
		}

		dsn = dsn[index+1:]
	}

	index = strings.Index(dsn, "?")

	if index > -1 && !strings.Contains(dsn, "/") {
		dsn = strings.Replace(dsn, "?", "/?", 1)
	}

	data := strings.Split(dsn, "/")
	err = obj.parseHosts(data[0])

	if err != nil {
		return nil, err
	}

	if len(data) == 2 {
		tmp := strings.Split(data[1], "?")

		if len(tmp[0]) > 0 {
			obj.Database = tmp[0]
		}

		if len(tmp) == 2 && len(tmp[1]) > 0 {
			obj.parseOptions(tmp[1])
		}
	}

	return obj, nil
}

func (m *DSN) parseProtocol(dsn string) (string, error) {
	index := strings.Index(dsn, "://")

	if index == -1 {
		return "", ErrorProtocolNotFound
	}

	protocol := dsn[0:index]

	if protocol == "" {
		return "", ErrorProtocolNotFound
	}

	m.Protocol = protocol
	return dsn[index+3:], nil
}

func (m *DSN) parseHosts(hosts string) error {
	if len(hosts) <= 0 {
		return ErrorHostsNotFound
	}

	splitHosts := strings.Split(hosts, ",")

	for _, v := range splitHosts {
		tmp := strings.Split(v, ":")

		if tmp[0] == "" {
			return ErrorHostNameCanNotBeEmpty
		}

		host := &Host{}

		if len(tmp) == 2 {
			host.Host = tmp[0]
			host.Port = tmp[1]
		} else {
			host.Host = tmp[0]
		}

		m.Hosts = append(m.Hosts, host)
	}

	return nil
}

func (m *DSN) parseOptions(opts string) {
	optsSplit := strings.Split(opts, "&")

	for _, v := range optsSplit {
		opt := strings.Split(v, "=")

		if len(opt) == 1 {
			m.Options[opt[0]] = ""
		} else {
			m.Options[opt[0]] = opt[1]
		}
	}
}
