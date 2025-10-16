package util

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vearne/autotest/internal/config"
	slog "github.com/vearne/simplelog"
)

// ReportGenerator 报告生成器
type ReportGenerator struct {
	config config.AutoTestConfig
}

// TestCaseResult 测试用例结果
type TestCaseResult struct {
	ID          uint64        `json:"id"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
	Duration    time.Duration `json:"duration"`
	ErrorMsg    string        `json:"error_message,omitempty"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
}

// ReportData 报告数据
type ReportData struct {
	Summary struct {
		TotalTests   int           `json:"total_tests"`
		PassedTests  int           `json:"passed_tests"`
		FailedTests  int           `json:"failed_tests"`
		SkippedTests int           `json:"skipped_tests"`
		Duration     time.Duration `json:"duration"`
		StartTime    time.Time     `json:"start_time"`
		EndTime      time.Time     `json:"end_time"`
		PassRate     float64       `json:"pass_rate"`
	} `json:"summary"`
	TestCases []TestCaseResult `json:"test_cases"`
}

// JUnitTestSuite JUnit XML格式
type JUnitTestSuite struct {
	XMLName   xml.Name        `xml:"testsuite"`
	Name      string          `xml:"name,attr"`
	Tests     int             `xml:"tests,attr"`
	Failures  int             `xml:"failures,attr"`
	Errors    int             `xml:"errors,attr"`
	Skipped   int             `xml:"skipped,attr"`
	Time      float64         `xml:"time,attr"`
	Timestamp string          `xml:"timestamp,attr"`
	TestCases []JUnitTestCase `xml:"testcase"`
}

// JUnitTestCase JUnit测试用例
type JUnitTestCase struct {
	Name      string        `xml:"name,attr"`
	ClassName string        `xml:"classname,attr"`
	Time      float64       `xml:"time,attr"`
	Failure   *JUnitFailure `xml:"failure,omitempty"`
	Error     *JUnitError   `xml:"error,omitempty"`
	Skipped   *JUnitSkipped `xml:"skipped,omitempty"`
}

// JUnitFailure JUnit失败信息
type JUnitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

// JUnitError JUnit错误信息
type JUnitError struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

// JUnitSkipped JUnit跳过信息
type JUnitSkipped struct {
	Message string `xml:"message,attr"`
}

// NewReportGenerator 创建报告生成器
func NewReportGenerator(cfg config.AutoTestConfig) *ReportGenerator {
	return &ReportGenerator{
		config: cfg,
	}
}

// GenerateReports 生成所有格式的报告
func (rg *ReportGenerator) GenerateReports(data ReportData) error {
	formats := rg.config.Global.Report.Formats
	if len(formats) == 0 {
		formats = []string{"html"} // 默认生成HTML报告
	}

	// 确保报告目录存在
	if err := os.MkdirAll(rg.config.Global.Report.DirPath, 0755); err != nil {
		return fmt.Errorf("failed to create report directory: %w", err)
	}

	for _, format := range formats {
		switch strings.ToLower(format) {
		case "html":
			if err := rg.generateHTMLReport(data); err != nil {
				slog.Error("Failed to generate HTML report: %v", err)
			}
		case "json":
			if err := rg.generateJSONReport(data); err != nil {
				slog.Error("Failed to generate JSON report: %v", err)
			}
		case "csv":
			if err := rg.generateCSVReport(data); err != nil {
				slog.Error("Failed to generate CSV report: %v", err)
			}
		case "junit":
			if err := rg.generateJUnitReport(data); err != nil {
				slog.Error("Failed to generate JUnit report: %v", err)
			}
		default:
			slog.Warn("Unsupported report format: %s", format)
		}
	}

	return nil
}

// generateHTMLReport 生成HTML报告
func (rg *ReportGenerator) generateHTMLReport(data ReportData) error {
	templatePath := rg.config.Global.Report.TemplatePath
	if templatePath == "" {
		templatePath = rg.getDefaultHTMLTemplate()
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		// 如果模板文件不存在，使用默认模板
		tmpl, err = template.New("report").Parse(defaultHTMLTemplate)
		if err != nil {
			return fmt.Errorf("failed to parse HTML template: %w", err)
		}
	}

	filePath := filepath.Join(rg.config.Global.Report.DirPath, "report.html")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create HTML report file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute HTML template: %w", err)
	}

	slog.Info("HTML report generated: %s", filePath)
	return nil
}

// generateJSONReport 生成JSON报告
func (rg *ReportGenerator) generateJSONReport(data ReportData) error {
	filePath := filepath.Join(rg.config.Global.Report.DirPath, "report.json")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create JSON report file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON report: %w", err)
	}

	slog.Info("JSON report generated: %s", filePath)
	return nil
}

