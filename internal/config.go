package internal

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const SAMPLE_CONFIG_FILE string = "sample_config.toml"

//go:embed sample_config.toml
var sample_config embed.FS

// Create the sample chatbot config file in the user's local
// filesystem or die trying.
func createConfig() (string, string) {
	home, err := os.UserHomeDir()
	PanicOnErr(err)
	chatbot_path := filepath.Join(home, ".chatbot")
	chatbot_config := filepath.Join(chatbot_path, "config")
	err = os.MkdirAll(chatbot_path, os.ModePerm)
	PanicOnErr(err)
	copySampleConfigFile(chatbot_config)
	return chatbot_path, chatbot_config
}

func copySampleConfigFile(out_path string) {
	config_file_contents, err := sample_config.ReadFile(SAMPLE_CONFIG_FILE)
	PanicOnErr(err)
	err = os.WriteFile(out_path, config_file_contents, os.ModePerm)
	PanicOnErr(err)
}

// Parse the config with Viper and handle errors
func ParseConfigWithViper() {
	err := viper.ReadInConfig()
	PanicOnErr(err)
}

func InitChatbot() {
	path, _ := createConfig()
	fmt.Println("Created Chatbot config file at ", path)
}

// Checks to see if --verbose is set by the user
// by checking the 'verbose' viper setting
func CheckDebug() bool {
	return viper.GetBool("verbose")
}

func GetViperValueWithDefault(config_value_name string, default_value string) string {
	val := viper.GetString(config_value_name)
	if val != "" {
		return val
	}
	return default_value
}
