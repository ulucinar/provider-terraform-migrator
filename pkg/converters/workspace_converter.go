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
	srcapis "github.com/crossplane-contrib/provider-terraform/apis"
	srcv1alpha1 "github.com/crossplane-contrib/provider-terraform/apis/v1alpha1"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	targetv1beta1 "github.com/upbound/provider-terraform/apis/v1beta1"
	"github.com/upbound/upjet/pkg/migration"
)

func workspaceResources(mg resource.Managed) ([]resource.Managed, error) {
	source := mg.(*srcv1alpha1.Workspace)
	target := &targetv1beta1.Workspace{}
	if _, err := migration.CopyInto(source, target, targetv1beta1.WorkspaceGroupVersionKind); err != nil {
		return nil, errors.Wrap(err, "failed to copy source into target")
	}
	return []resource.Managed{
		target,
	}, nil
}

func init() {
	if err := Registry.AddToScheme(srcapis.AddToScheme); err != nil {
		panic(err)
	}
	Registry.RegisterConversionFunctions(srcv1alpha1.WorkspaceGroupVersionKind, workspaceResources, nil)
}
