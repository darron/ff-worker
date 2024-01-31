package d1

import (
	"database/sql/driver"
	"errors"
	"io"
	"math"
	"sync"
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

type rows struct {
	rowsObj    js.Value
	currentRow int
	// columns is cached value of Columns method.
	// do not use this directly.
	_columns    []string
	onceColumns sync.Once
	// _rowsLen is cached value of rowsLen method.
	// do not use this directly.
	_rowsLen    int
	onceRowsLen sync.Once
	mu          sync.Mutex
}

var _ driver.Rows = (*rows)(nil)

// Columns returns column names retrieved from query result object's keys.
// If rows are empty, this returns nil.
func (r *rows) Columns() []string {
	r.onceColumns.Do(func() {
		if r.rowsObj.Length() == 0 {
			// return nothing when row count is zero.
			return
		}
		colsArray := jsutil.ObjectClass.Call("keys", r.rowsObj.Index(0))
		colsLen := colsArray.Length()
		cols := make([]string, colsLen)
		for i := 0; i < colsLen; i++ {
			cols[i] = colsArray.Index(i).String()
		}
		r._columns = cols
	})
	return r._columns
}

func (r *rows) Close() error {
	// do nothing
	return nil
}

// isIntegralNumber returns if given float64 value is integral value or not.
func isIntegralNumber(f float64) bool {
	// If the value is NaN or Inf, returns the value to avoid being mistakenly treated as an integral value.
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return false
	}
	return f == math.Trunc(f)
}

// convertRowColumnValueToDriverValue converts row column's value in JS to Go's driver.Value.
// row column value is `null | Number | String | ArrayBuffer`.
// see: https://developers.cloudflare.com/d1/platform/client-api/#type-conversion
func convertRowColumnValueToAny(v js.Value) (driver.Value, error) {
	switch v.Type() {
	case js.TypeNull:
		return nil, nil
	case js.TypeNumber:
		fv := v.Float()
		// if the value can be treated as integral value, return as int64.
		if isIntegralNumber(fv) {
			return int64(fv), nil
		}
		return fv, nil
	case js.TypeString:
		return v.String(), nil
	case js.TypeObject:
		// TODO: handle BLOB type (ArrayBuffer).
		// see: https://developers.cloudflare.com/d1/platform/client-api/#type-conversion
		return nil, errors.New("d1: row column value type object is not currently supported")
	}
	return nil, errors.New("d1: unexpected row column value type")
}

func (r *rows) Next(dest []driver.Value) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.currentRow == r.rowsLen() {
		return io.EOF
	}
	rowObj := r.rowsObj.Index(r.currentRow)
	cols := r.Columns()
	for i, col := range cols {
		v, err := convertRowColumnValueToAny(rowObj.Get(col))
		if err != nil {
			return err
		}
		dest[i] = v
	}
	r.currentRow++
	return nil
}

func (r *rows) rowsLen() int {
	r.onceRowsLen.Do(func() {
		r._rowsLen = r.rowsObj.Length()
	})
	return r._rowsLen
}
