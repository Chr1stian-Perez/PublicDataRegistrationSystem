
# Benchmark.go

## This module is a benchmarking tool designed to evaluate the performance of the integrated Blockchain-IPFS architecture.

* Controlled Injected Load (Rate Limiting): Uses Go channels and Tickers to inject an exact transactions-per-second (TPS) rate, enabling evaluation of system behavior under different load levels (Low, Medium, Stress, Saturation).

* Real IPFS Integration: This script performs actual PDF document uploads on every transaction, forcing the generation of new identifiers (CIDs) and evaluating the real-world performance of distributed storage.

* Resilience and Stabilization: Implements cool-down periods between rounds to ensure that CPU and memory usage on the (Docker) nodes returns to baseline, eliminating statistical noise.

* Data Generation: Automatically exports results to a formatted CSV file for subsequent statistical analysis and graphical visualization (Python/Matplotlib).

## Collected Metrics

* Actual Throughput: Effective number of confirmed transactions per second.
Average 
* Latency: Average time for the network to reach consensus and persist the data.
* Success Rate: Percentage of valid transactions versus failed attempts (critical for measuring stability under saturation).

## How to run

To run the benchmark, we have an option in the application menu within the civil registry organization, specifically option number 6, and this will generate a csv file with the metrics for your study.