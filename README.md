[![Tests][tests-badge]][tests-link]
[![GitHub Release][release-badge]][release-link]
[![Go Report Card][report-badge]][report-link]
[![License][license-badge]][license-link]

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

- Checking numbers from the `numbers_to_check.json` file without sending notifications:
    ```shell
    $ ./check-christmas-lottery-numbers/cmd -f numbers_to_check_test.json
    2020/12/21 18:17:27 Opening file /tmp/numbers_to_check_test.json with numbers to check
    2020/12/21 18:17:27 Successfully opened file /tmp/numbers_to_check_test.json with numbers to check
    2020/12/21 18:17:27 Numbers are going to be check from 2 different owners
    2020/12/21 18:17:27 1 numbers are going to be check from owner Xoan. The probabily to win a prize is 0.001%
    2020/12/21 18:17:27 Checking prize obtained for number 48335 with 5€ bet and origin My favourite restaurant
    2020/12/21 18:17:27 Prize obtained for number 48335 with 5 € bet and origin is 0€
    2020/12/21 18:17:27 1 numbers are going to be check from owner Manuel. The probabily to win a prize is 0.001%
    2020/12/21 18:17:27 Checking prize obtained for number 48334 with 5€ bet and origin Company number
    2020/12/21 18:17:27 Prize obtained for number 48334 with 5 € bet and origin is 0€
    ```

- Checking numbers from the `numbers_to_check.json` file with sending notifications activated:
    ```shell
    $ ./check-christmas-lottery-numbers/cmd -f numbers_to_check_test.json -n
    2020/12/21 18:29:02 Opening file /tmp/numbers_to_check_test.json with numbers to check
    2020/12/21 18:29:02 Successfully opened file /tmp/numbers_to_check_test.json with numbers to check
    2020/12/21 18:29:02 Numbers are going to be check from 2 different owners
    2020/12/21 18:29:02 1 numbers are going to be check from owner Xoan. The probabily to win a prize is 0.001%
    2020/12/21 18:29:02 Checking prize obtained for number 26590 with 5€ bet and origin My favourite restaurant
    2020/12/21 18:29:02 Prize obtained for number 26590 with 5 € bet and origin is 400000€
    2020/12/21 18:29:02 A notification is going to be send with pushOver
    2020/12/21 18:29:03 Notification send correctly with id 959d78e9-b8d1-4e66-8a79-0e105e799e90 to PushOver App
    2020/12/21 18:29:03 1 numbers are going to be check from owner Manuel. The probabily to win a prize is 0.001%
    2020/12/21 18:29:03 Checking prize obtained for number 48334 with 5€ bet and origin Company number
    2020/12/21 18:29:04 Prize obtained for number 48334 with 5 € bet and origin is 0€
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
      --notify, -n                             Activate notifications through PushOver (default: false)
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

- Create [cronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) that executes each <N> minutes check-christmas-lottery-numbers and necessary resources with sending notifications disabled:
  ```shell
  helm upgrade --install check-christmas-lottery-numbers ./helm/check-christmas-lottery-numbers \
  --set config.minutesSchedulePeriod=<N> \
  --set-file numbers_to_check=<number_to_check_dir> \
  -n <check-christmas-lottery-numbers-namespace-name>
  ```

- Create [cronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) that executes each ten minutes check-christmas-lottery-numbers and necessary resources with sending notifications enabled:
  ```shell
  helm upgrade --install check-christmas-lottery-numbers ./helm/check-christmas-lottery-numbers \
  --set config.push_over_notification_token=$(echo -n "<PUSH_OVER_NOTIFICATION_TOKEN>" | base64 | tr -d "\n") \
  --set config.push_over_notification_user=$(echo -n "<PUSH_OVER_NOTIFICATION_USER>" | base64 | tr -d "\n") \
  --set-file numbers_to_check=<number_to_check_file_dir> \
  -n <check-christmas-lottery-numbers-namespace-name>
  ```

- Create [cronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) that executes each ten minutes check-christmas-lottery-numbers and necessary resources with sending notifications enabled, also creates a [Secret](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) with dockerhub credentials to avoid problem with limit in pull requests:
  ```shell
  helm upgrade --install check-christmas-lottery-numbers ./helm/check-christmas-lottery-numbers \
  --set config.push_over_notification_token=$(echo -n "<PUSH_OVER_NOTIFICATION_TOKEN>" | base64 | tr -d "\n") \
  --set config.push_over_notification_user=$(echo -n "<PUSH_OVER_NOTIFICATION_USER>" | base64 | tr -d "\n") \
  --set imageCredentials.url=https://index.docker.io/v1/  \
  --set imageCredentials.username=<dockerhub_username> \
  --set imageCredentials.password=<dockerhub_password> \
  --set imageCredentials.name=<dockhub_pull_secret_name> \
  --set-file numbers_to_check=<number_to_check_file_dir> \
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