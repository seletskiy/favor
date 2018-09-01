package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/reconquest/karma-go"
)

func pick(program string, trees []*Tree, items []*ScanItem) (*ScanItem, error) {
	cmd := exec.Command(program)

	cmd.Stderr = os.Stderr
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to pipe stdin for picker",
		)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to start picker",
		)
	}

	for _, item := range items {
		_, err := stdin.Write(
			[]byte(
				item.tree.Name + ": " + item.dir + "\n",
			),
		)
		if err != nil {
			break
		}
	}

	stdin.Close()

	contents, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to read picker stdout",
		)
	}

	result := strings.TrimSpace((string(contents)))

	err = cmd.Wait()
	if err != nil {
		return nil, karma.Format(
			err,
			"picker execution failed",
		)
	}

	parts := strings.Split(result, ":")
	if len(parts) != 2 {
		return nil, karma.Describe("output", result).Format(
			err,
			"picker '%s' returned invalid output, "+
				"expected to get format 'name: dir'",
			program,
		)
	}

	name := strings.TrimSpace(parts[0])
	dir := strings.TrimSpace(parts[1])

	for _, item := range items {
		if item.tree.Name == name && item.dir == dir {
			return item, nil
		}
	}

	return nil, karma.Describe("name", name).
		Describe("dir", dir).
		Format(
			nil,
			"invalid picker output: unexpected choose",
		)
}
