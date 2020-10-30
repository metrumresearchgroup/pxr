package configlib

import (
	"bytes"
	"fmt"
	"github.com/metrumresearchgroup/rcmd"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// packrat uses R.version platform, which is not the same as the Platform
// as printed in R --version, at least on windows
func packratPlatform(p string) string {
	switch p {
	case "x86_64-w64-mingw32/x64":
		return "x86_64-w64-mingw32"
	default:
		return p
	}
}

// NewConfig initialize a Config passed in by caller
func NewConfig(cfgPath string) Config {
	var cfg Config

	viper.SetEnvPrefix("pxr")
	viper.AutomaticEnv()
	loadDefaultSettings()

	err := loadConfigFromPath(&cfg, cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	//if len(cfg.LibPaths) == 0 {
	//	rs := rcmd.NewRSettings(cfg.RPath)
	//}
	rs := rcmd.NewRSettings(cfg.RPath)
	for i := range cfg.LibPaths {
		cfg.LibPaths[i] = expandTilde(getLibraryPath(cfg.LibPaths[i], rs.Version, rs.Platform))
	}
	// For all cfg	values that can be repos, make sure that ~ is expanded to the home directory.
	cfg.RPath = expandTilde(cfg.RPath)
	return cfg
}

/// expand the ~ at the beginning of a path to the home directory.
/// consider any problems a fatal error.
func expandTilde(p string) string {
	expanded, err := homedir.Expand(p)
	if err != nil {
		log.WithFields(log.Fields{
			"path":  p,
			"error": err,
		}).Fatal("problem parsing config file -- could not expand path")
	}
	return expanded
}

/// For a list of repos, expand the ~ at the beginning of each path to the home directory.
/// consider any problems a fatal error.
func expandTildes(paths []string) []string {
	var expanded []string
	for _, p := range paths {
		newPath := expandTilde(p)
		expanded = append(expanded, newPath)
	}
	return expanded
}

func getLibraryPath(libOrType string, rversion rcmd.RVersion, platform string) string {
	switch strings.ToLower(libOrType) {
	case "packrat":
		libOrType = filepath.Join("packrat", "lib", packratPlatform(platform), rversion.ToFullString())
	case "renv":
		s := fmt.Sprintf("R-%s", rversion.ToString())
		libOrType = filepath.Join("renv", "library", s, packratPlatform(platform))
	case "pkgr":
	default:
	}
	return libOrType
}

// loadConfigFromPath loads pkc configuration into the global Viper
func loadConfigFromPath(cfg *Config, configFilename string) error {
	// for pxr should only need to run from the command line config just gives extra customization control
	if configFilename != "" {
		configFilename, _ = homedir.Expand(filepath.Clean(configFilename))
		viper.SetConfigFile(configFilename)
		b, err := ioutil.ReadFile(configFilename)
		// panic if can't find or parse config as this could be explicit to user expectations
		if err != nil {
			return fmt.Errorf("could not find a config file at path: %s", configFilename)
		}
		expb := []byte(os.ExpandEnv(string(b)))
		err = viper.ReadConfig(bytes.NewReader(expb))
		if err != nil {
			if _, ok := err.(viper.ConfigParseError); ok {
				// found config file but couldn't parse it, should error
				return fmt.Errorf("unable to parse config file with error (%s)", err)
			}
			// maybe could be more loose on this later, but for now will require a config file
			return fmt.Errorf("error during parsing config file with error (%s)", err)
		}
	}
	// this is where all the reconciliation of defaults, yaml, env, cmd line flags all hits
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return fmt.Errorf("error parsing pkgr.yml: %s\n", err)
	}
	return nil
}

// loadDefaultSettings load default settings
func loadDefaultSettings() {
	viper.SetDefault("debug", false)
	// should be one of trace,debug,info,warn,error,fatal,panic
	viper.SetDefault("loglevel", "info")
	// path to R on system, defaults to R in path
	viper.SetDefault("rpath", "R")
	viper.SetDefault("threads", runtime.NumCPU())
	// whether to introspect how a user would interactively use R
	// given a desire to run as a user, should not run with vanilla, or force override libpaths
	viper.SetDefault("asuser", false)
}

// IsCustomizationSet ...
func IsCustomizationSet(key string, elems []interface{}, elem string) bool {
	for _, v := range elems {
		for k, iv := range v.(map[interface{}]interface{}) {
			if k == elem {
				for k2 := range iv.(map[interface{}]interface{}) {
					if k2 == key {
						return true
					}
				}
			}
		}
	}
	return false
}

// SetCustomizations ... set ENV values in Rsettings
func SetCustomizations(rSettings rcmd.RSettings, cfg Config, pkg string) rcmd.RSettings {
	pkgCustomizations, isPresent := cfg.Customizations[pkg]
	if !isPresent {
		return rSettings
	}
	for _, envVars := range pkgCustomizations.EnvVars {
		rSettings.EnvVars = rcmd.NvpAppendPair(rSettings.EnvVars, envVars)
	}
	return rSettings
}

//
//func setViperCustomizations(cfg Config, pkgSettings []interface{}) {
//	for pkg, v := range cfg.Customizations {
//		if IsCustomizationSet("Suggests", pkgSettings, pkg) {
//			pkgDepTypes := dependencyConfigurations.Default
//			pkgDepTypes.Suggests = v.Suggests
//			dependencyConfigurations.Deps[pkg] = pkgDepTypes
//		}
//		if IsCustomizationSet("Repo", pkgSettings, pkg) {
//			err := pkgNexus.SetPackageRepo(pkg, v.Repo)
//			if err != nil {
//				log.WithFields(log.Fields{
//					"pkg":  pkg,
//					"repo": v.Repo,
//				}).Fatal("error finding custom repo to set")
//			}
//		}
//		if IsCustomizationSet("Type", pkgSettings, pkg) {
//			err := pkgNexus.SetPackageType(pkg, v.Type)
//			if err != nil {
//				log.WithFields(log.Fields{
//					"pkg":  pkg,
//					"repo": v.Repo,
//				}).Fatal("error finding custom repo to set")
//			}
//		}
//	}
//}
