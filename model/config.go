package model

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/94peter/log"
)

type Config struct {
	PersistanceFile string `env:"PERSISTANCE_FILE"`

	Log       log.Logger
	serialMgr SerialMgr
}

func GetModelCfgFromEnv() (*Config, error) {
	config := &Config{}
	err := GetFromEnv(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) Close() {
	if c.serialMgr != nil {
		c.serialMgr.Persistance()
	}
}

func (c *Config) NewSerial() SerialMgr {
	if c.serialMgr == nil {
		c.serialMgr = NewSerial(c)
	}
	return c.serialMgr
}

func GetFromEnv(obj any) error {
	// check obj is pointer
	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		return errors.New("obj must be a pointer")
	}

	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envValue := os.Getenv(envTag)
		if envValue == "" {
			return errors.New("environmental variable " + envTag + " must not be blank")
		}

		fieldValue := v.Field(i)
		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(envValue)
		case reflect.Int:
			intValue, err := strconv.Atoi(envValue)
			if err != nil {
				return errors.New("environmental variable " + envTag + " must be an integer")
			}
			fieldValue.SetInt(int64(intValue))
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(envValue)
			if err != nil {
				return errors.New("environmental variable " + envTag + " must be a boolean")
			}
			fieldValue.SetBool(boolValue)
		case reflect.TypeOf(time.Duration(0)).Kind():
			durationValue, err := time.ParseDuration(envValue)
			if err != nil {
				return errors.New("environmental variable " + envTag + " must be a duration")
			}
			fieldValue.Set(reflect.ValueOf(durationValue))
		default:
			return errors.New("unsupported type: " + fieldValue.Kind().String())
		}
	}
	return nil
}
