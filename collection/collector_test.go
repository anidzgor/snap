package collection_test

import (
	. "github.com/lynxbat/pulse/collection"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Foo() bool {
	return true
}

var _ = Describe("NewCollectorByType", func() {

	Context("with \"collectd\" type and config", func() {
		It("returns a Collector with CollectorType() = \"collectd\"", func() {
			c := NewCollectorByType("collectd", NewCollectDConfig())
			Expect(c.CollectorType()).To(Equal("collectd"))
		})
	})

	Context("with \"facter\" type and config", func() {
		It("returns a Collector with CollectorType() = \"facter\"", func() {
			c := NewCollectorByType("facter", NewFacterConfig())
			Expect(c.CollectorType()).To(Equal("facter"))
		})
	})

	// PContext("with \"libcontainer\" type and config", func() {
	// 	It("returns a Collector with CollectorType() = \"facter\"", func() {
	// 		c := NewCollectorByType("libcontainer", NewCollectDConfig())
	// 		Expect(c.CollectorType()).To(Equal("libcontainer"))
	// 	})
	// })
})