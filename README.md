# go-util

Shared Go utility functions for internal use at Maximumsoft.

```bash
go get github.com/Maximumsoft-Co-LTD/go-util
```

```go
import goutil "github.com/Maximumsoft-Co-LTD/go-util"
```

---

## Slice — Search & Index

| Function | Signature | Description |
|---|---|---|
| `Index` | `Index(limit int, predicate func(i int) bool) int` | คืน index แรกที่ predicate คืนค่า true, คืน -1 หากไม่พบ |
| `IndexOf` | `IndexOf[T comparable](collection []T, element T) int` | คืน index แรกของ element |
| `LastIndexOf` | `LastIndexOf[T comparable](collection []T, element T) int` | คืน index สุดท้ายของ element |
| `Contains` | `Contains[T comparable](collection []T, element T) bool` | ตรวจว่า element อยู่ใน slice หรือไม่ |
| `ContainsBy` | `ContainsBy[T any](collection []T, predicate func(T) bool) bool` | ตรวจว่ามี element ที่ตรงเงื่อนไขหรือไม่ |

<details>
<summary>ตัวอย่าง Search & Index</summary>

#### `Index`
```go
names := []string{"alice", "bob", "charlie"}
idx := goutil.Index(len(names), func(i int) bool {
    return names[i] == "bob"
})
// idx = 1
```

#### `IndexOf`
```go
goutil.IndexOf([]string{"a", "b", "c"}, "b")  // 1
goutil.IndexOf([]string{"a", "b", "c"}, "z")  // -1
```

#### `LastIndexOf`
```go
goutil.LastIndexOf([]string{"a", "b", "a", "c"}, "a")  // 2
```

#### `Contains`
```go
goutil.Contains([]int{1, 2, 3}, 2)  // true
goutil.Contains([]int{1, 2, 3}, 9)  // false
```

#### `ContainsBy`
```go
type User struct { Name string; Age int }
users := []User{{"Alice", 30}, {"Bob", 17}}
goutil.ContainsBy(users, func(u User) bool { return u.Age >= 18 })
// true
```

</details>

---

## Slice — Find

| Function | Signature | Description |
|---|---|---|
| `Find` | `Find[T any](collection []T, predicate func(T) bool) (T, bool)` | คืน element แรกที่ตรงเงื่อนไข |
| `FindIndexOf` | `FindIndexOf[T any](collection []T, predicate func(T) bool) (T, int, bool)` | คืน element, index, bool ตัวแรก |
| `FindLastIndexOf` | `FindLastIndexOf[T any](collection []T, predicate func(T) bool) (T, int, bool)` | เหมือน FindIndexOf แต่ค้นจากท้าย |
| `FindOrElse` | `FindOrElse[T any](collection []T, fallback T, predicate func(T) bool) T` | Find พร้อม fallback ถ้าไม่เจอ |
| `FindUniques` | `FindUniques[T comparable](collection []T) []T` | elements ที่ปรากฏครั้งเดียว |
| `FindUniquesBy` | `FindUniquesBy[T any, K comparable](collection []T, keyFn func(T) K) []T` | elements ที่ key ปรากฏครั้งเดียว |
| `FindDuplicates` | `FindDuplicates[T comparable](collection []T) []T` | elements ที่ซ้ำ (1 ชุดต่อค่า) |
| `FindDuplicatesBy` | `FindDuplicatesBy[T any, K comparable](collection []T, keyFn func(T) K) []T` | elements ที่ key ซ้ำ |

<details>
<summary>ตัวอย่าง Find</summary>

#### `Find`
```go
val, ok := goutil.Find([]int{1, 2, 3, 4}, func(n int) bool { return n > 2 })
// val = 3, ok = true
```

#### `FindIndexOf`
```go
val, idx, ok := goutil.FindIndexOf([]string{"a", "bb", "ccc"}, func(s string) bool {
    return len(s) > 1
})
// val = "bb", idx = 1, ok = true
```

