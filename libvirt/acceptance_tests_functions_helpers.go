package libvirt

import (
	"fmt"

	"github.com/hashicorp/terraform/terraform"
	libvirt "github.com/libvirt/libvirt-go"
	"github.com/libvirt/libvirt-go-xml"
)

// This file contain function helpers used for testsuite/testacc

// getResourceFromTerraformState get aresource by name
// from terraform states produced during testacc
// and return the resource
func getResourceFromTerraformState(resourceName string, state *terraform.State) (*terraform.ResourceState, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", resourceName)
	}

	if rs.Primary.ID == "" {
		return nil, fmt.Errorf("No libvirt resource key ID is set")
	}
	return rs, nil
}

// helper used in network tests for retrieve xml network definition.
func getNetworkDef(state *terraform.State, name string, virConn libvirt.Connect) (*libvirtxml.Network, error) {
	var network *libvirt.Network
	rs, err := getResourceFromTerraformState(name, state)
	if err != nil {
		return nil, err
	}
	network, err = virConn.LookupNetworkByUUIDString(rs.Primary.ID)
	if err != nil {
		return nil, err
	}
	networkDef, err := getXMLNetworkDefFromLibvirt(network)
	if err != nil {
		return nil, fmt.Errorf("Error reading libvirt network XML description: %s", err)
	}
	return &networkDef, nil
}
