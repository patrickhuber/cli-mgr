#!/bin/bash
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega/...
ginkgo -p -r -race -randomizeAllSpecs -randomizeSuites .