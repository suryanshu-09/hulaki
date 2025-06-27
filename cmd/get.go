/*
Package cmd
cli for hulaki
*/
package cmd

import (
	"errors"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/styles"
	"github.com/suryanshu-09/hulaki/utils"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Send an HTTP GET request to retrieve data",
	Long: `The 'get' command allows you to send an HTTP GET request to a specified URL.
This command is commonly used to fetch data from a server, such as retrieving a webpage or an API response.
You can customize the request by including query parameters and headers.`,
	Example: `Examples:
1. Basic GET request:
   hulaki http get https://example.com

2. GET request with query parameters:
   hulaki http get https://api.example.com/data --params=type=user,status=active

3. GET request with custom headers:
   hulaki http get https://api.example.com/data --headers=Authorization=BearerToken,Accept=application/json`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		p, err := cmd.Flags().GetString("params")
		params := make(map[string]string, 0)
		if err == nil {
			if i := strings.Index(p, ","); i != -1 {
				paramsArr := strings.SplitSeq(p, ",")
				for param := range paramsArr {
					if i := strings.Index(param, "="); i != -1 {
						temp := strings.Split(param, "=")
						params[temp[0]] = temp[1]
					}
				}
			}
		}

		h, _ := cmd.Flags().GetString("headers")
		headers := make(map[string]string, 0)
		if err == nil {
			if i := strings.Index(h, ","); i != -1 {
				headersArr := strings.SplitSeq(h, ",")
				for header := range headersArr {
					if i := strings.Index(header, "="); i != -1 {
						temp := strings.Split(header, "=")
						headers[temp[0]] = temp[1]
					}
				}
			}
		}

		url := cmd.Flags().Args()
		if len(url) < 1 {
			return errors.New("please provide a url")
		}
		resp, err := utils.HTTPGet(url[0], utils.WithHeaders(headers), utils.WithParams(params))
		if err != nil {
			return err
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// Output
		less, _ := cmd.Flags().GetBool("less")
		if !less {
			cmd.Println(styles.Heading.Render("HEADERS"))
			for key, values := range resp.Header {
				for _, value := range values {
					cmd.Printf("%s: %s\n", styles.Key.Render(key), value)
				}
			}

			cmd.Println(styles.Heading.Render("BODY"))
			cmd.Println(styles.Content.Render(string(body)))
			return nil
		}
		cmd.Println(string(body))
		return nil
	},
}

func init() {
	httpCmd.AddCommand(getCmd)

	getCmd.Flags().String("headers", "", "Specify custom headers for the HTTP request")
	getCmd.Flags().StringP("params", "p", "", "Specify query parameters for the HTTP request")
	getCmd.Flags().BoolP("less", "l", false, "Display only the response body in the output")
}
