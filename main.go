package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/gommon/log"
)

func main() {
	config := GetConfig()

	db, errcon := InitDB(config)
	if errcon != nil {
		fmt.Println(errcon.Error())
		panic("error database")
	}
	defer db.Close()

	handler := New(db)

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", handler.routeIndexGet)
	http.HandleFunc("/process", handler.routeSubmitPost)

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	fmt.Println("server started at localhost:" + port)

	http.ListenAndServe(port, nil)
}

func InitDB(config *AppConfig) (*sql.DB, error) {

	connectionString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		config.Database.Username,
		config.Database.Password,
		config.Database.Address,
		config.Database.Port,
		config.Database.Name,
	)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println(connectionString)
		return nil, err
	}

	InitMigrate(db)
	return db, nil
}
func InitMigrate(db *sql.DB) {
	// db.Migrator().DropTable(&entities.User{})
	// db.AutoMigrate(&entities.User{})
	var err error
	rows, err := db.Query("CREATE TABLE IF NOT EXISTS `zodiak` (`id` int NOT NULL AUTO_INCREMENT,`name` varchar(50) NOT NULL UNIQUE,`character` varchar(255) NOT NULL,`message` varchar(500) NOT NULL,`start_date` INT,end_date INT, PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=latin1;")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	rows.Close()
	zodiac := []map[string]interface{}{{"name": "Aries", "character": "Energik dan petualang , Pemberani , pelopor , Antusias , percaya diri, Dinamis", "message": "Ambil sebuah pekerjaan sampingan & dapatkan tambahan pendapatan. Sedikit uang lebih akan berguna di musin liburan ini. Temukan cara membuat atasan atau klien terkesan pada anda.", "start_date": 321, "end_date": 419},
		{"name": "Taurus", "character": "Sabar dan terpercaya, Mencintai , hangat, Punya(tujuan) dan Mencintai keamanan", "message": "Jangan bergantung pada org lain maka anda tak akan pernah kecewa. Bukan berarti anda tak bisa melangkah kedepan saat ini, tetapi semua harus anda lakukan dgn tangan sendiri.", "start_date": 420, "end_date": 520},
		{"name": "Gemini", "character": "Mudah beradaptasi , Komunikatif , Intelektual , Berjiwa muda dan hidup", "message": "Sesuatu terjadi tanpa sepengetahuan anda, tapi tidak perlu khawatir. Mungkin sebuah kejutan menanti anda, tetap diam & berpura-pura tak menyadari saat semua org menunjukkan perhatiannya.", "start_date": 521, "end_date": 622},
		{"name": "Cancer", "character": "Suasana Hati Tidak Menentu, Sentimentil, Setia, Penuh Perhatian, Sulit Memaafkan, Memiliki Daya Ingat Yang Kuat", "message": "Muncul pertanyaan, Dekatkah aku dengan Tuhan? Karier: Stres. Tugas belum selesai, suasana kantor sudah gerah. Rekan kerja keras kepala. Lebih produktif, jika fokus pada tugas sendiri daripada ngomel. Keuangan: Berhemat dululah. Asmara: Cancer lajang, cinta datang dengan cara menarik. Tes dulu air sebelum menyelam. Pria Cancer: Kesabaran Anda sedang diuji.", "start_date": 623, "end_date": 722},
		{"name": "Leo", "character": "Kreatif , antusias , Berpikiran terbuka , ekspansif , Setia dan mencintai", "message": "Akan sulit mencari jalan keluar jika anda berkeras kepala dlm masalah keluarga. Hati-hati, anda malah merugikan diri sendiri. Mengambil pendekatan baru bisa memberikan hasil lebih positif.", "start_date": 723, "end_date": 822},
		{"name": "Virgo", "character": "Rendah hati dan pemalu , Dapat(dipercaya) , Praktikal , rajin , Pandai dan analitikal", "message": "Semua rahasia berada di tangan & semua org ingin bergabung dgn anda. Aksi yg agresif akan membuat semua pesaing menoleh serta meninggikan status anda.", "start_date": 823, "end_date": 922},
		{"name": "Libra", "character": "Diplomatis ,Romantis , charming , Sosial , Idealistik dan senang kedamaian", "message": "Terlibatlah dgn persoalan yg menyangkut relasi yg lebih berumur. Anda harus selalu ingat untuk menjaga semua org yg telah lebih dahulu menjaga anda. Jangan pernah mundur.", "start_date": 923, "end_date": 1022},
		{"name": "Scorpius", "character": "Emosional , Intuitif , Bertenaga dan hasrat", "message": "Acara liburan kali ini akan membawa anda menemui cinta... Aktif berpartisipasi akan memperlihatkan betapa dinamis & menariknya anda. Bergerak lebih dulu & aksi yg terlihat jelas akan memberikan jawaban 'ya'.", "start_date": 1023, "end_date": 1121},
		{"name": "Sagittarius", "character": "Optimistik dan mencintai kebebasan , Selera humor yang tinggi , Jujur , terbuka dan Intellektual dan filosofikal", "message": "Sebuah pekerjaan yg diselesaikan dgn baik akan memberikan penghargaan. Semakin banyak yg anda kerjakan di hari ini, semakin baik untuk perkembangan masa depan.", "start_date": 1122, "end_date": 1221},
		{"name": "Capricornus", "character": "Pendiam, Rajin dan Ambisius, Materialis, Gengsi Tinggi, Suka Memerintah, Suka memperalat Orang Lain", "message": "Tak usah sedih atau kecewa, hanya kesalahpahaman saja, kok. Karier: Tak usah cari penyakit dengan perang melawan bos. Cari waktu yang tepat agar opini Anda bisa diterima.Sekarang, tampaknya Anda sendiri juga kurang siap atau pede. Keuangan: Cukup aman. Asmara: Kalau ingin diseriusi, kenali dulu dia lebih baik. Pria Capricorn: Waktu Anda cukup tersita untuk pekerjaan.", "start_date": 1222, "end_date": 119},
		{"name": "Aquarius", "character": "Humanis dan terbuka , Jujur , setia , Original , kreatif , Independent dan intelektual", "message": "Bukan hari yg tepat untuk membuka rahasia anda. Pihak lain mungkin akan tergganggu dgn sikap diam anda pada awalnya. Perubahan bisa membawa dampak baik, jadi tak perlu menahan yg datang.", "start_date": 120, "end_date": 218},
		{"name": "Pisces", "character": "Imaginatif dan sensitif , Baik , Tidak (egois), Intuitif dan simpatik", "message": "Menghabiskan waktu bersama teman atau tetangga akan memberikan informasi berguna. Berikan waktu untuk menikmati liburan dgn seseorang yg hati anda dambakan.", "start_date": 219, "end_date": 320}}
	for _, val := range zodiac {
		_, err := db.Exec("insert into zodiak (`name`,`character`,`message`,`start_date`,`end_date`) values (?, ?, ?, ?, ?)", val["name"].(string), val["character"].(string), val["message"].(string), val["start_date"].(int), val["end_date"].(int))
		if err != nil {
			errCheck := strings.Contains(err.Error(), "Duplicate entry")
			if !errCheck {
				fmt.Println(err.Error())
			}
		}
	}

}

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) routeIndexGet(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		var tmpl = template.Must(template.New("form").ParseFiles("views/view.html"))
		err := tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}
