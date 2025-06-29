package cmd

import (
	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/utils"
)

// optionsCmd represents the options command
var optionsCmd = &cobra.Command{
	Use:   "options",
	Short: "Send an HTTP OPTIONS request to discover server capabilities",
	Long: `The 'options' command sends an HTTP OPTIONS request to a specified URL.
It is used to discover the HTTP methods and other options supported by a server for a specific resource.

You can include query parameters and headers to customize the request.`,
	Example: `Examples:
1. Basic OPTIONS request:
   hulaki http options https://example.com/resource

2. OPTIONS request with query parameters:
   hulaki http options https://example.com/resource --params=type=user,status=active

3. OPTIONS request with custom headers:
   hulaki http options https://example.com/resource --headers=Authorization=BearerToken,Accept=application/json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		iBody, params, headers, err := HTTPIn(cmd, args)
		if err != nil {
			return err
		}
		url := args[0]
		resp, err := utils.HTTPOptions(url, utils.WithHeaders(headers), utils.WithParams(params), utils.WithBody(&iBody))
		if err != nil {
			return err
		}
		return HTTPOut(cmd, resp)
	},
}

func init() {
	httpCmd.AddCommand(optionsCmd)

	optionsCmd.Flags().String("headers", "", "Custom headers for the HTTP request, formatted as key=value pairs separated by commas")
	optionsCmd.Flags().StringP("params", "p", "", "Query parameters for the HTTP request, formatted as key=value pairs separated by commas")
	optionsCmd.Flags().StringP("body", "b", "", "Request body for the HTTP request, formatted as key=value pairs separated by commas")
	optionsCmd.Flags().BoolP("less", "l", false, "Show only the response body, omitting headers")
}
