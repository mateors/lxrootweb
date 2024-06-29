package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Process struct {
	Name             string //
	State            string //
	Pid              int    //
	PPid             int    //
	Uid              int
	Gid              int
	Groups           string
	Threads          int
	Executable       string
	WorkingDirectory string //cwd
	SocketList       []map[string]interface{}
}

func newProcess(pid int) (*Process, error) {

	p := &Process{Pid: pid}
	return p, p.getProcessStat()
}

func (p *Process) getProcessStat() error {

	kvMap := processStatus(p.Pid)
	p.Name = kvMap["Name"]
	p.State = kvMap["State"]
	p.PPid = str2int(kvMap["PPid"])

	uidSlc := strings.Fields(kvMap["Uid"])
	if len(uidSlc) == 4 {
		p.Uid = str2int(uidSlc[0])
	}
	gidSlc := strings.Fields(kvMap["Gid"])
	if len(gidSlc) == 4 {
		p.Gid = str2int(gidSlc[0])
	}
	p.Groups = kvMap["Groups"]
	p.Threads = str2int(kvMap["Threads"])

	readlink, _ := os.Readlink(fmt.Sprintf("/proc/%d/exe", p.Pid))
	cwd, _ := os.Readlink(fmt.Sprintf("/proc/%d/cwd", p.Pid))
	p.Executable = readlink
	p.WorkingDirectory = cwd
	p.SocketList = pidOpenSocketList(p.Pid)
	return err
}

// func (p *Process) getProcessStat() error {

// 	statFile := fmt.Sprintf("/proc/%d/stat", p.Pid)
// 	bs, err := os.ReadFile(statFile)
// 	if err != nil {
// 		return err
// 	}
// 	content := string(bs)
// 	start := strings.IndexRune(content, '(') + 1
// 	end := strings.IndexRune(content, ')')
// 	p.Name = content[start:end]
// 	data := content[end+2:]
// 	_, err = fmt.Sscanf(data, "%s %d", &p.State, &p.PPid)
// 	if err != nil {
// 		fmt.Println("err->", err, p.Name, data[0:1], data[2:3], "**", data)
// 	}
// 	readlink, _ := os.Readlink(fmt.Sprintf("/proc/%d/exe", p.Pid))
// 	cwd, _ := os.Readlink(fmt.Sprintf("/proc/%d/cwd", p.Pid))
// 	p.Executable = readlink
// 	p.WorkingDirectory = cwd
// 	p.SocketList = pidOpenSocketList(p.Pid)
// 	return err
// }

//proc/411/stat

func ProcessList() ([]Process, error) {

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

func cmdOutputSlc2(input string) []string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var zeroSlc = []string{}
	for scanner.Scan() {
		line := scanner.Text()
		zeroSlc = append(zeroSlc, line)
	}
	return zeroSlc
}

func hexaNumberToInteger(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	//output, err: = strconv.ParseInt(hexaNumberToInteger(hexaNumber), 16, 64)
	return numberStr
}

func hexAddressParse(hexA string) string {

	first := hexA[0:2]  //0,1
	second := hexA[2:4] //2,3
	third := hexA[4:6]  //4,5
	fourth := hexA[6:8] //6,7
	return fmt.Sprintf("%d.%d.%d.%d", hexToDecimal(fourth), hexToDecimal(third), hexToDecimal(second), hexToDecimal(first))
}

func hexToDecimal(hexa string) int {

	listeningPort, err := strconv.ParseInt(hexa, 16, 64)
	if err != nil {
		return 0
	}
	return int(listeningPort)
}

func processFileDescriptor(data string) map[string]map[string]interface{} {

	var irow = make(map[string]map[string]interface{})
	// 	data := `sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
	// 0: 00000000:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 23888 1 ffff92a00ebb8000 100 0 0 10 0
	// 1: 00000000:0016 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 20716 1 ffff92a002038000 100 0 0 10 0
	// 2: 00000000:0035 00000000:0000 0A 00000000:00000000 00:00000000 00000000   108        0 19855 1 ffff92a042db8000 100 0 0 10 0
	// 3: 00000000:01BB 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 23889 1 ffff92a00ebb88c0 100 0 0 10 0
	// 4: 0100007F:0CEA 00000000:0000 0A 00000000:00000000 00:00000000 00000000   110        0 685 1 ffff92a05d078000 100 0 0 10 0
	// `
	slc := cmdOutputSlc2(data)
	for i, line := range slc {

		if i == 0 {
			continue
		}

		fslc := strings.Fields(line)
		localA := fslc[1]
		remoteA := fslc[2]
		uid := fslc[7]
		inode := fslc[9]
		var row = make(map[string]interface{})
		var localPort int
		var localAddress string
		slc = strings.Split(localA, ":")
		if len(slc) == 2 {
			localAddress = hexAddressParse(slc[0])
			portHexa := slc[1]
			localPort = hexToDecimal(portHexa)
		}
		//fmt.Println(fslc, len(fslc), localA, localAddress, localPort, remoteA, uid, inode)
		row["inode"] = inode
		row["local_address"] = localA
		row["remote_address"] = remoteA
		row["address"] = localAddress
		row["port"] = localPort
		row["uid"] = uid
		irow[localA] = row
	}
	return irow
}

