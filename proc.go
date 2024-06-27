package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Process struct {
	Name             string //
	State            string //
	Pid              int    //
	PPid             int    //
	Executable       string
	WorkingDirectory string //cwd
}

func newProcess(pid int) (*Process, error) {

	p := &Process{Pid: pid}
	return p, p.getProcessStat()
}

func (p *Process) getProcessStat() error {

	statFile := fmt.Sprintf("/proc/%d/stat", p.Pid)
	bs, err := os.ReadFile(statFile)
	if err != nil {
		return err
	}
	content := string(bs)
	start := strings.IndexRune(content, '(') + 1
	end := strings.IndexRune(content, ')')
	p.Name = content[start:end]
	data := content[end+2:]
	_, err = fmt.Sscanf(data, "%s %d", &p.State, &p.PPid)
	if err != nil {
		fmt.Println("err->", err, p.Name, data[0:1], data[2:3], "**", data)
	}
	readlink, _ := os.Readlink(fmt.Sprintf("/proc/%d/exe", p.Pid))
	cwd, _ := os.Readlink(fmt.Sprintf("/proc/%d/cwd", p.Pid))
	p.Executable = readlink
	p.WorkingDirectory = cwd
	return err
}

//proc/411/stat

func processList() ([]Process, error) {

	d, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer d.Close()

	results := make([]Process, 0, 50)

	for {

		names, err := d.Readdirnames(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, name := range names {

			if name[0] < '0' || name[0] > '9' {
				continue
			}
			//fmt.Println(name)
			pn, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				return nil, err
			}

			//fmt.Println(pn)
			p, err := newProcess(int(pn))
			if err != nil {
				log.Println(err)
				continue
			}
			results = append(results, *p)
		}
	}
	return results, nil
}
