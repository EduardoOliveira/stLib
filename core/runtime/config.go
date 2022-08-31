package runtime

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	LibraryPath      string   `toml:"library_path"`
	Port             int      `toml:"port"`
	MaxRenderWorkers int      `toml:"max_render_workers"`
	FileBlacklist    []string `toml:"file_blacklist"`
	ModelRenderColor string   `toml:"model_render_color"`
	ModelBackgroundColor string   `toml:"model_background_color"`
}

var Cfg *Config

func init() {
	_, err := toml.DecodeFile("config.toml", &Cfg)
	if err != nil {
		log.Fatal("Unable to read config file: ", err)
	}

	Cfg.FileBlacklist = append(Cfg.FileBlacklist, ".gitignore", ".gitkeep", ".DS_Store", ".project.stlib", ".thumb.png")

	if Cfg.MaxRenderWorkers == 0 {
		Cfg.MaxRenderWorkers = 5
	}
	if Cfg.ModelRenderColor == "" {
		Cfg.ModelRenderColor = "#FFFFFF"
	}
	if Cfg.ModelRenderColor == "" {
		Cfg.ModelRenderColor = "#FFFFFF"
	}

}
