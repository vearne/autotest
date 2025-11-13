package command

import (
	"context"
	"embed"
	"fmt"
	"github.com/lianggaoqiang/progress"
	"github.com/vearne/autotest/internal/config"
	"github.com/vearne/autotest/internal/model"
	"github.com/vearne/autotest/internal/resource"
	"github.com/vearne/autotest/internal/util"
	"github.com/vearne/executor"
	slog "github.com/vearne/simplelog"
	"github.com/vearne/zaplog"
	"go.uber.org/zap"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func HttpAutomateTest(httpTestCases map[string][]*config.TestCaseHttp) *UnifiedTestResults {
	total := 0
	for _, testcases := range httpTestCases {
		total += len(testcases)
	}
	if total <= 0 {
		return &UnifiedTestResults{}
	}

	begin := time.Now()
	slog.Info("[start]HttpTestCases, total:%v", total)

	workerNum := resource.GlobalConfig.Global.WorkerNum

	finishCount := 0
	successCount := 0
	failedCount := 0
	var failedCases []string
	
	for filePath := range httpTestCases {
		// if ignore_testcase_fail is false and some testcases have failed.
		if resource.TerminationFlag.Load() {
			break
		}

		info, tcResultList := HandleSingleFileHttp(workerNum, filePath)
		finishCount += info.Total
		successCount += info.SuccessCount
		failedCount += info.FailedCount
		
		// 收集失败用例信息
		for _, tcResult := range tcResultList {
			if tcResult.State != model.StateSuccessFul {
				failedCases = append(failedCases, fmt.Sprintf("HTTP_%d: %s", tcResult.ID, tcResult.Desc))
			}
		}
		
		slog.Info("HttpTestCases, total:%v, finishCount:%v, successCount:%v, failedCount:%v",
			total, finishCount, successCount, failedCount)
		// generate report file (保留旧的报告生成作为备份)
		GenReportFileHttp(filePath, tcResultList, info)
	}
	slog.Info("[end]HttpTestCases, total:%v, cost:%v", total, time.Since(begin))
	
	return &UnifiedTestResults{
		TotalTests:  finishCount,
		PassedTests: successCount,
		FailedTests: failedCount,
		FailedCases: failedCases,
	}
}

func GenReportFileHttp(testCasefilePath string, tcResultList []HttpTestCaseResult, info *ResultInfo) {
	filename := filepath.Base(testCasefilePath)
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	filename = name + ".csv"

	reportDirPath := resource.GlobalConfig.Global.Report.DirPath
	reportPath := filepath.Join(reportDirPath, filename)
	sort.Slice(tcResultList, func(i, j int) bool {
		return tcResultList[i].ID < tcResultList[j].ID
	})
	// 1. csv file
	var records [][]string
	records = append(records, []string{"id", "description", "state", "reason"})
	for _, item := range tcResultList {
		reasonStr := item.Reason.String()
		if item.Reason == model.ReasonSuccess {
			reasonStr = ""
		}
		records = append(records, []string{strconv.Itoa(int(item.ID)),
			item.Desc, item.State.String(), reasonStr})
	}
	util.WriterCSV(reportPath, records)
	// 2. html file
	dirName := util.MD5(reportDirPath + name)

	var caseResults []CaseShow
	for _, item := range tcResultList {
		caseResults = append(caseResults, CaseShow{ID: item.ID, Description: item.Desc,
			State: item.State.String(), Reason: item.Reason.String(),
			Link: fmt.Sprintf("./%v/%v.html", dirName, item.ID)})
	}

	obj := map[string]any{
		"info":         info,
		"tcResultList": caseResults,
	}
	// index file
	err := RenderTpl(mytpl, "template/index.tpl", obj, filepath.Join(reportDirPath, name+".html"))
	if err != nil {
		slog.Error("RenderTpl, %v", err)
		return
	}

	// case file
	for _, item := range tcResultList {
		data := map[string]any{
			"Error":      item.Error,
			"reqDetail":  item.ReqDetail(),
			"respDetail": item.RespDetail(),
		}
		err := RenderTpl(mytpl, "template/case.tpl", data,
			filepath.Join(reportDirPath, dirName, strconv.Itoa(int(item.ID))+".html"))
		if err != nil {
			slog.Error("RenderTpl, %v", err)
			return
		}
	}
	slog.Info("write report:%v", reportDirPath)
}

func RenderTpl(fs embed.FS, key string, obj map[string]any, targetPath string) error {
	data, err := fs.ReadFile(key)
	if err != nil {
		slog.Error("mytpl.ReadFile, %v", err)
		return err
	}
	t, err := template.New("index").Parse(string(data))
	if err != nil {
		slog.Error("template Parse, %v", err)
		return err
	}
	dirPath := filepath.Dir(targetPath)
	if !pathExists(dirPath) {
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(targetPath)
	if err != nil {
		slog.Error("Create file, %v", err)
		return err
	}
	defer file.Close()
	return t.Execute(file, obj)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	// os.IsNotExist 判断错误是否为文件或目录不存在
	return !os.IsNotExist(err)
}

func HandleSingleFileHttp(workerNum int, filePath string) (*ResultInfo, []HttpTestCaseResult) {
	workerNum = min(workerNum, 10)
	testcases := resource.HttpTestCases[filePath]
	slog.Info("[start]HandleSingleFileHttp, filePath:%v, len(testcase):%v", filePath, len(testcases))

	futureChan := make(chan executor.Future, len(testcases))
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
		zaplog.Debug("future.Get",
			zap.Uint64("ID", tcResult.ID),
			zap.String("Desc", tcResult.Desc),
			zap.Any("request", tcResult.Request),
			zap.Any("response", tcResult.Response),
		)

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
		// process bar will override this line.
		fmt.Println()
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