#### `FindLastIndexOf`
```go
val, idx, ok := goutil.FindLastIndexOf([]int{1, 2, 3, 4}, func(n int) bool {
    return n%2 == 0
})
// val = 4, idx = 3, ok = true
```

#### `FindOrElse`
```go
result := goutil.FindOrElse([]int{1, 2, 3}, -1, func(n int) bool { return n > 5 })
// result = -1 (ไม่เจอ ใช้ fallback)
```

#### `FindUniques`
```go
goutil.FindUniques([]int{1, 2, 2, 3, 3, 4})
// [1, 4] — เฉพาะตัวที่ปรากฏครั้งเดียว
```

#### `FindUniquesBy`
```go
type Item struct { Category string; Name string }
items := []Item{{"A", "x"}, {"B", "y"}, {"A", "z"}}
goutil.FindUniquesBy(items, func(i Item) string { return i.Category })
// [{B y}]
```

#### `FindDuplicates`
```go
goutil.FindDuplicates([]int{1, 2, 2, 3, 3, 4})
// [2, 3] — ค่าที่ซ้ำ (คืน 1 ชุดต่อค่า)
```

#### `FindDuplicatesBy`
```go
type Item struct { Category string; Name string }
items := []Item{{"A", "x"}, {"B", "y"}, {"A", "z"}}
goutil.FindDuplicatesBy(items, func(i Item) string { return i.Category })
// [{A x}]
```

</details>

---

## Slice — Transform

| Function | Signature | Description |
|---|---|---|
| `Filter` | `Filter[T any](collection []T, predicate func(T) bool) []T` | elements ที่ตรงเงื่อนไข |
| `Reject` | `Reject[T any](collection []T, predicate func(T) bool) []T` | elements ที่ไม่ตรงเงื่อนไข (ตรงข้าม Filter) |
| `Map` | `Map[T any, U any](collection []T, iteratee func(T) U) []U` | แปลง element แต่ละตัว |
| `FlatMap` | `FlatMap[T any, U any](collection []T, iteratee func(T) []U) []U` | map แล้ว flatten |
| `Reduce` | `Reduce[T any, U any](collection []T, accumulator func(U, T) U, initial U) U` | รวมค่าเป็นค่าเดียว |
| `ForEach` | `ForEach[T any](collection []T, iteratee func(T, int))` | วนลูปทุก element |
| `Compact` | `Compact[T comparable](collection []T) []T` | ลบค่า zero-value |
| `Reverse` | `Reverse[T any](collection []T) []T` | กลับลำดับ |
| `Without` | `Without[T comparable](collection []T, exclude ...T) []T` | ลบ values ที่ระบุ |

<details>
<summary>ตัวอย่าง Transform</summary>

#### `Filter`
```go
goutil.Filter([]int{1, 2, 3, 4, 5}, func(n int) bool { return n > 3 })
// [4, 5]
```

#### `Reject`
```go
goutil.Reject([]int{1, 2, 3, 4, 5}, func(n int) bool { return n > 3 })
// [1, 2, 3] — ตรงข้าม Filter
```

#### `Map`
```go
goutil.Map([]int{1, 2, 3}, func(n int) string { return fmt.Sprintf("#%d", n) })
// ["#1", "#2", "#3"]
```

#### `FlatMap`
```go
goutil.FlatMap([]int{1, 2, 3}, func(n int) []int { return []int{n, n * 10} })
// [1, 10, 2, 20, 3, 30]
```

#### `Reduce`
```go
sum := goutil.Reduce([]int{1, 2, 3, 4, 5}, func(acc int, n int) int {
    return acc + n
}, 0)
// sum = 15
```

#### `ForEach`
```go
goutil.ForEach([]string{"a", "b", "c"}, func(val string, idx int) {
    fmt.Printf("%d: %s\n", idx, val)
})
// 0: a
// 1: b
// 2: c
```

#### `Compact`
```go
goutil.Compact([]string{"a", "", "b", "", "c"})
// ["a", "b", "c"]
```

#### `Reverse`
```go
goutil.Reverse([]int{1, 2, 3})
// [3, 2, 1]
```

