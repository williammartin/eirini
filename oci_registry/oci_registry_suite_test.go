package oci_registry_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOciRegistry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OciRegistry Suite")
}
