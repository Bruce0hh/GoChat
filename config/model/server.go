package model

// GlobalConfig 全局配置
type GlobalConfig struct {
	ServerConfig     ServerConfig `mapstructure:"serverConfig"`
	PostgreSQLConfig string       `mapstructure:"postgreSqlConfig"`
	MongoDBConfig    string       `mapstructure:"mongodbConfig"`
	RedisConfig      string       `mapstructure:"redisConfig"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

// PostgreSQLConfig PGSQL配置
type PostgreSQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	DB   int    `mapstructure:"db"`
}
