package cmd

import (
	"fmt"
	"os"

	"github.com/platoon-cc/platoon-cli/settings"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "platoon-cli",
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("Hello world")
	// },
}

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Platoon CLI v0.0.1 -- HEAD")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(settings.Config)
	cobra.OnFinalize(settings.Save)
	rootCmd.AddCommand(versionCmd)
}

// var (
// 	konfig     *koanf.Koanf
// 	configFile = "config.json"
// )
//
// func saveConfig() {
// 	b, _ := konfig.Marshal(json.Parser())
// 	os.WriteFile(configFile, b, 0755)
// }
//
// func initConfig() {
// 	konfig = koanf.New(".")
// 	if err := konfig.Load(file.Provider(configFile), json.Parser()); err != nil {
// 		fmt.Printf("error in config file (%s): %v\n", configFile, err)
// 	}
//
// 	konfig.Set("test.value", 123)
//
// 	// if err := konfig.Load(env.Provider("PL_", ".", func(s string) string {
// 	// 	return strings.ReplaceAll(strings.ToLower(
// 	// 		strings.TrimPrefix(s, "PL_")), "_", ".")
// 	// }), nil); err != nil {
// 	// 	panic(err)
// 	// }
//
// 	// flag := cmd.Flags().Lookup("config")
// 	//
// 	//
// 	// if err := c.konfig.Load(posflag.ProviderWithValue(cmd.Flags(), ".", c.konfig, func(key, value string) (string, interface{}) {
// 	// 	return strings.ReplaceAll(key, "_", "."), value
// 	// }), nil); err != nil {
// 	// 	return err
// 	// }
// }