func (h *Handler) routeSubmitPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var tmpl = template.Must(template.New("result").ParseFiles("views/view.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var name = r.FormValue("name")
		var date = r.FormValue("date")

		dateArray := strings.Split(date, "-")
		var year int
		var month int
		var day int
		var dateInt int
		if len(dateArray) == 3 {
			jakarta, _ := time.LoadLocation("Asia/Jakarta")
			timeNow := time.Now().In(jakarta)

			layoutFormat := "2006-01-02 15:04 MST"
			fTime, _ := time.Parse(layoutFormat, date+" 23:59 WIB")
			year, month, day = calcAge(fTime, timeNow)

			dateString := dateArray[2]

			monthString := dateArray[1]
			if string(monthString[0]) == "0" {
				monthString = string(monthString[1])
			}
			date, _ := strconv.Atoi(monthString + dateString)

			dateInt = date
		}
		var zodiacId int
		var zodiacStartDate int
		var zodiacEndDate int
		var zodiacName string
		var zodiacCharacter string
		var zodiacMessage string

		if dateInt > 0 {
			if dateInt > 1221 || dateInt < 119 {
				err := h.db.QueryRow("SELECT * FROM zodiak WHERE start_date > ?", 1221).Scan(&zodiacId, &zodiacName, &zodiacCharacter, &zodiacMessage, &zodiacStartDate, &zodiacEndDate)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				err := h.db.QueryRow("SELECT * FROM zodiak WHERE start_date < ? AND end_date > ?", dateInt, dateInt).Scan(&zodiacId, &zodiacName, &zodiacCharacter, &zodiacMessage, &zodiacStartDate, &zodiacEndDate)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

		}

		// var message = r.Form.Get("message")

		var data = map[string]interface{}{"name": name, "id": zodiacId, "zodiak": zodiacName, "character": zodiacCharacter, "message": zodiacMessage, "start_date": zodiacStartDate, "end_date": zodiacEndDate, "date": date, "year": year, "month": month, "day": day}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func calcAge(bdate, cdate time.Time) (int, int, int) {
	if cdate.Year() < bdate.Year() {
		return -1, -1, -1
	}
	by, bm, bd := bdate.Date()
	cy, cm, cd := cdate.Date()
	cMonth := int(cm)
	bMont := int(bm)
	if cd < bd {
		cd += 30
		cm--
	}
	if cMonth < bMont {
		cMonth += 12
		cy--
	}
	month := cMonth - bMont
	return cy - by, month, cd - bd
}

type AppConfig struct {
	Port     int
	Database struct {
		Driver   string
		Name     string
		Address  string
		Port     string
		Username string
		Password string
	}
	GoogleClientID     string
	GoogleClientSecret string
}

var lock = &sync.Mutex{}
var appConfig *AppConfig

func GetConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = initConfig()
	}

	return appConfig
}

func initConfig() *AppConfig {

	portDB := os.Getenv("DB_PORT")

	var defaultConfig AppConfig
	port, errPort := strconv.Atoi(os.Getenv("APP_PORT"))
	if errPort != nil {
		log.Warn(errPort)
	}

	defaultConfig.Port = port
	defaultConfig.Database.Driver = os.Getenv("DB_DRIVER")
	defaultConfig.Database.Name = os.Getenv("DB_NAME")
	defaultConfig.Database.Address = os.Getenv("DB_ADDRESS")
	defaultConfig.Database.Port = portDB
	defaultConfig.Database.Username = os.Getenv("DB_USERNAME")
	defaultConfig.Database.Password = os.Getenv("DB_PASSWORD")

	return &defaultConfig
}
