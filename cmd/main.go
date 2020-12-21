package main

import (
	"github.com/urfave/cli/v2"
	"github.com/xoanmm/check-christmas-lottery-numbers/pkg/results"
	"log"
	"os"
	"path/filepath"
	"time"
)

// lotteryDrawResultsAPIURL represents the URL of the API where the results of the
// Christmas draw will be obtained
const lotteryDrawResultsAPIURL = "http://api.elpais.com/ws/LoteriaNavidadPremiados"

var date = time.Now().Format(time.RFC3339)

func main() {
	cmd := buildCLI()
	if err := cmd.Run(os.Args); err != nil {
		log.Fatal(err)
	}
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
			&cli.BoolFlag{
				Name:    "notify",
				Usage:   "Activate notifications through PushOver",
				Value:   false,
				Aliases: []string{"n"},
			},
		},
		Action: func(c *cli.Context) error {
			fileNumbersToCheck, _ := filepath.Abs(c.String("file-numbers-to-check"))
			notify := c.Bool("notify")

			lotteryDrawStatus, err := results.GetAPILotteryDrawStatus(lotteryDrawResultsAPIURL)
			if err != nil {
				return err
			}
			message, err := results.CheckAPILOtteryDrawStatus(lotteryDrawStatus.Status)
			if err != nil {
				return err
			}
			log.Println(message)
			personNumbers, err := results.OpenFilePersonsNumbersToCheck(fileNumbersToCheck)
			if err != nil {
				return err
			}

			err = results.CheckPersonsNumbers(lotteryDrawResultsAPIURL, personNumbers, notify)
			if err != nil {
				return err
			}
			return err
		},
	}
}
