package main

import (
	"log"
	"google.golang.org/api/compute/v1"
)

func GetSnapShotSchedule(computeClient *compute.Service,project,region,SnapshotScheduleName string) (string,bool){
	resourcePoliciesService := compute.NewResourcePoliciesService(computeClient)
	resourcePolicyGet := resourcePoliciesService.Get(project, region, SnapshotScheduleName)
	resourcePolicyData, err := resourcePolicyGet.Do()
	if err == nil {
		return resourcePolicyData.SelfLink,true
	} else {
		return "",false
	}
}

func CreateSnapShotSchedule(computeClient *compute.Service,project,region,SnapshotScheduleName string, SnapshotRetentionDays,SnapshotFrequencyInHours int64) (string,bool){
	resourcePoliciesService := compute.NewResourcePoliciesService(computeClient)
	if selfLink,ok:=GetSnapShotSchedule(computeClient,project, region, SnapshotScheduleName); ok {
		log.Printf("SnapshotScheduleName: %v Already Exists", SnapshotScheduleName)
		return selfLink,ok
	}
	log.Println("Snapshot Schedule Doesnt Exist.Proceeding with Creation ..")
	SnapshotSchedulePolicy := &compute.ResourcePolicySnapshotSchedulePolicy{
		RetentionPolicy: &compute.ResourcePolicySnapshotSchedulePolicyRetentionPolicy{
			MaxRetentionDays:   SnapshotRetentionDays,
			OnSourceDiskDelete: "APPLY_RETENTION_POLICY",
		},
		Schedule: &compute.ResourcePolicySnapshotSchedulePolicySchedule{
			HourlySchedule: &compute.ResourcePolicyHourlyCycle{
				HoursInCycle: SnapshotFrequencyInHours,
				StartTime:    "00:00",
			},
		},
		SnapshotProperties: &compute.ResourcePolicySnapshotSchedulePolicySnapshotProperties{
			GuestFlush:       false,
			StorageLocations: []string{region},
		},
	}
	resourcePolicy := &compute.ResourcePolicy{
		Description:            "Snapshot Schedule",
		Name:                   SnapshotScheduleName,
		Region:                 region,
		SnapshotSchedulePolicy: SnapshotSchedulePolicy,
	}
	ResourcePoliciesInsert := resourcePoliciesService.Insert(
		project,
		region,
		resourcePolicy,
	)
	op, err := ResourcePoliciesInsert.Do()
	if err != nil {
		log.Println(op.StatusMessage)
		return "",false
	}
	log.Println("SnapShot Schedule has been created")
	if selfLink,ok:=GetSnapShotSchedule(computeClient,project, region, SnapshotScheduleName); ok {
		return selfLink,ok
	} else {
		return "",false
	}
}
