package api

type ClusterParams struct {
	Namespace                           string                 `json:"namespace"`
	ExternalAPIDNSName                  string                 `json:"externalAPIDNSName"`
	ExternalAPIPort                     uint                   `json:"externalAPIPort"`
	ExternalAPIIPAddress                string                 `json:"externalAPIAddress"`
	ExternalOpenVPNDNSName              string                 `json:"externalVPNDNSName"`
	ExternalOpenVPNPort                 uint                   `json:"externalVPNPort"`
	ExternalOauthPort                   uint                   `json:"externalOauthPort"`
	IdentityProviders                   string                 `json:"identityProviders"`
	ServiceCIDR                         string                 `json:"serviceCIDR"`
	NamedCerts                          []NamedCert            `json:"namedCerts,omitempty"`
	PodCIDR                             string                 `json:"podCIDR"`
	ReleaseImage                        string                 `json:"releaseImage"`
	APINodePort                         uint                   `json:"apiNodePort"`
	IngressSubdomain                    string                 `json:"ingressSubdomain"`
	OpenShiftAPIClusterIP               string                 `json:"openshiftAPIClusterIP"`
	ImageRegistryHTTPSecret             string                 `json:"imageRegistryHTTPSecret"`
	RouterNodePortHTTP                  string                 `json:"routerNodePortHTTP"`
	RouterNodePortHTTPS                 string                 `json:"routerNodePortHTTPS"`
	OpenVPNNodePort                     string                 `json:"openVPNNodePort"`
	BaseDomain                          string                 `json:"baseDomain"`
	NetworkType                         string                 `json:"networkType"`
	Replicas                            string                 `json:"replicas"`
	EtcdClientName                      string                 `json:"etcdClientName"`
	OriginReleasePrefix                 string                 `json:"originReleasePrefix"`
	OpenshiftAPIServerCABundle          string                 `json:"openshiftAPIServerCABundle"`
	CloudProvider                       string                 `json:"cloudProvider"`
	CVOSetupImage                       string                 `json:"cvoSetupImage"`
	InternalAPIPort                     uint                   `json:"internalAPIPort"`
	RouterServiceType                   string                 `json:"routerServiceType"`
	KubeAPIServerResources              []ResourceRequirements `json:"kubeAPIServerResources"`
	OpenshiftControllerManagerResources []ResourceRequirements `json:"openshiftControllerManagerResources"`
	ClusterVersionOperatorResources     []ResourceRequirements `json:"clusterVersionOperatorResources"`
	KubeControllerManagerResources      []ResourceRequirements `json:"kubeControllerManagerResources"`
	OpenshiftAPIServerResources         []ResourceRequirements `json:"openshiftAPIServerResources"`
	KubeSchedulerResources              []ResourceRequirements `json:"kubeSchedulerResources"`
	ControlPlaneOperatorResources       []ResourceRequirements `json:"controlPlaneOperatorResources"`
	OAuthServerResources                []ResourceRequirements `json:"oAuthServerResources"`
	ClusterPolicyControllerResources    []ResourceRequirements `json:"clusterPolicyControllerResources"`
	AutoApproverResources               []ResourceRequirements `json:"autoApproverResources"`
	OpenVPNClientResources              []ResourceRequirements `json:"openVPNClientResources"`
	OpenVPNServerResources              []ResourceRequirements `json:"openVPNServerResources"`
	APIServerAuditEnabled               bool                   `json:"apiServerAuditEnabled"`
	RestartDate                         string                 `json:"restartDate"`
	ControlPlaneOperatorControllers     []string               `json:"controlPlaneOperatorControllers"`
	ExtraFeatureGates                   []string               `json:"extraFeatureGates"`
	ApiserverLivenessPath               string                 `json:"apiserverLivenessPath"`
	DefaultFeatureGates                 []string
	PlatformType                        string `json:"platformType"`
}

type NamedCert struct {
	NamedCertPrefix string `json:"namedCertPrefix"`
	NamedCertDomain string `json:"namedCertDomain"`
}

type ResourceRequirements struct {
	ResourceLimit   []ResourceLimit   `json:"resourceLimit"`
	ResourceRequest []ResourceRequest `json:"resourceRequest"`
}

type ResourceLimit struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type ResourceRequest struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}
