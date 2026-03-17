# go-util

Shared Go utility functions for internal use at Maximumsoft.

```bash
go get github.com/Maximumsoft-Co-LTD/go-util
```

---

## แพ็กเกจ

### `jwt` — สร้างและตรวจสอบ JWT Token

```go
import "github.com/Maximumsoft-Co-LTD/go-util/jwt"
```

| Function | Signature | Description |
|---|---|---|
| `NewService` | `NewService(secret, issuer string) (Service, error)` | สร้าง JWT Service ใหม่ คืน error ถ้า secret ว่าง |

**Interface:**

```go
type Service interface {
    GenerateToken(username string, sessionID string, ttl time.Duration) (string, error)
    ValidateToken(token string) (*jwt.Token, error)
}
```

**ตัวอย่างการใช้งาน:**

```go
svc, err := jwt.NewService(os.Getenv("JWT_SECRET"), "my-app")
if err != nil {
    log.Fatal(err)
}

token, err := svc.GenerateToken("username", "session-123", 3*time.Hour)

parsed, err := svc.ValidateToken(token)
```

---

### `timex` — จัดการเวลา

```go
import "github.com/Maximumsoft-Co-LTD/go-util/timex"
```

| Function | Signature | Description |
|---|---|---|
| `StartOfDay` | `StartOfDay() (time.Time, error)` | คืนค่าต้นวันปัจจุบันในรูปแบบ UTC |
| `EndOfDay` | `EndOfDay() (time.Time, error)` | คืนค่าสิ้นวันปัจจุบันในรูปแบบ UTC |
| `TruncateDate` | `TruncateDate(now time.Time, date int) time.Time` | คืนค่า 00:00 ของวันที่ระบุในเดือนปัจจุบันหรือเดือนก่อนหน้า |
| `TruncateWeek` | `TruncateWeek(now time.Time, day int) time.Time` | คืนค่า 00:00 ของวันในสัปดาห์ที่ผ่านมาล่าสุด (0 = อาทิตย์, 6 = เสาร์) |
| `FormatDateTimeString` | `FormatDateTimeString(t time.Time) (string, error)` | แปลง time.Time เป็น string รูปแบบ ISO |
| `Bod` | `Bod(t time.Time) time.Time` | คืนค่าต้นวันแบบ UTC (17:00 UTC = 00:00 เวลาไทย) |
| `BodLocal` | `BodLocal(t time.Time) time.Time` | คืนค่าต้นวัน 00:00 ตาม timezone เครื่อง |
| `Eod` | `Eod(t time.Time) time.Time` | คืนค่าสิ้นวันแบบ exclusive UTC (17:00 UTC วันถัดไป) ใช้กับ `[Bod, Eod)` |
| `EodLocal` | `EodLocal(t time.Time) time.Time` | คืนค่าสิ้นวันแบบ exclusive local (00:00 วันถัดไป) ใช้กับ `[BodLocal, EodLocal)` |
| `TimeFormat` | `TimeFormat(t time.Time, filter string) string` | คืนค่าเวลาเป็น string ต่อกัน รับ `"day"`, `"hour"`, `"minute"`, `"second"` |
| `TimeNowOnly` | `TimeNowOnly(ref time.Time, t int64) int64` | ดึงชั่วโมง/นาที/วินาทีจาก Unix timestamp แล้วใช้กับวันที่ ref |
| `TimeInRange` | `TimeInRange(tTime, startTime, endTime time.Time) bool` | ตรวจสอบว่าเวลาอยู่ในช่วง `[startTime, endTime)` หรือไม่ |

**ค่าคงที่:** `LayoutISO` = `"2006-01-02 15:04:05"`

**ตัวอย่างการใช้งาน:**

```go
bod := timex.Bod(time.Now())
eod := timex.Eod(time.Now())
inRange := timex.TimeInRange(time.Now(), bod, eod) // true ถ้าอยู่ภายในวันนี้
```

---

### `slice` — ค้นหาใน Slice

```go
import "github.com/Maximumsoft-Co-LTD/go-util/slice"
```

| Function | Signature | Description |
|---|---|---|
| `Index` | `Index(limit int, predicate func(i int) bool) int` | คืนค่า index แรกที่ predicate คืนค่า true คืน -1 หากไม่พบ |

**ตัวอย่าง:**

```go
names := []string{"alice", "bob", "charlie"}
idx := slice.Index(len(names), func(i int) bool {
    return names[i] == "bob"
})
// idx = 1
```

---

### `random` — สุ่มตัวเลข/สตริง และ hashing

```go
import "github.com/Maximumsoft-Co-LTD/go-util/random"
```

| Function | Signature | Description |
|---|---|---|
| `UniqueInts` | `UniqueInts(count, max int) []int` | สุ่มเลขจำนวนเต็มไม่ซ้ำกันในช่วง [0, max) |
| `String` | `String(length int) (string, error)` | สุ่ม string ด้วย crypto/rand (alphanumeric) |
| `HashMD5` | `HashMD5(input string) string` | คืนค่า MD5 hex digest |
| `HashSHA1` | `HashSHA1(input string) string` | คืนค่า SHA1 hex digest |
| `Decimal` | `Decimal(precision, maxAttempts int) (float64, error)` | สุ่มทศนิยมในช่วง [0.01, 1.0) |

---

### `format` — จัดรูปแบบตัวเลข

```go
import "github.com/Maximumsoft-Co-LTD/go-util/format"
```

| Function | Signature | Description |
|---|---|---|
| `ToFixed` | `ToFixed(num float64, precision int) float64` | ปัดเศษ float64 ตามจำนวนทศนิยมที่กำหนด |
| `Decimal2` | `Decimal2(x float64) float64` | ปัดเศษ float64 เป็น 2 ตำแหน่ง |

---

## โครงสร้างโปรเจกต์

```
go-util/
├── go.mod
├── go.sum
├── README.md
├── jwt/
│   └── jwt.go         # สร้างและตรวจสอบ JWT Token (golang-jwt/jwt/v5)
├── timex/
│   └── time.go        # จัดการเวลา ตัดวัน/สัปดาห์/ฟอร์แมต/ตรวจสอบช่วงเวลา
├── slice/
│   └── slice.go       # ค้นหา index ใน slice
├── random/
│   └── random.go      # สุ่มตัวเลข/สตริง, MD5, SHA1
└── format/
    └── format.go      # ปัดเศษทศนิยม
```

---

## License

Internal use — Maximumsoft Co., LTD.
