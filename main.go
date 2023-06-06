package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Item struct {
	Id    int
	Name  string
	Stock int
	Price int
}

type Employee struct {
	Id      int
	Name    string
	Address string
	Role    string
	Status  string
}

// ========================= Fungsi Koneksi Database ========================
func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_store")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return db, nil
}

// ========================= Fungsi CRUD Data Barang ========================
func getAllItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		isSale := 1
		rows, err := db.Query("select id, name, stock, price from items where isSale = ?", isSale)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer rows.Close()

		var result []Item
		for rows.Next() {
			var item = Item{}
			var err = rows.Scan(&item.Id, &item.Name, &item.Stock, &item.Price)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			result = append(result, item)
		}

		if err = rows.Err(); err != nil {
			fmt.Println(err.Error())
		}

		var tmpl = template.Must(template.New("index").ParseFiles("view/index.html"))
		err = tmpl.Execute(w, result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Method not valid", http.StatusBadRequest) // jika method != GET
}

func insertItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("insert").ParseFiles("view/index.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "POST" {
		db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		var name = r.FormValue("name")
		var stock = r.FormValue("stock")
		var price = r.Form.Get("price")

		fmt.Println(name, stock, price)

		isSale := 1
		_, err = db.Exec("insert into items values (?,?,?,?,?)", nil, name, stock, price, isSale)
		if err != nil {
			http.Redirect(w, r, "/insert", http.StatusMovedPermanently)
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Tambah Data Sukses!")
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		db, err := connect()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var item = Item{}
		err = db.
			QueryRow("select id, name, stock, price from items where id = ?", id).
			Scan(&item.Id, &item.Name, &item.Stock, &item.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var tmpl = template.Must(template.New("update").ParseFiles("view/index.html"))
		err = tmpl.Execute(w, item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "POST" {
		id := r.URL.Query().Get("id")
		db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		var name = r.FormValue("name")
		var stock = r.FormValue("stock")
		var price = r.Form.Get("price")

		isSale := 1
		_, err = db.
			Exec("update items set name = ?, stock = ?, price = ?, isSale = ? where id = ?",
				name, stock, price, isSale, id)
		if err != nil {
			http.Redirect(w, r, "/update", http.StatusMovedPermanently)
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Update Data Sukses!")
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

func arsipItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		db, err := connect()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		isSale := 2
		_, err = db.
			Exec("update items set isSale = ? where id = ?", isSale, id)
		if err != nil {
			http.Redirect(w, r, "/arsip", http.StatusMovedPermanently)
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Arsip Data Sukses!")
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

func getArsipItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		isSale := 2
		rows, err := db.Query("select id, name, stock, price from items where isSale = ?", isSale)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer rows.Close()

		var result []Item
		for rows.Next() {
			var item = Item{}
			var err = rows.Scan(&item.Id, &item.Name, &item.Stock, &item.Price)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			result = append(result, item)
		}

		if err = rows.Err(); err != nil {
			fmt.Println(err.Error())
		}

		var tmpl = template.Must(template.New("arsip").ParseFiles("view/index.html"))

		err = tmpl.Execute(w, result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "POST" {
		id := r.URL.Query().Get("id")
		db, err := connect()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		isSale := 1
		_, err = db.
			Exec("update items set isSale = ? where id = ?", isSale, id)
		if err != nil {
			http.Redirect(w, r, "/arsip", http.StatusMovedPermanently)
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Tampilkan Data Sukses!")
		http.Redirect(w, r, "/arsip", http.StatusMovedPermanently)
		return
	}

	http.Error(w, "Method not valid", http.StatusBadRequest) // jika method != GET
}

// ========================= Fungsi CRUD Data Karyawan ========================
func getAllEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		isActive := "Aktif"
		rows, err := db.Query("select id, name, address, role from employees where isActive = ?", isActive)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var result []Employee
		for rows.Next() {
			var emply = Employee{}
			var err = rows.Scan(&emply.Id, &emply.Name, &emply.Address, &emply.Role)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			result = append(result, emply)
		}

		if err = rows.Err(); err != nil {
			fmt.Println(err.Error())
			return
		}

		var tmpl = template.Must(template.New("employee").ParseFiles("view/employee.html"))
		err = tmpl.Execute(w, result)
		if err != nil {
			fmt.Println()
			return
		}
	}
}

func insertEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("employee/insert").ParseFiles("view/employee.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if r.Method == "POST" {
		db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		var name = r.FormValue("name")
		var address = r.FormValue("address")
		var role = r.FormValue("role")

		isActive := "Aktif"
		_, err = db.Exec("insert into employees values (?,?,?,?,?)", nil, name, address, role, isActive)
		if err != nil {
			http.Redirect(w, r, "/employee/insert", http.StatusMovedPermanently)
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Simpan Data Karyawan Berhasil!")
		http.Redirect(w, r, "/employee", http.StatusMovedPermanently)
	}
}

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var id = r.URL.Query().Get("id")
		db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		var emply = Employee{}
		err = db.QueryRow("select name, address, role from employees where id = ?", id).Scan(&emply.Name, &emply.Address, &emply.Role)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var tmpl = template.Must(template.New("employee/update").ParseFiles("view/employee.html"))
		err = tmpl.Execute(w, emply)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "POST" {
		id := r.URL.Query().Get("id")
		db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		var name = r.FormValue("name")
		var address = r.FormValue("address")
		var role = r.FormValue("role")

		_, err = db.Exec("update employees set name = ?, address = ?, role = ? where id = ?", name, address, role, id)
		if err != nil {
			http.Redirect(w, r, "/emplyee/update", http.StatusMovedPermanently)
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Update Data Berhasil!")
		http.Redirect(w, r, "/employee", http.StatusMovedPermanently)
	}
}

func inactiveEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		inactive := "Nonaktif"
		_, err = db.Exec("update employees set isActive = ? where id = ?", inactive, id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Nonaktif karyawan berhasil!")
		http.Redirect(w, r, "/employee", http.StatusMovedPermanently)
	}
}

func arsipEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		rows, err := db.Query("select * from employees order by isActive desc")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer rows.Close()

		var result []Employee
		for rows.Next() {
			var emply = Employee{}
			var err = rows.Scan(&emply.Id, &emply.Name, &emply.Address, &emply.Role, &emply.Status)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			result = append(result, emply)
		}

		if err = rows.Err(); err != nil {
			fmt.Println(err.Error())
		}

		var tmpl = template.Must(template.New("employee/arsip").ParseFiles("view/employee.html"))
		err = tmpl.Execute(w, result)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}
}

func activeEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		active := "Aktif"
		_, err = db.Exec("update employees set isActive = ? where id = ?", active, id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Karyawan berhasil diaktifkan!")
		http.Redirect(w, r, "/employee/arsip", http.StatusMovedPermanently)
	}
}

func main() {
	//=============== handlefunc data barang ===============
	http.HandleFunc("/index", getAllItem)
	http.HandleFunc("/insert", insertItem)
	http.HandleFunc("/update", updateItem)
	http.HandleFunc("/arsipItemId", arsipItem)
	http.HandleFunc("/arsip", getArsipItem)

	//=============== handlefunc data karyawan ===============
	http.HandleFunc("/employee", getAllEmployee)
	http.HandleFunc("/employee/insert", insertEmployee)
	http.HandleFunc("/employee/update", updateEmployee)
	http.HandleFunc("/employee/inactive", inactiveEmployee)
	http.HandleFunc("/employee/arsip", arsipEmployee)
	http.HandleFunc("/employee/active", activeEmployee)

	//=============== handle start server ===============
	fmt.Println("Server started at localhost:9000")
	http.ListenAndServe(":9000", nil)
}
