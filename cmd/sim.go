package cmd

// func call(req *http.Request) ([]byte, int, error) {
// 	req.Header.Add("X-API-KEY", "plt_pOh1xjmIEDfE68zgFr7djsc2rvzjSMotlo2ZIJXdA")
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, 0, err
// 	}
//
// 	defer res.Body.Close()
// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, res.StatusCode, err
// 	}
//
// 	// var response map[string]any
// 	// if err := json.Unmarshal(body, &response); err != nil {
// 	// 	return nil, 0, err
// 	// }
//
// 	return body, res.StatusCode, nil
// }
//
// func callPost(url string, data any) ([]byte, int, error) {
// 	payload, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, 0, err
// 	}
//
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	return call(req)
// }
//
// func main() {
// 	// err := sim.SimulateForProject(1)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
//
// 	startTime := time.Now()
// 	for i := range 10 {
// 		events := sim.SimulateSessionForUser(i, startTime)
// 		resp, status, err := callPost("http://pl.localhost:9999/api/ingest", events)
//
// 		if err == nil {
// 			fmt.Printf("status %d - resp: %v\n", status, resp)
// 		}
//
// 	}
// }
