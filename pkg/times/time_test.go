package times_test

import (
	"testing"
	"time"

	"github.com/driif/echo-go-starter/pkg/times"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartOfMonth(t *testing.T) {
	d := times.Date(2020, 3, 12, time.UTC)
	expected := time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, times.StartOfMonth(d))

	d = times.Date(2020, 12, 35, time.UTC)
	expected = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, times.StartOfMonth(d))
}

func TestTimeFromString(t *testing.T) {
	expected := time.Date(2020, 3, 29, 12, 34, 54, 0, time.UTC)

	d, err := times.TimeFromString("2020-03-29T12:34:54Z")
	require.NoError(t, err)

	assert.Equal(t, expected, d)
}

func TestStartOfQuarter(t *testing.T) {
	d := times.Date(2020, 3, 31, time.UTC)
	expected := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, times.StartOfQuarter(d))

	d = times.Date(2020, 1, 1, time.UTC)
	expected = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, times.StartOfQuarter(d))

	d = times.Date(2020, 12, 1, time.UTC)
	expected = time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, times.StartOfQuarter(d))

	d = times.Date(2020, 12, 35, time.UTC)
	expected = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, times.StartOfQuarter(d))

	d = times.Date(2020, 4, 1, time.UTC)
	expected = time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, times.StartOfQuarter(d))
}

func TestStartOfWeek(t *testing.T) {
	d := times.Date(2020, 3, 12, time.UTC)
	expected := time.Date(2020, 3, 9, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, times.StartOfWeek(d))

	d = times.Date(2020, 6, 15, time.UTC)
	expected = time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, times.StartOfWeek(d))

	d = times.Date(2020, 6, 21, time.UTC)
	expected = time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, times.StartOfWeek(d))
}

func TestDateFromString(t *testing.T) {
	res, err := times.DateFromString("2020-01-03")
	require.NoError(t, err)

	require.True(t, res.Equal(time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC)))

	res, err = times.DateFromString("2020-xx-03")
	require.Error(t, err)
	assert.Empty(t, res)
}

func TestEndOfMonth(t *testing.T) {
	d := times.Date(2020, 3, 12, time.UTC)
	expected := time.Date(2020, 3, 31, 23, 59, 59, 999999999, time.UTC)
	assert.True(t, expected.Equal(times.EndOfMonth(d)))

	d = times.Date(2020, 12, 35, time.UTC)
	expected = time.Date(2021, 1, 31, 23, 59, 59, 999999999, time.UTC)
	res := times.EndOfMonth(d)
	assert.True(t, expected.Equal(res))

	expected = time.Date(2021, 1, 31, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, times.TruncateTime(res))

}

func TestEndOfDay(t *testing.T) {
	d := times.Date(2020, 3, 12, time.UTC)
	expected := time.Date(2020, 3, 12, 23, 59, 59, 999999999, time.UTC)
	assert.True(t, expected.Equal(times.EndOfDay(d)))

	d = times.Date(2020, 12, 35, time.UTC)
	expected = time.Date(2021, 1, 4, 23, 59, 59, 999999999, time.UTC)
	res := times.EndOfDay(d)
	assert.True(t, expected.Equal(res))

	expected = time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, times.TruncateTime(res))
}

func TestDateAdds(t *testing.T) {
	d := times.Date(2020, 3, 12, time.UTC)
	expected := time.Date(2022, 4, 12, 0, 0, 0, 0, time.UTC)
	res := times.AddMonths(d, 25)
	assert.True(t, expected.Equal(res))

	d = times.Date(2020, 1, 30, time.UTC)
	expected = time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)
	res = times.AddMonths(d, 1)
	assert.True(t, expected.Equal(res))

	d = times.Date(2020, 1, 30, time.UTC)
	expected = time.Date(2020, 3, 5, 0, 0, 0, 0, time.UTC)
	res = times.AddWeeks(d, 5)
	assert.True(t, expected.Equal(res))
}

func TestDayBefore(t *testing.T) {
	d := times.Date(2020, 3, 1, time.UTC)
	expected := time.Date(2020, 2, 29, 23, 59, 59, 999999999, time.UTC)
	res := times.DayBefore(d)
	assert.True(t, expected.Equal(res))

	expected = time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, times.TruncateTime(res))

}
