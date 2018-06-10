# Daily Google Cloud Platform Cost Report Batch
This batch sends cost of google cloud platform from BigQuery.

## Usage
Set environmental variables for BigQuery.
```
$ export GOOGLE_APPLICATION_CREDENTIALS=../key/credentials.json
```

Set environment variables in codes.
```
slack.go
 slackToken     = "" // your token
 slackChannel   = "" // your channel
 slackIconEmoji = "" // your icon emoji

biguery.go
 ProjectID = "" // your project ID
 TableName = "" // your table name
```

Build and run command with a parameter.
```
$ ./gcp-cost -targetYmd 20170725
```

## Parameters
|Parameter|Rule|Description|Example|
|---|---|---|---|
|targetYmd|required|Target date to collect billing of Google Cloud Platform.<br>When you set targetYmd ``20170725``, this batch aggregates GCP billing from 2017/07/25 00:00:00 +0000 UTC to 2017/07/26 07:00:00 +0000 UTC.|20170725|
