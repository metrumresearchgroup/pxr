package configlib

// PkgConfig provides information about custom settings during package installation
type Customization struct {
	EnvVars []string `mapstructure:"environment_variables"`
}

// Config provides a struct for all pkgr related configuration
type Config struct {
	TestCmd        string                   `mapstructure:"testcmd" yaml:"test_cmd"`
	Packages       []string                 `mapstructure:"packages"`
	LibPaths       []string                 `mapstructure:"libpaths"`
	Customizations map[string]Customization `mapstructure:"packages"`
	Threads        int                      `mapstructure:"threads"`
	RPath          string                   `mapstructure:"rpath"`
	LogLevel       string                   `mapstructure:"loglevel"`
	AsUser         bool                     `mapstructure:"asuser" yaml:"as_user"`
	EnvVars        []string                 `mapstructure:"environment_variables" yaml:"environment_variables"`
}
