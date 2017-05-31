package monitor

import (
	"github.com/labstack/echo"
	newrelic "github.com/newrelic/go-agent"
	"log"
	"os"
)

var (
	App     newrelic.Application
	runable bool
)

func init() {
	config := newrelic.NewConfig(os.Getenv("NEWRELIC_APP_NAME"), os.Getenv("NEWRELIC_LICENSE_KEY"))
	var err error
	App, err = newrelic.NewApplication(config)
	runable = true
	if err != nil {
		log.Println("Error init newrelic", err)
		log.Println("Will not report newrelic stat")
		runable = false
	}
}

func Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tx := App.StartTransaction(c.Request().URL.Path, c.Response(), c.Request())

		c.Set("newrelic", tx)

		err := next(c)

		tx.NoticeError(err)

		tx.End()
		return err
	}
}