// generateCSVReport 生成CSV报告
func (rg *ReportGenerator) generateCSVReport(data ReportData) error {
	filePath := filepath.Join(rg.config.Global.Report.DirPath, "report.csv")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV report file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入标题行
	headers := []string{"ID", "Description", "Status", "Duration", "Start Time", "End Time", "Error Message"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// 写入数据行
	for _, testCase := range data.TestCases {
		record := []string{
			fmt.Sprintf("%d", testCase.ID),
			testCase.Description,
			testCase.Status,
			testCase.Duration.String(),
			testCase.StartTime.Format("2006-01-02 15:04:05"),
			testCase.EndTime.Format("2006-01-02 15:04:05"),
			testCase.ErrorMsg,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	slog.Info("CSV report generated: %s", filePath)
	return nil
}

// generateJUnitReport 生成JUnit XML报告
func (rg *ReportGenerator) generateJUnitReport(data ReportData) error {
	testSuite := JUnitTestSuite{
		Name:      "AutoTest",
		Tests:     data.Summary.TotalTests,
		Failures:  data.Summary.FailedTests,
		Errors:    0,
		Skipped:   data.Summary.SkippedTests,
		Time:      data.Summary.Duration.Seconds(),
		Timestamp: data.Summary.StartTime.Format("2006-01-02T15:04:05"),
	}

	for _, testCase := range data.TestCases {
		junitCase := JUnitTestCase{
			Name:      fmt.Sprintf("TestCase_%d", testCase.ID),
			ClassName: "autotest",
			Time:      testCase.Duration.Seconds(),
		}

		switch testCase.Status {
		case "failed":
			junitCase.Failure = &JUnitFailure{
				Message: testCase.ErrorMsg,
				Type:    "AssertionError",
				Content: testCase.ErrorMsg,
			}
		case "skipped":
			junitCase.Skipped = &JUnitSkipped{
				Message: "Test skipped",
			}
		}

		testSuite.TestCases = append(testSuite.TestCases, junitCase)
	}

	filePath := filepath.Join(rg.config.Global.Report.DirPath, "junit.xml")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create JUnit report file: %w", err)
	}
	defer file.Close()

	file.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(testSuite); err != nil {
		return fmt.Errorf("failed to encode JUnit report: %w", err)
	}

	slog.Info("JUnit report generated: %s", filePath)
	return nil
}

// getDefaultHTMLTemplate 获取默认HTML模板路径
func (rg *ReportGenerator) getDefaultHTMLTemplate() string {
	return filepath.Join(rg.config.Global.Report.DirPath, "default_template.html")
}

// 默认HTML模板
const defaultHTMLTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>AutoTest Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .summary { background: #f5f5f5; padding: 15px; border-radius: 5px; margin-bottom: 20px; }
        .pass { color: green; }
        .fail { color: red; }
        .skip { color: orange; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .status-pass { background-color: #d4edda; }
        .status-fail { background-color: #f8d7da; }
        .status-skip { background-color: #fff3cd; }
    </style>
</head>
<body>
    <h1>AutoTest Report</h1>
    
    <div class="summary">
        <h2>Summary</h2>
        <p><strong>Total Tests:</strong> {{.Summary.TotalTests}}</p>
        <p><strong>Passed:</strong> <span class="pass">{{.Summary.PassedTests}}</span></p>
        <p><strong>Failed:</strong> <span class="fail">{{.Summary.FailedTests}}</span></p>
        <p><strong>Skipped:</strong> <span class="skip">{{.Summary.SkippedTests}}</span></p>
        <p><strong>Pass Rate:</strong> {{printf "%.2f" .Summary.PassRate}}%</p>
        <p><strong>Duration:</strong> {{.Summary.Duration}}</p>
        <p><strong>Start Time:</strong> {{.Summary.StartTime.Format "2006-01-02 15:04:05"}}</p>
        <p><strong>End Time:</strong> {{.Summary.EndTime.Format "2006-01-02 15:04:05"}}</p>
    </div>

    <h2>Test Cases</h2>
    <table>
        <thead>
            <tr>
                <th>ID</th>
                <th>Description</th>
                <th>Status</th>
                <th>Duration</th>
                <th>Start Time</th>
                <th>End Time</th>
                <th>Error Message</th>
            </tr>
        </thead>
        <tbody>
            {{range .TestCases}}
            <tr class="status-{{.Status}}">
                <td>{{.ID}}</td>
                <td>{{.Description}}</td>
                <td>{{.Status}}</td>
                <td>{{.Duration}}</td>
                <td>{{.StartTime.Format "2006-01-02 15:04:05"}}</td>
                <td>{{.EndTime.Format "2006-01-02 15:04:05"}}</td>
                <td>{{.ErrorMsg}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>`
