// Command hst 生成"上传交易订单文件"接口所需的 XLSX 测试文件，
// 演示流式 Excel 生成，并计算文件 SM3 哈希（申请上传凭证时须提交）。
package main

import (
	"fmt"
	"io"
	"iter"
	"math/rand/v2"
	"os"
	"slices"
	"time"

	"github.com/emmansun/gmsm/sm3"
	"github.com/xuri/excelize/v2"
)

const testRows = 1000

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() (err error) {
	filename := fmt.Sprintf("trade-%s.xlsx", time.Now().Format("20060102"))
	start := time.Now()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	// 生成文件的同时计算 SM3 哈希，避免二次读取文件
	hash := sm3.New()
	if err = GenerateDocTradeExcel(
		io.MultiWriter(file, hash),
		randomDocTradeOrders(testRows),
	); err != nil {
		return err
	}

	if err = verifyDocTradeExcel(filename, testRows); err != nil {
		return fmt.Errorf("verify: %w", err)
	}

	info, err := os.Stat(filename)
	if err != nil {
		return err
	}
	fmt.Printf("已生成 %s：%d 条随机订单，%d 字节，耗时 %s\n",
		filename, testRows, info.Size(), time.Since(start).Round(time.Millisecond))
	fmt.Printf("文件 SM3 哈希：%x\n", hash.Sum(nil))
	return nil
}

var docTradeSummaries = []string{
	"6月直播佣金",
	"直播带货分账",
	"达人坑位费结算",
	"直播打赏分成",
	"带货佣金结算",
	"品牌专场分账",
}

// randomDocTradeOrders 返回 n 条随机订单序列，逐条产出、
// 不预先占用内存，与流式写入配合。
func randomDocTradeOrders(n int) iter.Seq[DocTradeOrder] {
	return func(yield func(DocTradeOrder) bool) {
		for i := 0; i < n; i++ {
			if !yield(DocTradeOrder{
				// 以“分”为单位随机后除以 100，避免浮点误差产生三位小数
				Amount:  float64(rand.IntN(9_999_900)+100) / 100, // 1.00 ~ 99999.99
				Summary: randomSummary(),
				PayeeID: fmt.Sprintf("MC%d%08d", 2025+rand.IntN(2), rand.Int64N(1e8)),
			}) {
				return
			}
		}
	}
}

// randomSummary 订单摘要为选填字段，按 20% 概率留空以覆盖该场景。
func randomSummary() string {
	if rand.IntN(10) < 2 {
		return ""
	}
	return docTradeSummaries[rand.IntN(len(docTradeSummaries))]
}

// verifyDocTradeExcel 读回文件校验：表头列名与顺序完全一致、
// 数据行数正确（模拟服务端导入校验）。
func verifyDocTradeExcel(filename string, wantRows int) error {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	rows, err := f.Rows("Sheet1")
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	var count int
	headerChecked := false
	for rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			return err
		}
		if !headerChecked {
			if !slices.Equal(cols, docTradeHeaders) {
				return fmt.Errorf("表头不匹配: %v", cols)
			}
			headerChecked = true
			continue
		}
		count++
	}
	if count != wantRows {
		return fmt.Errorf("数据行数 %d != %d", count, wantRows)
	}
	return nil
}
