package gocd

import "time"

// PipelineConfig a GoCD structure that contains a pipeline and a group
type PipelineConfig struct {
	Group    string   `json:"group"`
	Pipeline Pipeline `json:"pipeline"`
}

// EnvironmentVariables a GoCD structure for an environment variable
type EnvironmentVariables struct {
	Secure         bool   `json:"secure"`
	Name           string `json:"name"`
	Value          string `json:"value,omitempty"`
	EncryptedValue string `json:"encrypted_value,omitempty"`
}

// Pipeline a GoCD structure that represents a pipeline
type Pipeline struct {
	LabelTemplate         string                 `json:"label_template"`
	EnablePipelineLocking bool                   `json:"enable_pipeline_locking"`
	Name                  string                 `json:"name"`
	Template              interface{}            `json:"template"`
	Parameters            []interface{}          `json:"parameters"`
	EnvironmentVariables  []EnvironmentVariables `json:"environment_variables"`
	Materials             []struct {
		Type       string `json:"type"`
		Attributes struct {
			URL             string      `json:"url"`
			Destination     string      `json:"destination"`
			Filter          interface{} `json:"filter"`
			InvertFilter    bool        `json:"invert_filter"`
			Name            interface{} `json:"name"`
			AutoUpdate      bool        `json:"auto_update"`
			Branch          string      `json:"branch"`
			SubmoduleFolder interface{} `json:"submodule_folder"`
			ShallowClone    bool        `json:"shallow_clone"`
		} `json:"attributes"`
	} `json:"materials"`
	Stages []struct {
		Name                  string `json:"name"`
		FetchMaterials        bool   `json:"fetch_materials"`
		CleanWorkingDirectory bool   `json:"clean_working_directory"`
		NeverCleanupArtifacts bool   `json:"never_cleanup_artifacts"`
		Approval              struct {
			Type          string `json:"type"`
			Authorization struct {
				Roles []interface{} `json:"roles"`
				Users []interface{} `json:"users"`
			} `json:"authorization"`
		} `json:"approval"`
		EnvironmentVariables []EnvironmentVariables `json:"environment_variables"`
		Jobs                 []struct {
			Name                 string                 `json:"name"`
			RunInstanceCount     interface{}            `json:"run_instance_count"`
			Timeout              string                 `json:"timeout"`
			EnvironmentVariables []EnvironmentVariables `json:"environment_variables"`
			Resources            []interface{}          `json:"resources"`
			Tasks                []struct {
				Type       string `json:"type"`
				Attributes struct {
					RunIf            []string    `json:"run_if"`
					OnCancel         interface{} `json:"on_cancel"`
					Command          string      `json:"command"`
					Arguments        []string    `json:"arguments"`
					WorkingDirectory string      `json:"working_directory"`
				} `json:"attributes"`
			} `json:"tasks"`
			Tabs       []interface{} `json:"tabs"`
			Artifacts  []interface{} `json:"artifacts"`
			Properties interface{}   `json:"properties"`
		} `json:"jobs"`
	} `json:"stages"`
	TrackingTool interface{} `json:"tracking_tool"`
	Timer        interface{} `json:"timer"`
}

// Server a representstion of a GoCD server
type Server struct {
	Host     string
	Port     string
	User     string
	Password string
	Timeout  time.Duration
}
