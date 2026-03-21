package main

import (
	"bytes"
	"crypto/rand"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"
)

type Result struct {
	SizeLabel      string
	UpThroughput   float64 // MB/s (System capacity)
	UpLatAvg       float64 // Seconds (User wait time)
	DownThroughput float64
	DownLatAvg     float64
}

func main() {
	sizes := []int64{100 * 1024, 500 * 1024, 1024 * 1024, 2048 * 1024, 5120 * 1024}
	labels := []string{"100KB", "500KB", "1MB", "2MB", "5MB"}
	concurrencia := 10 // Simulate 10 users going up/down at the same time
	fmt.Println("STARTING IPFS BENCHMARK: LATENCY AND THROUGHPUT")
	fmt.Println("=================================================")

	var allResults []Result

	for idx, size := range sizes {
		fmt.Printf("Evaluating size: %s... ", labels[idx])

		// --- PHASE 1: UPLOAD ---
		var wgUp sync.WaitGroup
		startUp := time.Now()//Start the clock BEFORE creating the threads.
		latUpChan := make(chan time.Duration, concurrencia)
		cidChan := make(chan string, concurrencia)

		for i := 0; i < concurrencia; i++ {
			wgUp.Add(1)
			go func() {
				defer wgUp.Done()
				t0 := time.Now() // Start Latency
				cid, _ := uploadToIPFS(size)
				latUpChan <- time.Since(t0) // End of Latency
				cidChan <- cid
			}()
		}
		wgUp.Wait()// Wait until the LAST one finishes
		close(latUpChan)
		close(cidChan)
		durationUp := time.Since(startUp).Seconds() //Total time for throughput

		// --- PHASE 2: DOWNLOAD ---
		cids := []string{}
		for c := range cidChan { cids = append(cids, c) }

		var wgDown sync.WaitGroup
		startDown := time.Now()
		latDownChan := make(chan time.Duration, concurrencia)

		for _, c := range cids {
			wgDown.Add(1)
			go func(cidStr string) {
				defer wgDown.Done()
				t0 := time.Now()
				downloadFromIPFS(cidStr)
				latDownChan <- time.Since(t0)
			}(c)
		}
		wgDown.Wait()
		close(latDownChan)
		durationDown := time.Since(startDown).Seconds()

		
		totalMB := (float64(size) * float64(concurrencia)) / (1024 * 1024)// Total Megabytes
		
		var sumLUp, sumLDown float64
		for l := range latUpChan { sumLUp += l.Seconds() }// Sum of all threads
		for l := range latDownChan { sumLDown += l.Seconds() }//Sum of all threads

		allResults = append(allResults, Result{
			SizeLabel:      labels[idx],
			UpThroughput:   totalMB / durationUp,// <--- THROUGHPUT CALCULATION
			UpLatAvg:       sumLUp / float64(concurrencia),// <--- FINAL AVERAGE CALCULATION
			DownThroughput: totalMB / durationDown,// <--- THROUGHPUT CALCULATION
			DownLatAvg:     sumLDown / float64(concurrencia),// <--- FINAL AVERAGE CALCULATION
		})
		fmt.Println("Finalized")
	}
	saveToCSV(allResults)
}

func uploadToIPFS(size int64) (string, error) {
	data := make([]byte, size); rand.Read(data)//Reserve a space in RAM of the size
	body := &bytes.Buffer{}//fill with random bytes
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.pdf")//Create the file label
	io.Copy(part, bytes.NewReader(data))
	writer.Close()
	resp, err := http.Post("http://localhost:5001/api/v0/add", writer.FormDataContentType(), body) //Send the packet to the endpoint
	if err != nil { return "", err }
	defer resp.Body.Close()
	var res struct{ Hash string }// Create a temporary structure to store the Hash
	json.NewDecoder(resp.Body).Decode(&res)// Read the IPFS JSON response and extract the "Hash" field
	return res.Hash, nil// Returns the Hash
}

func downloadFromIPFS(cid string) {
	// POST is used with the '/cat' command, passing the CID (Hash) as an argument to obtain the file contents.
	resp, err := http.Post("http://localhost:5001/api/v0/cat?arg="+cid, "", nil)
	if err != nil { return }
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)// READ all the bytes coming from the server.
}

func saveToCSV(res []Result) {
	f, _ := os.Create("resultados_finales_ipfs.csv")
	w := csv.NewWriter(f)
	w.Write([]string{"Tamanio", "Up_Throughput_MBps", "Up_Lat_Seg", "Down_Throughput_MBps", "Down_Lat_Seg"})
	for _, r := range res {
		w.Write([]string{r.SizeLabel, fmt.Sprintf("%.2f", r.UpThroughput), fmt.Sprintf("%.4f", r.UpLatAvg), fmt.Sprintf("%.2f", r.DownThroughput), fmt.Sprintf("%.4f", r.DownLatAvg)})
	}
	w.Flush()
	f.Close()
}
