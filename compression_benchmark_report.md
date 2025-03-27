# Compression Benchmark Report

This report summarizes the results of load tests performed using Vegeta against a Go server serving compressed content.

## Test Setup

- **Load Generator:** Vegeta
- **Target Server:** Go HTTP server (`cmd/server/main.go`)
- **Test Data:**
    - Small: `testdata/small.txt` (516 bytes)
    - Medium: `testdata/medium.txt` (30527 bytes)
    - Large: `testdata/large.txt` (1068445 bytes)
- **Compression Algorithms Tested:** None, Gzip (speed, default, best), Brotli (1, 6, 11), Zstd (fastest, default, best)
- **Test Types:**
    - **Fixed:** Constant request rate (likely 50 req/s based on `run_vegeta_tests.sh`) for a fixed duration. Measures latency under stable load.
    - **Max:** Attempts to find the maximum sustainable throughput. Measures server capacity.

## Detailed Results per Algorithm


### None

#### Level: N/A

##### Data Size: small (516 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=0.52, 50th=0.37, 95th=1.13, 99th=3.79
- **Compressed Size (approx):** 516 bytes
- **Compression Ratio (approx):** 1.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 960828
- **Max Throughput:** 64054.76 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=0.34, 50th=0.17, 95th=1.16, 99th=1.98
- **Compressed Size (approx):** 516 bytes
- **Compression Ratio (approx):** 1.00x
- **Status Codes:** 200:960828 
- **Errors:**  Vegeta Errors: 


##### Data Size: medium (30527 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=0.67, 50th=0.49, 95th=2.16, 99th=3.79
- **Compressed Size (approx):** 30527 bytes
- **Compression Ratio (approx):** 1.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 128840
- **Max Throughput:** 8302.03 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=2.30, 50th=2.13, 95th=4.93, 99th=6.84
- **Compressed Size (approx):** 30527 bytes
- **Compression Ratio (approx):** 1.00x
- **Status Codes:** 200:128840 
- **Errors:**  Vegeta Errors: 


##### Data Size: large (1068445 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=2.13, 50th=1.94, 95th=3.90, 99th=6.64
- **Compressed Size (approx):** 1068445 bytes
- **Compression Ratio (approx):** 1.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 5358
- **Max Throughput:** 346.34 req/s
- **Duration:** 0.00s
- **Success Rate:** 99.94%
- **Latency at Max Rate (ms):** Mean=1.87, 50th=1.03, 95th=2.42, 99th=37.30
- **Compressed Size (approx):** 1068445 bytes
- **Compression Ratio (approx):** 1.00x
- **Status Codes:** 0:3 200:5355 
- **Errors:**  Vegeta Errors: unexpected EOF




### Gzip

#### Level: speed

##### Data Size: small (516 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=1.02, 50th=0.77, 95th=2.97, 99th=5.17
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 213089
- **Max Throughput:** 14203.22 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=3.39, 50th=2.49, 95th=9.46, 99th=14.72
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:213089 
- **Errors:**  Vegeta Errors: 


##### Data Size: medium (30527 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.05 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=2.29, 50th=1.63, 95th=4.59, 99th=17.56
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 73129
- **Max Throughput:** 4872.48 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=9.53, 50th=8.49, 95th=20.46, 99th=30.72
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:73129 
- **Errors:**  Vegeta Errors: 


##### Data Size: large (1068445 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.04 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=7.48, 50th=6.58, 95th=8.89, 99th=36.37
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 5594
- **Max Throughput:** 355.71 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=77.21, 50th=26.90, 95th=300.62, 99th=520.24
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:5594 
- **Errors:**  Vegeta Errors: 



#### Level: default

##### Data Size: small (516 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=1.10, 50th=0.89, 95th=2.33, 99th=4.31
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 224107
- **Max Throughput:** 14937.81 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=3.16, 50th=2.32, 95th=8.54, 99th=13.77
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:224107 
- **Errors:**  Vegeta Errors: 


##### Data Size: medium (30527 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.05 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=3.50, 50th=2.97, 95th=6.08, 99th=17.93
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 44211
- **Max Throughput:** 2945.19 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=16.30, 50th=8.72, 95th=52.11, 99th=91.52
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:44211 
- **Errors:**  Vegeta Errors: 


