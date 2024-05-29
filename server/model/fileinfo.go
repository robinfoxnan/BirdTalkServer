package model

type FileInfo struct {
	HashCode  string
	Status    string
	StoreType string
	FileName  string
	UniqName  string
	FileSize  int64
	UserId    int64
	Gid       int64
	Tm        int64
	Tags      []string
}