#### `Without`
```go
goutil.Without([]int{1, 2, 3, 4, 5}, 2, 4)
// [1, 3, 5]
```

</details>

---

## Slice — Unique & Group

| Function | Signature | Description |
|---|---|---|
| `Unique` | `Unique[T comparable](collection []T) []T` | ลบค่าซ้ำ |
| `UniqueBy` | `UniqueBy[T any, K comparable](collection []T, keyFn func(T) K) []T` | ลบค่าซ้ำตาม key |
| `GroupBy` | `GroupBy[T any, K comparable](collection []T, keyFn func(T) K) map[K][]T` | จัดกลุ่มตาม key |
| `KeyBy` | `KeyBy[T any, K comparable](collection []T, keyFn func(T) K) map[K]T` | แปลง slice เป็น map |
| `Partition` | `Partition[T any](collection []T, predicate func(T) bool) ([]T, []T)` | แยก 2 กลุ่ม |
| `Chunk` | `Chunk[T any](collection []T, size int) [][]T` | แบ่ง chunk |
| `Flatten` | `Flatten[T any](collections [][]T) []T` | รวมหลาย slices |

<details>
<summary>ตัวอย่าง Unique & Group</summary>

#### `Unique`
```go
goutil.Unique([]string{"a", "b", "a", "c", "b"})
// ["a", "b", "c"]
```

#### `UniqueBy`
```go
type User struct { ID int; Name string }
users := []User{{1, "Alice"}, {2, "Bob"}, {1, "Alice2"}}
goutil.UniqueBy(users, func(u User) int { return u.ID })
// [{1 Alice} {2 Bob}]
```

#### `GroupBy`
```go
type User struct { Name string; Age int }
users := []User{{"A", 20}, {"B", 30}, {"C", 20}}
goutil.GroupBy(users, func(u User) int { return u.Age })
// map[20:[{A 20} {C 20}] 30:[{B 30}]]
```

#### `KeyBy`
```go
type User struct { ID int; Name string }
users := []User{{1, "Alice"}, {2, "Bob"}}
goutil.KeyBy(users, func(u User) int { return u.ID })
// map[1:{1 Alice} 2:{2 Bob}]
```

#### `Partition`
```go
pass, fail := goutil.Partition([]int{1, 2, 3, 4, 5}, func(n int) bool {
    return n > 3
})
// pass = [4, 5], fail = [1, 2, 3]
```

#### `Chunk`
```go
goutil.Chunk([]int{1, 2, 3, 4, 5}, 2)
// [[1, 2], [3, 4], [5]]
```

#### `Flatten`
```go
goutil.Flatten([][]int{{1, 2}, {3, 4}, {5}})
// [1, 2, 3, 4, 5]
```

</details>

---

## Slice — Set Operations

| Function | Signature | Description |
|---|---|---|
| `Intersect` | `Intersect[T comparable](list1, list2 []T) []T` | elements ที่มีทั้งสอง slices |
| `Difference` | `Difference[T comparable](list1, list2 []T) []T` | elements ใน list1 ที่ไม่อยู่ใน list2 |
| `Union` | `Union[T comparable](lists ...[]T) []T` | รวม unique elements |

<details>
<summary>ตัวอย่าง Set Operations</summary>

#### `Intersect`
```go
goutil.Intersect([]int{1, 2, 3}, []int{2, 3, 4})
// [2, 3]
```

#### `Difference`
```go
goutil.Difference([]int{1, 2, 3}, []int{2, 3, 4})
// [1]
```

#### `Union`
```go
goutil.Union([]int{1, 2, 3}, []int{2, 3, 4})
// [1, 2, 3, 4]
```

</details>

---

## Slice — Aggregate

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

<details>
<summary>ตัวอย่าง Aggregate</summary>

#### `Count`
```go
goutil.Count([]string{"a", "b", "a", "c", "a"}, "a")
// 3
```

#### `CountBy`
```go
goutil.CountBy([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
// 2
```

#### `Min`
```go
goutil.Min([]int{3, 1, 4, 1, 5})  // 1
```

