package util

import (
	"encoding/csv"
	slog "github.com/vearne/simplelog"
	"os"
)

func WriterCSV(path string, records [][]string) {
	File, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
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
	slog.Info("write report file:%v", path)
}
