package util

import (
	"crypto/md5"
	"encoding/csv"
	"fmt"
	slog "github.com/vearne/simplelog"
	"os"
)

func WriterCSV(path string, records [][]string) {
	File, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error("file open failed:%v", err)
		return
	}
	defer File.Close()

	writer := csv.NewWriter(File)

	err = writer.WriteAll(records)
	if err != nil {
		slog.Error("write csv file:%v", err)
		return
	}
	writer.Flush()
}

func MD5(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
