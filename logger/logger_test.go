/*
 * Copyright 2019 Eric Evans <eevans@wikimedia.org>, and Wikimedia
 * Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockWriter struct {
	data []byte
}

func (m *mockWriter) Write(data []byte) (n int, err error) {
	m.data = data
	return len(m.data), nil
}

func (m *mockWriter) ReadMessage() (msg *LogMessage, err error) {
	if err = json.Unmarshal(m.data, &msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func setUp(level string) (*mockWriter, *Logger) {
	writer := &mockWriter{}
	logger, _ := NewLogger(writer, "logtest", level)
	return writer, logger
}

func TestLogger(t *testing.T) {
	t.Run("Simple logging", func(t *testing.T) {
		testCases := []struct {
			format string
			arg    string
			level  Level
		}{
			{"Debug %s", "your bugs", DEBUG},
			{"Info %s", "wars", INFO},
			{"Consider yourself %s", "warned", WARNING},
			{"Errors are %s", "bad", ERROR},
			{"Fatal %s", "attraction", FATAL},
		}
		for _, tcase := range testCases {
			t.Run(LevelString(tcase.level), func(t *testing.T) {
				writer, logger := setUp("DEBUG")

				switch tcase.level {
				case DEBUG:
					logger.Debug(tcase.format, tcase.arg)
				case INFO:
					logger.Info(tcase.format, tcase.arg)
				case WARNING:
					logger.Warning(tcase.format, tcase.arg)
				case ERROR:
					logger.Error(tcase.format, tcase.arg)
				case FATAL:
					logger.Fatal(tcase.format, tcase.arg)
				default:
					t.Fatalf("Testcase has invalid level!")
				}

				r, err := writer.ReadMessage()
				if err != nil {
					t.Fatalf("Unable to deserialize JSON log message: %s", err)
				}

				assert.Equal(t, fmt.Sprintf(tcase.format, tcase.arg), r.Message, "Wrong message string attribute")
				assert.Equal(t, LevelString(tcase.level), r.Log.Level, "Wrong log level attribute")
				assert.Equal(t, "logtest", r.Service.Name, "Wrong appname attribute")
				assert.Equal(t, ecsVersion, r.ECS.Version)
			})
		}
	})

	// Logger is configured for INFO and above (DEBUG should be ignored)
	t.Run("Filtered", func(t *testing.T) {
		writer, logger := setUp("INFO")
		logger.Debug("Noisy log message")
		assert.Equal(t, 0, len(writer.data), "Unexpected log output")
	})

	t.Run("Using log module", func(t *testing.T) {
		writer, logger := setUp("INFO")
		log.SetFlags(0)
		log.SetOutput(logger)
		log.Println("Sent via log module")

		res, err := writer.ReadMessage()
		if err != nil {
			t.Fatalf("Unable to deserialize JSON log message: %s", err)
		}

		assert.Equal(t, "Sent via log module", res.Message, "Wrong message string attribute")
		assert.Equal(t, LevelString(WARNING), res.Log.Level, "Wrong log level attribute")
	})
}

func TestRequestScoped(t *testing.T) {
	writer, logger := setUp("DEBUG")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Request(r).Log(INFO, "In yer request, logging yer logs")
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	require.Nil(t, err)

	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	require.Nil(t, err)

	msg, err := writer.ReadMessage()
	require.Nil(t, err)
	assert.Equal(t, "In yer request, logging yer logs", msg.Message, "Message text")
	assert.Equal(t, LevelString(INFO), msg.Log.Level, "Logging level")
	assert.NotNil(t, msg.Client)
	assert.NotEmpty(t, msg.Client.IP)
	assert.NotEmpty(t, msg.Client.Port)
	assert.Equal(t, ecsVersion, msg.ECS.Version)

}
