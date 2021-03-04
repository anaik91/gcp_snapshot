package main

import (
	"strings"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	//"path/filepath"
	"gopkg.in/ini.v1"
	"context"
	"google.golang.org/api/option"
)

func main() {
	cfg, err := ini.Load("input.properties")
    if err != nil {
        log.Printf("Fail to read file: %v", err)
        os.Exit(1)
    }
	project:= cfg.Section("inputs").Key("project").String()
	gcp_project := cfg.Section("inputs").Key("gcp_project").String()
	region := cfg.Section("inputs").Key("region").String()
	zones := strings.Split(cfg.Section("inputs").Key("zones").String(),",")
	SnapshotRetentionDays,_ := strconv.ParseInt(cfg.Section("inputs").Key("SnapshotRetentionDays").String(), 10, 64)
	SnapshotFrequencyInHours,_ := strconv.ParseInt(cfg.Section("inputs").Key("SnapshotFrequencyInHours").String(), 10, 64)
	SnapshotScheduleName := project + "-" + region + "-snapshot-schedule"

	serviceAccountJson := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if len(serviceAccountJson) > 0 {
		log.Print("Found Env Variable : GOOGLE_APPLICATION_CREDENTIALS ")
	} else {
		log.Fatal("Env Variable GOOGLE_APPLICATION_CREDENTIALS Not Found")
		os.Exit(1)
	}
	ctx := context.Background()
	data, err := ioutil.ReadFile(serviceAccountJson)
	if err != nil {
		log.Fatal(err)
	}
	clientOptions := option.WithCredentialsJSON(data)
	computeClient := GetComputeClient(ctx, clientOptions)
	all_zones,_:= FetchComputeZones(computeClient,gcp_project,region)
	for _,v := range zones {
		flag:=false
		for _,z:= range all_zones {
			if z == region + "-" + v {
				flag=true
			}
		}
		if ! flag {
			log.Println("Invalid Input provided for zones List :" , zones,"Kindly recitfy.")
			os.Exit(1)
		}
	}
	log.Println("Zones Provided are valid.")
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
	for _,eachzone := range  zones {
		zone:=region+"-"+ eachzone
		instanceFilter:= map[string]string{"project":project,"backup":"backup"}
		log.Println("Listing VMs in Zone :", zone , "with filter  :",instanceFilter)
		instanceList:=ListComputeByLabel(computeClient,gcp_project,zone,instanceFilter)
		log.Println(instanceList)
		for _,v:= range instanceList {
			log.Println("Associating Snapshot Schedule : ", SnapshotScheduleName , "to Instance ",v)
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
	}
	log.Println("Done")
}