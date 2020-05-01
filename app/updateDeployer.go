package app

import (
	"encoding/json"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)


func (app *App) UpdateDeployer(c *cli.Context) error {
	type cfgSystem struct{
		ExecuteFileName string
		ArchiveName string
	}
	configBySystem := cfgSystem{"dep", "dep.linux-amd64.tar.gz"}
	if runtime.GOOS == "windows" {
		configBySystem.ExecuteFileName = "dep.exe"
		configBySystem.ArchiveName = "dep.windows-amd64.exe.tar.gz"
	}

	resp, err := http.Get("https://api.github.com/repos/evgeny-klyopov/dep/releases/latest")

	if err != nil  {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil  {
		return err
	}

	var response struct{
		TagName string `json:"tag_name"`
		Assets []struct{
			Name string `json:"name"`
			BrowserDownloadUrl string `json:"browser_download_url"`
		} `json:"assets"`
	}
	err = json.Unmarshal(body, &response)

	if err != nil  {
		return err
	}

	if app.GetVersion() == response.TagName {
		app.Color.Print(app.Color.Green, "You have the latest version.")
		return nil
	}

	tpmPath := os.TempDir()
	tpmArchivePath := tpmPath + "/" + configBySystem.ArchiveName
	link := "https://github.com/evgeny-klyopov/dep/releases/download/" + response.TagName + "/" + configBySystem.ArchiveName

	var tasks Tasks
	tasks = append(tasks, task{name: "download-release", remote: false, method: func(app *App) error {
		return app.cmd("curl", "-Ls", link, "-o", tpmArchivePath)
	}})

	tasks = append(tasks, task{name: "extract-release", remote: false, method: func(app *App) error {
		cmd := exec.Command("tar", "-xvzf", configBySystem.ArchiveName)
		cmd.Dir = tpmPath

		return  cmd.Run()
	}})

	tasks = append(tasks, task{name: "update-release", remote: false, method: func(app *App) error {
		cmd := exec.Command("mv", configBySystem.ExecuteFileName, os.Args[0])
		cmd.Dir = tpmPath
		return cmd.Run()
	}})

	tasks = append(tasks, task{name: "clear-release", remote: false, method: func(app *App) error {
		cmd := exec.Command("rm", configBySystem.ArchiveName)
		cmd.Dir = tpmPath
		return cmd.Run()
	}})

	_, err = app.run(tasks)

	if err != nil {
		return err
	}

	app.Color.Print(app.Color.Green, "Deployer has been updated to version " + response.TagName)

	return nil
}