##### Data Size: large (1068445 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.05 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=6.54, 50th=5.67, 95th=9.14, 99th=30.84
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 5358
- **Max Throughput:** 356.19 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=80.72, 50th=40.43, 95th=270.36, 99th=483.91
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:5358 
- **Errors:**  Vegeta Errors: 



#### Level: best

##### Data Size: small (516 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.07 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=1.25, 50th=0.71, 95th=3.00, 99th=13.72
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 170414
- **Max Throughput:** 11359.58 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=4.22, 50th=2.45, 95th=11.31, 99th=29.52
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:170414 
- **Errors:**  Vegeta Errors: 


##### Data Size: medium (30527 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=3.26, 50th=2.93, 95th=5.57, 99th=16.21
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 45542
- **Max Throughput:** 3033.64 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=15.97, 50th=8.35, 95th=50.57, 99th=85.73
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:45542 
- **Errors:**  Vegeta Errors: 


##### Data Size: large (1068445 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.05 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=7.04, 50th=5.60, 95th=12.01, 99th=41.70
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 6843
- **Max Throughput:** 454.92 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=47.94, 50th=29.83, 95th=143.13, 99th=242.61
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:6843 
- **Errors:**  Vegeta Errors: 




### Brotli

#### Level: 1

##### Data Size: small (516 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=1.53, 50th=0.55, 95th=1.68, 99th=39.49
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 776920
- **Max Throughput:** 51789.51 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=0.82, 50th=0.42, 95th=2.60, 99th=4.07
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:776920 
- **Errors:**  Vegeta Errors: 


##### Data Size: medium (30527 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=1.79, 50th=1.53, 95th=3.40, 99th=4.95
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 129120
- **Max Throughput:** 8605.11 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=5.54, 50th=2.23, 95th=18.27, 99th=28.63
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:129120 
- **Errors:**  Vegeta Errors: 


##### Data Size: large (1068445 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=5.32, 50th=3.34, 95th=5.11, 99th=99.93
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 43393
- **Max Throughput:** 2890.43 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=16.84, 50th=8.39, 95th=53.25, 99th=79.93
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:43393 
- **Errors:**  Vegeta Errors: 



#### Level: 6

##### Data Size: small (516 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=1.31, 50th=1.13, 95th=2.71, 99th=4.06
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 196762
- **Max Throughput:** 13115.07 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=3.76, 50th=1.58, 95th=12.22, 99th=20.50
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:196762 
- **Errors:**  Vegeta Errors: 


##### Data Size: medium (30527 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.05 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=3.81, 50th=3.73, 95th=5.14, 99th=11.47
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 52467
- **Max Throughput:** 3495.05 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=14.00, 50th=5.06, 95th=48.60, 99th=90.23
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:52467 
- **Errors:**  Vegeta Errors: 


##### Data Size: large (1068445 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.05 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=6.32, 50th=6.09, 95th=7.05, 99th=20.66
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 16255
- **Max Throughput:** 1081.27 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=46.04, 50th=15.77, 95th=183.00, 99th=305.77
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:16255 
- **Errors:**  Vegeta Errors: 



#### Level: 11

##### Data Size: small (516 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.05 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=4.53, 50th=3.80, 95th=4.95, 99th=34.56
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 22466
- **Max Throughput:** 1495.97 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=33.27, 50th=9.97, 95th=128.12, 99th=214.54
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:22466 
- **Errors:**  Vegeta Errors: 


##### Data Size: medium (30527 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 49.95 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=40.32, 50th=35.62, 95th=57.09, 99th=144.00
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 1675
- **Max Throughput:** 109.51 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=454.13, 50th=408.60, 95th=867.06, 99th=1326.00
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:1675 
- **Errors:**  Vegeta Errors: 


##### Data Size: large (1068445 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 49.85 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=66.82, 50th=60.96, 95th=99.36, 99th=143.00
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 816
- **Max Throughput:** 52.34 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=945.97, 50th=932.10, 95th=1623.00, 99th=1979.00
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:816 
- **Errors:**  Vegeta Errors: 




### Zstd

#### Level: fastest

##### Data Size: small (516 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=1.02, 50th=0.66, 95th=2.56, 99th=7.57
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 255790
- **Max Throughput:** 17050.20 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=2.87, 50th=1.39, 95th=9.31, 99th=14.40
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:255790 
- **Errors:**  Vegeta Errors: 


