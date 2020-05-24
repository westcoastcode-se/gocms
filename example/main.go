package main

import (
	"github.com/sirupsen/logrus"
	"github.com/westcoastcode-se/gocms/pkg/cms"
	"github.com/westcoastcode-se/gocms/pkg/config"
	"log"
)

func main() {
	logrus.SetReportCaller(true)

	log.SetFlags(log.Llongfile | log.LstdFlags)
	logrus.Info("Starting public web")

	// Create the server
	config := config.GetConfig()
	public := cms.NewServer(config)

	// Configure the server
	public.ContentRepository.RegisterModelType("models.News", ConvertToNews)

	// Start the server
	err := public.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	logrus.Info("Stopped public web")
}
