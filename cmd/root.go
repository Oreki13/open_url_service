package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"open_url_service/cmd/http"
)

func Start() {

	rootCmd := &cobra.Command{}

	//migrateCmd := &cobra.Command{
	//	Use:   "db:migrate",
	//	Short: "database migration",
	//	Run: func(c *cobra.Command, args []string) {
	//		migration.MigrateDatabase()
	//	},
	//}

	//migrateCmd.Flags().BoolP("version", "", false, "print version")
	//migrateCmd.Flags().StringP("dir", "", "database/migration/", "directory with migration files")
	//migrateCmd.Flags().StringP("table", "", "db", "migrations table name")
	//migrateCmd.Flags().BoolP("verbose", "", false, "enable verbose mode")
	//migrateCmd.Flags().BoolP("guide", "", false, "print help")

	//rabbitCmd := &cobra.Command{
	//	Use:   "rabbit",
	//	Short: "Run RabbitMQ Service",
	//	Run: func(cmd *cobra.Command, args []string) {
	//		broker.ServeRabbitMQ()
	//	},
	//}

	//rabbitCmd.Flags().StringP("name", "n", "", "queue and exchange name")
	//rabbitCmd.Flags().StringP("topics", "t", "", "topic to subscribe (separate with pipeline \"|\" if want multiple binding)")
	//rabbitCmd.Flags().BoolP("guide", "", false, "Print Help")
	//
	//if err := rabbitCmd.MarkFlagRequired("name"); err != nil {
	//	log.Fatal(err)
	//}
	//
	//if err := rabbitCmd.MarkFlagRequired("topics"); err != nil {
	//	log.Fatal(err)
	//}

	cmd := []*cobra.Command{
		{
			Use:   "http",
			Short: "Run HTTP Server",
			Run: func(cmd *cobra.Command, args []string) {
				http.Start()
			},
		},
		//rabbitCmd,
		//migrateCmd,
	}

	rootCmd.AddCommand(cmd...)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
