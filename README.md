# go-util

Shared Go utility functions for internal use at Maximumsoft.

```bash
go get github.com/Maximumsoft-Co-LTD/go-util
```

---

## แพ็กเกจ

### `slice` — Generic Slice Helpers (53 functions)

```go
import "github.com/Maximumsoft-Co-LTD/go-util/slice"
```

#### Search & Index

| Function | Signature | Description |
|---|---|---|
| `Index` | `Index(limit int, predicate func(i int) bool) int` | คืน index แรกที่ predicate คืนค่า true, คืน -1 หากไม่พบ |
| `IndexOf` | `IndexOf[T comparable](collection []T, element T) int` | คืน index แรกของ element, คืน -1 หากไม่พบ |
| `LastIndexOf` | `LastIndexOf[T comparable](collection []T, element T) int` | คืน index สุดท้ายของ element, คืน -1 หากไม่พบ |
| `Contains` | `Contains[T comparable](collection []T, element T) bool` | ตรวจว่า element อยู่ใน slice หรือไม่ |
| `ContainsBy` | `ContainsBy[T any](collection []T, predicate func(T) bool) bool` | ตรวจว่ามี element ที่ตรงเงื่อนไขหรือไม่ |

#### Find

| Function | Signature | Description |
|---|---|---|
| `Find` | `Find[T any](collection []T, predicate func(T) bool) (T, bool)` | คืน element แรกที่ตรงเงื่อนไข |
| `FindIndexOf` | `FindIndexOf[T any](collection []T, predicate func(T) bool) (T, int, bool)` | คืน element, index, และ bool ของตัวแรกที่ตรงเงื่อนไข |
| `FindLastIndexOf` | `FindLastIndexOf[T any](collection []T, predicate func(T) bool) (T, int, bool)` | เหมือน FindIndexOf แต่ค้นจากท้าย |
| `FindOrElse` | `FindOrElse[T any](collection []T, fallback T, predicate func(T) bool) T` | Find พร้อม fallback ถ้าไม่เจอ |
| `FindUniques` | `FindUniques[T comparable](collection []T) []T` | คืน elements ที่ปรากฏครั้งเดียว |
| `FindUniquesBy` | `FindUniquesBy[T any, K comparable](collection []T, keyFn func(T) K) []T` | คืน elements ที่ key ปรากฏครั้งเดียว |
| `FindDuplicates` | `FindDuplicates[T comparable](collection []T) []T` | คืน elements ที่ซ้ำ (1 ชุดต่อค่า) |
| `FindDuplicatesBy` | `FindDuplicatesBy[T any, K comparable](collection []T, keyFn func(T) K) []T` | คืน elements ที่ key ซ้ำ |

#### Transform

| Function | Signature | Description |
|---|---|---|
| `Filter` | `Filter[T any](collection []T, predicate func(T) bool) []T` | คืน elements ที่ตรงเงื่อนไข |
| `Reject` | `Reject[T any](collection []T, predicate func(T) bool) []T` | คืน elements ที่ไม่ตรงเงื่อนไข (ตรงข้าม Filter) |
| `Map` | `Map[T any, U any](collection []T, iteratee func(T) U) []U` | แปลง element แต่ละตัว |
| `FlatMap` | `FlatMap[T any, U any](collection []T, iteratee func(T) []U) []U` | แปลงแต่ละ element เป็น slice แล้ว flatten |
| `Reduce` | `Reduce[T any, U any](collection []T, accumulator func(U, T) U, initial U) U` | รวมค่าทั้ง slice เป็นค่าเดียว |
| `ForEach` | `ForEach[T any](collection []T, iteratee func(T, int))` | วนลูปทุก element |
| `Compact` | `Compact[T comparable](collection []T) []T` | ลบค่า zero-value ออก |
| `Reverse` | `Reverse[T any](collection []T) []T` | กลับลำดับ slice |
| `Without` | `Without[T comparable](collection []T, exclude ...T) []T` | คืน slice ที่ไม่รวม values ที่ระบุ |

