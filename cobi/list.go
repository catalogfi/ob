package cobi

import "github.com/spf13/cobra"

func List() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "list",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
		DisableAutoGenTag: true,
	}
	return cmd
}
