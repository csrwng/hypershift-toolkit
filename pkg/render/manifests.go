package render

import (
	"bytes"
	"path"
	"strings"
	"text/template"

	"github.com/openshift/hypershift-toolkit/pkg/api"
	assets "github.com/openshift/hypershift-toolkit/pkg/assets"
	"github.com/openshift/hypershift-toolkit/pkg/release"
)

// RenderClusterManifests renders manifests for a hosted control plane cluster
func RenderClusterManifests(params *api.ClusterParams, pullSecretFile, outputDir string, etcd bool, vpn bool, externalOauth bool, includeRegistry bool) error {
	releaseInfo, err := release.GetReleaseInfo(params.ReleaseImage, params.OriginReleasePrefix, pullSecretFile)
	if err != nil {
		return err
	}
	ctx := newClusterManifestContext(releaseInfo.Images, releaseInfo.Versions, params, outputDir, vpn)
	ctx.setupManifests(etcd, vpn, externalOauth, includeRegistry)
	return ctx.renderManifests()
}

type clusterManifestContext struct {
	*renderContext
	userManifestFiles []string
	userManifests     map[string]string
}

func newClusterManifestContext(images, versions map[string]string, params interface{}, outputDir string, includeVPN bool) *clusterManifestContext {
	ctx := &clusterManifestContext{
		renderContext: newRenderContext(params, outputDir),
		userManifests: make(map[string]string),
	}
	ctx.setFuncs(template.FuncMap{
		"version":           versionFunc(versions),
		"imageFor":          imageFunc(images),
		"base64String":      base64StringEncode,
		"indent":            indent,
		"address":           cidrAddress,
		"mask":              cidrMask,
		"include":           includeFileFunc(params, ctx.renderContext),
		"includeVPN":        includeVPNFunc(includeVPN),
		"randomString":      randomString,
		"includeData":       includeDataFunc(),
		"trimTrailingSpace": trimTrailingSpace,
	})
	return ctx
}

func (c *clusterManifestContext) setupManifests(etcd bool, vpn bool, externalOauth bool, includeRegistry bool) {
	if etcd {
		c.etcd()
	}
	c.kubeAPIServer(vpn)
	c.kubeControllerManager()
	c.kubeScheduler()
	c.clusterBootstrap()
	c.openshiftAPIServer()
	c.openshiftControllerManager()
	if externalOauth {
		c.oauthOpenshiftServer()
	}
	if vpn {
		c.openVPN()
	}
	c.clusterVersionOperator()
	if includeRegistry {
		c.registry()
	}
	c.userManifestsBootstrapper()
	c.controlPlaneOperator()
}

func (c *clusterManifestContext) etcd() {
	c.addManifestFiles(
		"etcd/etcd-cluster-crd.yaml",
		"etcd/etcd-cluster.yaml",
		"etcd/etcd-operator-cluster-role-binding.yaml",
		"etcd/etcd-operator-cluster-role.yaml",
		"etcd/etcd-operator.yaml",
	)

}

func (c *clusterManifestContext) oauthOpenshiftServer() {
	c.addManifestFiles(
		"oauth-openshift/oauth-browser-client.yaml",
		"oauth-openshift/oauth-challenging-client.yaml",
		"oauth-openshift/oauth-server-config-configmap.yaml",
		"oauth-openshift/oauth-server-deployment.yaml",
		"oauth-openshift/oauth-server-service.yaml",
		"oauth-openshift/v4-0-config-system-branding.yaml",
		"oauth-openshift/oauth-server-sessionsecret-secret.yaml",
	)
}

func (c *clusterManifestContext) kubeAPIServer(includeVPN bool) {
	c.addManifestFiles(
		"kube-apiserver/kube-apiserver-deployment.yaml",
		"kube-apiserver/kube-apiserver-service.yaml",
		"kube-apiserver/kube-apiserver-config-configmap.yaml",
		"kube-apiserver/kube-apiserver-oauth-metadata-configmap.yaml",
	)
	if includeVPN {
		c.addManifestFiles(
			"kube-apiserver/kube-apiserver-vpnclient-config.yaml",
		)
	}
}

func (c *clusterManifestContext) kubeControllerManager() {
	c.addManifestFiles(
		"kube-controller-manager/kube-controller-manager-deployment.yaml",
		"kube-controller-manager/kube-controller-manager-config-configmap.yaml",
	)
}

func (c *clusterManifestContext) kubeScheduler() {
	c.addManifestFiles(
		"kube-scheduler/kube-scheduler-deployment.yaml",
		"kube-scheduler/kube-scheduler-config-configmap.yaml",
	)
}

func (c *clusterManifestContext) registry() {
	c.addUserManifestFiles("registry/cluster-imageregistry-config.yaml")
}

