package file

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/funera1/gofmtal/internal/derror"
)

const chmodSupported = runtime.GOOS != "windows"

// backupFile writes data to a new file named filename<number> with permissions perm,
// with <number randomly chosen such that the file name is unique. backupFile returns
// the chosen file name.
func BackupFile(filename string, data []byte, perm fs.FileMode) (_ string, rerr error) {
	derror.Debug(&rerr, "backupFile(%q)", filename)

	// create backup file
	f, err := os.CreateTemp(filepath.Dir(filename), filepath.Base(filename))
	if err != nil {
		return "", err
	}
	bakname := f.Name()
	if chmodSupported {
		err = f.Chmod(perm)
		if err != nil {
			f.Close()
			os.Remove(bakname)
			return bakname, err
		}
	}

	// write data to backup file
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}

	return bakname, err
}
