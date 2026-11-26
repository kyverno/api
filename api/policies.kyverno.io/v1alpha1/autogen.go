package v1alpha1

type PodControllersGenerationConfiguration struct {
	// Enabled specifies whether to generate pod controller rules.
	// Optional. Defaults to "true" if not specified.
	Enabled *bool `json:"enabled,omitempty"`
	// Controllers specifies the list of pod controllers to generate rules for,
	// for example DaemonSet, Deployment, Job, StatefulSet, ReplicaSet, ReplicationController, CronJob.
	// Optional. Defaults to all supported pod controllers if not specified.
	Controllers []string `json:"controllers,omitempty"`
}

type Target struct {
	Group    string `json:"group,omitempty"`
	Version  string `json:"version"`
	Resource string `json:"resource"`
	Kind     string `json:"kind"`
}

type ValidatingPolicyAutogenStatus struct {
	Configs map[string]ValidatingPolicyAutogen `json:"configs,omitempty"`
}

type ImageValidatingPolicyAutogenStatus struct {
	Configs map[string]ImageValidatingPolicyAutogen `json:"configs,omitempty"`
}

type MutatingPolicyAutogenStatus struct {
	Configs map[string]MutatingPolicyAutogen `json:"configs,omitempty"`
}

type ValidatingPolicyAutogen struct {
	Targets []Target              `json:"targets"`
	Spec    *ValidatingPolicySpec `json:"spec"`
}

type ImageValidatingPolicyAutogen struct {
	Targets []Target                   `json:"targets"`
	Spec    *ImageValidatingPolicySpec `json:"spec"`
}

type MutatingPolicyAutogen struct {
	Targets []Target            `json:"targets"`
	Spec    *MutatingPolicySpec `json:"spec"`
}
