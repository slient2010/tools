package mokylin_front

import (
	"os/exec"
	"sort"
	"strings"
)

var (
	// ResRsyncCmd = "env RSYNC_PASSWORD=topsecret rsync -av --port 10873 gjqt@release.gjqt.com.cn::gjqt/pack_server/ | awk '{ print $5 }' | grep -E '^server_.+\\.tar$'"
	ResRsyncCmd = "env RSYNC_PASSWORD=topsecret rsync -av --port 8190 gjqt@180.149.146.30::gjqt/pack_server/ | awk '{ print $5 }' | grep -E '^server_.+\\.tar$'"
)

func ResList(daili, game []string) ([]string, error) {
	cmd := ResRsyncCmd
	data, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return nil, err
	}

	list := strings.Split(string(data), "\n")
	list = list[:len(list)-1] // 去除最后的空行
	sort.Strings(list)

	l := len(list)
	for i, j := 0, l/2; i < j; i++ {
		list[i], list[l-1-i] = list[l-1-i], list[i]
	}

	return list, nil
}
