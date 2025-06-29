package cmd

import (
	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/utils"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Send an HTTP PUT request to update resources",
	Long: `The 'put' command sends an HTTP PUT request to a specified URL. 
This is typically used to update or replace resources on a server, such as modifying user details or updating a file.

You can include query parameters, headers, and a request body to customize the request.`,
	Example: `Examples:
1. Basic PUT request:
   hulaki http put https://example.com/resource/123

2. PUT request with query parameters:
   hulaki http put https://example.com/resource --params="id=123,type=user"

3. PUT request with custom headers:
   hulaki http put https://example.com/resource --headers="Authorization=BearerToken,Content-Type=application/json"

4. PUT request with a body:
   hulaki http put https://example.com/resource --body="name=John,age=30"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		iBody, params, headers, err := HTTPIn(cmd, args)
		if err != nil {
			return err
		}
		url := args[0]
		resp, err := utils.HTTPPut(url, utils.WithHeaders(headers), utils.WithParams(params), utils.WithBody(&iBody))
		if err != nil {
			return err
		}
		return HTTPOut(cmd, resp)
	},
}

func init() {
	httpCmd.AddCommand(putCmd)

	putCmd.Flags().String("headers", "", "Custom headers for the HTTP request, formatted as key=value pairs separated by commas")
	putCmd.Flags().StringP("params", "p", "", "Query parameters for the HTTP request, formatted as key=value pairs separated by commas")
	putCmd.Flags().StringP("body", "b", "", "Request body for the HTTP request, formatted as key=value pairs separated by commas")
	putCmd.Flags().BoolP("less", "l", false, "Show only the response body, omitting headers")
}
