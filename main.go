package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"context"
	"google.golang.org/api/option"
)

func main() {
	project := "sap-abs-dev"
	region := "us-central1"
	//zone := "us-central1-a"
	var SnapshotRetentionDays int64 = 7
	var SnapshotFrequencyInHours int64 = 2
	SnapshotScheduleName := project + "-" + region + "-snapshot-schedule"

	serviceAccountJson := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if len(serviceAccountJson) > 0 {
		log.Print("Found Env Variable : GOOGLE_APPLICATION_CREDENTIALS ")
	} else {
		log.Fatal("Env Variable GOOGLE_APPLICATION_CREDENTIALS Not Found")
		os.Exit(1)
	}
	ctx := context.Background()
	data, err := ioutil.ReadFile(filepath.ToSlash(serviceAccountJson))
	if err != nil {
		log.Fatal(err)
	}
	clientOptions := option.WithCredentialsJSON(data)
	computeClient := getComputeClient(ctx, clientOptions)
	status:=CreateSnapShotSchedule(
		computeClient,
		project,
		region,
		SnapshotScheduleName,
		SnapshotRetentionDays,
		SnapshotFrequencyInHours,
	)
	if ! status {
		log.Fatal("Error Creating the Snapshot :" ,SnapshotScheduleName)
		os.Exit(1)
	}

}
