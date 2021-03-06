/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package terraform

// Resource represents a managed or data type that is created by the module
type Resource struct {
	Type           string
	ProviderName   string
	ProviderSource string
	Mode           string
	Version        string
}
