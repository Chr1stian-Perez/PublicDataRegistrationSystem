package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

const pdfPath = "./Certificate.pdf" 
// RoundResult stores the performance metrics collected in each round of testing
type RoundResult struct {
	Label       string// Round Name
	TargetTPS   int// Target Payload Sent to System
	ActualTPS   float64// Actual Throughput Processed by System
	MinLat      float64// Minimum Detected Latency (seconds)
	AvgLat      float64// Average Latency (seconds)
	MaxLat      float64// Maximum Detected Latency (seconds)
	Success     int// Number of Successful Transactions
	Fail        int// Number of Failed Transactions
	SuccessRate float64// Success Rate

}
// uploadRealPDF simulates uploading documents to IPFS.
// Generates a unique CID by modifying the file content with a timestamp and ID.
func uploadRealPDF(id int) (string, error) { 
	content, err := os.ReadFile(pdfPath)
	if err != nil {
		return "", fmt.Errorf("error leyendo PDF: %v", err)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, _ := writer.CreateFormFile("file", "certificate.pdf")
	part.Write(content)
	// Invisible modification to vary the CID on each upload
	fmt.Fprintf(part, "\n%%BENCH-ID-%d-%d", time.Now().UnixNano(), id)
	writer.Close()
	// Communication with the local IPFS node
	resp, err := http.Post("http://localhost:5001/api/v0/add", writer.FormDataContentType(), &body) 
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct{ Hash string }
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.Hash, nil
}
// runFullSuite manages the sequential flow of the injected payload rounds
func runFullSuite(contract *client.Contract, orgIdx int) {
	rondas := []struct {
		label string
		tps   int
		delay int // Tiempo de espera para estabilización de recursos
	}{
		{"Round 1 - 5 TPS (Low Load)", 5, 10},
		{"Round 2 - 20 TPS (Medium Load)", 20, 10},
		{"Round 3 - 50 TPS (Stress)", 50, 10},
		{"Round 4 - 100 TPS (Saturation)", 100, 10},
		{"Round 5 - 200 TPS (Saturation)", 200, 10},
	}

	os.Remove("resultados_tesis_final.csv")

	for _, r := range rondas {
		fmt.Printf("\n--- Starting: %s ---\n", r.label)
		fmt.Printf("Waiting %d sec to stabilize nodes...\n", r.delay)
		time.Sleep(time.Duration(r.delay) * time.Second)

		//We run each round for 20 seconds to avoid filling up the IPFS disk
		res := runControlledRound(contract, orgIdx, r.tps, 20, r.label)
		guardarEnCSV(res)
		
		fmt.Printf("✓ Completed. Successful: %d | Failed: %d | Avg Lat: %.3fs\n", res.Success, res.Fail, res.AvgLat)
	}
	fmt.Println("\n Benchmark completed! Results saved in resultados_tesis_final.csv")
}
// runControlledRound executes the concurrency and measurement logic.
func runControlledRound(contract *client.Contract, orgIdx int, targetTPS int, durationSecs int, label string) RoundResult {
	var wg sync.WaitGroup
	var mu sync.Mutex
	exitosas, fallidas := 0, 0
	latencias := make([]time.Duration, 0)
	// Ticker to control the rate of transaction injection per second (Rate Limiting)
	ticker := time.NewTicker(time.Second / time.Duration(targetTPS))
	defer ticker.Stop()

	stop := time.After(time.Duration(durationSecs) * time.Second)
	inicioRound := time.Now()

	count := 0
Loop:
	for {
		select {
		case <-stop:
			break Loop
		case <-ticker.C:
			count++
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				// Phase 1: Persistence in IPFS
				realCID, errIpfs := uploadRealPDF(id)
				if errIpfs != nil {
					mu.Lock()
					fallidas++
					mu.Unlock()
					return
				}

				// Phase 2: Registration in Hyperledger Fabric (Consensus and writing)
				citizenID := fmt.Sprintf("CI-BENCH-%d-%d", time.Now().UnixNano(), id)
				
				
				recordData := `{"first_names":"BENCH","last_names":"USER","birth_date":"2026-03-20","nationality":"ECUADORIAN"}`

				t0 := time.Now()
				// Call to Tx_Register
				_, err := contract.SubmitTransaction("Tx_Register", citizenID, "IDENTITY", recordData, realCID)
				lat := time.Since(t0)// Individual latency measurement (Round-trip time)

				mu.Lock()
				if err != nil {
					fallidas++
					
				} else {
					exitosas++
					latencias = append(latencias, lat)
				}
				mu.Unlock()
			}(count)
		}
	}

	wg.Wait()
	duracionReal := time.Since(inicioRound).Seconds()
	// ... (Calculation of statistics Min/Avg/Max) ...
	var sumaLat float64
	minLat := math.MaxFloat64
	maxLat := 0.0

	if len(latencias) > 0 {
		for _, l := range latencias {
			s := l.Seconds()
			sumaLat += s
			if s < minLat { minLat = s }
			if s > maxLat { maxLat = s }
		}
	} else {
		minLat = 0
	}

	avgLat := 0.0
	successRate := 0.0
	if (exitosas + fallidas) > 0 {
		if exitosas > 0 {
			avgLat = sumaLat / float64(exitosas)
		}
		successRate = (float64(exitosas) / float64(exitosas+fallidas)) * 100
	}

	return RoundResult{/* Final structure with calculations */
		Label: label, TargetTPS: targetTPS,
		ActualTPS:   float64(exitosas) / duracionReal,
		MinLat:      minLat,
		AvgLat:      avgLat,
		MaxLat:      maxLat,
		Success:     exitosas,
		Fail:        fallidas,
		SuccessRate: successRate,
	}
}

func guardarEnCSV(r RoundResult) {
	file, _ := os.OpenFile("resultados_tesis_final.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()

	info, _ := file.Stat()
	if info.Size() == 0 {
		w.Write([]string{"Round", "Target_TPS", "Actual_TPS", "Min_Lat_S", "Avg_Lat_S", "Max_Lat_S", "Success", "Fail", "Rate"})
	}
	w.Write([]string{
		r.Label,
		fmt.Sprintf("%d", r.TargetTPS),
		fmt.Sprintf("%.2f", r.ActualTPS),
		fmt.Sprintf("%.3f", r.MinLat),
		fmt.Sprintf("%.3f", r.AvgLat),
		fmt.Sprintf("%.3f", r.MaxLat),
		fmt.Sprintf("%d", r.Success),
		fmt.Sprintf("%d", r.Fail),
		fmt.Sprintf("%.2f%%", r.SuccessRate),
	})
}
