package lunarlauncher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const launchMetaUrl = "https://api.lunarclientprod.com/launcher/launch"

type LaunchMeta struct {
	Success        bool `json:"success"`
	LaunchTypeData struct {
		Artifacts []struct {
			Name string `json:"name"`
			Sha1 string `json:"sha1"`
			URL  string `json:"url"`
			Type string `json:"type"`
		} `json:"artifacts"`
		MainClass string `json:"mainClass"`
	} `json:"launchTypeData"`
	Licenses []struct {
		File string `json:"file"`
		URL  string `json:"url"`
		Sha1 string `json:"sha1"`
	} `json:"licenses"`
	Textures struct {
		IndexURL  string `json:"indexUrl"`
		IndexSha1 string `json:"indexSha1"`
		BaseURL   string `json:"baseUrl"`
	} `json:"textures"`
	Jre struct {
		Download struct {
			URL       string `json:"url"`
			Extension string `json:"extension"`
		} `json:"download"`
		ExecutablePathInArchive []string    `json:"executablePathInArchive"`
		CheckFiles              [][]string  `json:"checkFiles"`
		ExtraArguments          []string    `json:"extraArguments"`
		JavawDownload           interface{} `json:"javawDownload"`
		JavawExeChecksum        interface{} `json:"javawExeChecksum"`
		JavaExeChecksum         string      `json:"javaExeChecksum"`
	} `json:"jre"`
}

func FetchLaunchMeta(mcVersion McVersion, branch string) (LaunchMeta, error) {
	type LaunchRequest struct {
		Hwid            string `json:"hwid"`
		Os              string `json:"os"`
		Arch            string `json:"arch"`
		LauncherVersion string `json:"launcher_version"`
		Version         string `json:"version"`
		Branch          string `json:"branch"`
		LaunchType      string `json:"launch_type"`
		Classifier      string `json:"classifier"`
	}

	launchRequest := LaunchRequest{
		Hwid: "private",
		Os: "darwin",
		// todo: support apple chip
		Arch: "x64",
		LauncherVersion: "sandboxexec_LunarClient",
		Version: string(mcVersion),
		Branch: branch,
		LaunchType: "essa",
		Classifier: "sodium",
	}
	
	var launchMeta LaunchMeta

	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(launchRequest)
	if err != nil {
		return launchMeta, fmt.Errorf("encode launch request: %w", err)
	}
	r, err := http.Post(launchMetaUrl, "application/json", buffer)
	if err != nil {
		return launchMeta, fmt.Errorf("http get: %w", err)
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err == nil {
			return launchMeta, fmt.Errorf("invalid response code '%d' (%s)", r.StatusCode, string(bodyBytes))
		} else {
			return launchMeta, fmt.Errorf("invalid response code '%d'", r.StatusCode)
		}
	}
	modelDecoder := json.NewDecoder(r.Body)
	err = modelDecoder.Decode(&launchMeta)
	if err != nil {
		return launchMeta, fmt.Errorf("launch meta decode: %w", err)
	}

	return launchMeta, nil
}