#### `MinIndex`
```go
val, idx := goutil.MinIndex([]int{3, 1, 4, 1, 5})
// val = 1, idx = 1
```

#### `MinBy`
```go
type Product struct { Name string; Price float64 }
products := []Product{{"A", 100}, {"B", 50}, {"C", 200}}
goutil.MinBy(products, func(a, b Product) bool { return a.Price < b.Price })
// {B 50}
```

#### `MinIndexBy`
```go
val, idx := goutil.MinIndexBy(products, func(a, b Product) bool {
    return a.Price < b.Price
})
// val = {B 50}, idx = 1
```

#### `Max`
```go
goutil.Max([]int{3, 1, 4, 1, 5})  // 5
```

#### `MaxIndex`
```go
val, idx := goutil.MaxIndex([]int{3, 1, 4, 1, 5})
// val = 5, idx = 4
```

#### `MaxBy`
```go
goutil.MaxBy(products, func(a, b Product) bool { return a.Price > b.Price })
// {C 200}
```

#### `MaxIndexBy`
```go
val, idx := goutil.MaxIndexBy(products, func(a, b Product) bool {
    return a.Price > b.Price
})
// val = {C 200}, idx = 2
```

#### `Latest`
```go
t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
t2 := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
t3 := time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)
goutil.Latest([]time.Time{t1, t2, t3})
// 2024-06-15
```

#### `LatestBy`
```go
type Order struct { ID string; CreatedAt time.Time }
orders := []Order{
    {"A", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
    {"B", time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)},
}
goutil.LatestBy(orders, func(o Order) time.Time { return o.CreatedAt })
// {B 2024-06-15}
```

</details>

---

## Slice — Boolean Checks

| Function | Signature | Description |
|---|---|---|
| `Every` | `Every[T any](collection []T, predicate func(T) bool) bool` | ทุกตัวตรงเงื่อนไข |
| `Some` | `Some[T any](collection []T, predicate func(T) bool) bool` | มีอย่างน้อยหนึ่งตัว |
| `None` | `None[T any](collection []T, predicate func(T) bool) bool` | ไม่มีตัวใดตรง |

<details>
<summary>ตัวอย่าง Boolean Checks</summary>

#### `Every`
```go
goutil.Every([]int{2, 4, 6}, func(n int) bool { return n%2 == 0 })
// true — ทุกตัวเป็นเลขคู่
```

#### `Some`
```go
goutil.Some([]int{1, 3, 4}, func(n int) bool { return n%2 == 0 })
// true — มีอย่างน้อย 1 ตัวเป็นเลขคู่
```

#### `None`
```go
goutil.None([]int{1, 3, 5}, func(n int) bool { return n%2 == 0 })
// true — ไม่มีเลขคู่เลย
```

</details>

---

## Slice — Element Access

| Function | Signature | Description |
|---|---|---|
| `First` | `First[T any](collection []T) (T, bool)` | element แรก พร้อม bool |
| `FirstOrEmpty` | `FirstOrEmpty[T any](collection []T) T` | element แรก หรือ zero value |
| `Last` | `Last[T any](collection []T) (T, bool)` | element สุดท้าย พร้อม bool |
| `LastOrEmpty` | `LastOrEmpty[T any](collection []T) T` | element สุดท้าย หรือ zero value |

<details>
<summary>ตัวอย่าง Element Access</summary>

#### `First`
```go
val, ok := goutil.First([]int{10, 20, 30})
// val = 10, ok = true

val, ok = goutil.First([]int{})
// val = 0, ok = false
```

#### `FirstOrEmpty`
```go
goutil.FirstOrEmpty([]string{"hello", "world"})  // "hello"
goutil.FirstOrEmpty([]string{})                    // ""
```

#### `Last`
```go
val, ok := goutil.Last([]int{10, 20, 30})
// val = 30, ok = true
```

#### `LastOrEmpty`
```go
goutil.LastOrEmpty([]int{10, 20, 30})  // 30
goutil.LastOrEmpty([]int{})             // 0
```

