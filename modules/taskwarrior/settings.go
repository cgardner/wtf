package taskwarrior

import (
	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
)

const (
	defaultFocusable = true
	defaultTitle     = "TaskWarrior"
)

// Settings defines the configuration properties for this module
type Settings struct {
	common *cfg.Common

	filePaths            []interface{}
	maxDescriptionLength int
	maxProjectLength     int
}

// NewSettingsFromYAML creates a new settings instance from a YAML config block
func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {

	settings := Settings{
		common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		maxDescriptionLength: ymlConfig.UInt("maxDescriptionLength", 60),
		maxProjectLength:     ymlConfig.UInt("maxProjectLength", 30),
	}

	return &settings
}
