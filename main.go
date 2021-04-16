package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

//go:embed alpine.tar
var alpineImage []byte

const alpineImageTag = "public.ecr.aws/bacardi/alpine:3.13.0"

func main() {
	fmt.Println("Loading image...")

	dockerCli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	input := bytes.NewReader(alpineImage)
	imageLoadResponse, err := dockerCli.ImageLoad(context.Background(), input, true)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(imageLoadResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

	ctx, cancelFun := context.WithCancel(context.Background())
	defer cancelFun()

	config := &container.Config{
		Image: alpineImageTag,
		// Cmd:   []string{"cat", "/etc/hosts"},
		Tty: true,
	}

	resp, err := dockerCli.ContainerCreate(ctx, config, nil, nil, "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting container...")
	if err = dockerCli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatal(err)
	}

	command := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"echo", "hello world from a container!"},
	}

	execID, err := dockerCli.ContainerExecCreate(ctx, resp.ID, command)
	if err != nil {
		log.Fatal(err)
	}

	execAttachResp, err := dockerCli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		log.Fatal(err)
	}

	err = dockerCli.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		log.Fatal(err)
	}

	content, _, _ := execAttachResp.Reader.ReadLine()
	fmt.Println(string(content))


	fmt.Println("Stoping container...")
	timeoutValue := time.Duration(10) * time.Second
	timeout := &timeoutValue
	err = dockerCli.ContainerStop(ctx, resp.ID, timeout)
	if err != nil {
		log.Fatal(err)
	}
}
