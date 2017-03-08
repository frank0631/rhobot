package gocd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// NewServerConfig Create a Server object from a config
func NewServerConfig(host string, port string, user string, password string, timeoutStr string) *Server {

	// timeout casting to seconds
	timeout := time.Duration(120 * time.Second)
	i, err := strconv.Atoi(timeoutStr)
	if err == nil {
		timeout = time.Duration(i) * time.Second
	} else {
		log.Warn("Failed to convert timeout to seconds: ", err)
	}

	return &Server{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Timeout:  timeout,
	}
}

// URL returns the host of the GoCD server
func (server Server) URL() string {
	return fmt.Sprintf("%s:%s", server.Host, server.Port)
}

// History gets the run history of a pipeline of a given name exist. returns map
func History(server *Server, name string) (latestRuns map[string]int, err error) {

	//Get pipeline history if it exist
	latestRuns, err = server.historyGET(name)
	if err != nil {
		log.Fatalf("Could not find run history for pipeline: %v", name)
	}

	return
}

// Artifact gets an Areifact from a pipeline / stage / job
func Artifact(server *Server, pipelineName string, pipelineID int, stageName string, stageID int, jobName string, artifactPath string) (fileBytes *bytes.Buffer, err error) {
	fileBytes, err = server.artifactGET(
		pipelineName, pipelineID,
		stageName, stageID,
		jobName, artifactPath)
	return
}

// Push takes a pipeline from a file and sends it to GoCD
func Push(server *Server, path string, group string) (err error) {
	localPipeline, err := readPipelineJSONFromFile(path)
	if err != nil {
		return
	}

	etag, remotePipeline, err := Exist(server, localPipeline.Name)
	if err != nil {
		log.Info(err)
	}

	Compare(localPipeline, remotePipeline, path)

	if etag == "" {
		pipelineConfig := PipelineConfig{group, localPipeline}
		_, err = server.pipelineConfigPOST(pipelineConfig)
	} else {
		_, err = server.pipelineConfigPUT(localPipeline, etag)
	}
	return
}

// Pull reads pipeline from a file, finds it on GoCD, and updates the file
func Pull(server *Server, path string) (err error) {
	localPipeline, err := readPipelineJSONFromFile(path)
	if err != nil {
		return
	}

	name := localPipeline.Name
	remotePipeline, err := Clone(server, path, name)

	Compare(localPipeline, remotePipeline, path)

	return
}

// Exist checks if a pipeline of a given name exist, returns it's etag or an empty string
func Exist(server *Server, name string) (etag string, pipeline Pipeline, err error) {
	pipeline, etag, err = server.pipelineGET(name)
	return
}

// Clone finds a pipeline by name on GoCD and saves it to a file
func Clone(server *Server, path string, name string) (pipeline Pipeline, err error) {
	pipeline, _, err = server.pipelineGET(name)
	if err != nil {
		return
	}

	err = writePipeline(path, pipeline)
	return
}

// Compare saves copies of the local and remote pipeline if different
func Compare(localPipeline Pipeline, remotePipeline Pipeline, path string) {

	if !reflect.DeepEqual(localPipeline, remotePipeline) {
		log.Warn("Local and Remote are different")

		filepath := strings.TrimSuffix(path, filepath.Ext(path))
		localBakPath := filepath + ".local.bak.json"
		remoteBakPath := filepath + ".remote.bak.json"

		log.Info("Saving Local Backup: ", localBakPath)
		errLocal := writePipeline(localBakPath, localPipeline)
		log.Info("Saving Remote Backup: ", remoteBakPath)
		errRemote := writePipeline(remoteBakPath, remotePipeline)

		if errLocal != nil {
			log.Warn("Error while writing backup for local pipeline: ", errLocal)
		}
		if errRemote != nil {
			log.Warn("Error while writing backup for local pipeline: ", errLocal)
		}

	}
}
