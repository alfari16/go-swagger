package generator

import (
	"embed"
	"fmt"
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"path/filepath"
)

const DefaultConfigFile = "default_swagger_config_alfari16.yaml"

//go:embed swagger.yaml
var defaultConfig embed.FS

// LanguageDefinition in the configuration file.
type LanguageDefinition struct {
	Layout SectionOpts `mapstructure:"layout"`
}

// ConfigureOpts for generation
func (d *LanguageDefinition) ConfigureOpts(opts *GenOpts) error {
	opts.Sections = d.Layout
	if opts.LanguageOpts == nil {
		opts.LanguageOpts = GoLangOpts()
	}
	return nil
}

// LanguageConfig structure that is obtained from parsing a config file
type LanguageConfig map[string]LanguageDefinition

// ReadConfig at the specified path, when no path is specified it will look into
// the current directory and load a swagger.yaml.{yml,json,hcl,toml,properties} file
// Returns a viper config or an error
func ReadConfig(fpath string) (*viper.Viper, error) {
	v := viper.New()

	var file fs.File
	var err error
	if fpath == DefaultConfigFile {
		file, err = defaultConfig.Open("swagger.yaml")
	} else {
		abspath, err := filepath.Abs(fpath)
		if err != nil {
			return nil, err
		}
		file, err = os.Open(abspath)
		if err != nil {
			return nil, err
		}
		if !fileExists(fpath, "") {
			return nil, fmt.Errorf("can't find file for %q", fpath)
		}
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	ext := filepath.Ext(fpath)
	if len(ext) > 0 {
		ext = ext[1:]
	}
	v.SetConfigType(ext)
	if err := v.ReadConfig(file); err != nil {
		return nil, err
	}
	return v, nil
}
