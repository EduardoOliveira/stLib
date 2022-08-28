package runtime

type Config struct {
	LibraryPath string `toml:"library_path"`
	Port        int    `toml:"port"`
}

var Cfg *Config
