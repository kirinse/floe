package exe

import (
	"bufio"
	"io"
	"os/exec"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()

	var output []string
	out := make(chan string)
	rangeDone := make(chan bool)
	go func() {
		for t := range out {
			output = append(output, t)
		}
		rangeDone <- true
	}()

	status := Run(&tLog{t: t}, out, nil, ".", "echo", "hello world")

	if status != 0 {
		t.Error("echo failed", status)
	}
	<-rangeDone

	if output[2] != "hello world" {
		t.Error("bad output", output)
	}

	// confirm bad command fails no command found
	out = make(chan string, 100)
	status = Run(&tLog{t: t}, out, nil, "", "echop", `hello world`)
	if status != 1 {
		t.Error("status should have been 1", status)
	}
}

func TestRunOutput(t *testing.T) {
	t.Parallel()

	out, status := RunOutput(&tLog{t: t}, nil, "", "bash", "-c", `echo "hello world"`)
	if status != 0 {
		t.Error("echo failed", status)
	}
	for _, l := range out {
		t.Log(">>", l)
	}
	if out[2] != "hello world" {
		t.Errorf("bad output >%s<", out[2])
	}
}

func TestRunLongOutput(t *testing.T) {
	t.Parallel()

	out, status := RunOutput(&tLog{t: t}, nil, "", "bash", "-c", `for i in {1..50}; do echo "hello line number $i"; done`)
	if status != 0 {
		t.Error("echo failed", status)
	}

	for _, o := range out {
		t.Log(o)
	}
	if len(out) != 52 {
		t.Errorf("bad output: %d", len(out))
	}
}

type tLog struct {
	t *testing.T
}

func (l *tLog) Info(args ...interface{}) {
	args = append([]interface{}{"INFO"}, args...)
	l.t.Log(args...)
}
func (l *tLog) Debug(args ...interface{}) {
	args = append([]interface{}{"DEBUG"}, args...)
	l.t.Log(args...)
}
func (l *tLog) Error(args ...interface{}) {
	args = append([]interface{}{"ERROR"}, args...)
	l.t.Log(args...)
}
func (l *tLog) Infof(format string, args ...interface{}) {
	l.t.Logf("INFO - "+format, args...)
}

func TestPlay(t *testing.T) {

	eCmd := exec.Command("bash", "-c", "export")

	pr, pw := io.Pipe()
	eCmd.Stdout = pw
	eCmd.Stderr = pw

	var output []string
	scanDone := make(chan bool)
	go func() {
		scanner := bufio.NewScanner(pr)
		for scanner.Scan() {
			t := scanner.Text()
			output = append(output, t)
		}
		if e := scanner.Err(); e != nil {
			output = append(output, "scanning output failed with: "+e.Error())
		}
		scanDone <- true
	}()

	err := eCmd.Start()
	if err != nil {
		t.Error(err)
		return
	}

	err = eCmd.Wait()
	if err != nil {
		t.Error(err)
	}

	pr.Close()
	<-scanDone
	println(len(output))
}
