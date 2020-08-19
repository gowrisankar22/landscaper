// Copyright 2020 Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package installations

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	lsv1alpha1 "github.com/gardener/landscaper/pkg/apis/core/v1alpha1"
	"github.com/gardener/landscaper/pkg/landscaper/operation"
)

// Context contains the visible installations of a specific installation.
// This context is later used to validate and get import data
type Context struct {
	// Parent is the installation the installation is encompassed in.
	// Parents are handled separately as installation have access to the same imports as their parent.
	Parent *Installation

	// Siblings are installations with the same parent.
	// The installation has access to the exports of theses components
	Siblings []*Installation
}

// DetermineContext determines the visible context of a installation.
// The visible context consists of the installation's parent and siblings.
// The context is later used to validate and get imported data.
func (o *Operation) DetermineContext(ctx context.Context) (*Context, error) {
	if IsRootInstallation(o.Inst.Info) {
		// get all root object as siblings
		rootInstallations, err := o.GetRootInstallations(ctx, func(inst lsv1alpha1.Installation) bool { return inst.Name == o.Inst.Info.Name }, client.InNamespace(o.Inst.Info.Namespace))
		if err != nil {
			return nil, err
		}
		intInstallations, err := CreateInternalInstallations(ctx, o, rootInstallations...)
		if err != nil {
			return nil, err
		}
		return &Context{Siblings: intInstallations}, nil
	}

	// get the parent by owner reference
	parent, err := GetParent(ctx, o, o.Inst)
	if err != nil {
		return nil, err
	}

	// siblings are all encompassed installation of the parent installation
	subInstallations := make([]*lsv1alpha1.Installation, 0)
	for _, installationRef := range parent.Info.Status.InstallationReferences {
		if installationRef.Reference.Name == o.Inst.Info.Name {
			continue
		}
		subInst := &lsv1alpha1.Installation{}
		if err := o.Client().Get(ctx, installationRef.Reference.NamespacedName(), subInst); err != nil {
			return nil, err
		}
		subInstallations = append(subInstallations, subInst)
	}

	intSubInstallations, err := CreateInternalInstallations(ctx, o, subInstallations...)
	if err != nil {
		return nil, err
	}

	return &Context{Parent: parent, Siblings: intSubInstallations}, nil
}

// GetParent returns the parent of a installation.
// It returns nil if the installation has no parent
func GetParent(ctx context.Context, op operation.Interface, inst *Installation) (*Installation, error) {
	if IsRootInstallation(inst.Info) {
		return nil, nil
	}
	// get the parent by owner reference
	parentName := GetParentInstallationName(inst.Info)
	parent := &lsv1alpha1.Installation{}
	if err := op.Client().Get(ctx, client.ObjectKey{Name: parentName, Namespace: inst.Info.Namespace}, parent); err != nil {
		return nil, err
	}
	intParent, err := CreateInternalInstallation(ctx, op, parent)
	if err != nil {
		return nil, err
	}
	return intParent, err
}

// IsRoot returns if the current component is a root component
func (o *Operation) IsRoot() bool {
	return o.Context().Parent == nil
}