#### Unique & Group

| Function | Signature | Description |
|---|---|---|
| `Unique` | `Unique[T comparable](collection []T) []T` | ลบค่าซ้ำ |
| `UniqueBy` | `UniqueBy[T any, K comparable](collection []T, keyFn func(T) K) []T` | ลบค่าซ้ำตาม key function |
| `GroupBy` | `GroupBy[T any, K comparable](collection []T, keyFn func(T) K) map[K][]T` | จัดกลุ่มตาม key function |
| `KeyBy` | `KeyBy[T any, K comparable](collection []T, keyFn func(T) K) map[K]T` | แปลง slice เป็น map ตาม key function |
| `Partition` | `Partition[T any](collection []T, predicate func(T) bool) ([]T, []T)` | แยก slice เป็น 2 กลุ่ม (ตรง/ไม่ตรงเงื่อนไข) |
| `Chunk` | `Chunk[T any](collection []T, size int) [][]T` | แบ่ง slice เป็น chunk ตามขนาดที่กำหนด |
| `Flatten` | `Flatten[T any](collections [][]T) []T` | รวมหลาย slices เป็น slice เดียว |

#### Set Operations

| Function | Signature | Description |
|---|---|---|
| `Intersect` | `Intersect[T comparable](list1, list2 []T) []T` | คืน elements ที่มีอยู่ในทั้งสอง slices |
| `Difference` | `Difference[T comparable](list1, list2 []T) []T` | คืน elements ที่อยู่ใน list1 แต่ไม่อยู่ใน list2 |
| `Union` | `Union[T comparable](lists ...[]T) []T` | รวม unique elements จากหลาย slices |

#### Aggregate (Count / Min / Max)

| Function | Signature | Description |
|---|---|---|
| `Count` | `Count[T comparable](collection []T, value T) int` | นับจำนวน element ที่เท่ากับ value |
| `CountBy` | `CountBy[T any](collection []T, predicate func(T) bool) int` | นับจำนวนที่ตรงเงื่อนไข |
| `Min` | `Min[T cmp.Ordered](collection []T) T` | ค่าน้อยสุด |
| `MinIndex` | `MinIndex[T cmp.Ordered](collection []T) (T, int)` | ค่าน้อยสุดพร้อม index |
| `MinBy` | `MinBy[T any](collection []T, comparison func(a, b T) bool) T` | ค่าน้อยสุดด้วย comparison function |
| `MinIndexBy` | `MinIndexBy[T any](collection []T, comparison func(a, b T) bool) (T, int)` | ค่าน้อยสุดพร้อม index ด้วย comparison |
| `Max` | `Max[T cmp.Ordered](collection []T) T` | ค่ามากสุด |
| `MaxIndex` | `MaxIndex[T cmp.Ordered](collection []T) (T, int)` | ค่ามากสุดพร้อม index |
| `MaxBy` | `MaxBy[T any](collection []T, comparison func(a, b T) bool) T` | ค่ามากสุดด้วย comparison function |
| `MaxIndexBy` | `MaxIndexBy[T any](collection []T, comparison func(a, b T) bool) (T, int)` | ค่ามากสุดพร้อม index ด้วย comparison |

#### Time

| Function | Signature | Description |
|---|---|---|
| `Latest` | `Latest(collection []time.Time) time.Time` | เวลาล่าสุดจาก slice ของ time.Time |
| `LatestBy` | `LatestBy[T any](collection []T, timeFn func(T) time.Time) T` | element ที่มีเวลาล่าสุด |

#### Boolean Checks

| Function | Signature | Description |
|---|---|---|
| `Every` | `Every[T any](collection []T, predicate func(T) bool) bool` | ตรวจว่าทุกตัวตรงเงื่อนไข |
| `Some` | `Some[T any](collection []T, predicate func(T) bool) bool` | ตรวจว่ามีอย่างน้อยหนึ่งตัวตรงเงื่อนไข |
| `None` | `None[T any](collection []T, predicate func(T) bool) bool` | ตรวจว่าไม่มีตัวใดตรงเงื่อนไข |

