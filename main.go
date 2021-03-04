package main 

import (
	"os"
	"log"
	"context"
	"io/ioutil"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func getComputeClient(ctx context.Context,options option.ClientOption) *compute.Service {
	computeService, err := compute.NewService(ctx,options)
	if err != nil {
                log.Fatal(err)
        }
	return computeService
}

func main() {
	project:= "sap-abs-dev"
	region:= "us-central1"
	//zone := "us-central1-a"
	var SnapshotRetentionDays int64 =7
	var SnapshotFrequencyInHours int64 = 2
	SnapshotScheduleName:= project + "-" + region + "-snapshot-schedule"
	
	serviceAccountJson:= os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
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
	clientOptions:= option.WithCredentialsJSON(data)
	computeClient:=getComputeClient(ctx,clientOptions)
	resourcePoliciesService:=compute.NewResourcePoliciesService(computeClient)
	resourcePolicyGet:=resourcePoliciesService.Get(project,region,SnapshotScheduleName)
	_,err=resourcePolicyGet.Do()
	if err == nil {
		log.Printf("SnapshotScheduleName: %v Already Exists",SnapshotScheduleName)
		os.Exit(0)
	}
	log.Println("Snapshot Schedule Doesnt Exist.Proceeding with Creation ..")
	SnapshotSchedulePolicy:=&compute.ResourcePolicySnapshotSchedulePolicy{
		RetentionPolicy: &compute.ResourcePolicySnapshotSchedulePolicyRetentionPolicy{
			MaxRetentionDays: SnapshotRetentionDays,
			OnSourceDiskDelete: "APPLY_RETENTION_POLICY",
		},
		Schedule: &compute.ResourcePolicySnapshotSchedulePolicySchedule{
			HourlySchedule: &compute.ResourcePolicyHourlyCycle{
				HoursInCycle: SnapshotFrequencyInHours,
				StartTime: "00:00",
			},
		},
		SnapshotProperties: &compute.ResourcePolicySnapshotSchedulePolicySnapshotProperties{
			GuestFlush: false,
			StorageLocations: []string{region},
		},
	}
	resourcePolicy:=&compute.ResourcePolicy{
		Description: "Snapshot Schedule",
		Name: SnapshotScheduleName,
		Region : region,
		SnapshotSchedulePolicy: SnapshotSchedulePolicy,
	}
	ResourcePoliciesInsert:=resourcePoliciesService.Insert(
		project,
		region,
		resourcePolicy,
	)
	op,err:=ResourcePoliciesInsert.Do()
	log.Println("SnapShot Schedule has been created")

}
