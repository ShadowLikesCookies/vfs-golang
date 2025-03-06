package main

import "time"

type File struct {
	Name             string
	Content          string
	Size             int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ReadPermission   []int
	WritePermission  []int
	ModifyPermission []int
}

type CommandMap map[string]func([]string)
type UsageMap map[string]func()
type Directory struct {
	Name             string
	Files            map[string]*File
	SubDirs          map[string]*Directory
	Parent           *Directory
	CreatedAt        time.Time
	History          []string
	ModifyPermission []int
	ReadPermission   []int
	WritePermission  []int
}

type User struct {
	name       string
	groupPerms []int
}

type VFS struct {
	Root        *Directory
	CurrentDir  *Directory
	CurrentUser *User
}
