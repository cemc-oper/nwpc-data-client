package dataconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/embedded"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const ConfigFileBasename = ".yaml"

type PathItem struct {
	PathType  string `yaml:"type"`
	LevelType string `yaml:"level"`
	Path      string `yaml:"path"`
}

type DataConfig struct {
	// ConfigVersion is the reserved "_config_version" field declaring the
	// config format version. It is not used for path resolution; it exists
	// so that the field survives YAML parsing and strict structural checks.
	ConfigVersion int        `yaml:"_config_version,omitempty"`
	Default       string     `yaml:"default"`
	FileName      string     `yaml:"file_name"`
	FileNames     []string   `yaml:"file_names"`
	Paths         []PathItem `yaml:"paths"`
}

// LoadConfigContent read config content from a local file or from embedded config strings.
// The content is checked against the config format version this binary supports.
func LoadConfigContent(configDir string, dataType string) (string, error) {
	if len(configDir) == 0 {
		content, err := loadContentFromEmbedded(dataType)
		if err != nil {
			return "", fmt.Errorf("find embedded config for data type %s has error: %v", dataType, err)
		}
		if err := CheckConfigVersion(content); err != nil {
			return "", fmt.Errorf("embedded config for data type %s: %v", dataType, err)
		}
		return content, nil
	}

	content, err := loadContentFromLocal(configDir, dataType)
	if err != nil {
		return "", err
	}

	return content, nil
}

// loadContentFromLocal read config content from a local config file.
func loadContentFromLocal(configDir string, dataType string) (string, error) {
	configFilePath, err := findLocalConfig(configDir, dataType)
	if err != nil {
		return "", fmt.Errorf("model data type config is not found: %v", err)
	}

	return LoadConfigContentFromFile(configFilePath)
}

// LoadConfigContentFromFile read config content from a config file path.
// The content is checked against the config format version this binary supports.
// Tools that must accept every format version (such as "config lint") should
// read the file directly instead.
func LoadConfigContentFromFile(configFilePath string) (string, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return "", err
	}
	content := string(data)
	if err := CheckConfigVersion(content); err != nil {
		return "", fmt.Errorf("%s: %v", configFilePath, err)
	}
	return content, nil
}

// ParseConfigContent parse config content template with start time and forecast time.
func ParseConfigContent(content string, startTime time.Time, forecastTime time.Duration, member string) (DataConfig, error) {
	dataConfig := DataConfig{}
	tpVar := GenerateConfigTemplateVariable(startTime, forecastTime, member)

	forecastTimeStr := FormatForecastTimeShort(forecastTime)
	currentLog := log.WithFields(log.Fields{"forecastTime": forecastTimeStr})

	contentTemplate := template.Must(template.New("dataConfig").Funcs(ConfigTemplateFuncMap()).Parse(content))
	var configBuilder strings.Builder
	err := contentTemplate.Execute(&configBuilder, tpVar)
	if err != nil {
		currentLog.Errorf("file name template execute has error: %s", err)
		return dataConfig, err
	}
	parsedContent := configBuilder.String()

	dataConfig, err = LoadDataConfigFromContent(parsedContent)
	if err != nil {
		currentLog.Errorf("load dataConfig from content has error: %s", err)
		return dataConfig, err
	}
	return dataConfig, nil
}

// loadContentFromEmbedded find data config string in embedded configs.
func loadContentFromEmbedded(dataType string) (string, error) {
	for _, configItem := range embedded.EmbeddedConfigs {
		if configItem[0] == dataType {
			return configItem[1], nil
		}
	}
	return "", fmt.Errorf("can't find embedded config: %s", dataType)
}

// LoadDataConfigFromContent parse DataConfig from string content.
func LoadDataConfigFromContent(content string) (DataConfig, error) {
	dataConfig := DataConfig{}
	err := yaml.Unmarshal([]byte(content), &dataConfig)
	if err != nil {
		return dataConfig, err
	}
	return dataConfig, nil
}

// locationLevelPattern allows a single level ("runtime") or a combined level
// ("runtime/archive") as used by existing configs.
var locationLevelPattern = regexp.MustCompile(`^(runtime|archive|storage)(/(runtime|archive|storage))*$`)

// ValidateDataConfig validates the structure of a rendered (template-free)
// config document: strict YAML decoding (unknown fields are rejected) plus
// field-level checks. It mirrors schemas/config/data-config.schema.json.
func ValidateDataConfig(content string) error {
	dataConfig := DataConfig{}
	decoder := yaml.NewDecoder(strings.NewReader(content))
	decoder.KnownFields(true)
	if err := decoder.Decode(&dataConfig); err != nil {
		return fmt.Errorf("yaml parse: %v", err)
	}
	return ValidateDataConfigStruct(dataConfig)
}

// ValidateDataConfigStruct checks the structural rules of a parsed DataConfig.
func ValidateDataConfigStruct(dataConfig DataConfig) error {
	if len(dataConfig.Paths) == 0 {
		return fmt.Errorf("paths is required and must not be empty")
	}
	for i, item := range dataConfig.Paths {
		if item.PathType == "" {
			return fmt.Errorf("paths[%d]: type is required", i)
		}
		if item.PathType != "local" {
			return fmt.Errorf("paths[%d]: unsupported type %q (only \"local\" is supported)", i, item.PathType)
		}
		if item.LevelType == "" {
			return fmt.Errorf("paths[%d]: level is required", i)
		}
		if !locationLevelPattern.MatchString(item.LevelType) {
			return fmt.Errorf("paths[%d]: invalid level %q (allowed: runtime, archive, storage, or a \"/\"-joined combination)", i, item.LevelType)
		}
		if item.Path == "" {
			return fmt.Errorf("paths[%d]: path is required", i)
		}
	}
	return nil
}

func findLocalConfig(configDir string, dataType string) (string, error) {
	configDirPath, _ := filepath.Abs(configDir)
	configFilePath := filepath.Join(configDirPath, dataType+ConfigFileBasename)
	relPath, err := filepath.Rel(configDirPath, configFilePath)
	if err != nil {
		return "", fmt.Errorf("check directory match has error: %v", err)
	}
	if relPath[0:2] == ".." {
		return "", fmt.Errorf("check directory failed: dataType (%s) should be under configDir (%s)",
			dataType, configDirPath)
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return configFilePath, fmt.Errorf("file is not exist")
	}
	return configFilePath, nil
}
