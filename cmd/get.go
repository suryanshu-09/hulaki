/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/utils"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("get called")
		// url, err := cmd.Flags().GetString("url")
		// if err != nil {
		// 	log.Fatal("Need a url")
		// }
		url := cmd.Flags().Args()

		resp, err := utils.HTTPGet(url[0])
		if err != nil {
			log.Fatal(err.Error())
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Print(string(body))
	},
}

func init() {
	httpCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")
	// getCmd.PersistentFlags().String("url", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
