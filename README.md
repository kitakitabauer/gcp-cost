# Daily Google Cloud Platform Cost Report Batch
This batch sends cost of google cloud platform from BigQuery.

## Usage
Set environmental variables for BigQuery.
```
$ export GOOGLE_APPLICATION_CREDENTIALS=../key/gcp-cost-XXX.json
```

Run command with some parameters.
```
$ ./gcp-cost -targetYmd 2020/04/10 -channel #api-server
```

## Parameters
|Parameter|Rule|Description|Example|
|---|---|---|---|
|targetYmd|required|Target date to collect billing of Google Cloud Platform.<br>When you set targetYmd ``2020/04/10``, this batch aggregates GCP billing from 2020/04/10 00:00:00 +0000 UTC to 2020/04/10 23:59:59 +0000 UTC.|2020/04/10|
|channel|required|Set the channel name or channel ID of slack you want to notify.| #api-server |