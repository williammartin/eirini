package recipe_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "code.cloudfoundry.org/eirini/recipe"
)

var _ = Describe("Cfclient", func() {

	var client CfClient

	JustBeforeEach(func() {
		client = CfClient{}
	})

	FIt("works", func() {
		Expect(client.PushDroplet("", "")).ToNot(HaveOccurred())
	})
})
