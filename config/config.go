package config

// GlobalConfig 全局配置
type GlobalConfig struct {
	ServerConfig     serverConfig     `mapstructure:"server"`
	PostgreSQLConfig postgreSQLConfig `mapstructure:"postgres"`
	MongoDBConfig    mongoDBConfig    `mapstructure:"mongodb"`
	RedisConfig      redisConfig      `mapstructure:"redis"`
}

// ServerConfig 服务器配置
type serverConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

// PostgreSQLConfig PGSQL配置
type postgreSQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// MongoDBConfig MongoDB配置
type mongoDBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// RedisConfig Redis配置
type redisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}
