package main

import (
	"log"
	"strings"
	"context"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func GetComputeClient(ctx context.Context,options option.ClientOption) *compute.Service {
	computeService, err := compute.NewService(ctx,options)
	if err != nil {
                log.Fatal(err)
        }
	return computeService
}

func FetchComputeZones(computeClient *compute.Service, project,region string) ([]string,bool) {
	zones:= []string{}
	ZonesService:=compute.NewZonesService(computeClient)
	ZonesListCall:=ZonesService.List(project)
	zonesInfo,err:=ZonesListCall.Do()
	if err != nil {
		log.Println(err)
		return []string{},false
	}
	for _,v:=range zonesInfo.Items {
		if strings.Contains(v.Region,region) {
			zones = append(zones,v.Name)
		}
	}
	return zones,true
}

func AttachSnapshotSchedule(computeClient *compute.Service, project,zone ,disk,scheduleName string) bool {
	DisksService:=compute.NewDisksService(computeClient)
	DisksAddResourcePoliciesCall:=DisksService.AddResourcePolicies(
		project,
		zone,
		disk,
		&compute.DisksAddResourcePoliciesRequest{
			ResourcePolicies: []string{scheduleName},
		},
	)
	_,err:=DisksAddResourcePoliciesCall.Do()
	if err != nil {
		if strings.Contains(err.Error(), "Disk already has resource policy attached") {
			log.Println(disk," : already has resource policy attached")
			return true	
		}
		log.Println(err.Error())
		return false
	}
	return true	
}

func ListComputeByLabel(computeClient *compute.Service, project,zone string ,labels map[string]string) []string {
	var instancelist []string
	InstanceService:=compute.NewInstancesService(computeClient)
	InstanceList:=InstanceService.List(project,zone)
	instances,_:=InstanceList.Do()
	for _,v := range instances.Items {
		filterCheck:=func ( uLabel,iLabel map[string]string ) bool {
			for uk,uv:= range uLabel {
				if lv,ok:=iLabel[uk]; ok && lv == uv {
				} else {
					return false
				}
			}
			return true
		}
		if filterCheck(labels,v.Labels) {
			instancelist=append(instancelist,v.Name)
		}
	}
	return instancelist
}

func GetComputeDisks(computeClient *compute.Service, project,zone,vmName string ) []string {
	var instanceDisks []string
	InstanceService:=compute.NewInstancesService(computeClient)
	InstancesGetCall:=InstanceService.Get(project,zone,vmName)
	instanceData,err:=InstancesGetCall.Do()
	if err != nil {
		log.Fatal(err)
		return instanceDisks
	}
	for _,v := range instanceData.Disks {
		instanceDisks=append(instanceDisks,v.DeviceName)
	}
	return instanceDisks
}