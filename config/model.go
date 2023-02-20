package config

type Config struct {
	Hop   HopConfig   `ini:"hop"`
	Mysql MysqlConfig `ini:"mysql"`
}

type HopConfig struct {
	Host string `ini:"host"`
	Port string `ini:"port"`
}

type MysqlConfig struct {
	Username string `ini:"username"`
	Password string `ini:"password"`
	Host     string `ini:"host"`
	Port     string `ini:"port"`
	Database string `ini:"database"`
	Dsn      string `ini:"dsn"`
}
