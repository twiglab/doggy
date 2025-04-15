package conf

type PlatformConfig struct {
	BaseURL string `yaml:"base-url" mapstructure:"base-url"`
	Address string `yaml:"address" mapstructure:"address"`
	Port    int    `yaml:"port" mapstructure:"port"`
}

type CommonDeviceConfig struct {
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
}

type Web struct {
	Addr     string `yaml:"addr" mapstructure:"addr"`
	CertFile string `yaml:"cert-file" mapstructure:"cert-file"`
	KeyFile  string `yaml:"key-file" mapstructure:"key-file"`
}

type DB struct {
	DSN string `yaml:"dsn" mapstructure:"dsn"`
}

type App struct {
	ID                 string             `yaml:"id" mapstructure:"id"`
	PlatformConfig     PlatformConfig     `yaml:"platform" mapstructure:"platform"`
	CommonDeviceConfig CommonDeviceConfig `yaml:"device" mapstructure:"device"`
	Web                Web                `yaml:"server" mapstructure:"server"`
	DB                 DB                 `yaml:"db" mapstructure:"db"`
}
