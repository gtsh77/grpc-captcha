package entity

type HeathCheckResponse struct {
	ID         string `json:"runtime_id"`
	Name       string `json:"service_name"`
	Version    string `json:"version"`
	CompiledAt string `json:"compiled_at"`
}

type HealthOperableResponse struct {
	IsRdy bool `json:"is_ready"`
}
