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
   --file-numbers-to-check value, -f value  JSON file with numbers to check
   --notify, -n                             Activate notifications through PushOver (default: false)
   --help, -h                               show help (default: false)
 */
package main