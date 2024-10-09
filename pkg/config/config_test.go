package config

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/grd/FreePDM/pkg/util"
	"github.com/stretchr/testify/assert"
)

// Test GetUid
func TestGetUid(t *testing.T) {
	Conf.Users = map[string]int{
		"testuser": 1001,
		"admin":    1002,
	}

	t.Run("existing user", func(t *testing.T) {
		uid := GetUid("testuser")
		assert.Equal(t, 1001, uid)
	})

	t.Run("non-existing user", func(t *testing.T) {
		uid := GetUid("nonexistent")
		assert.Equal(t, -1, uid)
	})
}

// Test ReadConfig
func TestReadConfig(t *testing.T) {
	// Prepare a temporary config file
	tmpDir := t.TempDir()
	tmpFile := path.Join(tmpDir, "FreePDM.toml")
	ioutil.WriteFile(tmpFile, []byte(`
		StartupDirectory = "/tmp"
		LogFile = "log.txt"
		LogLevel = "info"
		[Users]
		"testuser" = 1001
		`), 0644)

	// Backup original values
	origConfigName := configName
	configName = tmpFile
	defer func() { configName = origConfigName }()

	t.Run("valid config", func(t *testing.T) {
		err := ReadConfig()
		assert.NoError(t, err)
		assert.Equal(t, "/tmp", Conf.StartupDirectory)
		assert.Equal(t, 1001, Conf.Users["testuser"])
	})

	t.Run("invalid config file", func(t *testing.T) {
		// Switch to a non-existing config file
		configName = path.Join(tmpDir, "nonexistent.toml")
		err := ReadConfig()
		assert.Error(t, err)
	})
}

// Test WriteConfig
func TestWriteConfig(t *testing.T) {
	// Prepare a temporary config file
	tmpDir := t.TempDir()
	tmpFile := path.Join(tmpDir, "FreePDM.toml")

	Conf = Config{
		StartupDirectory: "/home/test",
		LogFile:          "test.log",
		LogLevel:         "debug",
		Users: map[string]int{
			"testuser": 1001,
		},
	}

	configName = tmpFile

	err := WriteConfig()
	assert.NoError(t, err)

	// Check if the file was written
	_, err = os.Stat(tmpFile)
	assert.NoError(t, err)

	// Read and check contents
	content, err := ioutil.ReadFile(tmpFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "StartupDirectory = \"/home/test\"")
}

// Test Config.String()
func TestConfigString(t *testing.T) {
	Conf = Config{
		StartupDirectory: "/home/test",
		LogFile:          "test.log",
		LogLevel:         "debug",
		Users: map[string]int{
			"testuser": 1001,
		},
	}

	str := Conf.String()
	assert.Contains(t, str, "StartupDirectory = \"/home/test\"")
	assert.Contains(t, str, "testuser = 1001")
}

// Test InitializeConfig (indirectly tests init logic)
func TestInitializeConfig(t *testing.T) {
	// Prepare a temporary config and directory
	tmpDir := t.TempDir()
	configDir = tmpDir
	configName = path.Join(configDir, "FreePDM.toml")

	// Simulate initialization (should create directory and config file)
	InitializeConfig()

	// Check if directory was created
	_, err := os.Stat(configDir)
	assert.NoError(t, err)

	// Check if the config file was created
	_, err = os.Stat(configName)
	assert.NoError(t, err)
}

func InitializeConfig() {
	// create the new directory if it doesn't exist
	if !util.DirExists(configDir) {
		os.Mkdir(configDir, 0700)
	}

	// create a new config file when it doesn't exist
	if !util.FileExists(configName) {
		WriteConfig()
	}

	// Reading the configuration file
	ReadConfig()
}

func init() {
	// Call the initialization logic when the package is loaded
	InitializeConfig()
}
