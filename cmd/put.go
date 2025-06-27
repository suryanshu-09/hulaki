package cmd

import (
	"errors"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/styles"
	"github.com/suryanshu-09/hulaki/utils"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Perform an HTTP PUT request",
	Long: `The 'put' command allows you to perform an HTTP PUT request to a specified URL.
You can optionally include query parameters and headers to customize the request.

This command is commonly used to update resources on a server, such as modifying user details or updating a file.`,
	Example: `Examples:
1. Basic PUT request:
   hulaki http put https://example.com/resource/123

2. PUT request with query parameters:
   hulaki http put https://example.com/resource --params=id=123,type=user

3. PUT request with custom headers:
   hulaki http put https://example.com/resource --headers=Authorization=BearerToken,Content-Type=application/json`,
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
		resp, err := utils.HTTPPut(url[0], utils.WithHeaders(headers), utils.WithParams(params))
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
	httpCmd.AddCommand(putCmd)

	putCmd.Flags().String("headers", "", "Specify custom headers for the HTTP request")
	putCmd.Flags().StringP("params", "p", "", "Specify query parameters for the HTTP request")
	putCmd.Flags().BoolP("less", "l", false, "Display only the response body in the output")
}
