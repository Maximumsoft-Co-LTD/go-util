# go-util

Shared Go utility functions for internal use at Maximumsoft.

```bash
go get github.com/Maximumsoft-Co-LTD/go-util
```

```go
import goutil "github.com/Maximumsoft-Co-LTD/go-util"
```

Flat package — import ครั้งเดียว ใช้ได้ทุก function:

```go
goutil.Filter(names, func(s string) bool { return len(s) > 3 })
goutil.Keys(m)
goutil.IsNil(v)
goutil.Bod(time.Now())
```

---

## Slice Helpers (53 functions)

#### Search & Index

| Function | Signature | Description |
|---|---|---|
| `Index` | `Index(limit int, predicate func(i int) bool) int` | คืน index แรกที่ predicate คืนค่า true, คืน -1 หากไม่พบ |
| `IndexOf` | `IndexOf[T comparable](collection []T, element T) int` | คืน index แรกของ element |
| `LastIndexOf` | `LastIndexOf[T comparable](collection []T, element T) int` | คืน index สุดท้ายของ element |
| `Contains` | `Contains[T comparable](collection []T, element T) bool` | ตรวจว่า element อยู่ใน slice หรือไม่ |
| `ContainsBy` | `ContainsBy[T any](collection []T, predicate func(T) bool) bool` | ตรวจว่ามี element ที่ตรงเงื่อนไขหรือไม่ |

#### Find

| Function | Signature | Description |
|---|---|---|
| `Find` | `Find[T any](collection []T, predicate func(T) bool) (T, bool)` | คืน element แรกที่ตรงเงื่อนไข |
| `FindIndexOf` | `FindIndexOf[T any](collection []T, predicate func(T) bool) (T, int, bool)` | คืน element, index, bool ตัวแรก |
| `FindLastIndexOf` | `FindLastIndexOf[T any](collection []T, predicate func(T) bool) (T, int, bool)` | เหมือน FindIndexOf แต่ค้นจากท้าย |
| `FindOrElse` | `FindOrElse[T any](collection []T, fallback T, predicate func(T) bool) T` | Find พร้อม fallback |
| `FindUniques` | `FindUniques[T comparable](collection []T) []T` | elements ที่ปรากฏครั้งเดียว |
| `FindUniquesBy` | `FindUniquesBy[T any, K comparable](collection []T, keyFn func(T) K) []T` | elements ที่ key ปรากฏครั้งเดียว |
| `FindDuplicates` | `FindDuplicates[T comparable](collection []T) []T` | elements ที่ซ้ำ (1 ชุดต่อค่า) |
| `FindDuplicatesBy` | `FindDuplicatesBy[T any, K comparable](collection []T, keyFn func(T) K) []T` | elements ที่ key ซ้ำ |

#### Transform

| Function | Signature | Description |
|---|---|---|
| `Filter` | `Filter[T any](collection []T, predicate func(T) bool) []T` | elements ที่ตรงเงื่อนไข |
| `Reject` | `Reject[T any](collection []T, predicate func(T) bool) []T` | elements ที่ไม่ตรงเงื่อนไข |
| `Map` | `Map[T any, U any](collection []T, iteratee func(T) U) []U` | แปลง element แต่ละตัว |
| `FlatMap` | `FlatMap[T any, U any](collection []T, iteratee func(T) []U) []U` | map แล้ว flatten |
| `Reduce` | `Reduce[T any, U any](collection []T, accumulator func(U, T) U, initial U) U` | รวมค่าเป็นค่าเดียว |
| `ForEach` | `ForEach[T any](collection []T, iteratee func(T, int))` | วนลูปทุก element |
| `Compact` | `Compact[T comparable](collection []T) []T` | ลบค่า zero-value |
| `Reverse` | `Reverse[T any](collection []T) []T` | กลับลำดับ |
| `Without` | `Without[T comparable](collection []T, exclude ...T) []T` | ลบ values ที่ระบุ |

#### Unique & Group

