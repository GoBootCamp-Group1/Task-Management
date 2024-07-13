package config

type Config struct {
	Server Server `mapstructure:"server"`
	DB     DB     `mapstructure:"db"`
	Redis  Redis  `mapstructure:"redis"`
	Email  Email  `mapstructure:"email"`
}

type Server struct {
	HttpPort               int    `mapstructure:"http_port"`
	Host                   string `mapstructure:"host"`
	TokenExpMinutes        uint   `mapstructure:"token_exp_minutes"`
	RefreshTokenExpMinutes uint   `mapstructure:"refresh_token_exp_minute"`
	TokenSecret            string `mapstructure:"token_secret"`
	MaxRateLimit           int    `mapstructure:"max_rate_limit"`
	RateLimitExpiration    int    `mapstructure:"rate_limit_expiration"`
}

type DB struct {
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"pass"`
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	DBName string `mapstructure:"db_name"`
}

type Redis struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
}

type Email struct {
	SmtpHost        string `mapstructure:"smtp_host"`
	SmtpPort        int    `mapstructure:"smtp_port"`
	SmtpUsername    string `mapstructure:"smtp_username"`
	SmtpPassword    string `mapstructure:"smtp_password"`
	SmtpFromAddress string `mapstructure:"smtp_from_address"`
	SmtpEncryption  string `mapstructure:"smtp_encryption"`
	SmtpFromName    string `mapstructure:"smtp_from_name"`
}
