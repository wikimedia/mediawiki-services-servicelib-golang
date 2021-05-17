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
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestLogger(t *testing.T) {
	setUp := func(level Level) (*mockWriter, *Logger) {
		writer := &mockWriter{}
		logger, _ := NewLogger(writer, "logtest", "logger", LevelString(level))
		return writer, logger
	}

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
				writer, logger := setUp(DEBUG)

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
			})
		}
	})

	// Logger is configured for INFO and above (DEBUG should be ignored)
	t.Run("Filtered", func(t *testing.T) {
		writer, logger := setUp(INFO)
		logger.Debug("Noisy log message")
		assert.Equal(t, 0, len(writer.data), "Unexpected log output")
	})

	t.Run("Scoped", func(t *testing.T) {
		writer, logger := setUp(INFO)
		logger.RequestID("0000000a-000a-000a-000a-00000000000a").Log(WARNING, "Consider yourself %s", "warned")

		res, err := writer.ReadMessage()
		if err != nil {
			t.Fatalf("Unable to deserialize JSON log message: %s", err)
		}

		assert.Equal(t, "Consider yourself warned", res.Message, "Wrong message string attribute")
		assert.Equal(t, LevelString(WARNING), res.Log.Level, "Wrong log level attribute")
		assert.Equal(t, "logtest", res.Service.Name, "Wrong appname attribute")
		assert.Equal(t, "0000000a-000a-000a-000a-00000000000a", res.Trace.ID, "Wrong request_id attribute")
	})

	t.Run("Using log module", func(t *testing.T) {
		writer, logger := setUp(INFO)
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
