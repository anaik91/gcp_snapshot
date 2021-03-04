package main

import (
	"os"
	"log"
	"google.golang.org/api/compute/v1"
)

func CreateSnapShotSchedule(computeClient *compute.Service,project,region,SnapshotScheduleName string, SnapshotRetentionDays,SnapshotFrequencyInHours int64) bool {
	resourcePoliciesService := compute.NewResourcePoliciesService(computeClient)
	resourcePolicyGet := resourcePoliciesService.Get(project, region, SnapshotScheduleName)
	_, err := resourcePolicyGet.Do()
	if err == nil {
		log.Printf("SnapshotScheduleName: %v Already Exists", SnapshotScheduleName)
		os.Exit(0)
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
		return false
	}
	log.Println("SnapShot Schedule has been created")
	return true
}
