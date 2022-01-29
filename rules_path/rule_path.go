package rules_path

import (
	"regexp"
	"strings"
	"sync"
)

type PathStruct struct {
	Path    string
	Length  int
	IsRegex bool
	Role    string
}

type PathMap struct {
	Lock sync.RWMutex
	Path map[string]*PathStruct
}

func InitPathMap() *PathMap {
	return &PathMap{
		Lock: sync.RWMutex{},
		Path: map[string]*PathStruct{},
	}
}

func (p *PathMap) Put(roleId, key, value string) bool {
	node := p
	pathList := strings.Split(value, "/")
	flag := false
	path := ""
	for _, item := range pathList {
		if strings.HasPrefix(item, ":") {
			item = "*"
			flag = true
		}
		path += item + "/"
	}
	mapKey := roleId + key + strings.TrimRight(path, "/")
	_, ok := node.Path[mapKey]
	if ok {
		return false
	}
	node.Lock.Lock()
	defer node.Lock.Unlock()
	node.Path[mapKey] = &PathStruct{
		Path:    key + strings.TrimRight(path, "/"),
		Length:  len(pathList),
		IsRegex: flag,
		Role:    roleId,
	}
	return true
}

func (p *PathMap) Get(roleId string, value string) bool {
	node := p
	pathList := strings.Split(value, "/")
	for _, item := range node.Path {
		if item.Length != len(pathList) || item.Role != roleId {
			continue
		}
		if item.IsRegex {
			reg := regexp.MustCompile(item.Path)
			flag := reg.MatchString(value)
			if flag {
				return true
			}
		}
		if item.Path == value {
			return true
		}
	}
	return false
}

func (p *PathMap) Delete(roleId, key, value string) bool {
	node := p
	pathList := strings.Split(value, "/")
	path := ""
	for _, item := range pathList {
		if strings.HasPrefix(item, ":") {
			item = "*"
		}
		path += item + "/"
	}
	mapKey := roleId + key + strings.TrimRight(path, "/")
	_, ok := node.Path[mapKey]
	if !ok {
		return false
	}
	node.Lock.Lock()
	defer node.Lock.Unlock()
	delete(node.Path, mapKey)
	return true
}
