package main

import (
	"fmt"
	"io"
	"iter"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// DocTradeOrder 表示"上传交易订单文件"中的一行数据，
// 字段与文档《文件格式与字段》一一对应。
type DocTradeOrder struct {
	// Amount 交易金额，单位元，必填。服务端按 ×100 转为分存储，
	// 因此最多保留两位小数。
	Amount float64
	// Summary 订单摘要 / 标题，选填，可为空。
	Summary string
	// PayeeID 结算（收款）商户 ID，必填。
	PayeeID string
}

// docTradeHeaders 为表头列名，须与服务端要求完全一致，顺序不可调整。
// payChannel / payProductCode / solutionType / tradeScene 由平台导入时
// 自动补齐，不出现在文件中。
var docTradeHeaders = []string{"交易金额(元)", "订单摘要", "收款商户ID"}

// GenerateDocTradeExcel 以流式写入生成上传交易订单 XLSX 文件并写入 w。
//
// 首行为表头，orders 逐行产出数据（基于 excelize StreamWriter），
// 内存占用与数据量基本无关，适合大文件生成；orders 为 iter.Seq，
// 数据源可以是生成器或数据库游标，无需整体加载到内存。
func GenerateDocTradeExcel(w io.Writer, orders iter.Seq[DocTradeOrder]) error {
	f := excelize.NewFile()
	defer f.Close()

	const sheet = "Sheet1"
	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		return fmt.Errorf("hst: create stream writer: %w", err)
	}

	// 列宽与样式必须在写入任何行之前设置
	for i, width := range []float64{14, 28, 22} {
		if err = sw.SetColWidth(i+1, i+1, width); err != nil {
			return fmt.Errorf("hst: set column width: %w", err)
		}
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#E2EFDA"},
		},
	})
	if err != nil {
		return fmt.Errorf("hst: create header style: %w", err)
	}
	numFmt := "0.00" // 金额固定展示两位小数
	amountStyle, err := f.NewStyle(&excelize.Style{CustomNumFmt: &numFmt})
	if err != nil {
		return fmt.Errorf("hst: create amount style: %w", err)
	}

	header := make([]interface{}, len(docTradeHeaders))
	for i, name := range docTradeHeaders {
		header[i] = name
	}
	if err = sw.SetRow("A1", header, excelize.RowOpts{StyleID: headerStyle}); err != nil {
		return fmt.Errorf("hst: write header row: %w", err)
	}

	row := 1
	for order := range orders {
		row++
		if err = validateDocTradeOrder(order, row); err != nil {
			return err
		}
		cell, err := excelize.CoordinatesToCellName(1, row)
		if err != nil {
			return fmt.Errorf("hst: cell name of row %d: %w", row, err)
		}
		if err = sw.SetRow(cell, []interface{}{
			excelize.Cell{StyleID: amountStyle, Value: order.Amount},
			order.Summary,
			order.PayeeID,
		}); err != nil {
			return fmt.Errorf("hst: write row %d: %w", row, err)
		}
	}
	if err = sw.Flush(); err != nil {
		return fmt.Errorf("hst: flush: %w", err)
	}
	if err = f.Write(w); err != nil {
		return fmt.Errorf("hst: write xlsx: %w", err)
	}
	return nil
}

// validateDocTradeOrder 校验单个订单的必填字段与金额精度。
func validateDocTradeOrder(order DocTradeOrder, row int) error {
	if order.Amount <= 0 {
		return fmt.Errorf("hst: row %d: 交易金额(元) 必填且须大于 0", row)
	}
	if n := decimalPlaces(order.Amount); n > 2 {
		return fmt.Errorf("hst: row %d: 交易金额 %v 超过两位小数", row, order.Amount)
	}
	if order.PayeeID == "" {
		return fmt.Errorf("hst: row %d: 收款商户ID 必填", row)
	}
	return nil
}

// decimalPlaces 返回浮点数最短十进制表示的小数位数，
// 避免浮点乘法（×100）引入的精度误差。
func decimalPlaces(f float64) int {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	if i := strings.IndexByte(s, '.'); i >= 0 {
		return len(s) - i - 1
	}
	return 0
}