</details>

---

## Map Helpers

| Function | Signature | Description |
|---|---|---|
| `Keys` | `Keys[K comparable, V any](m map[K]V) []K` | keys ทั้งหมด |
| `Values` | `Values[K comparable, V any](m map[K]V) []V` | values ทั้งหมด |
| `ContainsKey` | `ContainsKey[K comparable, V any](m map[K]V, key K) bool` | ตรวจ key |
| `Pick` | `Pick[K comparable, V any](m map[K]V, keys []K) map[K]V` | เฉพาะ keys ที่ระบุ |
| `Omit` | `Omit[K comparable, V any](m map[K]V, keys []K) map[K]V` | ไม่รวม keys ที่ระบุ |
| `Merge` | `Merge[K comparable, V any](maps ...map[K]V) map[K]V` | รวมหลาย maps (ตัวหลังทับ) |
| `MapValues` | `MapValues[K comparable, V any, U any](m map[K]V, iteratee func(V) U) map[K]U` | แปลง values |
| `MapKeys` | `MapKeys[K comparable, V any, J comparable](m map[K]V, iteratee func(K) J) map[J]V` | แปลง keys |
| `FilterMap` | `FilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V` | กรอง map ตามเงื่อนไข |
| `ForEachMap` | `ForEachMap[K comparable, V any](m map[K]V, iteratee func(K, V))` | วนลูปทุก entry |
| `Invert` | `Invert[K comparable, V comparable](m map[K]V) map[V]K` | สลับ key/value |
| `Entries` | `Entries[K comparable, V any](m map[K]V) []Entry[K, V]` | map → key-value pairs |
| `FromEntries` | `FromEntries[K comparable, V any](entries []Entry[K, V]) map[K]V` | key-value pairs → map |
| `FindKey` | `FindKey[K comparable, V comparable](m map[K]V, value V) (K, bool)` | หา key จาก value |
| `FindKeyBy` | `FindKeyBy[K comparable, V any](m map[K]V, predicate func(K, V) bool) (K, bool)` | หา key ด้วย predicate |

<details>
<summary>ตัวอย่าง Map Helpers</summary>

#### `Keys`
```go
goutil.Keys(map[string]int{"a": 1, "b": 2, "c": 3})
// ["a", "b", "c"]
```

#### `Values`
```go
goutil.Values(map[string]int{"a": 1, "b": 2, "c": 3})
// [1, 2, 3]
```

#### `ContainsKey`
```go
goutil.ContainsKey(map[string]int{"a": 1, "b": 2}, "a")  // true
goutil.ContainsKey(map[string]int{"a": 1, "b": 2}, "z")  // false
```

#### `Pick`
```go
m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
goutil.Pick(m, []string{"a", "c"})
// map[a:1 c:3]
```

#### `Omit`
```go
m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
goutil.Omit(m, []string{"b", "d"})
// map[a:1 c:3]
```

#### `Merge`
```go
m1 := map[string]int{"a": 1, "b": 2}
m2 := map[string]int{"b": 99, "c": 3}
goutil.Merge(m1, m2)
// map[a:1 b:99 c:3]
```

#### `MapValues`
```go
m := map[string]int{"a": 1, "b": 2, "c": 3}
goutil.MapValues(m, func(v int) int { return v * 10 })
// map[a:10 b:20 c:30]
```

#### `MapKeys`
```go
m := map[string]int{"a": 1, "b": 2}
goutil.MapKeys(m, func(k string) string { return strings.ToUpper(k) })
// map[A:1 B:2]
```

#### `FilterMap`
```go
m := map[string]int{"a": 1, "b": 2, "c": 3}
goutil.FilterMap(m, func(k string, v int) bool { return v > 1 })
// map[b:2 c:3]
```

#### `ForEachMap`
```go
goutil.ForEachMap(map[string]int{"a": 1, "b": 2}, func(k string, v int) {
    fmt.Printf("%s=%d\n", k, v)
})
```

