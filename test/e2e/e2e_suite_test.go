package e2e_test

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

var _ = Describe("e2e", func() {
	It("", func() {
		bernyd := NewBernyd()
		bernyd.Start()
		bernyd.Wait()
	})
})

type Bernyd struct {
	cmd *exec.Cmd
}

func NewBernyd() *Bernyd {
	cmd := exec.Command("../../build/bernyd")

	stderr, err := cmd.StderrPipe()
	Expect(err).To(Succeed())
	stdout, err := cmd.StdoutPipe()
	Expect(err).To(Succeed())
	go LogWithPrefix(stderr, "bernyd# ")
	go LogWithPrefix(stdout, "bernyd> ")

	return &Bernyd{
		cmd: cmd,
	}
}

func (b *Bernyd) Start() {
	Expect(b.cmd.Start()).To(Succeed())
}

func (b *Bernyd) Wait() {
	err := b.cmd.Wait()
	if _, ok := err.(*exec.ExitError); ok {
		return
	}
	Expect(err).To(Succeed())
}

func LogWithPrefix(in io.Reader, prefix string) error {
	defer GinkgoRecover()

	buf := bufio.NewReader(in)
	for {
		line, _, err := buf.ReadLine()
		if err == io.ErrUnexpectedEOF || err == io.EOF || err == io.ErrClosedPipe || err == os.ErrClosed {
			return nil
		} else {
			Expect(err).To(Succeed())
		}
		_, err = GinkgoWriter.Write(append([]byte(prefix), line...))
		Expect(err).To(Succeed())
	}
}
