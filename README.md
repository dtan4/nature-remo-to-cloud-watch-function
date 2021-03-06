# nature-remo-to-cloud-watch-function

[![GitHub Actions](https://github.com/dtan4/nature-remo-to-cloud-watch-function/workflows/Test/badge.svg)](https://github.com/dtan4/nature-remo-to-cloud-watch-function/actions?query=workflow%3ATest+branch%3Amaster)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=dtan4/nature-remo-to-cloud-watch-function)](https://dependabot.com)
[![codecov](https://codecov.io/gh/dtan4/nature-remo-to-cloud-watch-function/branch/master/graph/badge.svg)](https://codecov.io/gh/dtan4/nature-remo-to-cloud-watch-function)

A Lambda function which fetches the room temperature from [Nature Remo](https://nature.global/en/top) [Cloud API](https://developer.nature.global/en/overview), then posts it to CloudWatch Metrics every minute.

![image](https://user-images.githubusercontent.com/680124/52900772-95da0d00-323d-11e9-98d4-6c3a64cd54dc.png)

### X-Ray

![image](https://user-images.githubusercontent.com/680124/52900804-2b759c80-323e-11e9-9b49-ff5136244896.png)

## Usage

### Parameters

These parameters must be set in [System Manager Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html) with encryption.

| key                                                     | description                                                                                                            |
|---------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------|
| `/natureRemoToCloudWatchFunction/natureRemoAccessToken` | Nature Remo Cloud API access token                                                                                     |
| `/natureRemoToCloudWatchFunction/deviceID`              | Device ID retrieved from [List Devices API](http://swagger.nature.global/#/default/get_1_devices) |
| `/natureRemoToCloudWatchFunction/sentryDsn`              | (Optional) Sentry DSN |

## Development

### Build a binary

```bash
make
```

### Run tests

```bash
make test
```

### Deploy

The following environment variables are required to deploy:

- AWS credentials (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_DEFAULT_REGION`)
- `AWS_S3_BUCKET`
  - S3 bucket to store Lambda artifact
- `AWS_CLOUDFORMATION_STACK_NAME`
  - CloudFormation stack name to use

```bash
make deploy
```

## Author

Daisuke Fujita ([@dtan4](https://github.com/dtan4))

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