| Function | Signature | Description |
|---|---|---|
| `Unique` | `Unique[T comparable](collection []T) []T` | ลบค่าซ้ำ |
| `UniqueBy` | `UniqueBy[T any, K comparable](collection []T, keyFn func(T) K) []T` | ลบค่าซ้ำตาม key |
| `GroupBy` | `GroupBy[T any, K comparable](collection []T, keyFn func(T) K) map[K][]T` | จัดกลุ่มตาม key |
| `KeyBy` | `KeyBy[T any, K comparable](collection []T, keyFn func(T) K) map[K]T` | แปลง slice เป็น map |
| `Partition` | `Partition[T any](collection []T, predicate func(T) bool) ([]T, []T)` | แยก 2 กลุ่ม |
| `Chunk` | `Chunk[T any](collection []T, size int) [][]T` | แบ่ง chunk |
| `Flatten` | `Flatten[T any](collections [][]T) []T` | รวมหลาย slices |

#### Set Operations

| Function | Signature | Description |
|---|---|---|
| `Intersect` | `Intersect[T comparable](list1, list2 []T) []T` | elements ที่มีทั้งสอง slices |
| `Difference` | `Difference[T comparable](list1, list2 []T) []T` | elements ใน list1 ที่ไม่อยู่ใน list2 |
| `Union` | `Union[T comparable](lists ...[]T) []T` | รวม unique elements |

#### Aggregate

| Function | Signature | Description |
|---|---|---|
| `Count` | `Count[T comparable](collection []T, value T) int` | นับจำนวน |
| `CountBy` | `CountBy[T any](collection []T, predicate func(T) bool) int` | นับตามเงื่อนไข |
| `Min` | `Min[T cmp.Ordered](collection []T) T` | ค่าน้อยสุด |
| `MinIndex` | `MinIndex[T cmp.Ordered](collection []T) (T, int)` | ค่าน้อยสุดพร้อม index |
| `MinBy` | `MinBy[T any](collection []T, comparison func(a, b T) bool) T` | ค่าน้อยสุดด้วย comparison |
| `MinIndexBy` | `MinIndexBy[T any](collection []T, comparison func(a, b T) bool) (T, int)` | ค่าน้อยสุดพร้อม index ด้วย comparison |
| `Max` | `Max[T cmp.Ordered](collection []T) T` | ค่ามากสุด |
| `MaxIndex` | `MaxIndex[T cmp.Ordered](collection []T) (T, int)` | ค่ามากสุดพร้อม index |
| `MaxBy` | `MaxBy[T any](collection []T, comparison func(a, b T) bool) T` | ค่ามากสุดด้วย comparison |
| `MaxIndexBy` | `MaxIndexBy[T any](collection []T, comparison func(a, b T) bool) (T, int)` | ค่ามากสุดพร้อม index ด้วย comparison |
| `Latest` | `Latest(collection []time.Time) time.Time` | เวลาล่าสุด |
| `LatestBy` | `LatestBy[T any](collection []T, timeFn func(T) time.Time) T` | element ที่มีเวลาล่าสุด |

#### Boolean Checks

| Function | Signature | Description |
|---|---|---|
| `Every` | `Every[T any](collection []T, predicate func(T) bool) bool` | ทุกตัวตรงเงื่อนไข |
| `Some` | `Some[T any](collection []T, predicate func(T) bool) bool` | มีอย่างน้อยหนึ่งตัว |
| `None` | `None[T any](collection []T, predicate func(T) bool) bool` | ไม่มีตัวใดตรง |

#### Element Access

| Function | Signature | Description |
|---|---|---|
| `First` | `First[T any](collection []T) (T, bool)` | element แรก พร้อม bool |
| `FirstOrEmpty` | `FirstOrEmpty[T any](collection []T) T` | element แรก หรือ zero value |
| `Last` | `Last[T any](collection []T) (T, bool)` | element สุดท้าย พร้อม bool |
| `LastOrEmpty` | `LastOrEmpty[T any](collection []T) T` | element สุดท้าย หรือ zero value |

---

## Map Helpers (15 functions)

