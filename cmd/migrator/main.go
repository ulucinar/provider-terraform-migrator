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

package main

import (
	"os"
	"path/filepath"

	"github.com/upbound/upjet/pkg/migration"
	"gopkg.in/alecthomas/kingpin.v2"
	"sigs.k8s.io/yaml"

	"github.com/ulucinar/migration/pkg/converters"
)

func main() {
	var (
		app      = kingpin.New(filepath.Base(os.Args[0]), "Upbound migration plan generator for migrating Kubernetes objects from community providers to official providers.").DefaultEnvars()
		planPath = app.Flag("plan-path", "Path where the generated migration plan will be stored").Short('p').Default("migration_plan.yaml").String()
		// sourcePath           = app.Flag("source-path", "Path of the root directory for the filesystem source").Short('s').Required().String()
		kubeconfigPath = app.Flag("kubeconfig", "Path of the kubernetes config file").String()
	)
	if len(*kubeconfigPath) == 0 {
		homeDir, err := os.UserHomeDir()
		kingpin.FatalIfError(err, "Failed to get user's home directory")
		*kubeconfigPath = filepath.Join(homeDir, ".kube/config")
	}
	kingpin.MustParse(app.Parse(os.Args[1:]))
	dc, err := migration.InitializeDynamicClient(*kubeconfigPath)
	kingpin.FatalIfError(err, "Failed to initialize Kubernetes dynamic client")
	source, err := migration.NewKubernetesSource(converters.Registry, dc)
	kingpin.FatalIfError(err, "Failed to initialize a Kubernetes source")
	target := migration.NewFileSystemTarget()
	pg := migration.NewPlanGenerator(converters.Registry, source, target)
	err = pg.GeneratePlan()
	kingpin.FatalIfError(err, "Failed to generate the migration plan")
	buff, err := yaml.Marshal(pg.Plan)
	kingpin.FatalIfError(err, "Failed to marshal the migration plan into YAML")
	kingpin.FatalIfError(os.WriteFile(*planPath, buff, 0o600), "Failed to store the migration plan: %s", planPath)
}
