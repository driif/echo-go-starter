package cmd

import (
	"fmt"
	"os"

	"github.com/driif/echo-go-starter/internal/server/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version: config.GetFormattedBuildArgs(),
	Use:     "app",
	Short:   config.ModuleName,
	Long: fmt.Sprintf(`%v
A stateless RESTful JSON service written in Go.
Requires configuration through environment variables.`, config.ModuleName),
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init sets the version template for the root command.
func init() {
	rootCmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)
}
