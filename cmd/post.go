package cmd

import (
	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/utils"
)

// postCmd represents the post command
var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Send an HTTP POST request to a specified URL",
	Long: `The 'post' command sends an HTTP POST request to a specified URL. 
This command is typically used to submit data to a server, such as form data or file uploads. 
You can customize the request by including headers, query parameters, and a request body.`,
	Example: `Examples:
1. Send a basic POST request:
   hulaki http post https://example.com/resource

2. Send a POST request with a JSON body:
   hulaki http post https://example.com/resource --body=name=John,age=30

3. Send a POST request with custom headers:
   hulaki http post https://example.com/resource --headers=Authorization=BearerToken,Content-Type=application/json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		iBody, params, headers, err := HTTPIn(cmd, args)
		if err != nil {
			return err
		}
		url := args[0]
		resp, err := utils.HTTPPost(url, utils.WithHeaders(headers), utils.WithParams(params), utils.WithBody(&iBody))
		if err != nil {
			return err
		}
		return HTTPOut(cmd, resp)
	},
}

func init() {
	httpCmd.AddCommand(postCmd)

	postCmd.Flags().String("headers", "", "Custom headers for the HTTP request, formatted as key=value pairs separated by commas")
	postCmd.Flags().StringP("params", "p", "", "Query parameters for the HTTP request, formatted as key=value pairs separated by commas")
	postCmd.Flags().StringP("body", "b", "", "Request body for the HTTP POST request, formatted as key=value pairs separated by commas")
	postCmd.Flags().BoolP("less", "l", false, "Show only the response body, omitting headers in the output")
}
