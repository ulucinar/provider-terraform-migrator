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
	"fmt"
	"strings"

	srcapis "github.com/crossplane-contrib/provider-aws/apis"
	srcv1beta1 "github.com/crossplane-contrib/provider-aws/apis/ec2/v1beta1"
	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	v1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	"github.com/pkg/errors"
	targetv1beta1 "github.com/upbound/provider-aws/apis/ec2/v1beta1"
	"github.com/upbound/upjet/pkg/migration"
)

func subnetResources(mg resource.Managed) ([]resource.Managed, error) {
	source := mg.(*srcv1beta1.Subnet)
	target := &targetv1beta1.Subnet{}
	if _, err := migration.CopyInto(source, target, targetv1beta1.Subnet_GroupVersionKind, "spec.forProvider.tags"); err != nil {
		return nil, errors.Wrap(err, "failed to copy source into target")
	}
	target.Spec.ForProvider.Tags = make(map[string]*string, len(source.Spec.ForProvider.Tags))
	for _, t := range source.Spec.ForProvider.Tags {
		v := t.Value
		target.Spec.ForProvider.Tags[t.Key] = &v
	}
	return []resource.Managed{
		target,
	}, nil
}

func subnetComposedTemplates(cmp v1.ComposedTemplate, convertedBase ...*v1.ComposedTemplate) error {
	for i, cb := range convertedBase {
		for j, p := range cb.Patches {
			if p.ToFieldPath == nil || !strings.HasPrefix(*p.ToFieldPath, "spec.forProvider.tags") {
				continue
			}
			u, err := migration.FromRawExtension(cmp.Base)
			if err != nil {
				return errors.Wrap(err, "failed to convert ComposedTemplate with Subnet base")
			}
			paved := fieldpath.Pave(u.Object)
			key, err := paved.GetString(strings.ReplaceAll(*p.ToFieldPath, ".value", ".key"))
			if err != nil {
				return errors.Wrap(err, "failed to get value from paved")
			}
			s := fmt.Sprintf(`spec.forProvider.tags["%s"]`, key)
			convertedBase[i].Patches[j].ToFieldPath = &s
		}
	}
	return nil
}

func init() {
	if err := Registry.AddToScheme(srcapis.AddToScheme); err != nil {
		panic(err)
	}
	Registry.RegisterConversionFunctions(srcv1beta1.SubnetGroupVersionKind, subnetResources, subnetComposedTemplates)
}
