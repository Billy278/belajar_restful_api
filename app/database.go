package app

import (
	"belajar-golang-api/helper"
	"database/sql"
	"time"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/belajar_golang_restful_api")
	helper.PanicIfError(err)
	//set minimal jumlah koneksi ke databse disni 10
	db.SetMaxIdleConns(10)
	//mengeset maksimal koneksi yg dijinkan disi contohnya 100
	db.SetMaxOpenConns(20)
	//disini di set apabila lagi bengong dan koneksi yg digunakan lebih dari minimum
	//maka koneksi akan dimatikan /koneksi yg dinukan tidak digunakan lagi
	db.SetConnMaxIdleTime(10 * time.Minute)
	// disini diset bahwa apabila sudah mencapai batas
	//dimana di contohkan 60 menit akan di hapus semua koneksi
	//beserta koneksi minimumnya dan di buat lagi koneksi minimun nya
	//jadi merest semua koneksi dan membuat koneksi miimun yg baru
	db.SetConnMaxLifetime(60 * time.Minute)
	return db

}
