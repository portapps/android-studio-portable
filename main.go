//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"
	"strings"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/utl"
)

var (
	app *portapps.App
)

const (
	vmOptionsFile = "studio.vmoptions"
)

func init() {
	var err error

	// Init app
	if app, err = portapps.New("android-studio-portable", "Android Studio"); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "bin", "studio64.exe")
	app.WorkingDir = utl.PathJoin(app.AppPath, "bin")

	// override idea.properties
	studioPropContent := strings.Replace(`# DO NOT EDIT! AUTOMATICALLY GENERATED BY PORTAPPS.
idea.config.path={{ DATA_PATH }}/config
idea.system.path={{ DATA_PATH }}/system
idea.plugins.path={{ DATA_PATH }}/plugins
idea.log.path={{ DATA_PATH }}/log`, "{{ DATA_PATH }}", utl.FormatUnixPath(app.DataPath), -1)

	studioPropPath := utl.PathJoin(app.DataPath, "idea.properties")
	if err := utl.CreateFile(studioPropPath, studioPropContent); err != nil {
		log.Fatal().Err(err).Msg("Cannot write idea.properties")
	}

	// https://developer.android.com/studio/command-line/variables
	os.Setenv("ANDROID_HOME", utl.PathJoin(app.DataPath, "sdk"))
	os.Setenv("ANDROID_SDK_ROOT", utl.PathJoin(app.DataPath, "sdk"))
	os.Setenv("ANDROID_SDK_HOME", utl.PathJoin(app.DataPath, ".android"))
	os.Setenv("GRADLE_USER_HOME", utl.PathJoin(app.DataPath, ".gradle"))

	// https://developer.android.com/studio/intro/studio-config
	os.Setenv("STUDIO_PROPERTIES", studioPropPath)
	os.Setenv("STUDIO_VM_OPTIONS", utl.PathJoin(app.DataPath, vmOptionsFile))
	if !utl.Exists(utl.PathJoin(app.DataPath, vmOptionsFile)) {
		utl.CopyFile(utl.PathJoin(app.AppPath, "bin", "studio64.exe.vmoptions"), utl.PathJoin(app.DataPath, vmOptionsFile))
	} else {
		utl.CopyFile(utl.PathJoin(app.DataPath, vmOptionsFile), utl.PathJoin(app.AppPath, "bin", "studio64.exe.vmoptions"))
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