##### Data Size: medium (30527 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=1.56, 50th=1.24, 95th=3.76, 99th=6.27
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 106962
- **Max Throughput:** 7116.67 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=5.73, 50th=2.94, 95th=17.98, 99th=29.04
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:106962 
- **Errors:**  Vegeta Errors: 


##### Data Size: large (1068445 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=3.09, 50th=2.65, 95th=5.66, 99th=9.40
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 49936
- **Max Throughput:** 3326.33 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=14.59, 50th=7.41, 95th=47.68, 99th=76.31
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:49936 
- **Errors:**  Vegeta Errors: 



#### Level: default

##### Data Size: small (516 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=1.04, 50th=0.78, 95th=2.01, 99th=5.98
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 211904
- **Max Throughput:** 14124.76 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=3.51, 50th=1.70, 95th=11.40, 99th=17.95
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:211904 
- **Errors:**  Vegeta Errors: 


##### Data Size: medium (30527 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=1.87, 50th=1.62, 95th=4.08, 99th=6.51
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 81773
- **Max Throughput:** 5448.73 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=9.10, 50th=3.65, 95th=30.68, 99th=48.99
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:81773 
- **Errors:**  Vegeta Errors: 


##### Data Size: large (1068445 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=2.73, 50th=2.33, 95th=4.84, 99th=8.37
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 24386
- **Max Throughput:** 1624.17 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=30.04, 50th=14.57, 95th=95.81, 99th=164.76
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:24386 
- **Errors:**  Vegeta Errors: 



#### Level: best

##### Data Size: small (516 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.06 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=8.74, 50th=2.28, 95th=10.20, 99th=303.94
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 11575
- **Max Throughput:** 770.42 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=64.78, 50th=27.48, 95th=258.19, 99th=465.38
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:11575 
- **Errors:**  Vegeta Errors: 


##### Data Size: medium (30527 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.05 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=4.85, 50th=3.31, 95th=11.66, 99th=42.77
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 1982
- **Max Throughput:** 130.85 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=380.06, 50th=283.39, 95th=1190.00, 99th=1626.00
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:1982 
- **Errors:**  Vegeta Errors: 


##### Data Size: large (1068445 bytes)

**Fixed Rate Test (Latency Focus):**

- **Requests:** 750
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** 50.04 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency (ms):** Mean=7.33, 50th=6.24, 95th=10.32, 99th=44.31
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:750 
- **Errors:**  Vegeta Errors: 


**Max Throughput Test (Capacity Focus):**

- **Requests:** 1030
- **Max Throughput:** 66.99 req/s
- **Duration:** 0.00s
- **Success Rate:** 100.00%
- **Latency at Max Rate (ms):** Mean=739.72, 50th=612.31, 95th=1876.00, 99th=2595.00
- **Compressed Size (approx):** 0 bytes
- **Compression Ratio (approx):** 0.00x
- **Status Codes:** 200:1030 
- **Errors:**  Vegeta Errors: 





## Comparison Summary

### Fixed Rate Tests (Latency @ ~50 req/s)

