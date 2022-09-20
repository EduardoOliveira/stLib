package runtime

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	LibraryPath          string   `toml:"library_path"`
	Port                 int      `toml:"port"`
	MaxRenderWorkers     int      `toml:"max_render_workers"`
	FileBlacklist        []string `toml:"file_blacklist"`
	ModelRenderColor     string   `toml:"model_render_color"`
	ModelBackgroundColor string   `toml:"model_background_color"`
	ThingiverseToken     string   `toml:"thingiverse_token"`
	LogPath              string   `toml:"log_path"`
}

var Cfg *Config

func init() {

	viper.SetDefault("port", 8000)
	viper.SetDefault("library_path", "./library")
	viper.SetDefault("max_render_workers", 5)
	viper.SetDefault("file_blacklist", []string{})
	viper.SetDefault("model_render_color", "#ffffff")
	viper.SetDefault("model_background_color", "#000000")
	viper.SetDefault("thingiverse_token", "")
	viper.SetDefault("log_path", "./log")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	viper.SetConfigType("toml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("error config file: %w", err)
	}

	Cfg = &Config{
		LibraryPath:          viper.GetString("library_path"),
		Port:                 viper.GetInt("port"),
		MaxRenderWorkers:     viper.GetInt("max_render_workers"),
		FileBlacklist:        append(viper.GetStringSlice("file_blacklist"), ".gitignore", ".gitkeep", ".DS_Store", ".project.stlib", ".thumb.png", ".gitkeep"),
		ModelRenderColor:     viper.GetString("model_render_color"),
		ModelBackgroundColor: viper.GetString("model_background_color"),
		ThingiverseToken:     viper.GetString("thingiverse_token"),
		LogPath:              viper.GetString("log_path"),
	}
}
