// Copyright 2022 Upbound Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package converters

import (
	"github.com/upbound/upjet/pkg/migration"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	Registry = migration.NewRegistry(runtime.NewScheme())
)

func init() {
	if err := Registry.AddCompositionTypes(); err != nil {
		panic(err)
	}
	Registry.AddClaimType(
		// claim GVK
		schema.GroupVersionKind{
			Group:   "aws.platformref.upbound.io",
			Version: "v1alpha1",
			Kind:    "Subnet",
		},
	)
	Registry.AddCompositeType(
		// composite resource GVK
		schema.GroupVersionKind{
			Group:   "aws.platformref.upbound.io",
			Version: "v1alpha1",
			Kind:    "XSubnet",
		},
	)
}
