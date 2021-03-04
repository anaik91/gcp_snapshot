package main

import (
	"log"
	"context"
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