#### `Invert`
```go
goutil.Invert(map[string]int{"a": 1, "b": 2, "c": 3})
// map[1:a 2:b 3:c]
```

#### `Entries`
```go
goutil.Entries(map[string]int{"a": 1, "b": 2})
// [{Key:"a" Value:1}, {Key:"b" Value:2}]
```

#### `FromEntries`
```go
entries := []goutil.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}}
goutil.FromEntries(entries)
// map[a:1 b:2]
```

#### `FindKey`
```go
key, ok := goutil.FindKey(map[string]int{"a": 1, "b": 2, "c": 3}, 2)
// key = "b", ok = true
```

#### `FindKeyBy`
```go
key, ok := goutil.FindKeyBy(map[string]int{"a": 1, "b": 20, "c": 3}, func(k string, v int) bool {
    return v > 10
})
// key = "b", ok = true
```

</details>

---

## Value & Nil Helpers

| Function | Signature | Description |
|---|---|---|
| `IsNil` | `IsNil(v any) bool` | ตรวจ nil (รองรับ typed nil interface) |
| `IsNotNil` | `IsNotNil(v any) bool` | ตรวจไม่ใช่ nil |
| `Empty` | `Empty[T any]() T` | zero value ของ type T |
| `IsEmpty` | `IsEmpty[T comparable](v T) bool` | ตรวจ zero value |
| `IsNotEmpty` | `IsNotEmpty[T comparable](v T) bool` | ตรวจไม่ใช่ zero value |

<details>
<summary>ตัวอย่าง Value & Nil</summary>

#### `IsNil`
```go
goutil.IsNil(nil)  // true

var p *int = nil
var i interface{} = p
i == nil            // false (Go gotcha!)
goutil.IsNil(i)     // true  (ถูกต้อง)
```

#### `IsNotNil`
```go
goutil.IsNotNil("hello")  // true
goutil.IsNotNil(nil)       // false
```

#### `Empty`
```go
goutil.Empty[string]()  // ""
goutil.Empty[int]()     // 0
goutil.Empty[bool]()    // false
```

#### `IsEmpty`
```go
goutil.IsEmpty("")       // true
goutil.IsEmpty(0)        // true
goutil.IsEmpty("hello")  // false
```

#### `IsNotEmpty`
```go
goutil.IsNotEmpty("hello")  // true
goutil.IsNotEmpty(0)        // false
```

</details>

---

## Time Helpers

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

<details>
<summary>ตัวอย่าง Time Helpers</summary>

#### `StartOfDay`
```go
goutil.StartOfDay()
// 2024-03-17 00:00:00 +0000 UTC
```

#### `EndOfDay`
```go
goutil.EndOfDay()
// 2024-03-17 23:59:59.999999999 +0000 UTC
```

#### `TruncateDate`
```go
now := time.Date(2024, 3, 17, 14, 0, 0, 0, time.Local)
goutil.TruncateDate(now, 10)
// 2024-03-10 00:00:00

goutil.TruncateDate(now, 25)
// 2024-02-25 00:00:00 — ยังไม่ถึงวันที่ 25 → ย้อนเดือนก่อน
```

#### `TruncateWeek`
```go
now := time.Date(2024, 3, 17, 14, 0, 0, 0, time.Local) // อาทิตย์
goutil.TruncateWeek(now, 1) // 1 = Monday
// 2024-03-11 00:00:00
```

#### `FormatDateTimeString`
```go
t := time.Date(2024, 3, 17, 14, 30, 45, 0, time.UTC)
goutil.FormatDateTimeString(t)
// "2024-03-17 14:30:45"
```

#### `Bod`
```go
goutil.Bod(time.Now())
// ต้นวัน UTC (17:00 UTC = 00:00 เวลาไทย)
```

#### `BodLocal`
```go
goutil.BodLocal(time.Now())
// 2024-03-17 00:00:00 +0700 ICT
```

#### `Eod`
```go
goutil.Eod(time.Now())
// สิ้นวัน exclusive UTC — ใช้กับ [Bod, Eod)
```

