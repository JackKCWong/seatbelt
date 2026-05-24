package scanner

type SecretDetector interface {
	Detect(prompt string) ([]Finding, error)
	Name() string
}

type Finding struct {
	Type     string
	Value    string
	Location string
}

type Registry struct {
	detectors []SecretDetector
}

func NewRegistry(detectors ...SecretDetector) *Registry {
	return &Registry{detectors: detectors}
}

func (r *Registry) Detect(prompt string) ([]Finding, error) {
	var allFindings []Finding
	for _, d := range r.detectors {
		findings, err := d.Detect(prompt)
		if err != nil {
			return nil, err
		}
		allFindings = append(allFindings, findings...)
	}
	return allFindings, nil
}