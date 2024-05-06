package assets

import "embed"

//go:embed apps aws elastic firebase gcp generic k8s oci openstack programming saas
var Embedded embed.FS
