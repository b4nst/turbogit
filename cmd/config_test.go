package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseValue(t *testing.T) {
	// Parse bool
	v, err := parseValue("true")
	assert.NoError(t, err)
	assert.IsType(t, false, v)
	assert.Equal(t, true, v)
	v, err = parseValue("False")
	assert.NoError(t, err)
	assert.IsType(t, false, v)
	assert.Equal(t, false, v)

	// Parse int
	v, err = parseValue("113")
	assert.NoError(t, err)
	assert.IsType(t, int64(113), v)
	assert.EqualValues(t, 113, v)

	// Parse float
	v, err = parseValue("113.42")
	assert.NoError(t, err)
	assert.IsType(t, 113.42, v)
	assert.InDelta(t, 113.42, v, 1e-5)

	// Parse array
	v, err = parseValue("[foo, bar]")
	assert.NoError(t, err)
	assert.Len(t, v, 2)
	assert.Equal(t, []interface{}([]interface{}{"foo", "bar"}), v)
}

func TestConfigure(t *testing.T) {
	cmd := &cobra.Command{}
	delete := cmd.Flags().BoolP("delete", "d", false, "Delete config.")

	stdout := os.Stdout
	restore := func() {
		os.Stdout = stdout
	}
	defer restore()

	devnull, err := ioutil.TempFile("", "dev-null")
	require.NoError(t, err)
	defer os.RemoveAll(devnull.Name())
	os.Stdout = devnull

	cfg, err := ioutil.TempFile("", "*.toml")
	require.NoError(t, err)
	defer os.RemoveAll(cfg.Name())
	viper.SetConfigFile(cfg.Name())

	*delete = false
	err = configure(cmd, []string{"key", "value"})
	assert.NoError(t, err)
	assert.Equal(t, viper.Get("key"), "value")
	err = configure(cmd, []string{"key"})
	assert.NoError(t, err)

	*delete = true
	err = configure(cmd, []string{"key"})
	assert.NoError(t, err)
	assert.Equal(t, viper.Get("key"), "")
}
