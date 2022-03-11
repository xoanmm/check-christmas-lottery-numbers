[![Tests][tests-badge]][tests-link]
[![GitHub Release][release-badge]][release-link]
[![Go Report Card][report-badge]][report-link]
[![License][license-badge]][license-link]
[![Coverage][coverage-badge]][coverage-link]

# check-christmas-lottery-numbers
CLI to check the lottery numbers with the possibility to send notifications for the winning numbers using [PushOver](https://pushover.net/).


This CLI has been created with the intention of providing a lottery number and through the API available to check the numbers of the draw to know if any of the numbers provided have been awarded.

It must be taken into account that the CLI will check the status of the draw and in case it has not started yet it will not carry out the check of the numbers, showing a message like the following one:
```shell
2020/12/21 18:21:25 Checking the status of the Christmas lottery draw
2020/12/21 18:21:26 the draw has not started yet, it is not possible to get the results for the numbers
```

## Installation

Go to [release page](https://github.com/xoanmm/check-christmas-lottery-numbers/releases) and download the binary you need.

### Requirements
If the user wants to send notifications, the following variables must be exported:
- `PUSH_OVER_NOTIFICATION_TOKEN`: [PushOver](https://pushover.net/) account token for sending notifications

- `PUSH_OVER_NOFITICATION_USER`: [PushOver](https://pushover.net/) Application/API token registered in account for receive notifications

- `MONGO_INITDB_ROOT_USERNAME`: MongoDB Root Username of MongoDB host where the notifications results are going to be saved.

- `MONGO_INITDB_ROOT_PASSWORD`: MongoDB Root Password of MongoDB host where the notifications results are going to be saved.

## Example of usage

A JSON file must be provided with the people whose numbers we want to check, following the format shown below:

```sh
{
    "numbers_to_check":  [
        {
            "owner": "<owner_name_1>",
            "numbers": [
                {
                    "number": "<number_to_check>",
                    "bet_amount": <bet_amount>,
                    "origin": "<number_origin>"
                },
                . . .
                {
                    "number": "<number_to_check>",
                    "bet_amount": <bet_amount>,
                    "origin": "<number_origin>"
                }
            ]
        },
        . . .
        {
            "owner": "<owner_name_N>",
            "numbers": [
                {
                    "number": "<number_to_check>",
                    "bet_amount": <bet_amount>,
                    "origin": "<number_origin>"
                },
                . . .
                {
                    "number": "<number_to_check>",
                    "bet_amount": <bet_amount>,
                    "origin": "<number_origin>"
                }
            ]
        }
    ]
}
```

The following are different examples of execution:

- Checking numbers from the `numbers_to_check.json` file without sending notifications for christmas lottery draw:
    ```shell
    $ ./check-christmas-lottery-numbers -d christmas -f /tmp/numbers_to_check.json
    2020/12/27 18:36:21 Checking the status of the christmas lottery draw
    2020/12/27 18:36:21 The draw is over and the list of numbers is official
    2020/12/27 18:36:21 Opening file /tmp/numbers_to_check.json with numbers to check
    2020/12/27 18:36:21 Successfully opened file /tmp/numbers_to_check.json  with numbers to check
    2020/12/27 18:36:21 Numbers are going to be check from 1 different owners
    2020/12/27 18:36:21 2 numbers are going to be check from owner Xoan. The probability to win a prize is 0.002%
    2020/12/27 18:36:21 Checking prize obtained for number 26967 with 5€ bet and origin My favourite restaurant
    2020/12/27 18:36:22 Xoan won 5.00€ in number 26967 with 5 € bet and origin My favourite restaurant
    2020/12/27 18:36:22 Checking prize obtained for number 33972 with 20€ bet and origin My company
    2020/12/27 18:36:22 Xoan won 100.00€ in number 33972 with 20 € bet and origin My company
    ```

- Checking numbers from the `numbers_to_check.json` file with sending notifications activated and without use mongodb to store notification for christmas lottery draw:
    ```shell
    $ ./check-christmas-lottery-numbers -d="christmas" -f /tmp/numbers_to_check.json -n
    2020/12/27 18:37:02 Checking the status of the christmas lottery draw
    2020/12/27 18:37:02 The draw is over and the list of numbers is official
    2020/12/27 18:37:02 Opening file /tmp/numbers_to_check.json with numbers to check
    2020/12/27 18:37:02 Successfully opened file /tmp/numbers_to_check.json  with numbers to check
    2020/12/27 18:37:02 Numbers are going to be check from 1 different owners
    2020/12/27 18:37:02 2 numbers are going to be check from owner Xoan. The probability to win a prize is 0.002%
    2020/12/27 18:37:02 Checking prize obtained for number 26967 with 5€ bet and origin My favourite restaurant
    2020/12/27 18:37:03 Xoan won 5.00€ in number 26967 with 5 € bet and origin My favourite restaurant
    2020/12/27 18:37:03 A notification is going to be send with pushOver
    2020/12/27 18:37:03 Notification send correctly with id a600ed98-a640-4317-8c04-631fcdfdbccf to PushOver App
    2020/12/27 18:37:03 Checking prize obtained for number 33972 with 20€ bet and origin My company
    2020/12/27 18:37:03 Xoan won 100.00€ in number 33972 with 20 € bet and origin My company
    2020/12/27 18:37:03 A notification is going to be send with pushOver
    2020/12/27 18:37:04 Notification send correctly with id 8205a8e5-e11c-4741-8cac-51debc6fe201 to PushOver App
    ```

- Checking numbers from the `numbers_to_check.json` file with sending notifications activated and using mongodb for store notifications for christmas lottery draw:
    ```shell
    $ ./check-christmas-lottery-numbers -d="christmas" -m="localhost:27017" -f /tmp/numbers_to_check.json -n -s -u=<MONGO_INITDB_ROOT_USERNAME> -p=<MONGO_INITDB_ROOT_PASSWORD>
    2020/12/27 18:43:34 Checking the status of the christmas lottery draw
    2020/12/27 18:43:35 The draw is over and the list of numbers is official
    2020/12/27 18:43:35 Opening file /tmp/numbers_to_check.json with numbers to check
    2020/12/27 18:43:35 Successfully opened file /tmp/numbers_to_check.json  with numbers to check
    2020/12/27 18:43:35 Numbers are going to be check from 1 different owners
    2020/12/27 18:43:35 2 numbers are going to be check from owner Xoan. The probability to win a prize is 0.002%
    2020/12/27 18:43:35 Checking prize obtained for number 26967 with 5€ bet and origin My favourite restaurant
    2020/12/27 18:43:35 Xoan won 5.00€ in number 26967 with 5 € bet and origin My favourite restaurant
    2020/12/27 18:43:35 Notification is not inserted in collection christmas yet
    2020/12/27 18:43:35 A notification is going to be send with pushOver
    2020/12/27 18:43:36 Notification send correctly with id 13e58739-ea10-4c04-aba1-3c78a8acfa21 to PushOver App
    2020/12/27 18:43:36 Notification inserted sucessfully in mongodb collection with id ObjectID("5fe8c7c873fa0e4c8595e1ef")
    2020/12/27 18:43:36 Checking prize obtained for number 33972 with 20€ bet and origin My company
    2020/12/27 18:43:36 Xoan won 100.00€ in number 33972 with 20 € bet and origin My company
    2020/12/27 18:43:36 Notification is not inserted in collection christmas yet
    2020/12/27 18:43:36 A notification is going to be send with pushOver
    2020/12/27 18:43:37 Notification send correctly with id 4332ea36-cce6-4831-8d3f-af40226366fd to PushOver App
    2020/12/27 18:43:37 Notification inserted sucessfully in mongodb collection with id ObjectID("5fe8c7c973fa0e4c8595e1f0")
    ```

- Checking numbers from the `numbers_to_check.json` file with notifications activated for childhoods lottery draw:
   ```shell
    $ ./check-christmas-lottery-numbers -d="childhoods" -f /tmp/numbers_to_check.json -n
    2020/12/27 17:48:10 Checking the status of the childhoods lottery draw
    2020/12/27 17:48:11 The draw is over and the list of numbers is official
    2020/12/27 17:48:11 Opening file /tmp/numbers_to_check.json with numbers to check
    2020/12/27 17:48:11 Successfully opened file /tmp/numbers_to_check.json  with numbers to check
    2020/12/27 17:48:11 Numbers are going to be check from 1 different owners
    2020/12/27 17:48:11 2 numbers are going to be check from owner Xoan. The probability to win a prize is 0.002%
    2020/12/27 17:48:11 Checking prize obtained for number 26967 with 5€ bet and origin My favourite restaurant
    2020/12/27 17:48:12 Xoan won 0.00€ in number 26967 with 5 € bet and origin My favourite restaurant
    2020/12/27 17:48:12 Checking prize obtained for number 33972 with 20€ bet and origin My company
    2020/12/27 17:48:13 Xoan won 20.00€ in number 33972 with 20 € bet and origin My company
    2020/12/27 17:48:13 A notification is going to be send with pushOver
    2020/12/27 17:48:14 Notification send correctly with id 59389f39-28ba-4f90-8edb-86e4f66432c9 to PushOver App
    ```

## Usage

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

## Helm chart

A [helm chart](./helm/check-christmas-lottery-numbers) has been created in order to run the CLI as a [cronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) in Kubernetes.

It is recommended to create a unique namespace to use the chart, this can be done through the following command:
```shell
kubectl create namespace <check-christmas-lottery-numbers-namespace-name> --dry-run -o yaml | kubectl apply -f -
```

Below are examples of how to use the chart:

- Create [cronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) that executes each ten minutes check-christmas-lottery-numbers and necessary resources with sending notifications disabled:
  ```shell
  helm upgrade --install check-christmas-lottery-numbers ./helm/check-christmas-lottery-numbers \
  --set-file numbers_to_check=<number_to_check_dir> \
  -n <check-christmas-lottery-numbers-namespace-name>
  ```

- Create [cronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) that executes each `<N>` minutes check-christmas-lottery-numbers and necessary resources with sending notifications disabled:
  ```shell
  helm upgrade --install check-christmas-lottery-numbers ./helm/check-christmas-lottery-numbers \
  --set config.minutesSchedulePeriod=<N> \
  --set-file numbers_to_check=<number_to_check_dir> \
  --set-string config.notify="false" \
  -n <check-christmas-lottery-numbers-namespace-name>
  ```

- Create [cronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) that executes each ten minutes check-christmas-lottery-numbers and necessary resources with sending notifications enabled and store notifications enabled with a [statefulset](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) with a mongo host to store notifications results:
  ```shell
  helm upgrade --install check-christmas-lottery-numbers ./helm/check-christmas-lottery-numbers \
  --set config.push_over_notification_token=$$PUSH_OVER_NOTIFICATION_TOKEN \
  --set config.push_over_notification_user=$PUSH_OVER_NOTIFICATION_USER \
  --set-file numbers_to_check=<number_to_check_file_dir> \
  --set-string config.notify="true" \
  --set-string config.storeNotifications="true" \
  --set config.mongodb.rootCredentials.username="user" \
  --set config.mongodb.rootCredentials.password="password" \
  --set config.mongodb.storage.className="standard" \
  -n <check-christmas-lottery-numbers-namespace-name>
  ```

- Create [cronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) that executes each ten minutes check-christmas-lottery-numbers and necessary resources with sending notifications enabled and store notifications enabled with a [statefulset](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) with a mongo host to store notifications results, also creates a [Secret](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) with dockerhub credentials to avoid problem with limit in pull requests:
  ```shell
  helm upgrade --install check-christmas-lottery-numbers ./helm/check-christmas-lottery-numbers \
  --set config.push_over_notification_token=$PUSH_OVER_NOTIFICATION_TOKEN \
  --set config.push_over_notification_user=$PUSH_OVER_NOTIFICATION_USER \
  --set imageCredentials.url=https://index.docker.io/v1/  \
  --set imageCredentials.username=<dockerhub_username> \
  --set imageCredentials.password=<dockerhub_password> \
  --set imageCredentials.name=<dockhub_pull_secret_name> \
  --set-file numbers_to_check=<number_to_check_file_dir> \
  --set-string config.notify="true" \
  --set-string config.storeNotifications="true" \
  --set config.mongodb.rootCredentials.username="user" \
  --set config.mongodb.rootCredentials.password="password" \
  --set config.mongodb.storage.className="standard" \
  -n <check-christmas-lottery-numbers-namespace-name>
  ```

### Running the tests
Due to being an application with a single entry point, it does not make sense to perform unit tests, but rather integration tests that check that the expected actions are performed based on the input parameters provided.

#### Tests requirements

A series of variables must be provided in order to carry out the execution of the integration tests mentioned, these variables are mentioned previously in [requirements](#requirements) section.

### Dependencies & Refs

- [urfave/cli](https://github.com/urfave/cli)
- [joho/godotenv](https://github.com/joho/godotenv)
- [google/go-cmp](https://github.com/google/go-cmp)

### LICENSE

[MIT license](LICENSE)

### Author(s)

- [xoanmm](https://github.com/xoanmm)

<!-- JUST BADGES & LINKS -->
[tests-badge]: https://img.shields.io/github/workflow/status/xoanmm/check-christmas-lottery-numbers/Test
[tests-link]: https://github.com/xoanmm/check-christmas-lottery-numbers/actions?query=workflow%3ATest

[release-badge]: https://img.shields.io/github/release/xoanmm/check-christmas-lottery-numbers.svg?logo=github&labelColor=262b30
[release-link]: https://github.com/xoanmm/check-christmas-lottery-numbers/releases

[report-badge]: https://goreportcard.com/badge/github.com/xoanmm/check-christmas-lottery-numbers
[report-link]: https://goreportcard.com/report/github.com/xoanmm/check-christmas-lottery-numbers

[license-badge]: https://img.shields.io/github/license/xoanmm/check-christmas-lottery-numbers
[license-link]: https://github.com/xoanmm/check-christmas-lottery-numbers/blob/master/LICENSE

[coverage-badge]: https://sonarcloud.io/api/project_badges/measure?project=xoanmm_check-christmas-lottery-numbers&metric=coverage
[coverage-link]: https://sonarcloud.io/dashboard?id=xoanmm_check-christmas-lottery-numbers
