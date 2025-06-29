package cmd

import (
	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/utils"
)

// headCmd represents the head command
var headCmd = &cobra.Command{
	Use:   "head",
	Short: "Send an HTTP HEAD request to retrieve headers",
	Long: `The 'head' command sends an HTTP HEAD request to a specified URL. 
This command retrieves only the headers of the response, without the body, making it useful for checking metadata, resource availability, or server capabilities.

You can include query parameters and headers to customize the request.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		iBody, params, headers, err := HTTPIn(cmd, args)
		if err != nil {
			return err
		}
		url := args[0]
		resp, err := utils.HTTPHead(url, utils.WithHeaders(headers), utils.WithParams(params), utils.WithBody(&iBody))
		if err != nil {
			return err
		}
		return HTTPOut(cmd, resp)
	},
}

func init() {
	httpCmd.AddCommand(headCmd)

	headCmd.Flags().String("headers", "", "Custom headers for the HTTP request, formatted as key=value pairs separated by commas")
	headCmd.Flags().StringP("params", "p", "", "Query parameters for the HTTP request, formatted as key=value pairs separated by commas")
	headCmd.Flags().StringP("body", "b", "", "Request body for the HTTP request, formatted as key=value pairs separated by commas")
	headCmd.Flags().BoolP("less", "l", false, "Display only the response headers, omitting additional formatting")
}