| Function | Signature | Description |
|---|---|---|
| `Keys` | `Keys[K comparable, V any](m map[K]V) []K` | keys ทั้งหมด |
| `Values` | `Values[K comparable, V any](m map[K]V) []V` | values ทั้งหมด |
| `ContainsKey` | `ContainsKey[K comparable, V any](m map[K]V, key K) bool` | ตรวจ key |
| `Pick` | `Pick[K comparable, V any](m map[K]V, keys []K) map[K]V` | เฉพาะ keys ที่ระบุ |
| `Omit` | `Omit[K comparable, V any](m map[K]V, keys []K) map[K]V` | ไม่รวม keys ที่ระบุ |
| `Merge` | `Merge[K comparable, V any](maps ...map[K]V) map[K]V` | รวมหลาย maps |
| `MapValues` | `MapValues[K comparable, V any, U any](m map[K]V, iteratee func(V) U) map[K]U` | แปลง values |
| `MapKeys` | `MapKeys[K comparable, V any, J comparable](m map[K]V, iteratee func(K) J) map[J]V` | แปลง keys |
| `FilterMap` | `FilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V` | กรอง map ตามเงื่อนไข |
| `ForEachMap` | `ForEachMap[K comparable, V any](m map[K]V, iteratee func(K, V))` | วนลูปทุก entry |
| `Invert` | `Invert[K comparable, V comparable](m map[K]V) map[V]K` | สลับ key/value |
| `Entries` | `Entries[K comparable, V any](m map[K]V) []Entry[K, V]` | map → key-value pairs |
| `FromEntries` | `FromEntries[K comparable, V any](entries []Entry[K, V]) map[K]V` | key-value pairs → map |
| `FindKey` | `FindKey[K comparable, V comparable](m map[K]V, value V) (K, bool)` | หา key จาก value |
| `FindKeyBy` | `FindKeyBy[K comparable, V any](m map[K]V, predicate func(K, V) bool) (K, bool)` | หา key ด้วย predicate |

---

## Value & Nil Helpers (5 functions)

| Function | Signature | Description |
|---|---|---|
| `IsNil` | `IsNil(v any) bool` | ตรวจ nil (รองรับ typed nil interface) |
| `IsNotNil` | `IsNotNil(v any) bool` | ตรวจไม่ใช่ nil |
| `Empty` | `Empty[T any]() T` | zero value ของ type T |
| `IsEmpty` | `IsEmpty[T comparable](v T) bool` | ตรวจ zero value |
| `IsNotEmpty` | `IsNotEmpty[T comparable](v T) bool` | ตรวจไม่ใช่ zero value |

---

## Time Helpers (12 functions)

| Function | Signature | Description |
|---|---|---|
| `StartOfDay` | `StartOfDay() time.Time` | ต้นวันปัจจุบัน (UTC) |
| `EndOfDay` | `EndOfDay() time.Time` | สิ้นวันปัจจุบัน (UTC) |
| `TruncateDate` | `TruncateDate(now time.Time, date int) time.Time` | 00:00 ของวันที่ระบุ |
| `TruncateWeek` | `TruncateWeek(now time.Time, day int) time.Time` | 00:00 ของวันในสัปดาห์ที่ผ่านมา |
| `FormatDateTimeString` | `FormatDateTimeString(t time.Time) string` | แปลง time → ISO string |
| `Bod` | `Bod(t time.Time) time.Time` | ต้นวัน UTC (17:00 UTC = 00:00 ไทย) |
| `BodLocal` | `BodLocal(t time.Time) time.Time` | ต้นวัน 00:00 local |
| `Eod` | `Eod(t time.Time) time.Time` | สิ้นวัน exclusive UTC `[Bod, Eod)` |
| `EodLocal` | `EodLocal(t time.Time) time.Time` | สิ้นวัน exclusive local `[BodLocal, EodLocal)` |
| `TimeFormat` | `TimeFormat(t time.Time, filter string) string` | เวลาเป็น string |
| `TimeNowOnly` | `TimeNowOnly(ref time.Time, t int64) int64` | ดึง h/m/s จาก Unix timestamp |
| `TimeInRange` | `TimeInRange(tTime, startTime, endTime time.Time) bool` | ตรวจช่วง `[start, end)` |

**ค่าคงที่:** `LayoutISO` = `"2006-01-02 15:04:05"`

---

## Random & Hashing (5 functions)

