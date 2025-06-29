package cmd

import (
	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/utils"
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Send an HTTP PATCH request to update resources",
	Long: `The 'patch' command sends an HTTP PATCH request to a specified URL. 
This command is typically used to partially update resources on a server, such as modifying specific fields of a user profile. 
You can include query parameters, headers, and a request body to customize the request.`,
	Example: `Examples:
1. Send a basic PATCH request:
   hulaki http patch https://example.com/resource/123

2. Include query parameters in the PATCH request:
   hulaki http patch https://example.com/resource --params="id=123,type=user"

3. Add custom headers to the PATCH request:
   hulaki http patch https://example.com/resource --headers="Authorization=BearerToken,Content-Type=application/json"

4. Include a request body in the PATCH request:
   hulaki http patch https://example.com/resource --body="field1=value1,field2=value2"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		iBody, params, headers, err := HTTPIn(cmd, args)
		if err != nil {
			return err
		}
		url := args[0]
		resp, err := utils.HTTPPatch(url, utils.WithHeaders(headers), utils.WithParams(params), utils.WithBody(&iBody))
		if err != nil {
			return err
		}
		return HTTPOut(cmd, resp)
	},
}

func init() {
	httpCmd.AddCommand(patchCmd)

	patchCmd.Flags().String("headers", "", "Custom headers for the HTTP request, formatted as key=value pairs separated by commas")
	patchCmd.Flags().StringP("params", "p", "", "Query parameters for the HTTP request, formatted as key=value pairs separated by commas")
	patchCmd.Flags().StringP("body", "b", "", "Request body for the HTTP PATCH request, formatted as key=value pairs separated by commas")
	patchCmd.Flags().BoolP("less", "l", false, "Show only the response body, omitting headers")
}
