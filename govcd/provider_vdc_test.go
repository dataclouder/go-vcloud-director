//go:build providervdc || functional || ALL

package govcd

import (
	"fmt"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	. "gopkg.in/check.v1"
	"strings"
)

func init() {
	testingTags["providervdc"] = "provider_vdc_test.go"
}

func (vcd *TestVCD) Test_GetProviderVdc(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	var providerVdcs []*ProviderVdc
	providerVdc, err := vcd.client.GetProviderVdcByName(vcd.config.VCD.NsxtProviderVdc.Name)
	check.Assert(err, IsNil)
	providerVdcs = append(providerVdcs, providerVdc)
	providerVdc, err = vcd.client.GetProviderVdcById(providerVdc.ProviderVdc.ID)
	check.Assert(err, IsNil)
	providerVdcs = append(providerVdcs, providerVdc)
	providerVdc, err = vcd.client.GetProviderVdcByHref(providerVdc.ProviderVdc.HREF)
	check.Assert(err, IsNil)
	providerVdcs = append(providerVdcs, providerVdc)

	// Common asserts
	for _, providerVdc := range providerVdcs {
		check.Assert(providerVdc.ProviderVdc.Name, Equals, vcd.config.VCD.NsxtProviderVdc.Name)
		foundStorageProfile := false
		for _, storageProfile := range providerVdc.ProviderVdc.StorageProfiles.ProviderVdcStorageProfile {
			if storageProfile.Name == vcd.config.VCD.NsxtProviderVdc.StorageProfile {
				foundStorageProfile = true
				break
			}
		}
		check.Assert(foundStorageProfile, Equals, true)
		check.Assert(*providerVdc.ProviderVdc.IsEnabled, Equals, true)
		check.Assert(providerVdc.ProviderVdc.ComputeCapacity, NotNil)
		check.Assert(providerVdc.ProviderVdc.Status, Equals, 1)
		check.Assert(len(providerVdc.ProviderVdc.NetworkPoolReferences.NetworkPoolReference), Equals, 1)
		foundNetworkPool := false
		for _, networkPool := range providerVdc.ProviderVdc.NetworkPoolReferences.NetworkPoolReference {
			if networkPool.Name == vcd.config.VCD.NsxtProviderVdc.NetworkPool {
				foundNetworkPool = true
				break
			}
		}
		check.Assert(foundNetworkPool, Equals, true)
		check.Assert(providerVdc.ProviderVdc.Link, NotNil)
	}
}

func (vcd *TestVCD) Test_GetProviderVdcExtended(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	var providerVdcsExtended []*ProviderVdcExtended
	providerVdcExtended, err := vcd.client.GetProviderVdcExtendedByName(vcd.config.VCD.NsxtProviderVdc.Name)
	check.Assert(err, IsNil)
	providerVdcsExtended = append(providerVdcsExtended, providerVdcExtended)
	providerVdcExtended, err = vcd.client.GetProviderVdcExtendedById(providerVdcExtended.VMWProviderVdc.ID)
	check.Assert(err, IsNil)
	providerVdcsExtended = append(providerVdcsExtended, providerVdcExtended)
	providerVdcExtended, err = vcd.client.GetProviderVdcExtendedByHref(providerVdcExtended.VMWProviderVdc.HREF)
	check.Assert(err, IsNil)
	providerVdcsExtended = append(providerVdcsExtended, providerVdcExtended)

	// Common asserts
	for _, providerVdcExtended := range providerVdcsExtended {
		// Basic PVDC asserts
		check.Assert(providerVdcExtended.VMWProviderVdc.Name, Equals, vcd.config.VCD.NsxtProviderVdc.Name)
		foundStorageProfile := false
		for _, storageProfile := range providerVdcExtended.VMWProviderVdc.StorageProfiles.ProviderVdcStorageProfile {
			if storageProfile.Name == vcd.config.VCD.NsxtProviderVdc.StorageProfile {
				foundStorageProfile = true
				break
			}
		}
		check.Assert(foundStorageProfile, Equals, true)
		check.Assert(*providerVdcExtended.VMWProviderVdc.IsEnabled, Equals, true)
		check.Assert(providerVdcExtended.VMWProviderVdc.ComputeCapacity, NotNil)
		check.Assert(providerVdcExtended.VMWProviderVdc.Status, Equals, 1)
		check.Assert(len(providerVdcExtended.VMWProviderVdc.NetworkPoolReferences.NetworkPoolReference), Equals, 1)
		foundNetworkPool := false
		for _, networkPool := range providerVdcExtended.VMWProviderVdc.NetworkPoolReferences.NetworkPoolReference {
			if networkPool.Name == vcd.config.VCD.NsxtProviderVdc.NetworkPool {
				foundNetworkPool = true
				break
			}
		}
		check.Assert(foundNetworkPool, Equals, true)
		check.Assert(providerVdcExtended.VMWProviderVdc.Link, NotNil)
		// Extended PVDC asserts
		check.Assert(providerVdcExtended.VMWProviderVdc.ComputeProviderScope, Equals, "vc1")
		check.Assert(len(providerVdcExtended.VMWProviderVdc.DataStoreRefs.VimObjectRef), Equals, 4)
		check.Assert(strings.HasPrefix(providerVdcExtended.VMWProviderVdc.HighestSupportedHardwareVersion, "vmx-"), Equals, true)
		check.Assert(providerVdcExtended.VMWProviderVdc.HostReferences, NotNil)
		check.Assert(providerVdcExtended.VMWProviderVdc.NsxTManagerReference, NotNil)
		check.Assert(providerVdcExtended.VMWProviderVdc.NsxTManagerReference.Name, Equals, vcd.config.VCD.Nsxt.Manager)
		check.Assert(providerVdcExtended.VMWProviderVdc.ResourcePoolRefs, NotNil)
		check.Assert(providerVdcExtended.VMWProviderVdc.VimServer, NotNil)
	}
}