#### Element Access

| Function | Signature | Description |
|---|---|---|
| `First` | `First[T any](collection []T) (T, bool)` | คืน element แรก พร้อม bool |
| `FirstOrEmpty` | `FirstOrEmpty[T any](collection []T) T` | element แรก หรือ zero value |
| `Last` | `Last[T any](collection []T) (T, bool)` | คืน element สุดท้าย พร้อม bool |
| `LastOrEmpty` | `LastOrEmpty[T any](collection []T) T` | element สุดท้าย หรือ zero value |

**ตัวอย่าง:**

```go
names := []string{"alice", "bob", "charlie", "bob"}

slice.Contains(names, "bob")        // true
slice.IndexOf(names, "bob")         // 1
slice.LastIndexOf(names, "bob")     // 3
slice.Unique(names)                  // ["alice", "bob", "charlie"]
slice.Filter(names, func(s string) bool { return len(s) > 3 })
                                     // ["alice", "charlie"]

nums := []int{1, 2, 3, 4, 5}
doubled := slice.Map(nums, func(n int) int { return n * 2 })
                                     // [2, 4, 6, 8, 10]
sum := slice.Reduce(nums, func(acc int, n int) int { return acc + n }, 0)
                                     // 15
slice.Min(nums)                      // 1
slice.Max(nums)                      // 5
maxVal, maxIdx := slice.MaxIndex(nums) // 5, 4

// FindDuplicates / FindUniques
data := []int{1, 2, 2, 3, 3, 4}
slice.FindDuplicates(data)           // [2, 3]
slice.FindUniques(data)              // [1, 4]

// FindOrElse — ค้นหาพร้อม fallback
user := slice.FindOrElse(users, defaultUser, func(u User) bool {
    return u.ID == targetID
})

// MinBy / MaxBy — เปรียบเทียบ struct ด้วย custom logic
type Product struct { Name string; Price float64 }
products := []Product{{Name: "A", Price: 100}, {Name: "B", Price: 50}}
cheapest := slice.MinBy(products, func(a, b Product) bool { return a.Price < b.Price })
// {Name: "B", Price: 50}

// LatestBy — หา record ล่าสุดตาม field เวลา
type Order struct { ID string; CreatedAt time.Time }
latest := slice.LatestBy(orders, func(o Order) time.Time { return o.CreatedAt })

// Set operations
allowed := []string{"read", "write", "admin"}
userPerms := []string{"read", "delete"}
slice.Intersect(allowed, userPerms)   // ["read"]
slice.Difference(userPerms, allowed)  // ["delete"]

// Partition — แยก pass/fail ในรอบเดียว
pass, fail := slice.Partition(nums, func(n int) bool { return n > 3 })
// pass: [4, 5], fail: [1, 2, 3]
```

---

### `mapx` — Generic Map Helpers (15 functions)

```go
import "github.com/Maximumsoft-Co-LTD/go-util/mapx"
```