#### `EodLocal`
```go
goutil.EodLocal(time.Now())
// 2024-03-18 00:00:00 +0700 ICT — ใช้กับ [BodLocal, EodLocal)
```

#### `TimeFormat`
```go
t := time.Date(2024, 3, 17, 14, 30, 45, 0, time.UTC)
goutil.TimeFormat(t, "day")     // "2024March17"
goutil.TimeFormat(t, "hour")    // "2024March1714"
goutil.TimeFormat(t, "minute")  // "2024March171430"
goutil.TimeFormat(t, "second")  // "2024March17143045"
```

#### `TimeNowOnly`
```go
ref := time.Date(2024, 3, 17, 0, 0, 0, 0, time.Local)
ts := time.Date(2000, 1, 1, 14, 30, 0, 0, time.Local).Unix()
goutil.TimeNowOnly(ref, ts)
// Unix timestamp ของ 2024-03-17 14:30:00
```

#### `TimeInRange`
```go
bod := goutil.Bod(time.Now())
eod := goutil.Eod(time.Now())
goutil.TimeInRange(time.Now(), bod, eod)
// true
```

</details>

---

## Random & Hashing

| Function | Signature | Description |
|---|---|---|
| `UniqueInts` | `UniqueInts(count, max int) []int` | สุ่มเลขไม่ซ้ำ [0, max) |
| `RandString` | `RandString(length int) (string, error)` | สุ่ม string (crypto/rand) |
| `HashMD5` | `HashMD5(input string) string` | MD5 hex digest |
| `HashSHA1` | `HashSHA1(input string) string` | SHA1 hex digest |
| `RandDecimal` | `RandDecimal(precision, maxAttempts int) (float64, error)` | สุ่มทศนิยม [0.01, 1.0) |

<details>
<summary>ตัวอย่าง Random & Hashing</summary>

#### `UniqueInts`
```go
goutil.UniqueInts(3, 10)
// [7, 2, 5] — สุ่ม 3 เลขไม่ซ้ำจาก [0, 10)
```

#### `RandString`
```go
s, err := goutil.RandString(16)
// s = "aB3kL9mNpQrStUvW"
```

#### `HashMD5`
```go
goutil.HashMD5("hello")
// "5d41402abc4b2a76b9719d911017c592"
```

#### `HashSHA1`
```go
goutil.HashSHA1("hello")
// "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"
```

#### `RandDecimal`
```go
val, err := goutil.RandDecimal(2, 1000)
// val = 0.47
```

</details>

---

## Number Formatting

| Function | Signature | Description |
|---|---|---|
| `ToFixed` | `ToFixed(num float64, precision int) float64` | ปัดเศษตามจำนวนทศนิยม |
| `Decimal2` | `Decimal2(x float64) float64` | ปัดเศษ 2 ตำแหน่ง |

<details>
<summary>ตัวอย่าง Number Formatting</summary>

#### `ToFixed`
```go
goutil.ToFixed(3.14159, 2)  // 3.14
goutil.ToFixed(3.14159, 4)  // 3.1416
goutil.ToFixed(3.1, 0)      // 3
```

#### `Decimal2`
```go
goutil.Decimal2(3.14159)  // 3.14
goutil.Decimal2(2.999)    // 3
```

</details>

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

<details>
<summary>ตัวอย่าง JWT Service</summary>

#### `NewJWTService`
```go
svc, err := goutil.NewJWTService(os.Getenv("JWT_SECRET"), "my-app")
if err != nil {
    log.Fatal(err)
}
```

#### `GenerateToken`
```go
token, err := svc.GenerateToken("admin", "session-abc", 3*time.Hour)
// token = "eyJhbGciOiJIUzI1NiIs..."
```

#### `ValidateToken`
```go
parsed, err := svc.ValidateToken(token)
if err != nil {
    log.Fatal("invalid token")
}
claims := parsed.Claims.(jwt.MapClaims)
fmt.Println(claims["username"]) // "admin"
```

</details>

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

## License

Internal use — Maximumsoft Co., LTD.