func (vcd *TestVCD) Test_GetNonExistentProviderVdc(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	providerVdcExtended, err := vcd.client.GetProviderVdcExtendedByName("non-existent-pvdc")
	check.Assert(providerVdcExtended, IsNil)
	check.Assert(err, NotNil)
	providerVdcExtended, err = vcd.client.GetProviderVdcExtendedById("non-existent-pvdc")
	check.Assert(providerVdcExtended, IsNil)
	check.Assert(err, NotNil)
	providerVdcExtended, err = vcd.client.GetProviderVdcExtendedByHref("non-existent-pvdc")
	check.Assert(providerVdcExtended, IsNil)
	check.Assert(err, NotNil)
	providerVdc, err := vcd.client.GetProviderVdcByName("non-existent-pvdc")
	check.Assert(providerVdc, IsNil)
	check.Assert(err, NotNil)
	providerVdc, err = vcd.client.GetProviderVdcById("non-existent-pvdc")
	check.Assert(providerVdc, IsNil)
	check.Assert(err, NotNil)
	providerVdc, err = vcd.client.GetProviderVdcByHref("non-existent-pvdc")
	check.Assert(providerVdc, IsNil)
	check.Assert(err, NotNil)
}

func (vcd *TestVCD) Test_GetProviderVdcConvertFromExtendedToNormal(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}

	providerVdcExtended, err := vcd.client.GetProviderVdcExtendedByName(vcd.config.VCD.NsxtProviderVdc.Name)
	check.Assert(err, IsNil)
	providerVdc, err := providerVdcExtended.ToProviderVdc()
	check.Assert(err, IsNil)
	check.Assert(providerVdc.ProviderVdc.Name, Equals, vcd.config.VCD.NsxtProviderVdc.Name)
	foundStorageProfile := false
	for _, storageProfile := range providerVdc.ProviderVdc.StorageProfiles.ProviderVdcStorageProfile {
		if storageProfile.Name == vcd.config.VCD.NsxtProviderVdc.StorageProfile {
			foundStorageProfile = true
			break
		}
	}
	check.Assert(foundStorageProfile, Equals, true)
	check.Assert(*providerVdc.ProviderVdc.IsEnabled, Equals, true)
	check.Assert(providerVdc.ProviderVdc.Status, Equals, 1)
	check.Assert(len(providerVdc.ProviderVdc.NetworkPoolReferences.NetworkPoolReference), Equals, 1)
	check.Assert(providerVdc.ProviderVdc.NetworkPoolReferences.NetworkPoolReference[0].Name, Equals, vcd.config.VCD.NsxtProviderVdc.NetworkPool)
	check.Assert(providerVdc.ProviderVdc.Link, NotNil)
}

