package mokylin_common

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// game1@1: server 1
// game1@1_2_5: server 1,2,5
// game1@1-5: server 1,2,3,4,5
// game1@1-3_5: 1,2,3,5
// game1@1#2@3-5#4@6
const (
	prefix       = "game"
	sepAt        = "@"
	sepHyphen    = "*"
	sepUnderline = "_"
	sepPound     = "#"
	sepDbUrl     = ":"
)

func IsServerName(s string) bool {
	return strings.HasPrefix(s, prefix)
}

func NameToIds(name string) (map[int][]int, error) {
	errNameToIds := errors.New("游戏服名称解析非法! " + name)

	if !strings.HasPrefix(name, prefix) {
		return nil, errNameToIds
	}

	rMap := make(map[int][]int)
	for _, seg := range strings.Split(name[len(prefix):], sepPound) {
		l1 := strings.Split(seg, sepAt)
		if len(l1) != 2 {
			return nil, errNameToIds
		}

		pid, err := strconv.Atoi(l1[0])
		if err != nil {
			return nil, errNameToIds
		}

		if _, ok := rMap[pid]; ok {
			return nil, errNameToIds
		}

		sids := []int{}
		l2 := strings.Split(l1[1], sepUnderline)
		for _, v1 := range l2 {
			l3 := strings.Split(v1, sepHyphen)
			switch len(l3) {
			case 1:
				a, e := strconv.Atoi(l3[0])
				if e != nil {
					return nil, errNameToIds
				}
				sids = append(sids, a)
			case 2:
				a, e := strconv.Atoi(l3[0])
				if e != nil {
					return nil, errNameToIds
				}
				b, e := strconv.Atoi(l3[1])
				if e != nil {
					return nil, errNameToIds
				}
				if a >= b {
					return nil, errNameToIds
				}
				for i := a; i <= b; i++ {
					sids = append(sids, i)
				}
			default:
				return nil, errNameToIds
			}
		}

		if len(sids) == 0 {
			return nil, errNameToIds
		}

		sort.Ints(sids)
		rMap[pid] = sids
	}

	return rMap, nil
}

func IdsToName(idsMap map[int][]int) (string, error) {
	errIdsToName := errors.New(fmt.Sprintf("游戏服名称合并非法! %v", idsMap))

	i, l := 0, len(idsMap)
	if l == 0 {
		return "", errIdsToName
	}
	pids := make([]int, l) // platformIds
	for pid, _ := range idsMap {
		pids[i] = pid
		i++
	}
	sort.Ints(pids)

	buf := bytes.Buffer{}
	buf.WriteString(prefix)

	for _, pid := range pids {
		sids := idsMap[pid]
		sort.Ints(sids)

		prev := sids[0]
		if prev <= 0 {
			return "", errIdsToName
		}

		buf.WriteString(strconv.Itoa(pid))
		buf.WriteString(sepAt)
		buf.WriteString(strconv.Itoa(prev))
		hasHyphen := false

		for i, l := 1, len(sids); i < l; i++ {
			curr := sids[i]
			switch {
			case curr == prev:
				return "", errIdsToName
			case curr == prev+1:
				hasHyphen = true
			default:
				if hasHyphen {
					buf.WriteString(sepHyphen)
					buf.WriteString(strconv.Itoa(prev))
					hasHyphen = false
				}
				buf.WriteString(sepUnderline)
				buf.WriteString(strconv.Itoa(curr))
			}
			prev = curr
		}

		if hasHyphen {
			buf.WriteString(sepHyphen)
			buf.WriteString(strconv.Itoa(prev))
		}

		buf.WriteString(sepPound)
	}

	name := buf.String()
	name = name[:len(name)-1]

	return name, nil
}

func HasIntersection(n1, n2 string) (bool, error) {
	idsMap1, err := NameToIds(n1)
	if err != nil {
		return false, err
	}

	idsMap2, err := NameToIds(n2)
	if err != nil {
		return false, err
	}

	for pid1, sids1 := range idsMap1 {
		sids2, ok := idsMap2[pid1]
		if !ok {
			continue
		}

		for _, sid1 := range sids1 {
			for _, sid2 := range sids2 {
				if sid1 == sid2 {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
