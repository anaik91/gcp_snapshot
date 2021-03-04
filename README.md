# Enable Snapshots Based on Project & Backup Label

## Prerequisites

* GoLang version go1.15.8 linux/amd64
* GCP IAM Service Account JSON

## Inputs

```
[inputs]
project = i501950
gcp_project = sap-abs-dev
region = us-central1
zones = a,b,c
SnapshotRetentionDays = 7
SnapshotFrequencyInHours = 2
```

## Inputs Description

* `project`                   : **Concourse Project**
* `gcp_project`               : **Google Cloud Project Name**
* `region`                    : **Google Cloud Region**
* `zones`                     : **Google Cloud Project Zones (, seperated)**
* `SnapshotRetentionDays`     : **Number of Days to retain Snapshot**
* `SnapshotFrequencyInHours`  : **Snapshot Frequence in hours**


## Run Instructions

Before you can run the Automation . Kindly export the below environment variable to point to GCP IAM Service Account JSON File .

ON Linux/Mac
`export GOOGLE_APPLICATION_CREDENTIALS="/tmp/sap-abs-dev.json"`
ON Windows
`export GOOGLE_APPLICATION_CREDENTIALS=C:\tmp\sap-abs-dev.json"`

TO run the Automation

```
git clone https://github.tools.sap/I501950/gcp_snapshot.git
cd gcp_snapshot
go run .
```
