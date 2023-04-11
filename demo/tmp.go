package main

import (
	"encoding/json"
	"fmt"
)

// # 01 map 中的数组如何添加元素
// 错误
/*array := testMap["1"]

array = append(array, "10")*/

// 正确
/*
	testMap := map[string][]string{
		"1": []string{"1", "2"},
	}

	testMap["1"] = append(testMap["1"], "10")

*/

var reqJson = `{
	"eventId": 3704372516,
	"returnData": {
		"codeDesc": "Success",
		"code": 0,
		"data": {
			"credentials": {
				"tmpSecretKey": "NDNNy/Ft1q7Q9PaRXnbSjvjF+dk/FPbfS86wNeqsWME=",
				"tmpSecretId": "AKIDyOvI8aojvMk-37DU8i-HqNootR_-6GD003AyK7I0w-tmWuajapF0yB7IO3XsFIVN",
				"sessionToken": "4kaTz24zZaE8xdmTGre05CavwxcbQcRae7cdb5c8395bb3b1e3ef69467eb62986ni7DcYE5QijCUf_6gAFzvy-Y_SCUk4WFyFiog6OWTEOiyAtY7EvYEAGnETxnjPWt9rrcZnGs4pRNx4cVdR1Nof2-GnjRGnQ5II05XfIce_GaV7INmvgBoL3Ghp1aLrSA27EUyUNjG_cbAjMXw_NlXpUZ0GJG4r9rNRqkKX_xZ_3XtUfkF8QE6MivjkyZceOoLnHgO9rA3omx8DMX_RObVYV0uYELvicEj2e-2qayqmrp3_-OnnqbhLFyMUv31f5L6OAAbIy76erDQaiVe2s_TtBatVfnOJRn7zxl13vY2B-c4KMEPaNQIO9gw2MV66-v-BzpiU-qpxBavS7VPw6rD2-_fxeqaPWmvbFW11s3LhE"
			},
			"expiredTime": 1652091132,
			"expiration": "2022-05-09T10:12:12Z"
		}
	},
	"version": "1.0",
	"returnValue": 0,
	"returnMsg": "success",
	"timestamp": 1652087131,
	"caller": "cloudprovider",
	"callee": "NORM"
}`

type NormGetAgentCredentialRsp struct {
	Data struct {
		Credentials struct {
			Token        string `json:"sessionToken"`
			TmpSecretId  string `json:"tmpSecretId"`
			TmpSecretKey string `json:"tmpSecretKey"`
		} `json:"credentials"`
		ExpiredTime int `json:"expiredTime"`
	} `json:"data"`
}

type NormResponse struct {
	Code     int         `json:"returnValue"`
	Msg      string      `json:"returnMsg"`
	Version  string      `json:"version"`
	Password string      `json:"password"`
	Data     interface{} `json:"returnData"`
}

func main() {
	/*storage, _ := tstorage.NewStorage(
		tstorage.WithTimestampPrecision(tstorage.Seconds),
	)
	defer storage.Close()

	tm := time.Now()

	_ = storage.InsertRows([]tstorage.Row{
		{
			Metric:    "metric1",
			Labels:    []tstorage.Label{{"l1", "v1"}, {"l2", "v2"}},
			DataPoint: tstorage.DataPoint{Timestamp: tm.UnixMilli(), Value: 0.1},
		},
		{
			Metric:    "metric1",
			Labels:    []tstorage.Label{{"l1", "v1"}, {"l2", "v2"}},
			DataPoint: tstorage.DataPoint{Timestamp: tm.Add(time.Second).UnixMilli(), Value: 0.2},
		},
		{
			Metric:    "metric1",
			Labels:    []tstorage.Label{{"l1", "v1"}, {"l2", "v2"}},
			DataPoint: tstorage.DataPoint{Timestamp: tm.Add(2 * time.Second).UnixMilli(), Value: 0.3},
		},
		{
			Metric:    "metric1",
			Labels:    []tstorage.Label{{"l3", "v3"}, {"l2", "v2"}},
			DataPoint: tstorage.DataPoint{Timestamp: tm.Add(3 * time.Second).UnixMilli(), Value: 0.1},
		},
	})
	points, _ := storage.Select("metric1", []tstorage.Label{{"l1", "v1"}, {"l2", "v2"}}, tm.Unix(), tm.Add(3*time.Second).Unix())
	for _, p := range points {
		fmt.Printf("timestamp: %v, value: %v\n", p.Timestamp, p.Value)
		// => timestamp: 1600000000, value: 0.1
	}*/
	var res = &NormGetAgentCredentialRsp{}
	var resp NormResponse
	resp.Data = res
	err := unPackageResponse([]byte(reqJson), &resp)
	if err != nil {
		fmt.Printf("error: unPackageResponse failed. err:%v\n", err)
		return
	}
	data, _ := json.Marshal(resp)

	fmt.Printf("%s", string(data))
}

func unPackageResponse(data []byte, responseData *NormResponse) error {
	err := json.Unmarshal(data, responseData)
	if err != nil {
		fmt.Printf("error: unPackageResponse: json.Unmarshal failed. err:%v\n", err)
		return err
	}
	return nil
}
