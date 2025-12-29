## Pembuat
*Aleron Maulana Firjatullah*  
*NIM:* 434231040  
*Kelas:* TI-C5  
UAS Pemrograman Backend Lanjut  
Teknik Informatika – Universitas Airlangga

## Rules
1. jadi saat mahasiswa melaporan prestasi maka laporan tersebut masuk kedalam draft (status uploadnya menjadi "draft"), pada saat ini dosen wali bisa menerima atau menolak prestasi.

2. Mahasiswa dapat menghapus prestasi draft, untuk status ketika melakukan softdelete ditambahkan status deleted.

3. Dosen wali berelasi dengan mahasiswa, karena setiap dosen wali hanya bisa melihat mahasiswa bimbingannya saja.

4. Mahasiswa hanya bisa melihat dirinya sendiri (seperti pada nomor 2.3 Karakteristik Pengguna, yang menjelaskan terkait role, deskripsi, dan hak akses setiap role. Tertulis dalam PDF SRS).

5. Data sudah disediakan (dosen) tidak boleh merubah tabel (tabel berada pada nomor 3. Arsitektur Database) anda cukup mengikuti yang sudah ada di PDF SRS.

6. Pada Tabel achievement_references (3.1.7 Tabel achievement_references) terdapat mongo_achievement_id. Jadi terdapat Prestasi dan detail prestasi dengan relasi ke MongoDB.

7. Terdapat Functional Reporting & Analytics (pada nomor 4.5 Reporting & Analytics) yang bisa dilihat setiap role seperti mahasiswa hanya bisa melihat Reporting & Analytics miliknya sendiri, dosen wali bisa melihat mahasiswa bimbingannya, dan Admin bisa melihat seluruh Reporting & Analytics milik mahasiswa. Output Top mahasiswa berprestasi hanya tampil untuk dosen wali dan admin saja.

8. Semua endpoint harus ada dokumentasinya (swagger).

9. Kita juga perlu melakukan testing pada method post (wajib di test), method put, dan delete (yang memiliki parameter) juga perlu di test.

10. Kalau bisa ada validator untuk mengantisipasi input Password: or 0 = 0, jangan sampai terjadi SQL Injection, lakukan Penetration Testing (Pentesting) sebagai pencegahan.

11. Tidak Menggunakan ORM, Hanya Boleh Menggunakan Query Row.

12. Pada nomor 5. API Endpoints tepatnya pada 5.4 Achievements terdapat endpoint GET /api/v1/achievements/:id/history // Status history, pada endpoint ini menampilkan history semua achievement yang pernah diupload.

13. Pada route hanya diperbolehkan menggunakan Fiber, tidak boleh yang lain.

14. notification dihilangkan dari SRS, dosen saya memohon maaf terlewat menghilangkan notification dalam SRS (tidak perlu menggunakan notification), jika tidak salah atau kurang lebih (tolong anda cek ulang) dalam SRS yang tertulis terdapat notification yaitu pada bagian FR-004: Submit untuk Verifikas, dan FR-008: Reject Prestasi