func (c *clusterManifestContext) clusterBootstrap() {
	manifests, err := assets.AssetDir("cluster-bootstrap")
	if err != nil {
		panic(err.Error())
	}
	for _, m := range manifests {
		c.addUserManifestFiles("cluster-bootstrap/" + m)
	}
}

func (c *clusterManifestContext) openshiftAPIServer() {
	c.addManifestFiles(
		"openshift-apiserver/openshift-apiserver-deployment.yaml",
		"openshift-apiserver/openshift-apiserver-service.yaml",
		"openshift-apiserver/openshift-apiserver-config-configmap.yaml",
	)
	c.addUserManifestFiles(
		"openshift-apiserver/openshift-apiserver-user-service.yaml",
		"openshift-apiserver/openshift-apiserver-user-endpoint.yaml",
	)
	apiServices := &bytes.Buffer{}
	for _, apiService := range []string{
		"v1.apps.openshift.io",
		"v1.authorization.openshift.io",
		"v1.build.openshift.io",
		"v1.image.openshift.io",
		"v1.oauth.openshift.io",
		"v1.project.openshift.io",
		"v1.quota.openshift.io",
		"v1.route.openshift.io",
		"v1.security.openshift.io",
		"v1.template.openshift.io",
		"v1.user.openshift.io"} {

		params := map[string]string{
			"APIService":                 apiService,
			"APIServiceGroup":            trimFirstSegment(apiService),
			"OpenshiftAPIServerCABundle": c.params.(*api.ClusterParams).OpenshiftAPIServerCABundle,
		}
		entry, err := c.substituteParams(params, "openshift-apiserver/service-template.yaml")
		if err != nil {
			panic(err.Error())
		}
		apiServices.WriteString(entry)
	}
	c.addUserManifest("openshift-apiserver-apiservices.yaml", apiServices.String())
}

func (c *clusterManifestContext) openshiftControllerManager() {
	c.addManifestFiles(
		"openshift-controller-manager/openshift-controller-manager-deployment.yaml",
		"openshift-controller-manager/openshift-controller-manager-config-configmap.yaml",
		"openshift-controller-manager/cluster-policy-controller-deployment.yaml",
	)
	c.addUserManifestFiles(
		"openshift-controller-manager/00-openshift-controller-manager-namespace.yaml",
		"openshift-controller-manager/openshift-controller-manager-service-ca.yaml",
	)
}

func (c *clusterManifestContext) controlPlaneOperator() {
	c.addManifestFiles(
		"control-plane-operator/cp-operator-deployment.yaml",
	)
}

func (c *clusterManifestContext) openVPN() {
	c.addManifestFiles(
		"openvpn/openvpn-server-deployment.yaml",
		"openvpn/openvpn-server-service.yaml",
		"openvpn/openvpn-ccd-configmap.yaml",
		"openvpn/openvpn-server-configmap.yaml",
	)
	c.addUserManifestFiles(
		"openvpn/openvpn-client-deployment.yaml",
		"openvpn/openvpn-client-configmap.yaml",
	)
}

func (c *clusterManifestContext) clusterVersionOperator() {
	c.addManifestFiles(
		"cluster-version-operator/cluster-version-operator-deployment.yaml",
	)
}

func (c *clusterManifestContext) userManifestsBootstrapper() {
	c.addManifestFiles(
		"user-manifests-bootstrapper/user-manifests-bootstrapper-pod.yaml",
	)
	for _, file := range c.userManifestFiles {
		data, err := c.substituteParams(c.params, file)
		if err != nil {
			panic(err.Error())
		}
		name := path.Base(file)
		params := map[string]string{
			"data": data,
			"name": userConfigMapName(name),
		}
		manifest, err := c.substituteParams(params, "user-manifests-bootstrapper/user-manifest-template.yaml")
		if err != nil {
			panic(err.Error())
		}
		c.addManifest("user-manifest-"+name, manifest)
	}

	for name, data := range c.userManifests {
		params := map[string]string{
			"data": data,
			"name": userConfigMapName(name),
		}
		manifest, err := c.substituteParams(params, "user-manifests-bootstrapper/user-manifest-template.yaml")
		if err != nil {
			panic(err.Error())
		}
		c.addManifest("user-manifest-"+name, manifest)
	}
}

func (c *clusterManifestContext) addUserManifestFiles(name ...string) {
	c.userManifestFiles = append(c.userManifestFiles, name...)
}

func (c *clusterManifestContext) addUserManifest(name, content string) {
	c.userManifests[name] = content
}

func trimFirstSegment(s string) string {
	parts := strings.Split(s, ".")
	return strings.Join(parts[1:], ".")
}

func userConfigMapName(file string) string {
	parts := strings.Split(file, ".")
	return "user-manifest-" + strings.ReplaceAll(parts[0], "_", "-")
}
