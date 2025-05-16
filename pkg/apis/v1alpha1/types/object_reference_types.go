package v1alpha1

type LocalObjectReference struct {
	Group Group      `json:"group"`
	Kind  Kind       `json:"kind"`
	Name  ObjectName `json:"name"`
}

type SecretObjectReference struct {
	Group     *Group     `json:"group"`
	Kind      *Kind      `json:"kind"`
	Name      ObjectName `json:"name"`
	Namespace *Namespace `json:"namespace,omitempty"`
}