| Function | Signature | Description |
|---|---|---|
| `Keys` | `Keys[K comparable, V any](m map[K]V) []K` | คืนค่า keys ทั้งหมด |
| `Values` | `Values[K comparable, V any](m map[K]V) []V` | คืนค่า values ทั้งหมด |
| `ContainsKey` | `ContainsKey[K comparable, V any](m map[K]V, key K) bool` | ตรวจว่า key มีอยู่หรือไม่ |
| `Pick` | `Pick[K comparable, V any](m map[K]V, keys []K) map[K]V` | คืน map เฉพาะ keys ที่ระบุ |
| `Omit` | `Omit[K comparable, V any](m map[K]V, keys []K) map[K]V` | คืน map ที่ไม่รวม keys ที่ระบุ |
| `Merge` | `Merge[K comparable, V any](maps ...map[K]V) map[K]V` | รวมหลาย maps เข้าด้วยกัน (ตัวหลังทับตัวก่อน) |
| `MapValues` | `MapValues[K comparable, V any, U any](m map[K]V, iteratee func(V) U) map[K]U` | แปลง values ด้วย function |
| `MapKeys` | `MapKeys[K comparable, V any, J comparable](m map[K]V, iteratee func(K) J) map[J]V` | แปลง keys ด้วย function |
| `Filter` | `Filter[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V` | คืน map ที่ตรงเงื่อนไข |
| `ForEach` | `ForEach[K comparable, V any](m map[K]V, iteratee func(K, V))` | วนลูปทุก entry |
| `Invert` | `Invert[K comparable, V comparable](m map[K]V) map[V]K` | สลับ key กับ value |
| `Entries` | `Entries[K comparable, V any](m map[K]V) []Entry[K, V]` | แปลง map เป็น slice ของ key-value pairs |
| `FromEntries` | `FromEntries[K comparable, V any](entries []Entry[K, V]) map[K]V` | แปลง key-value pairs กลับเป็น map |
| `FindKey` | `FindKey[K comparable, V comparable](m map[K]V, value V) (K, bool)` | หา key จาก value ที่ตรงกัน |
| `FindKeyBy` | `FindKeyBy[K comparable, V any](m map[K]V, predicate func(K, V) bool) (K, bool)` | หา key ด้วย predicate |

**ตัวอย่าง:**

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}

mapx.Keys(m)                              // ["a", "b", "c"]
mapx.Pick(m, []string{"a", "c"})          // map[a:1 c:3]
mapx.Omit(m, []string{"b"})               // map[a:1 c:3]

m2 := map[string]int{"c": 99, "d": 4}
mapx.Merge(m, m2)                         // map[a:1 b:2 c:99 d:4]

doubled := mapx.MapValues(m, func(v int) int { return v * 2 })
                                           // map[a:2 b:4 c:6]

// FindKey — หา key จาก value
key, ok := mapx.FindKey(m, 2)              // "b", true

// FindKeyBy — หา key ด้วย predicate
key, ok = mapx.FindKeyBy(m, func(k string, v int) bool { return v > 2 })
                                           // "c", true
```

---

### `val` — Value & Nil Checking Helpers (5 functions)

```go
import "github.com/Maximumsoft-Co-LTD/go-util/val"
```

| Function | Signature | Description |
|---|---|---|
| `IsNil` | `IsNil(v any) bool` | ตรวจว่าเป็น nil หรือไม่ (รองรับ typed nil interface) |
| `IsNotNil` | `IsNotNil(v any) bool` | ตรวจว่าไม่ใช่ nil |
| `Empty` | `Empty[T any]() T` | คืน zero value ของ type T |
| `IsEmpty` | `IsEmpty[T comparable](v T) bool` | ตรวจว่าเป็น zero value หรือไม่ |
| `IsNotEmpty` | `IsNotEmpty[T comparable](v T) bool` | ตรวจว่าไม่ใช่ zero value |

**ตัวอย่าง:**

```go
// IsNil — จัดการ Go interface nil gotcha ได้ถูกต้อง
var p *User = nil
var i interface{} = p
i == nil                // false (Go gotcha!)
val.IsNil(i)            // true  (ถูกต้อง)
val.IsNotNil(i)         // false

// Empty — สร้าง zero value
val.Empty[string]()     // ""
val.Empty[int]()        // 0
val.Empty[*User]()      // nil

