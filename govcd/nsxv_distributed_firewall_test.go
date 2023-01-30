//go:build functional || ALL

/*
 * Copyright 2022 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */
package govcd

import (
	"fmt"
	"strings"

	"github.com/kr/pretty"
	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_NsxvDistributedFirewall(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())

	org, err := vcd.client.GetAdminOrgByName(vcd.config.VCD.Org)
	check.Assert(err, IsNil)
	check.Assert(org, NotNil)

	if vcd.config.VCD.Nsxt.Vdc != "" {
		if testVerbose {
			fmt.Println("Testing attempted access to NSX-T VDC")
		}
		// Retrieve a NSX-T VDC
		nsxtVdc, err := org.GetAdminVDCByName(vcd.config.VCD.Nsxt.Vdc, false)
		check.Assert(err, IsNil)

		notWorkingDfw := NewNsxvDistributedFirewall(nsxtVdc.client, nsxtVdc.AdminVdc.ID)
		check.Assert(notWorkingDfw, NotNil)

		isEnabled, err := notWorkingDfw.IsEnabled()
		check.Assert(err, IsNil)
		check.Assert(isEnabled, Equals, false)

		err = notWorkingDfw.Enable()
		// NSX-T VDCs don't support NSX-V distributed firewalls. We expect an error here.
		check.Assert(err, NotNil)
		check.Assert(strings.Contains(err.Error(), "Forbidden"), Equals, true)
		if testVerbose {
			fmt.Printf("notWorkingDfw: %s\n", err)
		}
		config, err := notWorkingDfw.GetConfiguration()
		// Also this operation should fail
		check.Assert(err, NotNil)
		check.Assert(config, IsNil)
	}

	// Retrieve a NSX-V VDC
	vdc, err := org.GetAdminVDCByName(vcd.config.VCD.Vdc, false)
	check.Assert(err, IsNil)
	check.Assert(vdc, NotNil)

	dfw := NewNsxvDistributedFirewall(vdc.client, vdc.AdminVdc.ID)
	check.Assert(dfw, NotNil)

	// dfw.Enable is an idempotent operation. It can be repeated on an already enabled DFW
	// without errors.
	err = dfw.Enable()
	check.Assert(err, IsNil)

	enabled, err := dfw.IsEnabled()
	check.Assert(err, IsNil)
	check.Assert(enabled, Equals, true)

	config, err := dfw.GetConfiguration()
	check.Assert(err, IsNil)
	check.Assert(config, NotNil)
	if testVerbose {
		fmt.Printf("%# v\n", pretty.Formatter(config))
	}

	// Repeat enable operation
	err = dfw.Enable()
	check.Assert(err, IsNil)

	enabled, err = dfw.IsEnabled()
	check.Assert(err, IsNil)
	check.Assert(enabled, Equals, true)

	err = dfw.Disable()
	check.Assert(err, IsNil)
	enabled, err = dfw.IsEnabled()
	check.Assert(err, IsNil)
	check.Assert(enabled, Equals, false)

	// Also dfw.Disable is idempotent
	err = dfw.Disable()
	check.Assert(err, IsNil)

	enabled, err = dfw.IsEnabled()
	check.Assert(err, IsNil)
	check.Assert(enabled, Equals, false)
}

/*
// ----------------------------------------------------------------------------------------------
// methods from here till the end of the file will be removed if we decide we don't need services
// ----------------------------------------------------------------------------------------------

func (vcd *TestVCD) Test_NsxvServices(check *C) {
	fmt.Printf("Running: %s\n", check.TestName())

	org, err := vcd.client.GetAdminOrgByName(vcd.config.VCD.Org)
	check.Assert(err, IsNil)
	check.Assert(org, NotNil)

	// Retrieve a NSX-V VDC
	vdc, err := org.GetAdminVDCByName(vcd.config.VCD.Vdc, false)
	check.Assert(err, IsNil)
	check.Assert(vdc, NotNil)

	dfw := NewNsxvDistributedFirewall(vdc.client, vdc.AdminVdc.ID)
	check.Assert(dfw, NotNil)

	services, err := dfw.GetServices(false)
	check.Assert(err, IsNil)
	check.Assert(services, NotNil)
	check.Assert(len(services) > 0, Equals, true)

	if testVerbose {
		fmt.Printf("services: %d\n", len(services))
		fmt.Printf("%# v\n", pretty.Formatter(services[0]))
	}

	serviceName := "PostgreSQL"
	serviceByName, err := dfw.GetServiceByName(serviceName)

	check.Assert(err, IsNil)
	check.Assert(serviceByName, NotNil)
	check.Assert(serviceByName.Name, Equals, serviceName)

	serviceById, err := dfw.GetServiceById(serviceByName.ObjectID)
	check.Assert(err, IsNil)
	check.Assert(serviceById.Name, Equals, serviceName)

	searchRegex := "M.SQL" // Finds, among others, names containing "MySQL" or "MSSQL"
	servicesByRegex, err := dfw.GetServicesByRegex(searchRegex)
	check.Assert(err, IsNil)
	check.Assert(len(servicesByRegex) > 1, Equals, true)

	searchRegex = "." // Finds all services
	servicesByRegex, err = dfw.GetServicesByRegex(searchRegex)
	check.Assert(err, IsNil)
	check.Assert(len(servicesByRegex), Equals, len(services))

	searchRegex = "^####--no-such-service--####$" // Finds no services
	servicesByRegex, err = dfw.GetServicesByRegex(searchRegex)
	check.Assert(err, IsNil)
	check.Assert(len(servicesByRegex), Equals, 0)

	serviceGroups, err := dfw.GetServiceGroups(false)
	check.Assert(err, IsNil)
	check.Assert(serviceGroups, NotNil)
	check.Assert(len(serviceGroups) > 0, Equals, true)

	serviceGroupName := "Orchestrator"
	serviceGroupByName, err := dfw.GetServiceGroupByName(serviceGroupName)
	check.Assert(err, IsNil)
	check.Assert(serviceGroupByName, NotNil)
	check.Assert(serviceGroupByName.Name, Equals, serviceGroupName)

	serviceGroupById, err := dfw.GetServiceGroupById(serviceGroupByName.ObjectID)
	check.Assert(err, IsNil)
	check.Assert(serviceGroupById, NotNil)
	check.Assert(serviceGroupById.Name, Equals, serviceGroupName)

	searchRegex = "Oracle"
	serviceGroupsByRegex, err := dfw.GetServiceGroupsByRegex(searchRegex)
	check.Assert(err, IsNil)
	check.Assert(len(serviceGroupsByRegex) > 1, Equals, true)

	searchRegex = "."
	serviceGroupsByRegex, err = dfw.GetServiceGroupsByRegex(searchRegex)
	check.Assert(err, IsNil)
	check.Assert(len(serviceGroupsByRegex), Equals, len(serviceGroups))

	searchRegex = "^####--no-such-service-group--####$"
	serviceGroupsByRegex, err = dfw.GetServiceGroupsByRegex(searchRegex)
	check.Assert(err, IsNil)
	check.Assert(len(serviceGroupsByRegex), Equals, 0)
}


*/
