package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// JSON 구조와 매핑될 구조체 정의
type DataStore struct {
	National []string `json:"national"`
	Ulsan    []string `json:"ulsan"`
}

type Response struct {
	Destination string `json:"destination"`
}

var districts DataStore

func main() {
	rand.Seed(time.Now().UnixNano())

	// 1. 외부 JSON 파일 로드
	if err := loadDistrictsJSON(); err != nil {
		fmt.Printf("❌ 데이터 로드 실패: %v\n", err)
		return
	}

	// 2. 라우터 설정
	http.HandleFunc("/api/draw", drawHandler)
	http.HandleFunc("/", indexHandler)

	fmt.Println("🚀 [지역 추첨기] 서버가 작동 중입니다. (포트: 58000)")
	if err := http.ListenAndServe(":58000", nil); err != nil {
		fmt.Printf("서버 실행 에러: %v\n", err)
	}
}

// districts.json 파일을 읽어서 메모리에 적재하는 함수
func loadDistrictsJSON() error {
	file, err := os.ReadFile("districts.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &districts)
}

// REST API 핸들러
func drawHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	option := r.URL.Query().Get("option")
	var targetList []string

	if option == "ulsan" {
		targetList = districts.Ulsan
	} else {
		targetList = districts.National
	}

	// 데이터 예외 처리
	if len(targetList) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Destination: "데이터가 없습니다."})
		return
	}

	randomIndex := rand.Intn(len(targetList))
	res := Response{Destination: targetList[randomIndex]}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// 외부 HTML 파일을 읽어 사용자에게 서빙하는 핸들러
func indexHandler(w http.ResponseWriter, r *http.Request) {
	htmlData, err := os.ReadFile("index.html")
	if err != nil {
		http.Error(w, "HTML 파일을 찾을 수 없습니다.", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(htmlData)
}
