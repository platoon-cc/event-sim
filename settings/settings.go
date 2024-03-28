package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	koanfJson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	konfig       *koanf.Koanf
	settingsFile string
)

func Config() {
	dir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("error getting config folder: %v\n", err)
		os.Exit(1)
	}

	path := filepath.Join(dir, "platoon")
	os.MkdirAll(path, 0755)
	settingsFile = filepath.Join(path, "settings.json")

	konfig = koanf.New(".")
	if err := konfig.Load(file.Provider(settingsFile), koanfJson.Parser()); err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("error in config file (%s): %v\n", settingsFile, err)
		os.Exit(1)
	}

	// if err := konfig.Load(env.Provider("PL_", ".", func(s string) string {
	// 	return strings.ReplaceAll(strings.ToLower(
	// 		strings.TrimPrefix(s, "PL_")), "_", ".")
	// }), nil); err != nil {
	// 	panic(err)
	// }

	// flag := cmd.Flags().Lookup("config")
	//
	//
	// if err := c.konfig.Load(posflag.ProviderWithValue(cmd.Flags(), ".", c.konfig, func(key, value string) (string, interface{}) {
	// 	return strings.ReplaceAll(key, "_", "."), value
	// }), nil); err != nil {
	// 	return err
	// }
}

func Save() {
	// b, _ := konfig.Marshal(koanfJson.Parser())

	b, _ := json.MarshalIndent(konfig.Raw(), "", "\t")
	os.WriteFile(settingsFile, b, 0755)
}

func GetAuthToken() string {
	return konfig.String("auth.token")
}

func SetAuthToken(token string) {
	konfig.Set("auth.token", token)
}
