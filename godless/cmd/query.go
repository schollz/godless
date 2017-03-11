// Copyright © 2017 Johnny Morrice <john@functorama.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/johnny-morrice/godless"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query a godless server",
	Long: `Send a query to a godless server over HTTP.`,
	Run: func(cmd *cobra.Command, args []string) {
		var query *godless.Query

		if source != "" {
			q, err := godless.CompileQuery(source)

			if err != nil {
				die(err)
			}

			query = q
		}

		if analyse {
			fmt.Printf("Query analysis:\n\n%s\n\n", query.Analyse())
			fmt.Println("Syntax tree:\n\n")
			query.Parser.PrintSyntaxTree()
		}
	},
}

var source string
var analyse bool

func init() {
	RootCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringVar(&source, "query", "", "Godless NOSQL query text")
	queryCmd.Flags().BoolVar(&analyse, "analyse", false, "Analyse query")
}