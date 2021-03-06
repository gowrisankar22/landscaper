// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package installations

import (
	"context"

	"github.com/gardener/component-cli/ociclient"
	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/mandelsoft/vfs/pkg/osfs"
	corev1 "k8s.io/api/core/v1"

	"github.com/gardener/component-cli/ociclient/credentials"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	componentsregistry "github.com/gardener/landscaper/pkg/landscaper/registry/components"
	"github.com/gardener/landscaper/pkg/utils"
)

// SetupRegistries sets up components and blueprints registries for the current reconcile
func (a *actuator) SetupRegistries(ctx context.Context, pullSecrets []lsv1alpha1.ObjectReference, installation *lsv1alpha1.Installation) error {
	// resolve all pull secrets
	secrets, err := a.resolveSecrets(ctx, pullSecrets)
	if err != nil {
		return err
	}

	if a.lsConfig.Registry.Local != nil {
		componentsOCIRegistry, err := componentsregistry.NewLocalClient(a.Log(), a.lsConfig.Registry.Local.RootPath)
		if err != nil {
			return err
		}
		if err := a.componentsRegistryMgr.Set(componentsOCIRegistry); err != nil {
			return err
		}
	}

	// always add a oci client to support unauthenticated requests
	ociConfigFiles := make([]string, 0)
	if a.lsConfig.Registry.OCI != nil {
		ociConfigFiles = a.lsConfig.Registry.OCI.ConfigFiles
	}
	ociKeyring, err := credentials.NewBuilder(a.Log()).DisableDefaultConfig().
		WithFS(osfs.New()).
		FromConfigFiles(ociConfigFiles...).
		FromPullSecrets(secrets...).
		Build()
	if err != nil {
		return err
	}
	ociClient, err := ociclient.NewClient(a.Log(),
		utils.WithConfiguration(a.lsConfig.Registry.OCI),
		ociclient.WithResolver{Resolver: ociKeyring},
		ociclient.WithCache{Cache: a.componentsRegistryMgr.SharedCache()},
	)
	if err != nil {
		return err
	}

	var inlineCd *cdv2.ComponentDescriptor = nil
	if installation.Spec.ComponentDescriptor != nil {
		inlineCd = installation.Spec.ComponentDescriptor.Inline
	}

	componentsOCIRegistry, err := componentsregistry.NewOCIRegistryWithOCIClient(ociClient, inlineCd)
	if err != nil {
		return err
	}
	if err := a.componentsRegistryMgr.Set(componentsOCIRegistry); err != nil {
		return err
	}

	return nil
}

func (a *actuator) resolveSecrets(ctx context.Context, secretRefs []lsv1alpha1.ObjectReference) ([]corev1.Secret, error) {
	secrets := make([]corev1.Secret, len(secretRefs))
	for i, secretRef := range secretRefs {
		secret := corev1.Secret{}
		// todo: check for cache
		if err := a.Client().Get(ctx, secretRef.NamespacedName(), &secret); err != nil {
			return nil, err
		}
		secrets[i] = secret
	}
	return secrets, nil
}
