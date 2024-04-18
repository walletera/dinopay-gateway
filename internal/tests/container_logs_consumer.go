package tests

import (
    "bufio"
    "fmt"
    "os"

    "github.com/testcontainers/testcontainers-go"
)

type ContainerLogConsumer struct {
    file *os.File
}

func NewContainerLogConsumer(containerName string) *ContainerLogConsumer {
    wd, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    file, err := os.Create(fmt.Sprintf("%s/containerlogs/%s.log", wd, containerName))
    if err != nil {
        panic(err)
    }
    return &ContainerLogConsumer{
        file: file,
    }
}

func (c *ContainerLogConsumer) Accept(log testcontainers.Log) {
    w := bufio.NewWriter(c.file)
    _, err := w.Write(log.Content)
    if err != nil {
        panic(err)
    }
    err = w.Flush()
    if err != nil {
        panic(err)
    }
}
