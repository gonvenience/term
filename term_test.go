// Copyright © 2019 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package term_test

import (
	"bytes"
	"io"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gonvenience/term"
)

func WithCustomEnvVars(envVars map[string]string, f func()) {
	var tmp = make(map[string]string, len(envVars))
	for key, value := range envVars {
		tmp[key] = os.Getenv(key)
		os.Setenv(key, value)
	}

	defer func() {
		for key, value := range tmp {
			os.Setenv(key, value)
		}
	}()

	f()
}

func CaptureStdout(f func()) string {
	r, w, err := os.Pipe()
	Expect(err).ToNot(HaveOccurred())

	tmp := os.Stdout
	defer func() {
		os.Stdout = tmp
	}()

	os.Stdout = w
	f()
	w.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	Expect(err).ToNot(HaveOccurred())

	return buf.String()
}

var _ = Describe("term package", func() {
	Context("detection functions that base on environment variables", func() {
		It("should return that it is a dumb terminal, if the environment is set to dumb", func() {
			WithCustomEnvVars(map[string]string{"TERM": "dumb"}, func() {
				Expect(IsDumbTerminal()).To(BeTrue())
			})
		})

		It("should detect truecolor support", func() {
			WithCustomEnvVars(map[string]string{"COLORTERM": "truecolor"}, func() {
				Expect(IsTrueColor()).To(BeTrue())
			})

			WithCustomEnvVars(map[string]string{"COLORTERM": "24bit"}, func() {
				Expect(IsTrueColor()).To(BeTrue())
			})

			WithCustomEnvVars(map[string]string{"COLORTERM": "foobar"}, func() {
				Expect(IsTrueColor()).To(BeFalse())
			})

			WithCustomEnvVars(map[string]string{"COLORTERM": ""}, func() {
				Expect(IsTrueColor()).To(BeFalse())
			})
		})
	})

	Context("Concourse/Garden specific checks", func() {
		It("should at least not panic when trying to detect the Garden init process", func() {
			// Please note: Not sure why, but if one would run this test case
			// in Concourse, it would fail. So, please do not.
			Expect(IsGardenContainer()).To(BeFalse())
		})
	})

	Context("cursor hide and show convenience", func() {
		It("should print the right hide sequence to the Stdout", func() {
			Expect(CaptureStdout(HideCursor)).To(Equal("\x1b[?25l"))
		})

		It("should print the right show sequence to the Stdout", func() {
			Expect(CaptureStdout(ShowCursor)).To(Equal("\x1b[?25h"))
		})
	})
})
