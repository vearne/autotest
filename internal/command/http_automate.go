package command

import (
	"context"
	"github.com/lianggaoqiang/progress"
	"github.com/vearne/autotest/internal/config"
	"github.com/vearne/autotest/internal/model"
	"github.com/vearne/autotest/internal/resource"
	"github.com/vearne/autotest/internal/util"
	"github.com/vearne/executor"
	slog "github.com/vearne/simplelog"
	"github.com/vearne/zaplog"
	"go.uber.org/zap"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ResultInfo struct {
	Total        int
	SuccessCount int
	FailedCount  int
}

func HttpAutomateTest(httpTestCases map[string][]*config.TestCaseHttp) {
	total := 0
	for _, testcases := range httpTestCases {
		total += len(testcases)
	}
	if total <= 0 {
		return
	}

	begin := time.Now()
	slog.Info("[start]HttpTestCases, total:%v", total)

	workerNum := resource.GlobalConfig.Global.WorkerNum

	finishCount := 0
	successCount := 0
	failedCount := 0
	for filePath := range httpTestCases {
		// if ignore_testcase_fail is false and some testcases have failed.
		if resource.TerminationFlag.Load() {
			break
		}

		info, tcResultList := HandleSingleFileHttp(workerNum, filePath)
		finishCount += info.Total
		successCount += info.SuccessCount
		failedCount += info.FailedCount
		slog.Info("HttpTestCases, total:%v, finishCount:%v, successCount:%v, failedCount:%v",
			total, finishCount, successCount, failedCount)
		// generate report file
		GenReportFileHttp(filePath, tcResultList)
	}
	slog.Info("[end]HttpTestCases, total:%v, cost:%v", total, time.Since(begin))
}

func GenReportFileHttp(testCasefilePath string, tcResultList []HttpTestCaseResult) {
	filename := filepath.Base(testCasefilePath)
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	filename = name + ".csv"

	reportDirPath := resource.GlobalConfig.Global.Report.DirPath
	reportPath := filepath.Join(reportDirPath, filename)
	sort.Slice(tcResultList, func(i, j int) bool {
		return tcResultList[i].ID < tcResultList[j].ID
	})
	var records [][]string
	records = append(records, []string{"id", "desc", "state", "reason"})
	for _, item := range tcResultList {
		reasonStr := item.Reason.String()
		if item.Reason == model.ReasonSuccess {
			reasonStr = ""
		}
		records = append(records, []string{strconv.Itoa(int(item.ID)),
			item.Desc, item.State.String(), reasonStr})
	}
	util.WriterCSV(reportPath, records)
}

func HandleSingleFileHttp(workerNum int, filePath string) (*ResultInfo, []HttpTestCaseResult) {
	workerNum = min(workerNum, 10)
	testcases := resource.HttpTestCases[filePath]
	slog.Info("[start]HandleSingleFileHttp, filePath:%v, len(testcase):%v", filePath, len(testcases))

	futureChan := make(chan executor.Future, 10)
	pool := executor.NewFixedGPool(context.Background(), workerNum)
	defer pool.WaitTerminate()

	stateGroup := model.NewStateGroup()
	for _, testcase := range testcases {
		stateGroup.SetState(testcase.GetID(), model.StateNotExecuted)
	}
	// producer
	go func() {
		for i := 0; i < len(testcases); i++ {
			tc := testcases[i]
			f, err := pool.Submit(&HttpTestCallable{testcase: tc, stateGroup: stateGroup})
			if err != nil {
				zaplog.Error("pool.Submit", zap.Any("testcase", tc), zap.Error(err))
			}
			futureChan <- f
		}
	}()

	p := progress.Start()
	bar := progress.NewBar().Custom(
		progress.BarSetting{
			StartText:       "[",
			EndText:         "]",
			PassedText:      "-",
			FirstPassedText: ">",
			NotPassedText:   "=",
		},
	)
	p.AddBar(bar)

	finishCount := 0
	successCount := 0
	failedCount := 0

	var tcResultList []HttpTestCaseResult
	for future := range futureChan {
		gResult := future.Get()
		tcResult := gResult.Value.(HttpTestCaseResult)
		zaplog.Debug("future.Get", zap.Any("tcResult", tcResult))

		if tcResult.State == model.StateNotExecuted {
			time.Sleep(200 * time.Millisecond)
			// wait for a while
			f, err := pool.Submit(&HttpTestCallable{testcase: tcResult.TestCase, stateGroup: stateGroup})
			if err != nil {
				zaplog.Error("pool.Submit", zap.Any("testcase", tcResult.TestCase), zap.Error(err))
			} else {
				futureChan <- f
			}
			continue
		}

		if tcResult.State == model.StateSuccessFul {
			successCount++
		} else {
			failedCount++
			// terminate subsequent testcases
			if !resource.GlobalConfig.Global.IgnoreTestCaseFail {
				resource.TerminationFlag.Store(true)
				break
			}
		}

		stateGroup.SetState(tcResult.ID, tcResult.State)
		// prepare to write report
		tcResultList = append(tcResultList, tcResult)

		// process the variables generated when testcase is run
		for key, value := range tcResult.KeyValues {
			resource.CustomerVars.Store(key, value)
		}

		finishCount++
		//nolint: errcheck
		bar.Percent(float64(finishCount) / float64(len(testcases)) * 100)
		if finishCount >= len(testcases) {
			// finish all test cases
			break
		}
	}
	slog.Info("[end]HandleSingleFile, filePath:%v", filePath)
	return &ResultInfo{Total: finishCount, SuccessCount: successCount,
		FailedCount: failedCount}, tcResultList
}
