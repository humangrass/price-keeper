package xtype

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Interval time.Duration

// Scan реализует интерфейс sql.Scanner для чтения из БД
func (i *Interval) Scan(value interface{}) error {
	if value == nil {
		*i = Interval(0)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return i.parsePostgresInterval(string(v))
	case string:
		return i.parsePostgresInterval(v)
	case time.Duration:
		*i = Interval(v)
		return nil
	default:
		return fmt.Errorf("unsupported type for Interval: %T", value)
	}
}

// parsePostgresInterval парсит PostgreSQL INTERVAL формат
func (i *Interval) parsePostgresInterval(s string) error {
	var hours, minutes, seconds int
	_, err := fmt.Sscanf(s, "%02d:%02d:%02d", &hours, &minutes, &seconds)
	if err != nil {
		return fmt.Errorf("failed to parse PostgreSQL interval: %w", err)
	}

	dur := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second

	*i = Interval(dur)
	return nil
}

// Value реализует интерфейс driver.Valuer для записи в БД
func (i Interval) Value() (driver.Value, error) {
	d := time.Duration(i)
	return fmt.Sprintf("%02d:%02d:%02d",
			int(d.Hours())%24,
			int(d.Minutes())%60,
			int(d.Seconds())%60),
		nil
}

// String возвращает человеко-читаемое представление
func (i Interval) String() string {
	d := time.Duration(i)
	if d < time.Minute {
		return d.Round(time.Second).String()
	}
	if d < time.Hour {
		return d.Round(time.Minute).String()
	}
	return d.Round(time.Hour).String()
}

// ToDuration конвертирует в стандартный time.Duration
func (i Interval) ToDuration() time.Duration {
	return time.Duration(i)
}

// FromDuration создает Interval из time.Duration
func FromDuration(d time.Duration) Interval {
	return Interval(d)
}

func (i Interval) Hours() float64 {
	return time.Duration(i).Hours()
}

func (i Interval) Minutes() float64 {
	return time.Duration(i).Minutes()
}

func (i Interval) Seconds() float64 {
	return time.Duration(i).Seconds()
}
