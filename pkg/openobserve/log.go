package openobserve

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tranhuyducseven/goem-template/pkg/utils"
)

type LogWriter struct {
	io.Writer
	Level    zerolog.Level
	password string
	username string

	ch      chan []byte
	buf     [][]byte
	buf_idx int
}

func wrapData(b []byte) *bytes.Buffer {
	data := make([]byte, len(b)+2)
	copy(data, []byte("[")[0:1])
	copy(data[1:], b)
	copy(data[len(b)+1:], []byte("]")[0:1])

	return bytes.NewBuffer(data)
}

func NewLogWriter(level zerolog.Level) *LogWriter {
	b, err := base64.StdEncoding.DecodeString(config.Credential)
	if err != nil {
		log.Fatal().Err(err).Msg("Error decoding credential")
		return nil
	}
	rawUsernamePassword := strings.Split(string(b), ":")

	lw := &LogWriter{
		Level:    level,
		username: rawUsernamePassword[0],
		password: rawUsernamePassword[1],
		ch:       make(chan []byte),
	}
	go lw.collectLogs()
	return lw
}

func (w *LogWriter) collectLogs() {
	const BATCH_SIZE = 5
	w.buf = make([][]byte, BATCH_SIZE)
	w.buf_idx = 0

	for logItem := range w.ch {
		w.buf[w.buf_idx] = logItem
		w.buf_idx++

		// Send the logs if the buffer is full
		if w.buf_idx == BATCH_SIZE {
			w.send()
		}
	}
}

func (w *LogWriter) send() {
	defer func() {
		w.buf_idx = 0
	}()

	b := bytes.Join(w.buf, []byte(","))

	payload := wrapData(b)

	req, err := http.NewRequest("POST", config.logEndpoint, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.SetBasicAuth(w.username, w.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v", err)
		}

		now := time.Now().Format("3:04PM")
		// DON'T use log.Error() here, it will cause infinite loop
		fmt.Printf("%s%s Failed to send logs to OpenObservability: %s\n", now, utils.YellowMsg("WRN"), string(bodyBytes))
	}
}

func (w *LogWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level >= w.Level {
		b := make([]byte, len(p))

		// Write must not modify the slice data, even temporarily. Implementations must not retain p
		// Need to copy the data to a new slice, otherwise the data will be overwritten by the next log
		// See https://stackoverflow.com/a/41250437
		copy(b, p)
		w.ch <- b
	}
	return len(p), nil
}