func tcpIpv4() map[string]map[string]interface{} {

	pmap := make(map[string]map[string]interface{})
	tcpFile := "/proc/net/tcp"
	bs, err := os.ReadFile(tcpFile)
	if err != nil {
		return nil
	}
	tcpContent := string(bs)
	irow := processFileDescriptor(tcpContent)
	for _, vmap := range irow {
		//fmt.Println(">>", localA, vmap)
		inode := vmap["inode"].(string)
		pmap[inode] = vmap
	}
	return pmap
}

func tcpIpv6() map[string]map[string]interface{} {

	pmap := make(map[string]map[string]interface{})
	tcpFile := "/proc/net/tcp6"
	bs, err := os.ReadFile(tcpFile)
	if err != nil {
		return nil
	}
	tcpContent := string(bs)
	irow := processFileDescriptor(tcpContent)
	for _, vmap := range irow {
		//fmt.Println(">>", localA, vmap)
		inode := vmap["inode"].(string)
		pmap[inode] = vmap
	}
	return pmap
}

func fileDescriptorList(pid int) []map[string]string {

	var dlist = make([]map[string]string, 0)
	///proc/723/task/723/fd
	taskFd := fmt.Sprintf("/proc/%d/task/%d/fd", pid, pid)
	fsd, err := os.ReadDir(taskFd)
	if err != nil {
		return nil
	}

	for _, fd := range fsd {

		//fmt.Println(fd.Name(), fd.Type() == fs.ModeSymlink)
		if fd.Type() == fs.ModeSymlink {

			fname := filepath.Join(taskFd, fd.Name())
			readlink, err := os.Readlink(fname)
			if err == nil {
				rmap := make(map[string]string)
				//rmap[fd.Name()] = readlink
				rmap["name"] = fd.Name()
				rmap["value"] = readlink
				dlist = append(dlist, rmap)
			}
			//fmt.Println(err, fd.Name(), readlink)
			// var ltype string
			// var inode int
			// spaceSeparatedStr := strings.Replace(readlink, ":", " ", -1)       //socket:[491]
			// n, err := fmt.Sscanf(spaceSeparatedStr, "%s [%d]", &ltype, &inode) //socket [491]
			// if err == nil {
			// 	fmt.Println("**", n, readlink, ltype, inode)
			// }
		}
	}
	return dlist
}

func linkValueToSocketInode(linkValue string) (inode int) {

	//var lmap = make(map[string]int)
	var ltype string
	//var inode int
	spaceSeparatedStr := strings.Replace(linkValue, ":", " ", -1)      //socket:[491]
	_, err := fmt.Sscanf(spaceSeparatedStr, "%s [%d]", &ltype, &inode) //socket [491]
	if err == nil {
		//fmt.Println("**", n, linkValue, ltype, inode)
		//lmap["ltype"] = ltype
		//lmap["node"] = inode
		if ltype == "socket" {
			return inode
		}
	}
	return
}

// pid to socket list
func fileDescriptorSocketList(pid int) []map[string]interface{} {

	var rows = make([]map[string]interface{}, 0)
	fdList := fileDescriptorList(pid)
	for _, fdMap := range fdList {

		row := make(map[string]interface{})
		name := fdMap["name"]
		linkvalue := fdMap["value"]
		inode := linkValueToSocketInode(linkvalue)
		if inode > 0 {
			row["fd"] = name            //fd
			row["inode"] = inode        //socket:[491]
			row["link_type"] = "socket" //socket:[491]
			rows = append(rows, row)
			//fmt.Println(name, inode)
		}
	}
	return rows
}

func pidOpenSocketList(pid int) []map[string]interface{} {

	slist := make([]map[string]interface{}, 0)
	rows := fileDescriptorSocketList(pid)
	for _, row := range rows {

		fd, _ := row["fd"].(string)
		inode, _ := row["inode"].(int)

		//check ipv4 socket
		pmap := tcpIpv4()
		mrow, isOk := pmap[fmt.Sprint(inode)]
		if isOk {
			//fmt.Println(fd, inode, mrow)
			mrow["version"] = "ipv4"
			mrow["fd"] = fd
			slist = append(slist, mrow)
		}

		//check ipv6 socket
		pmap = tcpIpv6()
		mrow, isOk = pmap[fmt.Sprint(inode)]
		if isOk {
			//fmt.Println(fd, inode, mrow)
			mrow["version"] = "ipv6"
			mrow["fd"] = fd
			slist = append(slist, mrow)
		}

	}
	return slist
}

func processStatus(pid int) map[string]string {

	//cat /proc/490/status
	kvMap := make(map[string]string)
	statusFile := fmt.Sprintf("/proc/%d/status", pid)
	bs, err := os.ReadFile(statusFile)
	if err != nil {
		log.Println("processStatusERR:", err)
		return nil
	}
	content := string(bs)
	slc := cmdOutputSlc2(content)
	for _, line := range slc {
		slc = strings.Split(line, ":")
		if len(slc) == 2 {
			key := slc[0]
			kvMap[key] = strings.TrimSpace(slc[1])
		}
	}
	return kvMap
}