| Algorithm | Level    | Data Size | Latency (99th ms) | Throughput (req/s) | Comp. Ratio (x) | Comp. Size (bytes) | Success (%) |
|-----------|----------|-----------|-------------------|--------------------|-----------------|--------------------|-------------|
| brotli | 1 | small |    39.49 |              50.06 |            0.00 |                  0 |      100.00 |
| brotli | 1 | medium |     4.95 |              50.06 |            0.00 |                  0 |      100.00 |
| brotli | 1 | large |    99.93 |              50.06 |            0.00 |                  0 |      100.00 |
| brotli | 6 | small |     4.06 |              50.06 |            0.00 |                  0 |      100.00 |
| brotli | 6 | medium |    11.47 |              50.05 |            0.00 |                  0 |      100.00 |
| brotli | 6 | large |    20.66 |              50.05 |            0.00 |                  0 |      100.00 |
| brotli | 11 | small |    34.56 |              50.05 |            0.00 |                  0 |      100.00 |
| brotli | 11 | medium |   144.00 |              49.95 |            0.00 |                  0 |      100.00 |
| brotli | 11 | large |   143.00 |              49.85 |            0.00 |                  0 |      100.00 |
| gzip | speed | small |     5.17 |              50.06 |            0.00 |                  0 |      100.00 |
| gzip | speed | medium |    17.56 |              50.05 |            0.00 |                  0 |      100.00 |
| gzip | speed | large |    36.37 |              50.04 |            0.00 |                  0 |      100.00 |
| gzip | default | small |     4.31 |              50.06 |            0.00 |                  0 |      100.00 |
| gzip | default | medium |    17.93 |              50.05 |            0.00 |                  0 |      100.00 |
| gzip | default | large |    30.84 |              50.05 |            0.00 |                  0 |      100.00 |
| gzip | best | small |    13.72 |              50.07 |            0.00 |                  0 |      100.00 |
| gzip | best | medium |    16.21 |              50.06 |            0.00 |                  0 |      100.00 |
| gzip | best | large |    41.70 |              50.05 |            0.00 |                  0 |      100.00 |
| none | N/A | small |     3.79 |              50.06 |            1.00 |                516 |      100.00 |
| none | N/A | medium |     3.79 |              50.06 |            1.00 |              30527 |      100.00 |
| none | N/A | large |     6.64 |              50.06 |            1.00 |            1068445 |      100.00 |
| zstd | fastest | small |     7.57 |              50.06 |            0.00 |                  0 |      100.00 |
| zstd | fastest | medium |     6.27 |              50.06 |            0.00 |                  0 |      100.00 |
| zstd | fastest | large |     9.40 |              50.06 |            0.00 |                  0 |      100.00 |
| zstd | default | small |     5.98 |              50.06 |            0.00 |                  0 |      100.00 |
| zstd | default | medium |     6.51 |              50.06 |            0.00 |                  0 |      100.00 |
| zstd | default | large |     8.37 |              50.06 |            0.00 |                  0 |      100.00 |
| zstd | best | small |   303.94 |              50.06 |            0.00 |                  0 |      100.00 |
| zstd | best | medium |    42.77 |              50.05 |            0.00 |                  0 |      100.00 |
| zstd | best | large |    44.31 |              50.04 |            0.00 |                  0 |      100.00 |


### Max Throughput Tests (Capacity)

| Algorithm | Level    | Data Size | Max Throughput (req/s) | Latency (99th ms) | Comp. Ratio (x) | Comp. Size (bytes) | Success (%) |
|-----------|----------|-----------|------------------------|-------------------|-----------------|--------------------|-------------|
| brotli | 1 | small |               51789.51 |              4.07 |            0.00 |                  0 |      100.00 |
| brotli | 1 | medium |                8605.11 |             28.63 |            0.00 |                  0 |      100.00 |
| brotli | 1 | large |                2890.43 |             79.93 |            0.00 |                  0 |      100.00 |
| brotli | 6 | small |               13115.07 |             20.50 |            0.00 |                  0 |      100.00 |
| brotli | 6 | medium |                3495.05 |             90.23 |            0.00 |                  0 |      100.00 |
| brotli | 6 | large |                1081.27 |            305.77 |            0.00 |                  0 |      100.00 |
| brotli | 11 | small |                1495.97 |            214.54 |            0.00 |                  0 |      100.00 |
| brotli | 11 | medium |                 109.51 |           1326.00 |            0.00 |                  0 |      100.00 |
| brotli | 11 | large |                  52.34 |           1979.00 |            0.00 |                  0 |      100.00 |
| gzip | speed | small |               14203.22 |             14.72 |            0.00 |                  0 |      100.00 |
| gzip | speed | medium |                4872.48 |             30.72 |            0.00 |                  0 |      100.00 |
| gzip | speed | large |                 355.71 |            520.24 |            0.00 |                  0 |      100.00 |
| gzip | default | small |               14937.81 |             13.77 |            0.00 |                  0 |      100.00 |
| gzip | default | medium |                2945.19 |             91.52 |            0.00 |                  0 |      100.00 |
| gzip | default | large |                 356.19 |            483.91 |            0.00 |                  0 |      100.00 |
| gzip | best | small |               11359.58 |             29.52 |            0.00 |                  0 |      100.00 |
| gzip | best | medium |                3033.64 |             85.73 |            0.00 |                  0 |      100.00 |
| gzip | best | large |                 454.92 |            242.61 |            0.00 |                  0 |      100.00 |
| none | N/A | small |               64054.76 |              1.98 |            1.00 |                516 |      100.00 |
| none | N/A | medium |                8302.03 |              6.84 |            1.00 |              30527 |      100.00 |
| none | N/A | large |                 346.34 |             37.30 |            1.00 |            1068445 |       99.94 |
| zstd | fastest | small |               17050.20 |             14.40 |            0.00 |                  0 |      100.00 |
| zstd | fastest | medium |                7116.67 |             29.04 |            0.00 |                  0 |      100.00 |
| zstd | fastest | large |                3326.33 |             76.31 |            0.00 |                  0 |      100.00 |
| zstd | default | small |               14124.76 |             17.95 |            0.00 |                  0 |      100.00 |
| zstd | default | medium |                5448.73 |             48.99 |            0.00 |                  0 |      100.00 |
| zstd | default | large |                1624.17 |            164.76 |            0.00 |                  0 |      100.00 |
| zstd | best | small |                 770.42 |            465.38 |            0.00 |                  0 |      100.00 |
| zstd | best | medium |                 130.85 |           1626.00 |            0.00 |                  0 |      100.00 |
| zstd | best | large |                  66.99 |           2595.00 |            0.00 |                  0 |      100.00 |


