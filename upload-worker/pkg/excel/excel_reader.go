package excel

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

func ReadExcel(file io.Reader) ([][]string, error) {
	f, err := excelize.OpenReader(file)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	sheets := f.GetSheetList()

	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in excel file")
	}

	rows, err := f.GetRows(sheets[0])

	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %v", err)
	}

	return rows, nil
}
