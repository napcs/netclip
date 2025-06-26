package netclip_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"netclip"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigTailscale(t *testing.T) {
	configContent := `port: "8080"
tailscale:
  enabled: true
  hostname: "test-netclip"
  use_tls: true`

	// Create temporary config file
	tmpfile, err := os.CreateTemp("", "netclip-config-*.yml")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte(configContent))
	assert.NoError(t, err)
	err = tmpfile.Close()
	assert.NoError(t, err)

	// Load config
	config, err := netclip.LoadConfig(tmpfile.Name())
	assert.NoError(t, err)

	// Verify values
	assert.Equal(t, "8080", config.Port)
	assert.True(t, config.Tailscale.Enabled)
	assert.Equal(t, "test-netclip", config.Tailscale.Hostname)
	assert.True(t, config.Tailscale.UseTLS)
}

func TestLoadConfigTailscaleWithEnvAuthKey(t *testing.T) {
	configContent := `tailscale:
  enabled: true
  hostname: "test-netclip"`

	// Set environment variable
	os.Setenv("TS_AUTHKEY", "test-auth-key")
	defer os.Unsetenv("TS_AUTHKEY")

	// Create temporary config file
	tmpfile, err := os.CreateTemp("", "netclip-config-*.yml")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte(configContent))
	assert.NoError(t, err)
	err = tmpfile.Close()
	assert.NoError(t, err)

	// Load config
	config, err := netclip.LoadConfig(tmpfile.Name())
	assert.NoError(t, err)

	// Verify Tailscale config loaded correctly
	assert.True(t, config.Tailscale.Enabled)
	assert.Equal(t, "test-netclip", config.Tailscale.Hostname)
}

func TestLoadConfigFileNotFound(t *testing.T) {
	// Try to load a non-existent config file
	_, err := netclip.LoadConfig("/nonexistent/path/config.yml")
	assert.Error(t, err)
}

func TestLoadConfigMalformedYAML(t *testing.T) {
	malformedContent := `port: "8080"
tailscale:
  enabled: true
  hostname: [invalid yaml structure
  use_tls: missing closing bracket`

	// Create temporary config file with malformed YAML
	tmpfile, err := os.CreateTemp("", "netclip-malformed-*.yml")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte(malformedContent))
	assert.NoError(t, err)
	err = tmpfile.Close()
	assert.NoError(t, err)

	// Try to load malformed config
	_, err = netclip.LoadConfig(tmpfile.Name())
	assert.Error(t, err)
}

// ApplyFlags tests
func TestApplyFlagsPortOverridesConfig(t *testing.T) {
	config := netclip.Config{Port: "4000"}
	result := netclip.ApplyFlags(config, "8080", "", "", "", false, false)
	assert.Equal(t, "8080", result.Port)
}

func TestApplyFlagsPortDefault(t *testing.T) {
	config := netclip.Config{}
	result := netclip.ApplyFlags(config, "", "", "", "", false, false)
	assert.Equal(t, "9999", result.Port)
}

func TestApplyFlagsCertKeyOverrideConfig(t *testing.T) {
	config := netclip.Config{
		CertFile: "config.crt",
		KeyFile:  "config.key",
	}
	result := netclip.ApplyFlags(config, "", "flag.crt", "flag.key", "", false, false)
	assert.Equal(t, "flag.crt", result.CertFile)
	assert.Equal(t, "flag.key", result.KeyFile)
}

func TestApplyFlagsTailscaleEnablesOverConfig(t *testing.T) {
	config := netclip.Config{
		Tailscale: netclip.TailscaleConfig{Enabled: false},
	}
	result := netclip.ApplyFlags(config, "", "", "", "", true, false)
	assert.True(t, result.Tailscale.Enabled)
}

func TestApplyFlagsTailscaleHostnameOverridesConfig(t *testing.T) {
	config := netclip.Config{
		Tailscale: netclip.TailscaleConfig{
			Enabled:  true,
			Hostname: "config-host",
		},
	}
	result := netclip.ApplyFlags(config, "", "", "", "flag-host", false, false)
	assert.Equal(t, "flag-host", result.Tailscale.Hostname)
}

