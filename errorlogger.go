/*
Logger interface for sendcloud
*/
package mailgun

import (
    "fmt"
)

type ErrorLogger interface {
    ErrorLog(source string, code int, msg string) error
}

type FmtErrorLogger struct {
}

func (l FmtErrorLogger) ErrorLog(source string, code int, msg string) error {
    if code < 200 || code >= 300 {
        return fmt.Errorf("%s: code=%d, msg=%s", source, code, msg)
    }
    fmt.Printf("%s: code=%d, msg=%s", source, code, msg)
    return nil
}
