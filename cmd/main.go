package main

import (
	"github.com/urfave/cli/v2"
	"github.com/xoanmm/check-christmas-lottery-numbers/pkg/results"
	"log"
	"os"
	"path/filepath"
	"time"
)

var date = time.Now().Format(time.RFC3339)
var year = time.Now().Year()

func main() {
	cmd := buildCLI()
	if err := cmd.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getCLINecessaryArgs(c *cli.Context) (string, bool, bool, string, string, string, string) {
	fileNumbersToCheck, _ := filepath.Abs(c.String("file-numbers-to-check"))
	notify := c.Bool("notify")
	storeNotifications := c.Bool("store-notifications")
	draw := c.String("draw")
	mongoHostURL := c.String("mongoHost")
	mongoRootUsername := c.String("mongoRootUsername")
	mongoRootPassword := c.String("mongoRootPassword")

	return fileNumbersToCheck, notify, storeNotifications, draw, mongoHostURL, mongoRootUsername, mongoRootPassword
}

// buildCLI creates a CLI app
func buildCLI() *cli.App {
	d, _ := time.Parse(time.RFC3339, date)

	return &cli.App{
		Name: "check-christmas-lottery-numbers",
		Usage: "Interacts with the Christmas lottery API to obtain the results " +
			"of the numbers provided and send notifications with PushOver if you win a prize",
		Compiled:  d,
		UsageText: "check-lottery-results --file-numbers-to-check <file-numbers-to-check> [-n]",
		Authors: []*cli.Author{
			{
				Name:  "Xoan Mallon",
				Email: "xoanmallon@gmail.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file-numbers-to-check",
				Usage:   "JSON file with numbers to check",
				Value:   "/tmp/numbers_to_check.json",
				Aliases: []string{"f"},
			},
			&cli.StringFlag{
				Name:    "draw",
				Usage:   "Draw for which the numbers will be checked",
				Value:   "",
				Aliases: []string{"d"},
			},
			&cli.BoolFlag{
				Name:    "notify",
				Usage:   "Activate notifications through PushOver",
				Value:   false,
				Aliases: []string{"n"},
			},
			&cli.BoolFlag{
				Name:    "store-notifications",
				Usage:   "Activate notifications store in mongodb",
				Value:   false,
				Aliases: []string{"s"},
			},
			&cli.StringFlag{
				Name:    "mongoHost",
				Usage:   "Mongo Host URL used to store notifications",
				Value:   "localhost:27017",
				Aliases: []string{"m"},
			},
			&cli.StringFlag{
				Name:    "mongoRootUsername",
				Usage:   "Root username for mongo host",
				Value:   "",
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    "mongoRootPassword",
				Usage:   "Root password for mongo host",
				Value:   "",
				Aliases: []string{"p"},
			},
		},
		Action: func(c *cli.Context) error {
			fileNumbersToCheck, notify, storeNotifications, draw, mongoHostURL, mongoRootUsername, mongoRootPassword := getCLINecessaryArgs(c)
			lotteryDrawAPIURL, err := results.GetDrawAPIURLToCheckNumbers(draw)
			if err != nil {
				return err
			}
			lotteryDrawStatus, err := results.GetAPILotteryDrawStatus(draw, *lotteryDrawAPIURL)
			if err != nil {
				return err
			}
			message, err := results.CheckAPILOtteryDrawStatus(lotteryDrawStatus.Status)
			if err != nil {
				return err
			}
			log.Println(*message)
			personNumbers, err := results.OpenFilePersonsNumbersToCheck(fileNumbersToCheck)
			if err != nil {
				return err
			}

			err = results.CheckPersonsNumbers(*lotteryDrawAPIURL, personNumbers, draw, notify, storeNotifications, mongoHostURL, mongoRootUsername, mongoRootPassword, year)
			if err != nil {
				return err
			}
			return err
		},
	}
}