func (vcd *TestVCD) Test_ProviderVdcCRUD(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}
	if vcd.config.Vsphere.ResourcePoolForVcd1 == "" {
		check.Skip("no resource pool defined for this VCD")
	}
	providerVdcName := check.TestName()
	providerVdcDescription := check.TestName()
	storageProfileList, err := vcd.client.Client.QueryAllProviderVdcStorageProfiles()
	check.Assert(err, IsNil)
	check.Assert(len(storageProfileList) > 0, Equals, true)
	var storageProfile types.QueryResultProviderVdcStorageProfileRecordType
	for _, sp := range storageProfileList {
		if sp.Name == vcd.config.VCD.NsxtProviderVdc.StorageProfile {
			storageProfile = *sp
		}
	}
	check.Assert(storageProfile.HREF, Not(Equals), "")

	vcenter, err := vcd.client.GetVcenterByName(vcd.config.VCD.VimServer)
	check.Assert(err, IsNil)
	check.Assert(vcenter, NotNil)

	resourcePool, err := vcenter.GetResourcePoolByName(vcd.config.Vsphere.ResourcePoolForVcd1)
	check.Assert(err, IsNil)
	check.Assert(resourcePool, NotNil)

	nsxtManagers, err := vcd.client.QueryNsxtManagerByName(vcd.config.VCD.Nsxt.Manager)
	check.Assert(err, IsNil)
	check.Assert(len(nsxtManagers), Equals, 1)

	hwVersion, err := resourcePool.GetDefaultHardwareVersion()
	check.Assert(err, IsNil)

	vcenterUrl, err := vcenter.GetVimServerUrl()
	check.Assert(err, IsNil)

	networkPool, err := vcd.client.GetNetworkPoolByName(vcd.config.VCD.NsxtProviderVdc.NetworkPool)
	check.Assert(err, IsNil)
	networkPoolHref, err := networkPool.GetOpenApiUrl()
	check.Assert(err, IsNil)

	providerVdcCreation := types.ProviderVdcCreation{
		Name:                            providerVdcName,
		Description:                     providerVdcDescription,
		HighestSupportedHardwareVersion: hwVersion,
		IsEnabled:                       true,
		VimServer: []*types.Reference{
			{
				HREF: vcenterUrl,
				ID:   extractUuid(vcenter.VSphereVcenter.VcId),
				Name: vcenter.VSphereVcenter.Name,
			},
		},
		ResourcePoolRefs: &types.VimObjectRefs{
			VimObjectRef: []*types.VimObjectRef{
				{
					VimServerRef: &types.Reference{
						HREF: vcenterUrl,
						ID:   extractUuid(vcenter.VSphereVcenter.VcId),
						Name: vcenter.VSphereVcenter.Name,
					},
					MoRef:         resourcePool.ResourcePool.Moref,
					VimObjectType: "RESOURCE_POOL",
				},
			},
		},
		StorageProfile: []string{storageProfile.Name},
		NsxTManagerReference: types.Reference{
			HREF: nsxtManagers[0].HREF,
			ID:   extractUuid(nsxtManagers[0].HREF),
			Name: nsxtManagers[0].Name,
		},
		NetworkPool: types.Reference{
			HREF: networkPoolHref,
			Name: networkPool.NetworkPool.Name,
			ID:   extractUuid(networkPool.NetworkPool.Id),
			Type: networkPool.NetworkPool.PoolType,
		},
		AutoCreateNetworkPool: false,
	}
	printVerbose("  creating provider VDC '%s' using resource pool '%s' and storage profile '%s'\n",
		providerVdcName, resourcePool.ResourcePool.Name, storageProfile.Name)
	providerVdcJson, err := vcd.client.CreateProviderVdc(&providerVdcCreation)
	check.Assert(err, IsNil)
	check.Assert(providerVdcJson, NotNil)
	check.Assert(providerVdcJson.VMWProviderVdc.Name, Equals, providerVdcName)

	AddToCleanupList(providerVdcName, "provider_vdc", "", check.TestName())
	retrievedPvdc, err := vcd.client.GetProviderVdcExtendedByName(providerVdcName)
	check.Assert(err, IsNil)

	err = retrievedPvdc.Disable()
	check.Assert(err, IsNil)
	check.Assert(retrievedPvdc.VMWProviderVdc.IsEnabled, NotNil)
	check.Assert(*retrievedPvdc.VMWProviderVdc.IsEnabled, Equals, false)

	err = retrievedPvdc.Enable()
	check.Assert(err, IsNil)
	check.Assert(retrievedPvdc.VMWProviderVdc.IsEnabled, NotNil)
	check.Assert(*retrievedPvdc.VMWProviderVdc.IsEnabled, Equals, true)

	newProviderVdcName := "TestNewName"
	newProviderVdcDescription := "Test New provider VDC description"
	printVerbose("  renaming provider VDC to '%s'\n", newProviderVdcName)
	err = retrievedPvdc.Rename(newProviderVdcName, newProviderVdcDescription)
	check.Assert(err, IsNil)
	check.Assert(retrievedPvdc.VMWProviderVdc.Name, Equals, newProviderVdcName)
	check.Assert(retrievedPvdc.VMWProviderVdc.Description, Equals, newProviderVdcDescription)

	printVerbose("  renaming back provider VDC to '%s'\n", providerVdcName)
	err = retrievedPvdc.Rename(providerVdcName, providerVdcDescription)
	check.Assert(err, IsNil)
	check.Assert(retrievedPvdc.VMWProviderVdc.Name, Equals, providerVdcName)
	check.Assert(retrievedPvdc.VMWProviderVdc.Description, Equals, providerVdcDescription)

	secondResourcePoolName := vcd.config.Vsphere.ResourcePoolForVcd2
	if secondResourcePoolName != "" {
		printVerbose("  adding resource pool '%s' to provider VDC\n", secondResourcePoolName)
		secondResourcePool, err := vcenter.GetResourcePoolByName(secondResourcePoolName)
		check.Assert(err, IsNil)
		check.Assert(secondResourcePool, NotNil)
		err = retrievedPvdc.AddResourcePools([]*ResourcePool{secondResourcePool})
		check.Assert(err, IsNil)
		err = retrievedPvdc.Refresh()
		check.Assert(err, IsNil)
		check.Assert(len(retrievedPvdc.VMWProviderVdc.ResourcePoolRefs.VimObjectRef), Equals, 2)

		printVerbose("  removing resource pool '%s' from provider VDC\n", secondResourcePoolName)
		err = retrievedPvdc.DeleteResourcePools([]*ResourcePool{secondResourcePool})
		check.Assert(err, IsNil)
		err = retrievedPvdc.Refresh()
		check.Assert(err, IsNil)
		check.Assert(len(retrievedPvdc.VMWProviderVdc.ResourcePoolRefs.VimObjectRef), Equals, 1)
	}

	secondStorageProfile := vcd.config.VCD.NsxtProviderVdc.StorageProfile2
	if secondStorageProfile != "" {
		printVerbose("  adding storage profile '%s' to provider VDC\n", secondStorageProfile)
		// Adds a storage profile
		err = retrievedPvdc.AddStorageProfiles([]string{secondStorageProfile})
		check.Assert(err, IsNil)
		check.Assert(len(retrievedPvdc.VMWProviderVdc.StorageProfiles.ProviderVdcStorageProfile), Equals, 2)

		fmt.Printf("  removing storage profile '%s' from provider VDC\n", secondStorageProfile)
		// Remove a storage profile
		err = retrievedPvdc.DeleteStorageProfiles([]string{secondStorageProfile})
		check.Assert(err, IsNil)
		check.Assert(len(retrievedPvdc.VMWProviderVdc.StorageProfiles.ProviderVdcStorageProfile), Equals, 1)
	}

	// Deleting while the Provider VDC is still enabled will fail
	task, err := retrievedPvdc.Delete()
	check.Assert(err, NotNil)

	// Properly deleting provider VDC: first disabling, then removing
	printVerbose("  disabling provider VDC '%s'\n", providerVdcName)
	err = retrievedPvdc.Disable()
	check.Assert(err, IsNil)
	check.Assert(retrievedPvdc.VMWProviderVdc.IsEnabled, NotNil)
	check.Assert(*retrievedPvdc.VMWProviderVdc.IsEnabled, Equals, false)

	printVerbose("  removing provider VDC '%s'\n", providerVdcName)
	task, err = retrievedPvdc.Delete()
	check.Assert(err, IsNil)
	err = task.WaitTaskCompletion()
	check.Assert(err, IsNil)
}
