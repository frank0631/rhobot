package gocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// PipelineConfig a GoCD structure that contains a pipeline and a group
type PipelineConfig struct {
	Group    string   `json:"group"`
	Pipeline Pipeline `json:"pipeline"`
}

// Pipeline a GoCD structure that represents a pipeline
type Pipeline struct {
	LabelTemplate         string        `json:"label_template"`
	EnablePipelineLocking bool          `json:"enable_pipeline_locking"`
	Name                  string        `json:"name"`
	Template              interface{}   `json:"template"`
	Parameters            []interface{} `json:"parameters"`
	EnvironmentVariables  []struct {
		Secure bool   `json:"secure"`
		Name   string `json:"name"`
		Value  string `json:"value"`
	} `json:"environment_variables"`
	Materials []struct {
		Type       string `json:"type"`
		Attributes struct {
			URL             string      `json:"url"`
			Destination     string      `json:"destination"`
			Filter          interface{} `json:"filter"`
			Name            interface{} `json:"name"`
			AutoUpdate      bool        `json:"auto_update"`
			Branch          string      `json:"branch"`
			SubmoduleFolder interface{} `json:"submodule_folder"`
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
		EnvironmentVariables []interface{} `json:"environment_variables"`
		Jobs                 []struct {
			Name                 string        `json:"name"`
			RunInstanceCount     interface{}   `json:"run_instance_count"`
			Timeout              interface{}   `json:"timeout"`
			EnvironmentVariables []interface{} `json:"environment_variables"`
			Resources            []interface{} `json:"resources"`
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

func unmarshalPipeline(data []byte) (pipeline Pipeline, err error) {
	err = json.Unmarshal(data, &pipeline)
	if err != nil {
		return pipeline, err
	}
	return pipeline, nil
}

func unmarshalPipelineConfig(data []byte) (pipelineConfig PipelineConfig, err error) {
	err = json.Unmarshal(data, &pipelineConfig)
	if err != nil {
		return pipelineConfig, err
	}
	return pipelineConfig, nil
}

// ReadPipelineJSONFromFile reads a GoCD structure from a json file
func ReadPipelineJSONFromFile(path string) (pipeline Pipeline, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return pipeline, err
	}
	return unmarshalPipeline(data)
}

func pipelineConfigPUT(gocdURL string, pipeline Pipeline, etag string) (pipelineResult Pipeline, err error) {

	pipelineName := pipeline.Name
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	payloadBytes, err := json.Marshal(pipeline)
	if err != nil {
		return pipelineResult, err
	}
	payloadBody := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("PUT", gocdURL+"/go/api/admin/pipelines/"+pipelineName, payloadBody)
	if err != nil {
		return pipelineResult, err
	}
	user := os.Getenv("GOCDUSER")
	pass := os.Getenv("GOCDPASSWORD")
	req.SetBasicAuth(user, pass)
	req.Header.Set("Accept", "application/vnd.go.cd.v1+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("If-Match", etag)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return pipelineResult, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pipelineResult, err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return pipelineResult, err
	}
	//fmt.Println("pipelineConfigPOST:", string(prettyJSON.Bytes()))
	pipelineResult, err = unmarshalPipeline(body)
	return pipelineResult, err
}

func pipelineConfigPOST(gocdURL string, pipelineConfig PipelineConfig) (pipeline Pipeline, err error) {

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	payloadBytes, err := json.Marshal(pipelineConfig)
	if err != nil {
		return pipeline, err
	}
	payloadBody := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", gocdURL+"/go/api/admin/pipelines", payloadBody)
	if err != nil {
		return pipeline, err
	}
	user := os.Getenv("GOCDUSER")
	pass := os.Getenv("GOCDPASSWORD")
	req.SetBasicAuth(user, pass)
	req.Header.Set("Accept", "application/vnd.go.cd.v1+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return pipeline, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pipeline, err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return pipeline, err
	}
	//fmt.Println("pipelineConfigPOST:", string(prettyJSON.Bytes()))

	return unmarshalPipeline(body)
}

func pipelineGET(gocdURL string, pipelineName string) (pipeline Pipeline, etag string, err error) {

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	req, err := http.NewRequest("GET", gocdURL+"/go/api/admin/pipelines/"+pipelineName, nil)
	if err != nil {
		return pipeline, etag, err
	}
	user := os.Getenv("GOCDUSER")
	pass := os.Getenv("GOCDPASSWORD")
	req.SetBasicAuth(user, pass)
	req.Header.Set("Accept", "application/vnd.go.cd.v1+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return pipeline, etag, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pipeline, etag, err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return pipeline, etag, err
	}
	//fmt.Println("pipelineConfigGET:", string(prettyJSON.Bytes()))
	etag = resp.Header.Get("ETag")
	pipeline, err = unmarshalPipeline(body)
	return pipeline, etag, err
}

// Push takes a pipeline from a file and sends it to GoCD
func Push(gocdURL string, path string, group string) {

	pipeline, err := ReadPipelineJSONFromFile(path)
	check(err)

	etag := Exist(gocdURL, pipeline.Name)
	if etag == "" {
		pipelineConfig := PipelineConfig{group, pipeline}
		_, err = pipelineConfigPOST(gocdURL, pipelineConfig)
		check(err)
	} else {
		_, err = pipelineConfigPUT(gocdURL, pipeline, etag)
		check(err)
	}

}

// Pull reads pipeline from a file, finds it on GoCD, and updates the file
func Pull(gocdURL string, path string) {

	pipeline, err := ReadPipelineJSONFromFile(path)
	check(err)
	name := pipeline.Name
	Clone(gocdURL, path, name)

}

// Exist checks if a pipeline of a given name exist, returns it's etag or an empty string
func Exist(gocdURL string, name string) (etag string) {

	_, etag, err := pipelineGET(gocdURL, name)
	check(err)
	return etag
}

// Clone finds a pipeline by name on GoCD and saves it to a file
func Clone(gocdURL string, path string, name string) {

	pipelineFetched, _, err := pipelineGET(gocdURL, name)
	check(err)
	pipelineJSON, _ := json.MarshalIndent(pipelineFetched, "", "    ")
	err = ioutil.WriteFile(path, pipelineJSON, 0666)
	check(err)
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}