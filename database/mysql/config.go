package mysql

// Config defines mysql configuration.
type Config struct {
	Master       Conn   `yaml:"master"`
	Slaves       []Conn `yaml:"slaves"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	LogLevel     string `yaml:"log_level"`
}

type Conn struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}
