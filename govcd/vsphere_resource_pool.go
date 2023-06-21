/*
 * Copyright 2023 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"net/url"
)

type ResourcePool struct {
	ResourcePool *types.ResourcePool
	vcenter      *VCenter
	client       *VCDClient
}

func NewResourcePool(client *VCDClient, vcenter *VCenter) *ResourcePool {
	return &ResourcePool{
		ResourcePool: new(types.ResourcePool),
		vcenter:      vcenter,
		client:       client,
	}
}

func (vcenter VCenter) GetAllAvailableResourcePools(queryParams url.Values) ([]*ResourcePool, error) {
	client := vcenter.client.Client
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointResourcePools
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(fmt.Sprintf(endpoint, vcenter.VSphereVcenter.VcId))
	if err != nil {
		return nil, err
	}

	var retrieved []*types.ResourcePool

	err = client.OpenApiGetAllItems(minimumApiVersion, urlRef, queryParams, &retrieved, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting resource pool list: %s", err)
	}

	if len(retrieved) == 0 {
		return nil, nil
	}
	var returnList []*ResourcePool

	for _, r := range retrieved {
		if !r.Eligible {
			newRp, err := vcenter.GetAvailableResourcePoolById(r.Moref)
			if err != nil {
				return nil, fmt.Errorf("error while retrieving child resource pool for %s: %s", r.Name, err)
			}
			r = newRp.ResourcePool
		}
		returnList = append(returnList, &ResourcePool{
			ResourcePool: r,
			vcenter:      &vcenter,
			client:       vcenter.client,
		})
	}
	return returnList, nil
}

func (vcenter VCenter) GetAvailableResourcePoolByName(name string) (*ResourcePool, error) {

	resourcePools, err := vcenter.GetAllAvailableResourcePools(nil)
	if err != nil {
		return nil, err
	}
	for _, rp := range resourcePools {
		if rp.ResourcePool.Name == name {
			return rp, nil
		}
	}
	return nil, fmt.Errorf("resource pool '%s' not found: %s", name, ErrorEntityNotFound)
}

func (vcenter VCenter) GetAvailableResourcePoolById(id string) (*ResourcePool, error) {

	client := vcenter.client.Client
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointResourcePools
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(fmt.Sprintf(endpoint, vcenter.VSphereVcenter.VcId) + "/" + id)
	if err != nil {
		return nil, err
	}

	retrieved := []*types.ResourcePool{}
	err = client.OpenApiGetAllItems(minimumApiVersion, urlRef, nil, &retrieved, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting resource pool: %s", err)
	}

	if len(retrieved) == 0 {
		return nil, fmt.Errorf("resource pool %s not found: %s", id, ErrorEntityNotFound)
	}

	return &ResourcePool{
		ResourcePool: retrieved[0],
		vcenter:      &vcenter,
		client:       vcenter.client,
	}, nil
}

func (rp ResourcePool) GetAvailableHardwareVersions() (*types.OpenApiSupportedHardwareVersions, error) {

	client := rp.client.Client
	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointResourcePoolHardware
	minimumApiVersion, err := client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(fmt.Sprintf(endpoint, rp.vcenter.VSphereVcenter.VcId, rp.ResourcePool.Moref))
	if err != nil {
		return nil, err
	}

	retrieved := types.OpenApiSupportedHardwareVersions{}
	err = client.OpenApiGetItem(minimumApiVersion, urlRef, nil, &retrieved, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting resource pool hardware versions: %s", err)
	}

	return &retrieved, nil
}

func (rp ResourcePool) GetDefaultHardwareVersion() (string, error) {

	versions, err := rp.GetAvailableHardwareVersions()
	if err != nil {
		return "", err
	}

	for _, v := range versions.SupportedVersions {
		if v.IsDefault {
			return v.Name, nil
		}
	}
	return "", fmt.Errorf("no default hardware version found for resource pool %s", rp.ResourcePool.Name)
}