// IsEmpty / IsNotEmpty
val.IsEmpty("")          // true
val.IsEmpty(0)           // true
val.IsEmpty("hello")     // false
val.IsNotEmpty(42)       // true
```

---

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

### `timex` — จัดการเวลา (12 functions)

```go
import "github.com/Maximumsoft-Co-LTD/go-util/timex"
```

| Function | Signature | Description |
|---|---|---|
| `StartOfDay` | `StartOfDay() time.Time` | คืนค่าต้นวันปัจจุบัน (UTC) |
| `EndOfDay` | `EndOfDay() time.Time` | คืนค่าสิ้นวันปัจจุบัน (UTC) |
| `TruncateDate` | `TruncateDate(now time.Time, date int) time.Time` | คืนค่า 00:00 ของวันที่ระบุ |
| `TruncateWeek` | `TruncateWeek(now time.Time, day int) time.Time` | คืนค่า 00:00 ของวันในสัปดาห์ที่ผ่านมาล่าสุด |
| `FormatDateTimeString` | `FormatDateTimeString(t time.Time) string` | แปลง time.Time เป็น string ISO |
| `Bod` | `Bod(t time.Time) time.Time` | ต้นวัน UTC (17:00 UTC = 00:00 ไทย) |
| `BodLocal` | `BodLocal(t time.Time) time.Time` | ต้นวัน 00:00 ตาม timezone เครื่อง |
| `Eod` | `Eod(t time.Time) time.Time` | สิ้นวัน exclusive UTC ใช้กับ `[Bod, Eod)` |
| `EodLocal` | `EodLocal(t time.Time) time.Time` | สิ้นวัน exclusive local ใช้กับ `[BodLocal, EodLocal)` |
| `TimeFormat` | `TimeFormat(t time.Time, filter string) string` | คืนค่าเวลาเป็น string |
| `TimeNowOnly` | `TimeNowOnly(ref time.Time, t int64) int64` | ดึงชั่วโมง/นาที/วินาทีจาก Unix timestamp |
| `TimeInRange` | `TimeInRange(tTime, startTime, endTime time.Time) bool` | ตรวจว่าอยู่ในช่วง `[start, end)` |

**ค่าคงที่:** `LayoutISO` = `"2006-01-02 15:04:05"`

**ตัวอย่าง:**

```go
bod := timex.Bod(time.Now())
eod := timex.Eod(time.Now())
inRange := timex.TimeInRange(time.Now(), bod, eod) // true ถ้าอยู่ภายในวันนี้
```

---

### `random` — สุ่มตัวเลข/สตริง และ hashing (5 functions)

```go
import "github.com/Maximumsoft-Co-LTD/go-util/random"
```

| Function | Signature | Description |
|---|---|---|
| `UniqueInts` | `UniqueInts(count, max int) []int` | สุ่มเลขจำนวนเต็มไม่ซ้ำกัน [0, max) |
| `String` | `String(length int) (string, error)` | สุ่ม string ด้วย crypto/rand |
| `HashMD5` | `HashMD5(input string) string` | คืนค่า MD5 hex digest |
| `HashSHA1` | `HashSHA1(input string) string` | คืนค่า SHA1 hex digest |
| `Decimal` | `Decimal(precision, maxAttempts int) (float64, error)` | สุ่มทศนิยม [0.01, 1.0) |

---

### `format` — จัดรูปแบบตัวเลข (2 functions)

```go
import "github.com/Maximumsoft-Co-LTD/go-util/format"
```

| Function | Signature | Description |
|---|---|---|
| `ToFixed` | `ToFixed(num float64, precision int) float64` | ปัดเศษ float64 ตามจำนวนทศนิยม |
| `Decimal2` | `Decimal2(x float64) float64` | ปัดเศษเป็น 2 ตำแหน่ง |

---

## โครงสร้างโปรเจกต์

```
go-util/
├── go.mod
├── go.sum
├── README.md
├── slice/
│   └── slice.go       # Generic slice helpers (53 functions)
├── mapx/
│   └── mapx.go        # Generic map helpers (15 functions)
├── val/
│   └── val.go         # Value & nil checking helpers (5 functions)
├── jwt/
│   └── jwt.go         # JWT Token service (golang-jwt/jwt/v5)
├── timex/
│   └── time.go        # Time helpers (12 functions)
├── random/
│   └── random.go      # Random & hashing (5 functions)
└── format/
    └── format.go      # Number formatting (2 functions)
```

---

## License

Internal use — Maximumsoft Co., LTD.