func TestApplyFlagsTailscaleHostnameDefaultWhenEnabled(t *testing.T) {
	config := netclip.Config{}
	result := netclip.ApplyFlags(config, "", "", "", "", true, false)
	assert.Equal(t, "netclip", result.Tailscale.Hostname)
}

func TestApplyFlagsTailscaleTLSOverridesConfig(t *testing.T) {
	config := netclip.Config{
		Tailscale: netclip.TailscaleConfig{UseTLS: false},
	}
	result := netclip.ApplyFlags(config, "", "", "", "", false, true)
	assert.True(t, result.Tailscale.UseTLS)
}

func TestApplyFlagsConfigPreservedWhenNoFlags(t *testing.T) {
	config := netclip.Config{
		Port:     "4000",
		CertFile: "test.crt",
		KeyFile:  "test.key",
		Tailscale: netclip.TailscaleConfig{
			Enabled:  true,
			Hostname: "test-host",
			UseTLS:   true,
		},
	}
	result := netclip.ApplyFlags(config, "", "", "", "", false, false)
	assert.Equal(t, "4000", result.Port)
	assert.Equal(t, "test.crt", result.CertFile)
	assert.Equal(t, "test.key", result.KeyFile)
	assert.True(t, result.Tailscale.Enabled)
	assert.Equal(t, "test-host", result.Tailscale.Hostname)
	assert.True(t, result.Tailscale.UseTLS)
}

// GetConfigPaths tests
func TestGetConfigPaths(t *testing.T) {
	paths := netclip.GetConfigPaths()
	
	// Should have multiple paths
	assert.Greater(t, len(paths), 2)
	
	// All paths should end with netclip.yml
	for _, path := range paths {
		assert.True(t, strings.HasSuffix(path, "netclip.yml"))
	}
	
	// Should include current working directory
	cwd, _ := os.Getwd()
	expectedCwdPath := filepath.Join(cwd, "netclip.yml")
	assert.Contains(t, paths, expectedCwdPath)
	
	// Should include home directory paths
	if homeDir, err := os.UserHomeDir(); err == nil {
		expectedHomePath := filepath.Join(homeDir, ".netclip.yml")
		assert.Contains(t, paths, expectedHomePath)
	}
	
	// Platform-specific path checks
	switch runtime.GOOS {
	case "linux":
		assert.Contains(t, paths, "/etc/netclip/netclip.yml")
		assert.Contains(t, paths, "/usr/local/etc/netclip/netclip.yml")
		assert.Contains(t, paths, "/opt/netclip/netclip.yml")
	case "darwin":
		assert.Contains(t, paths, "/etc/netclip/netclip.yml")
		assert.Contains(t, paths, "/usr/local/etc/netclip/netclip.yml")
		assert.Contains(t, paths, "/Library/Application Support/netclip/netclip.yml")
	case "windows":
		// Check for environment variable paths
		if programData := os.Getenv("PROGRAMDATA"); programData != "" {
			expectedPath := filepath.Join(programData, "netclip", "netclip.yml")
			assert.Contains(t, paths, expectedPath)
		}
	}
}

func TestGetConfigPathsXDGConfigHome(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG_CONFIG_HOME not applicable on Windows")
	}
	
	// Save original value
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		if originalXDG == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
		} else {
			os.Setenv("XDG_CONFIG_HOME", originalXDG)
		}
	}()
	
	// Test with custom XDG_CONFIG_HOME
	testConfigDir := "/tmp/test-xdg"
	os.Setenv("XDG_CONFIG_HOME", testConfigDir)
	
	paths := netclip.GetConfigPaths()
	expectedPath := filepath.Join(testConfigDir, "netclip", "netclip.yml")
	assert.Contains(t, paths, expectedPath)
}

