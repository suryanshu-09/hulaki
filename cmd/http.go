package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/styles"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Make an HTTP request",
	Long:  `Make an HTTP request`,
}

func init() {
	rootCmd.AddCommand(httpCmd)
}

func HTTPIn(cmd *cobra.Command, args []string) (iBody bytes.Buffer, params map[string]string, headers map[string]string, err error) {
	b, err := cmd.Flags().GetString("body")
	jBody := make(map[string]string, 0)
	if err == nil {
		if b == "-" {
			io.Copy(&iBody, os.Stdin)
		} else {
			if i := strings.Index(b, ","); i != -1 {
				bodyArr := strings.SplitSeq(b, ",")
				for bo := range bodyArr {
					if i := strings.Index(bo, "="); i != -1 {
						temp := strings.Split(bo, "=")
						jBody[temp[0]] = temp[1]
					}
				}
			} else {
				if i := strings.Index(b, "="); i != -1 {
					temp := strings.Split(b, "=")
					jBody[temp[0]] = temp[1]
				}
			}
			if len(b) > 0 {
				marshalled, err := json.Marshal(jBody)
				if err != nil {
					return iBody, nil, nil, err
				}
				iBody = *bytes.NewBuffer(marshalled)
			}
		}
	}

	p, err := cmd.Flags().GetString("params")
	params = make(map[string]string, 0)
	if err == nil {
		if i := strings.Index(p, ","); i != -1 {
			paramsArr := strings.SplitSeq(p, ",")
			for param := range paramsArr {
				if i := strings.Index(param, "="); i != -1 {
					temp := strings.Split(param, "=")
					params[temp[0]] = temp[1]
				}
			}
		} else {
			if i := strings.Index(p, "="); i != -1 {
				temp := strings.Split(p, "=")
				params[temp[0]] = temp[1]
			}
		}
	}

	h, _ := cmd.Flags().GetString("headers")
	headers = make(map[string]string, 0)
	if err == nil {
		if i := strings.Index(h, ","); i != -1 {
			headersArr := strings.SplitSeq(h, ",")
			for header := range headersArr {
				if i := strings.Index(header, "="); i != -1 {
					temp := strings.Split(header, "=")
					headers[temp[0]] = temp[1]
				}
			}
		} else {
			if i := strings.Index(h, "="); i != -1 {
				temp := strings.Split(h, "=")
				headers[temp[0]] = temp[1]
			}
		}
	}

	if len(args) < 1 {
		return iBody, nil, nil, errors.New("please provide a url")
	}
	return
}

func HTTPOut(cmd *cobra.Command, resp *http.Response) error {
	body := bytes.Buffer{}
	_, err := io.Copy(&body, resp.Body)
	if err != nil {
		return err
	}

	less, _ := cmd.Flags().GetBool("less")
	if !less {
		fmt.Fprintf(os.Stdout, "%s\n", styles.Heading.Render("HEADERS"))
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Fprintf(os.Stdout, "%s: %s\n", styles.Key.Render(key), value)
			}
		}

		fmt.Fprintf(os.Stdout, "%s\n", styles.Heading.Render("BODY"))
		fmt.Fprintf(os.Stdout, "%s\n", styles.Content.Render(body.String()))
		return nil
	}
	io.Copy(os.Stdout, bytes.NewBuffer(body.Bytes()))

	return nil
}
