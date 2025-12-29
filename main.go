package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Candidate struct {
	ID     int
	Height float64
}

type Result struct {
	Size          int
	IterativeTime float64
	RecursiveTime float64
	PassedCount   int
	TotalCount    int
}

func generateCandidates(n int) []Candidate {
	candidates := make([]Candidate, n)
	for i := 0; i < n; i++ {
		candidates[i] = Candidate{
			ID:     i + 1,
			Height: 150.0 + rand.Float64()*50.0,
		}
	}
	return candidates
}

func countPassedIterative(candidates []Candidate, minHeight float64) int {
	count := 0
	for i := 0; i < len(candidates); i++ {
		if candidates[i].Height >= minHeight {
			count++
		}
	}
	return count
}

func countPassedRecursive(candidates []Candidate, minHeight float64, index int) int {
	if index >= len(candidates) {
		return 0
	}

	if candidates[index].Height >= minHeight {
		return 1 + countPassedRecursive(candidates, minHeight, index+1)
	}
	return countPassedRecursive(candidates, minHeight, index+1)
}

func benchmarkSingle(size int, minHeight float64) Result {
	candidates := generateCandidates(size)

	startIter := time.Now()
	passedCount := countPassedIterative(candidates, minHeight)
	iterTime := time.Since(startIter).Microseconds()

	startRec := time.Now()
	countPassedRecursive(candidates, minHeight, 0)
	recTime := time.Since(startRec).Microseconds()

	return Result{
		Size:          size,
		IterativeTime: float64(iterTime),
		RecursiveTime: float64(recTime),
		PassedCount:   passedCount,
		TotalCount:    size,
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
<title>Analisis Kompleksitas Algoritma</title>
<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
<style>
body{font-family:Arial;margin:20px;background:#f0f0f0}
.container{max-width:1100px;margin:0 auto;background:white;padding:20px;border-radius:10px;box-shadow:0 2px 10px rgba(0,0,0,0.1)}
h1{color:#333;text-align:center;margin-bottom:5px}
h3{color:#666;text-align:center;margin-top:0}
.form{background:#f9f9f9;padding:15px;border-radius:5px;margin:20px 0;display:flex;gap:10px;align-items:center;justify-content:center;flex-wrap:wrap}
input{padding:8px;width:150px;border:1px solid #ddd;border-radius:3px}
button{padding:10px 20px;background:#4CAF50;color:white;border:none;cursor:pointer;border-radius:5px;font-weight:bold}
button:hover{background:#45a049}
button:disabled{background:#ccc;cursor:not-allowed}
.charts-container{display:grid;grid-template-columns:1fr 1fr;gap:20px;margin-top:20px}
.chart-box{background:#f9f9f9;padding:15px;border-radius:5px}
.chart-box h4{margin-top:0;color:#333;text-align:center}
.stats{background:#e8f5e9;padding:15px;border-radius:5px;margin:20px 0}
.stats h4{margin-top:0;color:#2e7d32}
.stat-row{display:flex;justify-content:space-between;margin:10px 0;padding:8px;background:white;border-radius:3px}
.stat-label{font-weight:bold;color:#555}
.stat-value{color:#2e7d32;font-weight:bold}
@media (max-width: 768px){
.charts-container{grid-template-columns:1fr}
}
</style>
</head>
<body>
<div class="container">
<h1> Analisis Kompleksitas Algoritma</h1>
<h3>Seleksi Tinggi Badan: Iteratif vs Rekursif</h3>
<div class="form">
<label>Tinggi Min (cm): <input type="number" id="minHeight" value="170"></label>
<label>Max Data: <input type="number" id="maxSize" value="10000"></label>
<button id="btnRun" onclick="runBenchmark()">Jalankan Benchmark</button>
</div>

<div id="stats" style="display:none" class="stats">
<h4> Ringkasan Statistik</h4>
<div class="stat-row">
<span class="stat-label">Rata-rata Iteratif:</span>
<span class="stat-value" id="avgIter">-</span>
</div>
<div class="stat-row">
<span class="stat-label">Rata-rata Rekursif:</span>
<span class="stat-value" id="avgRec">-</span>
</div>
<div class="stat-row">
<span class="stat-label">Rata-rata Peserta Lolos:</span>
<span class="stat-value" id="avgPassed">-</span>
</div>
<div class="stat-row">
<span class="stat-label">Tinggi Minimum:</span>
<span class="stat-value" id="minHeightDisplay">-</span>
</div>
</div>

<div class="charts-container">
<div class="chart-box">
<h4> Perbandingan Waktu Eksekusi</h4>
<canvas id="timeChart"></canvas>
</div>
<div class="chart-box">
<h4> Jumlah Peserta Lolos per Ukuran Data</h4>
<canvas id="passedChart"></canvas>
</div>
</div>
</div>

<script>
let timeChart=null;
let passedChart=null;

async function runBenchmark(){
const minHeight=document.getElementById('minHeight').value;
const maxSize=document.getElementById('maxSize').value;

const btn=document.getElementById('btnRun');
btn.disabled=true;
btn.textContent=' Memproses...';

try {
const res=await fetch('/benchmark',{
method:'POST',
headers:{'Content-Type':'application/x-www-form-urlencoded'},
body:'minHeight='+minHeight+'&maxSize='+maxSize
});
const text=await res.text();
const lines=text.trim().split('\n');
const data=[];

for(let i=1;i<lines.length;i++){
const p=lines[i].split(',');
data.push({
size:parseInt(p[0]),
iter:parseFloat(p[1]),
rec:parseFloat(p[2]),
passed:parseInt(p[3]),
total:parseInt(p[4])
});
}

// Update Time Chart
if(timeChart)timeChart.destroy();
const ctx1=document.getElementById('timeChart').getContext('2d');
timeChart=new Chart(ctx1,{
type:'line',
data:{
labels:data.map(d=>d.size),
datasets:[
{label:'Iteratif (μs)',data:data.map(d=>d.iter),borderColor:'#4CAF50',backgroundColor:'rgba(76,175,80,0.1)',fill:true,tension:0.4},
{label:'Rekursif (μs)',data:data.map(d=>d.rec),borderColor:'#f44336',backgroundColor:'rgba(244,67,54,0.1)',fill:true,tension:0.4}
]},
options:{responsive:true,plugins:{legend:{position:'top'}},scales:{y:{beginAtZero:true,title:{display:true,text:'Waktu (μs)'}}}}
});

// Update Passed Chart
if(passedChart)passedChart.destroy();
const ctx2=document.getElementById('passedChart').getContext('2d');
passedChart=new Chart(ctx2,{
type:'bar',
data:{
labels:data.map(d=>d.size),
datasets:[
{label:'Peserta Lolos',data:data.map(d=>d.passed),backgroundColor:'#2196F3'},
{label:'Tidak Lolos',data:data.map(d=>d.total-d.passed),backgroundColor:'#FF9800'}
]},
options:{responsive:true,plugins:{legend:{position:'top'}},scales:{x:{stacked:true},y:{stacked:true,beginAtZero:true,title:{display:true,text:'Jumlah Peserta'}}}}
});

// Calculate statistics
const iterAvg=data.reduce((s,d)=>s+d.iter,0)/data.length;
const recAvg=data.reduce((s,d)=>s+d.rec,0)/data.length;
const passedAvg=data.reduce((s,d)=>s+d.passed,0)/data.length;

document.getElementById('avgIter').textContent=iterAvg.toFixed(2)+' μs';
document.getElementById('avgRec').textContent=recAvg.toFixed(2)+' μs';
document.getElementById('avgPassed').textContent=passedAvg.toFixed(0)+' peserta';
document.getElementById('minHeightDisplay').textContent=minHeight+' cm';
document.getElementById('stats').style.display='block';

} catch(error) {
console.error('Error:', error);
alert('Terjadi error saat menjalankan benchmark');
}

btn.disabled=false;
btn.textContent='Jalankan Benchmark';
}
</script>
</body>
</html>`

	t, _ := template.New("index").Parse(tmpl)
	t.Execute(w, nil)
}

func benchmarkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	minHeight, _ := strconv.ParseFloat(r.FormValue("minHeight"), 64)
	maxSize, _ := strconv.Atoi(r.FormValue("maxSize"))

	sizes := []int{}
	step := maxSize / 10
	for i := step; i <= maxSize; i += step {
		sizes = append(sizes, i)
	}

	var results []Result
	for _, size := range sizes {
		result := benchmarkSingle(size, minHeight)
		results = append(results, result)
	}

	w.Header().Set("Content-Type", "text/csv")
	var sb strings.Builder
	sb.WriteString("Size,IterativeTime,RecursiveTime,PassedCount,TotalCount\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("%d,%.2f,%.2f,%d,%d\n",
			r.Size, r.IterativeTime, r.RecursiveTime, r.PassedCount, r.TotalCount))
	}
	fmt.Fprint(w, sb.String())
}

func main() {
	fmt.Println("\nStarting web server on http://localhost:8080")

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/benchmark", benchmarkHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
