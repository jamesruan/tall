package tall

import ()

// implements File
type file struct {
	name     string    // use as key for Dir.Entries
	rAgent   Reader    // nil unless opened
	wAgent   Writer    // nil unless created
	fileStat *FileStat // nil unless directory being read
}
