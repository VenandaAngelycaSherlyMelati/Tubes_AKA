# Tubes_AKA
Tugas Besar Analisis Kompleksitas Algoritma. Algoritma seleksi: iteratif dan rekursif
Tugas besar ini disusun untuk menentukan kelulusan pada sebuah seleksi dengan kriteria tinggi badan tertentu sebagai standar passing grade.

- Algoritma Iteratif: Menggunakan LOOP FOR untuk menghitung kandidat yang lolos
- Algoritma Rekursif: Menggunakan pemanggilan fungsi rekursif untuk menghitung kandidat yang lolos

Aplikasi menghasilkan visualisasi grafik interaktif yang menampilkan perbandingan waktu eksekusi dan statistik hasil seleksi.

Fitur:
- Benchmark otomatis dengan berbagai ukuran input data
- Visualisasi grafik perbandingan waktu eksekusi (Chart.js)
- Grafik jumlah peserta yang lolos vs tidak lolos
- konfigurasi tinggi minimum dan maksimal ukuran data

Penggunaan:
1. Atur Parameter:
   - Tinggi Min (cm): Masukkan tinggi minimum untuk seleksi 
   - Max Data: Masukkan ukuran data maksimum untuk benchmark

2. Jalankan Benchmark:
   - Klik tombol "Jalankan Benchmark"
   - Aplikasi akan melakukan 10 iterasi benchmark dengan ukuran data yang bertambah secara bertahap
   - Tunggu hingga proses selesai

3. Analisis Hasil:
   - Lihat grafik perbandingan waktu eksekusi
   - Analisis grafik jumlah peserta yang lolos
   - Periksa statistik ringkasan untuk data agregat
  
Interpretasi Hasil:

Grafik Waktu Eksekusi
- Sumbu X: Ukuran data (jumlah kandidat)
- Sumbu Y: Waktu eksekusi dalam mikrodetik (μs)
- Garis Hijau: Waktu eksekusi algoritma iteratif
- Garis Merah: Waktu eksekusi algoritma rekursif

Grafik Peserta Lolos
- Bar Biru: Jumlah peserta yang memenuhi kriteria tinggi minimum
- Bar Oranye: Jumlah peserta yang tidak memenuhi kriteria

Statistik Ringkasan
- Rata-rata waktu iteratif dan rekursif
- Rata-rata jumlah peserta yang lolos
- Parameter tinggi minimum yang digunakan

Backend: Go (Golang)
Frontend: HTML5, CSS3, JavaScript
Visualisasi: Chart.js
Server: HTTP native Go
