package main

import "time"

type File struct {
	Name      string
	Content   string
	Size      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Directory struct {
	Name      string
	Files     map[string]*File
	SubDirs   map[string]*Directory
	Parent    *Directory
	CreatedAt time.Time
}

type VFS struct {
	Root       *Directory
	CurrentDir *Directory
}