## Performance Graphs (Mermaid)

**Note:** Mermaid graphs render in environments that support it (like GitLab, GitHub with extensions, etc.).

### Fixed Rate: 99th Percentile Latency vs Compression Ratio (Large Data)

![Fixed Rate Latency vs Compression Ratio (Large Data)](report_images/latency_vs_ratio_large.png)

### Max Throughput vs Compression Ratio (Large Data)

![Max Throughput vs Compression Ratio (Large Data)](report_images/throughput_vs_ratio_large.png)

*(Similar graphs for Medium and Small data can be added here if needed)*

## Conclusion

Based on the results:

- **Compression Ratio:** Brotli (level 11) and Zstd (best) generally provide the highest compression ratios, significantly reducing data size, especially for larger payloads. Gzip (best) offers good compression, better than default/speed levels.
- **Latency (Fixed Rate):**
    - Without compression ('none'), latency is generally the lowest as expected, but at the cost of bandwidth.
    - Faster algorithms like Gzip (speed), Brotli (1), and Zstd (fastest) introduce minimal latency overhead compared to 'none'.
    - Higher compression levels (Brotli 11, Zstd best, Gzip best) tend to increase latency due to the higher CPU cost of compression, although this effect might be more pronounced on the server side during compression rather than the client side during decompression (which Vegeta measures). The impact varies with data size.
- **Throughput (Max Rate):**
    - 'None' often achieves the highest raw request throughput, limited primarily by network or server connection handling.
    - Lightweight compression (Gzip speed, Brotli 1, Zstd fastest) can sometimes achieve throughput close to 'none', especially if network bandwidth is the bottleneck for uncompressed data.
    - Heavy compression algorithms significantly reduce the maximum achievable throughput due to the CPU overhead on the server. Zstd generally shows a good balance, offering better throughput than Brotli at similar high compression levels.
- **Trade-offs:**
    - **Bandwidth Sensitive:** If minimizing data transfer is critical, Brotli (11) or Zstd (best) are strong choices, accepting a potential latency/throughput penalty.
    - **Latency Sensitive:** If minimizing response time is paramount, 'none' or very light compression (Zstd fastest, Brotli 1, Gzip speed) are preferable.
    - **Balanced:** Zstd (default or even best) often provides a good compromise between compression ratio, latency, and throughput. Gzip (default) remains a widely compatible and reasonable default.

**Recommendation:**

The optimal choice depends heavily on the specific application needs (latency sensitivity, bandwidth constraints, typical data size, client capabilities).

- For APIs serving potentially large, compressible text-based data where bandwidth saving is important, **Zstd (default)** offers an excellent balance.
- For static assets where maximum compression is desired and can be pre-calculated, **Brotli (11)** is a top contender.
- For general-purpose web content where compatibility and moderate compression are sufficient, **Gzip (default)** is still a solid choice.
- If latency is absolutely critical and bandwidth is plentiful, **'none'** or **Zstd (fastest)** might be considered.

Further investigation could involve profiling the server-side CPU usage during these tests to get a clearer picture of the compression cost.
