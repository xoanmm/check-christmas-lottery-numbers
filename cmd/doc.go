/*
NAME:
   check-christmas-lottery-numbers - Interacts with the Christmas lottery API to obtain the results of the numbers provided and send notifications with PushOver if you win a prize

USAGE:
   check-lottery-results --file-numbers-to-check <file-numbers-to-check> [-n]

AUTHOR:
   Xoan Mallon <xoanmallon@gmail.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file-numbers-to-check value, -f value  JSON file with numbers to check (default: "/tmp/numbers_to_check.json")
   --draw value, -d value                   Draw for which the numbers will be checked
   --notify, -n                             Activate notifications through PushOver (default: false)
   --store-notifications, -s                Activate notifications store in mongodb (default: false)
   --mongoHost value, -m value              Mongo Host URL used to store notifications (default: "localhost:27017")
   --mongoRootUsername value, -u value      Root username for mongo host
   --mongoRootPassword value, -p value      Root password for mongo host
   --help, -h                               show help (default: false)
*/
package main
