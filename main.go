package main

import (
	"strings"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"context"
	"google.golang.org/api/option"
)

func main() {
	project:= "i501950"
	gcp_project := "sap-abs-dev"
	region := "us-central1"
	zone := "us-central1-a"
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
	computeClient := GetComputeClient(ctx, clientOptions)
	SnapShotScheduleSelf,status:=CreateSnapShotSchedule(
		computeClient,
		gcp_project,
		region,
		SnapshotScheduleName,
		SnapshotRetentionDays,
		SnapshotFrequencyInHours,
	)
	if ! status {
		log.Fatal("Error Creating the Snapshot :" ,SnapshotScheduleName)
		os.Exit(1)
	}
	instanceFilter:= map[string]string{"project":project,"backup":"backup"}
	instanceList:=ListComputeByLabel(computeClient,gcp_project,zone,instanceFilter)
	log.Println(instanceList)
	for _,v:= range instanceList {
		log.Println("Associating Snapshot Schedule : ", SnapShotScheduleSelf , "to Instance ",v)
		disks:=GetComputeDisks(computeClient,gcp_project,zone,v)
		log.Printf("Disks attached to %v are %v\n",v,strings.Join(disks,","))
		for _,d := range disks {
			log.Println("Setting Snapshot Schedule : ", SnapshotScheduleName , " to disk :" ,d)
			if ok:=AttachSnapshotSchedule(
				computeClient,
				gcp_project,
				zone,
				d,
				SnapShotScheduleSelf,
			); ok {
				log.Println("Snapshot Schedule : ", SnapshotScheduleName , " has been associated to disk :" ,d)
			}
		}
	}
	log.Println("Done")
}