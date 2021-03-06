package templates_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/templates"
)

var _ = Describe("Tokenizer", func() {
	Describe("Next", func() {
		Context("when no variable", func() {
			It("can tokenize", func() {
				t := templates.NewVariableTokenizer("abc 123")
				token := t.Next()
				Expect(token).ToNot(BeNil())
				Expect(token.TokenType).To(Equal(templates.VariableAstText))
				Expect(t.Next()).To(BeNil())
			})
		})
		Context("when single variable", func() {
			It("can tokenize", func() {
				t := templates.NewVariableTokenizer("((test))")

				token := t.Next()
				Expect(token).ToNot(BeNil())
				Expect(token.TokenType).To(Equal(templates.VariableAstOpen))

				token = t.Next()
				Expect(token).ToNot(BeNil())
				Expect(token.TokenType).To(Equal(templates.VariableAstText))
				Expect(token.Capture).To(Equal("test"))

				token = t.Next()
				Expect(token).ToNot(BeNil())
				Expect(token.TokenType).To(Equal(templates.VariableAstClose))

				Expect(t.Next()).To(BeNil())
			})
		})
		Context("when nested", func() {
			It("can tokenize", func() {
				t := templates.NewVariableTokenizer("((test((nested))))")

				token := t.Next()
				Expect(token).ToNot(BeNil())
				Expect(token.TokenType).To(Equal(templates.VariableAstOpen))

				token = t.Next()
				Expect(token).ToNot(BeNil())
				Expect(token.TokenType).To(Equal(templates.VariableAstText))
				Expect(token.Capture).To(Equal("test"))

				token = t.Next()
				Expect(token).ToNot(BeNil())
				Expect(token.TokenType).To(Equal(templates.VariableAstOpen))

				token = t.Next()
				Expect(token).ToNot(BeNil())
				Expect(token.TokenType).To(Equal(templates.VariableAstText))
				Expect(token.Capture).To(Equal("nested"))

				token = t.Next()
				Expect(token).ToNot(BeNil())
				Expect(token.TokenType).To(Equal(templates.VariableAstClose))

				token = t.Next()
				Expect(token).ToNot(BeNil())
				Expect(token.TokenType).To(Equal(templates.VariableAstClose))
			})
		})
	})
	Describe("Peek", func() {
		It("does not consume token", func() {
			t := templates.NewVariableTokenizer("test")

			Expect(t.Peek()).ToNot(BeNil())
			Expect(t.Peek()).ToNot(BeNil())
			Expect(t.Next()).ToNot(BeNil())
			Expect(t.Peek()).To(BeNil())
		})
	})
})
