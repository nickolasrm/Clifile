package file_test

import (
	"testing"

	"github.com/nickolasrm/clifile/internal/util/file"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "File Util Suite")
}

var _ = Describe("ReadRunes", func() {
	Context("when a valid file is passed", func() {
		It("should return a channel of runes", func() {
			ch, err := file.ReadRunes("file.go")
			Expect(err).To(BeNil())
			Expect(ch).ToNot(BeNil())

		})
		It("channel should return runes", func() {
			ch, err := file.ReadRunes("file.go")
			Expect(err).To(BeNil())
			Expect(<-ch).To(Equal('p'))
			Expect(<-ch).To(Equal('a'))
			Expect(<-ch).To(Equal('c'))
			Expect(<-ch).To(Equal('k'))
			Expect(<-ch).To(Equal('a'))
			Expect(<-ch).To(Equal('g'))
			Expect(<-ch).To(Equal('e'))
		})
		It("channel should be closed after all runes are read", func() {
			ch, err := file.ReadRunes("file.go")
			defer Expect(ch).To(BeClosed())
			Expect(err).To(BeNil())
			for run := range ch {
				Expect(run).ToNot(BeNil())
			}
		})
	})

	Context("when an invalid file is passed", func() {
		It("should return an error", func() {
			ch, err := file.ReadRunes("file.go1")
			Expect(err).ToNot(BeNil())
			Expect(ch).To(BeNil())
		})
	})
})
