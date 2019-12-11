package influxdb

import (
	"fmt"
	"strings"

	v2 "github.com/influxdata/influxdb1-client/v2"
	"github.com/pkg/errors"
)

const (
	TagKeyName    = "tagKey"
	FieldKeyName  = "fieldKey"
	FieldTypeName = "fieldType"
)

type Client struct {
	influxClient v2.Client
}

func NewClient(cfg Config) (Client, error) {

	client, err := v2.NewHTTPClient(v2.HTTPConfig{
		Addr:               cfg.Address,
		Username:           cfg.Username,
		Password:           cfg.Password,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return Client{}, errors.Wrap(err, "new influxdb http client failed")
	}

	return Client{
		influxClient: client,
	}, nil
}

func (c Client) ListMeasurementsInDatabase(db string) ([]string, error) {

	cmd := "SHOW MEASUREMENTS"

	q := v2.NewQuery(cmd, db, "")
	resp, err := c.influxClient.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "query influxdb failed")
	} else if err := resp.Error(); err != nil {
		return nil, errors.Wrap(err, "query influxdb failed")
	}

	if len(resp.Results) != 1 {
		return nil, errors.New("length of results must be 1")
	} else if len(resp.Results[0].Series) > 1 {
		return nil, errors.New("length of series must not be greater 1")
	}

	measurements := make([]string, 0)
	for _, serie := range resp.Results[0].Series {
		for _, value := range serie.Values {
			if len(value) != 1 {
				return nil, errors.New("length of values must be 1")
			}
			m, ok := value[0].(string)
			if !ok {
				return nil, errors.Wrap(err, "convert value from interface{} to string failed")
			}
			measurements = append(measurements, m)
		}
	}

	return measurements, nil
}

func (c Client) ListMeasurements(db string, measurements []string) ([]Measurement, error) {

	showTAGCMD := "SHOW TAG KEYS FROM "
	for _, m := range measurements {
		showTAGCMD += fmt.Sprintf(`%s,`, m)
	}
	showTAGCMD = strings.TrimSuffix(showTAGCMD, ",")

	showFieldCMD := "SHOW FIELD KEYS FROM "
	for _, m := range measurements {
		showFieldCMD += fmt.Sprintf(`%s,`, m)
	}
	showFieldCMD = strings.TrimSuffix(showFieldCMD, ",")

	cmd := fmt.Sprintf(`%s;%s`, showTAGCMD, showFieldCMD)

	q := v2.NewQuery(cmd, db, "")
	resp, err := c.influxClient.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "query influxdb failed")
	} else if err := resp.Error(); err != nil {
		return nil, errors.Wrap(err, "query influxdb failed")
	}

	measurementMap := map[string]Measurement{}

	showTagResult := resp.Results[0]
	for _, serie := range showTagResult.Series {
		if len(serie.Columns) != 1 {
			return nil, errors.Errorf(`length of column must be "1", get "%d"`, len(serie.Columns))
		} else if serie.Columns[0] != TagKeyName {
			return nil, errors.Errorf(`columns[0] in serie must be "%s", get "%s"`, TagKeyName, serie.Columns[0])
		}

		measurementName := serie.Name
		measurement, exist := measurementMap[measurementName]
		if !exist {
			measurement = Measurement{Database: db, Name: measurementName}
		}

		tags := make([]string, len(serie.Values))
		for i, value := range serie.Values {
			if len(value) != 1 {
				return nil, errors.New("each value's length in values should be 1")
			}
			tag, ok := value[0].(string)
			if !ok {
				return nil, errors.New("convert tag to string failed")
			}
			tags[i] = tag
		}
		measurement.Tags = tags
		measurementMap[measurementName] = measurement
	}

	showFieldResult := resp.Results[1]
	for _, serie := range showFieldResult.Series {
		if len(serie.Columns) != 2 {
			return nil, errors.Errorf(`length of column must be "2", get "%d"`, len(serie.Columns))
		} else if serie.Columns[0] != FieldKeyName || serie.Columns[1] != FieldTypeName {
			return nil, errors.Errorf(`columns[0] and columns[1] in serie must be "%s" and "%s", get "%s" and "%s"`,
				TagKeyName, FieldTypeName, serie.Columns[0], serie.Columns[1])
		}

		measurementName := serie.Name
		measurement, exist := measurementMap[measurementName]
		if !exist {
			measurement = Measurement{Database: db, Name: measurementName}
		}

		fields := make([]Field, len(serie.Values))
		for i, value := range serie.Values {
			if len(value) != 2 {
				return nil, errors.New("length should be 2")
			}

			name, ok := value[0].(string)
			if !ok {
				return nil, errors.New("convert field name to string failed")
			}
			fieldType, ok := value[1].(string)
			if !ok {
				return nil, errors.New("convert field type to string failed")
			}

			fields[i] = Field{Name: name, Type: fieldType}
		}
		measurement.Fields = fields
		measurementMap[measurementName] = measurement
	}

	res := make([]Measurement, 0, len(measurementMap))
	for _, m := range measurementMap {
		res = append(res, m)
	}
	return res, nil
}

func (c Client) DropDatabases(databases []string) error {

	cmd := ""
	for _, database := range databases {
		cmd += fmt.Sprintf(`DROP DATABASE %s;`, database)
	}

	q := v2.NewQuery(cmd, "", "")
	resp, err := c.influxClient.Query(q)
	if err != nil {
		return errors.Wrap(err, "query influxdb failed")
	} else if err := resp.Error(); err != nil {
		return errors.Wrap(err, "query influxdb failed")
	}

	return nil
}
