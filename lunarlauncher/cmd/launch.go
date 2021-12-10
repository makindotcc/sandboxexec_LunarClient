package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/makindotcc/lunarlauncher"
)

func downloadTextures(texturesBaseUrl string, textures []lunarlauncher.TextureMeta) {
	for processedCount, texture := range textures {
		log.Printf("Downloading texture (%d/%d): %s\n", (processedCount + 1),
			len(textures), texture)
		err := lunarlauncher.DownloadTexture(texturesBaseUrl, texture)
		if errors.Is(err, lunarlauncher.ErrAlreadyDownloaded) {
			log.Printf("Texture %v is already downloaded. Skipping...\n", texture)
		} else if err != nil {
			log.Printf("Texture %v download error: %s\n", texture, err)
		}
	}
}

func downloadLaunchTextures(launchMeta lunarlauncher.LaunchMeta) {
	textures, err := lunarlauncher.FetchTextures(launchMeta.Textures.IndexURL)
	if err != nil {
		panic(err)
	}
	downloadTextures(launchMeta.Textures.BaseURL, textures)
}

func downloadLaunchArtifacts(mcVersion lunarlauncher.McVersion,
	artifacts []lunarlauncher.Artifact) {
	for processedCount, artifact := range artifacts {
		log.Printf("Downloading artifact (%d/%d): %s\n", (processedCount + 1),
			len(artifacts), artifact)
		lunarlauncher.DownloadArtifact(mcVersion, artifact)
	}
}

func launchSandboxed(mcVersion lunarlauncher.McVersion) {
	cmd := exec.Command("sandbox-exec", "-f", "../../../lunar.sb", "java", "-Dlog4j2.formatMsgNoLookups=true", "--add-modules", "jdk.naming.dns", "--add-exports",
	"jdk.naming.dns/com.sun.jndi.dns=java.naming", "-Djna.boot.library.path=natives",
	"--add-opens", "java.base/java.io=ALL-UNNAMED", "-XstartOnFirstThread", "-Xms1024m", "-Xmx1024m", "-Djava.library.path=natives",
	"-XX:+DisableAttachMechanism", 
	"-cp", "vpatcher-prod.jar:lunar-prod-optifine.jar:lunar-libs.jar:lunar-assets-prod-1-optifine.jar:"+
	"lunar-assets-prod-2-optifine.jar:lunar-assets-prod-3-optifine.jar:OptiFine.jar",
	"com.moonsworth.lunar.patcher.LunarMain", "--version", string(mcVersion), "--accessToken", "0", "--assetIndex", string(mcVersion),
	"--userProperties", "{}", "--gameDir", ".minecraft", "--texturesDir", "../../textures", "--launcherVersion", "2.7.5",
	"--width", "1280", "--height", "720")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = path.Join("offline", string(mcVersion))
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	const mcVersion = lunarlauncher.Mc1_16

	log.Println("sandbox:", os.Environ())

	log.Println("Preparing lunar assets...")
	log.Println("Fetching launch meta...")
	launchMeta, err := lunarlauncher.FetchLaunchMeta(mcVersion, "master")
	if err != nil {
		panic(err)
	}
	log.Println("Got launch meta:", launchMeta)

	downloadLaunchTextures(launchMeta)
	downloadLaunchArtifacts(mcVersion, launchMeta.LaunchTypeData.Artifacts)

	launchSandboxed(mcVersion)
}