| Function | Signature | Description |
|---|---|---|
| `UniqueInts` | `UniqueInts(count, max int) []int` | สุ่มเลขไม่ซ้ำ [0, max) |
| `RandString` | `RandString(length int) (string, error)` | สุ่ม string (crypto/rand) |
| `HashMD5` | `HashMD5(input string) string` | MD5 hex digest |
| `HashSHA1` | `HashSHA1(input string) string` | SHA1 hex digest |
| `RandDecimal` | `RandDecimal(precision, maxAttempts int) (float64, error)` | สุ่มทศนิยม [0.01, 1.0) |

---

## Number Formatting (2 functions)

| Function | Signature | Description |
|---|---|---|
| `ToFixed` | `ToFixed(num float64, precision int) float64` | ปัดเศษตามจำนวนทศนิยม |
| `Decimal2` | `Decimal2(x float64) float64` | ปัดเศษ 2 ตำแหน่ง |

---

## JWT Service

| Function | Signature | Description |
|---|---|---|
| `NewJWTService` | `NewJWTService(secret, issuer string) (JWTService, error)` | สร้าง JWT Service |

**Interface:**

```go
type JWTService interface {
    GenerateToken(username string, sessionID string, ttl time.Duration) (string, error)
    ValidateToken(token string) (*jwt.Token, error)
}
```

**ตัวอย่าง:**

```go
svc, err := goutil.NewJWTService(os.Getenv("JWT_SECRET"), "my-app")
if err != nil {
    log.Fatal(err)
}
token, err := svc.GenerateToken("username", "session-123", 3*time.Hour)
parsed, err := svc.ValidateToken(token)
```

---

## ตัวอย่างการใช้งาน

```go
import goutil "github.com/Maximumsoft-Co-LTD/go-util"

// Slice
names := []string{"alice", "bob", "charlie", "bob"}
goutil.Contains(names, "bob")         // true
goutil.IndexOf(names, "bob")          // 1
goutil.LastIndexOf(names, "bob")      // 3
goutil.Unique(names)                   // ["alice", "bob", "charlie"]
goutil.Filter(names, func(s string) bool { return len(s) > 3 })
                                       // ["alice", "charlie"]

nums := []int{1, 2, 3, 4, 5}
goutil.Min(nums)                       // 1
goutil.Max(nums)                       // 5
goutil.Map(nums, func(n int) int { return n * 2 })
                                       // [2, 4, 6, 8, 10]

// FindDuplicates / FindUniques
data := []int{1, 2, 2, 3, 3, 4}
goutil.FindDuplicates(data)            // [2, 3]
goutil.FindUniques(data)               // [1, 4]

// Map
m := map[string]int{"a": 1, "b": 2, "c": 3}
goutil.Keys(m)                         // ["a", "b", "c"]
goutil.Pick(m, []string{"a", "c"})     // map[a:1 c:3]
goutil.FilterMap(m, func(k string, v int) bool { return v > 1 })
                                       // map[b:2 c:3]

// Value checking
goutil.IsEmpty("")                     // true
goutil.IsNotEmpty(42)                  // true
goutil.IsNil(nil)                      // true

// Time
bod := goutil.Bod(time.Now())
eod := goutil.Eod(time.Now())
goutil.TimeInRange(time.Now(), bod, eod)
```

---

## โครงสร้างโปรเจกต์

```
go-util/
├── go.mod
├── go.sum
├── README.md
├── slice.go       # Slice helpers (53 functions)
├── map.go         # Map helpers (15 functions)
├── val.go         # Value & nil checking (5 functions)
├── time.go        # Time helpers (12 functions)
├── random.go      # Random & hashing (5 functions)
├── format.go      # Number formatting (2 functions)
└── jwt.go         # JWT service
```

**Package:** `goutil` — flat package, ไม่มี sub-folder

---

## Dependencies

| Dependency | เหตุผล |
|---|---|
| `github.com/golang-jwt/jwt/v5` | JWT signing/validation |

functions อื่นทั้งหมด **ไม่มี external dependency** — ใช้เฉพาะ Go standard library

**รวม: 97 functions**

---

## License

Internal use — Maximumsoft Co., LTD.
