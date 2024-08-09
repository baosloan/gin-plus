package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "success",
				"data": "pong",
			})
		})
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
