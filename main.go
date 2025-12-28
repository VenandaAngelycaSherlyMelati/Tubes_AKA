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
}

func generateCandidates(n int) []Candidate {
	rand.Seed(time.Now().UnixNano())
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
	countPassedIterative(candidates, minHeight)
	iterTime := time.Since(startIter).Microseconds()

	startRec := time.Now()
	countPassedRecursive(candidates, minHeight, 0)
	recTime := time.Since(startRec).Microseconds()

	return Result{
		Size:          size,
		IterativeTime: float64(iterTime),
		RecursiveTime: float64(recTime),
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
.container{max-width:900px;margin:0 auto;background:white;padding:20px;border-radius:10px}
h1{color:#333;text-align:center}
.form{background:#f9f9f9;padding:15px;border-radius:5px;margin:20px 0}
input{padding:8px;margin:5px;width:200px}
button{padding:10px 20px;background:#4CAF50;color:white;border:none;cursor:pointer;border-radius:5px}
button:hover{background:#45a049}
.chart{margin-top:20px}
</style>
</head>
<body>
<div class="container">
<h1>📊 Analisis Kompleksitas Algoritma</h1>
<h3>Seleksi Tinggi Badan: Iteratif vs Rekursif</h3>
<div class="form">
<label>Tinggi Min (cm): <input type="number" id="minHeight" value="170"></label>
<label>Max Data: <input type="number" id="maxSize" value="10000"></label>
<button onclick="runBenchmark()">Jalankan Benchmark</button>
</div>
<div class="chart"><canvas id="chart"></canvas></div>
<div id="stats"></div>
</div>
<script>
let chart=null;
async function runBenchmark(){
const minHeight=document.getElementById('minHeight').value;
const maxSize=document.getElementById('maxSize').value;
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
data.push({size:parseInt(p[0]),iter:parseFloat(p[1]),rec:parseFloat(p[2])});
}
if(chart)chart.destroy();
const ctx=document.getElementById('chart').getContext('2d');
chart=new Chart(ctx,{
type:'line',
data:{
labels:data.map(d=>d.size),
datasets:[
{label:'Iteratif (μs)',data:data.map(d=>d.iter),borderColor:'#4CAF50',fill:false},
{label:'Rekursif (μs)',data:data.map(d=>d.rec),borderColor:'#f44336',fill:false}
]},
options:{responsive:true,scales:{y:{beginAtZero:true}}}
});
const iterAvg=data.reduce((s,d)=>s+d.iter,0)/data.length;
const recAvg=data.reduce((s,d)=>s+d.rec,0)/data.length;
document.getElementById('stats').innerHTML='<p><b>Rata-rata Iteratif:</b> '+iterAvg.toFixed(2)+' μs</p><p><b>Rata-rata Rekursif:</b> '+recAvg.toFixed(2)+' μs</p>';
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
	sb.WriteString("Size,IterativeTime,RecursiveTime\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("%d,%.2f,%.2f\n", r.Size, r.IterativeTime, r.RecursiveTime))
	}
	fmt.Fprint(w, sb.String())
}

func main() {
	minHeight := 170.0
	var results []Result

	result := benchmarkSingle(1000, minHeight)
	results = append(results, result)

	fmt.Println("\n Starting web server on http://localhost:8080")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/benchmark", benchmarkHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
