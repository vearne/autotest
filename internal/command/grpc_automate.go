package command

import (
	"github.com/vearne/autotest/internal/config"
	"github.com/vearne/autotest/internal/resource"
	slog "github.com/vearne/simplelog"
	"time"
)

func GrpcAutomateTest(grpcTestCases map[string][]*config.TestCaseGrpc) {
	total := 0
	for _, testcases := range grpcTestCases {
		total += len(testcases)
	}
	if total <= 0 {
		return
	}

	begin := time.Now()
	slog.Info("[start]GrpcTestCases, total:%v", total)

	workerNum := resource.GlobalConfig.Global.WorkerNum

	finishCount := 0
	successCount := 0
	failedCount := 0
	for filePath := range grpcTestCases {
		// if ignore_testcase_fail is false and some testcases have failed.
		if resource.TerminationFlag.Load() {
			break
		}

		info, tcResultList := HandleSingleFileGrpc(workerNum, filePath)
		finishCount += info.Total
		successCount += info.SuccessCount
		failedCount += info.FailedCount
		slog.Info("GrpcTestCases, total:%v, finishCount:%v, successCount:%v, failedCount:%v",
			total, finishCount, successCount, failedCount)
		// generate report file
		GenReportFileGrpc(filePath, tcResultList)
	}
	slog.Info("[end]GrpcTestCases, total:%v, cost:%v", total, time.Since(begin))
}

func HandleSingleFileGrpc(workerNum int, filePath string) (*ResultInfo, []HttpTestCaseResult) {
	return nil, nil
}

func GenReportFileGrpc(testCasefilePath string, tcResultList []HttpTestCaseResult) {

}
