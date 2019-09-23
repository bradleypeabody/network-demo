// +build wasm

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/vugu/vugu"
)

func main() {

	mountPoint := flag.String("mount-point", "#vugu_mount_point", "The query selector for the mount point for the root component, if it is not a full HTML component")
	flag.Parse()

	fmt.Printf("Entering main(), -mount-point=%q\n", *mountPoint)
	defer fmt.Printf("Exiting main()\n")

	rootBuilder := &Root{}

	buildEnv, err := vugu.NewBuildEnv()
	if err != nil {
		log.Fatal(err)
	}

	renderer, err := vugu.NewJSRenderer(*mountPoint)
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Release()

	// STARTUP HTTP REQUEST CODE ADDED
	go func() {

		resp, err := http.Get("/api/startup")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		env := renderer.EventEnv()
		env.Lock()
		defer env.UnlockRender()
		rootBuilder.StartupData = string(body)

	}()

	for ok := true; ok; ok = renderer.EventWait() {

		buildResults := buildEnv.RunBuild(rootBuilder)

		err = renderer.Render(buildResults)
		if err != nil {
			panic(err)
		}
	}

}
