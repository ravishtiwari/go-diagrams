package openstack

import "github.com/emarais-godaddy/go-diagrams/diagram"

type operationsContainer struct {
	path string
	opts []diagram.NodeOption
}

var Operations = &operationsContainer{
	opts: diagram.OptionSet{diagram.Provider("openstack"), diagram.NodeShape("none")},
	path: "assets/openstack/operations",
}
