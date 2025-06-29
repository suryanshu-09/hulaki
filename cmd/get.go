/*
Package cmd
cli for hulaki
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/utils"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Perform an HTTP GET request",
	Long: `The 'get' command sends an HTTP GET request to a specified URL. 
It is commonly used to retrieve data from a server, such as fetching a webpage or an API response. 
You can include query parameters and headers to customize the request.`,
	Example: `Examples:
1. Perform a basic GET request:
   hulaki http get https://example.com

2. Perform a GET request with query parameters:
   hulaki http get https://api.example.com/data --params=type=user,status=active

3. Perform a GET request with custom headers:
   hulaki http get https://api.example.com/data --headers=Authorization=BearerToken,Accept=application/json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		iBody, params, headers, err := HTTPIn(cmd, args)
		if err != nil {
			return err
		}
		url := args[0]
		resp, err := utils.HTTPGet(url, utils.WithHeaders(headers), utils.WithParams(params), utils.WithBody(&iBody))
		if err != nil {
			return err
		}
		return HTTPOut(cmd, resp)
	},
}

func init() {
	httpCmd.AddCommand(getCmd)

	getCmd.Flags().String("headers", "", "Custom headers for the HTTP request, formatted as key=value pairs separated by commas")
	getCmd.Flags().StringP("params", "p", "", "Query parameters for the HTTP request, formatted as key=value pairs separated by commas")
	getCmd.Flags().StringP("body", "b", "", "Request body for the HTTP request, formatted as key=value pairs separated by commas")
	getCmd.Flags().BoolP("less", "l", false, "Show only the response body, omitting headers in the output")
}
