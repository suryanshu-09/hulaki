package cmd

import (
	"errors"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/styles"
	"github.com/suryanshu-09/hulaki/utils"
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Perform an HTTP PATCH request",
	Long: `The 'patch' command allows you to perform an HTTP PATCH request to a specified URL.
You can optionally include query parameters and headers to customize the request.

This command is commonly used to partially update resources on a server, such as modifying specific fields of a user profile.`,
	Example: `Examples:
1. Basic PATCH request:
   hulaki http patch https://example.com/resource/123

2. PATCH request with query parameters:
   hulaki http patch https://example.com/resource --params=id=123,type=user

3. PATCH request with custom headers:
   hulaki http patch https://example.com/resource --headers=Authorization=BearerToken,Content-Type=application/json`,
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
		resp, err := utils.HTTPPatch(url[0], utils.WithHeaders(headers), utils.WithParams(params))
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
	httpCmd.AddCommand(patchCmd)

	patchCmd.Flags().String("headers", "", "Specify custom headers for the request")
	patchCmd.Flags().StringP("params", "p", "", "Specify query parameters for the request")
	patchCmd.Flags().BoolP("less", "l", false, "Display only the response body")
}
