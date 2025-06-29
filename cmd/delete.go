package cmd

import (
	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/utils"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Send an HTTP DELETE request to remove a resource",
	Long: `The 'delete' command sends an HTTP DELETE request to a specified URL. 
You can include query parameters, headers, and a request body to customize the request.

This command is commonly used to delete resources on a server, such as removing a user or deleting a file.`,
	Example: `Examples:
1. Basic DELETE request:
   hulaki http delete https://example.com/resource/123

2. DELETE request with query parameters:
   hulaki http delete https://example.com/resource --params id=123,type=user

3. DELETE request with custom headers:
   hulaki http delete https://example.com/resource --headers Authorization=BearerToken

4. DELETE request with a body:
   hulaki http delete https://example.com/resource --body key=value`,
	RunE: func(cmd *cobra.Command, args []string) error {
		iBody, params, headers, err := HTTPIn(cmd, args)
		if err != nil {
			return err
		}
		url := args[0]
		resp, err := utils.HTTPDelete(url, utils.WithHeaders(headers), utils.WithParams(params), utils.WithBody(&iBody))
		if err != nil {
			return err
		}
		return HTTPOut(cmd, resp)
	},
}

func init() {
	httpCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().String("headers", "", "Custom headers for the HTTP request, formatted as key=value pairs separated by commas")
	deleteCmd.Flags().StringP("params", "p", "", "Query parameters for the HTTP request, formatted as key=value pairs separated by commas")
	deleteCmd.Flags().StringP("body", "b", "", "Request body for the HTTP request, formatted as key=value pairs separated by commas")
	deleteCmd.Flags().BoolP("less", "l", false, "Show only the response body, omitting headers")
}
