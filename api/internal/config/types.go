package config

type Server struct {
	HTTP HTTP `yaml:"http"`
}

type HTTP struct {
	Address      string `yaml:"address"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
}

type Database struct {
	Driver          string `yaml:"driver"`
	MaxOpenConns    string `yaml:"max_open_conns"`
	MaxIdleConns    string `yaml:"max_idle_conns"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`
	ConnMaxIdletime string `yaml:"conn_max_idletime"`

	MasterDsn         string
	FollowerDsn       string
	MasterDsnNoCred   string `yaml:"master_dsn_no_cred"`
	FollowerDsnNoCred string `yaml:"follower_dsn_no_cred"`
}
