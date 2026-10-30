package main

import (
        "reflect"
        "testing"
)

func TestGenerateRandomElements(t *testing.T) {
        tests := []struct {
                name string
                size int
        }{
                {"размер 0", 0},
                {"отрицательный размер", -5},
                {"один элемент", 1},
                {"сто элементов", 100},
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        data := generateRandomElements(tt.size)

                        expectedLen := tt.size
                        if expectedLen < 0 {
                                expectedLen = 0
                        }
                        if len(data) != expectedLen {
                                t.Fatalf("длина слайса = %d, хотим %d", len(data), expectedLen)
                        }

                        // Все элементы должны быть строго положительными.
                        for i, value := range data {
                                if value <= 0 {
                                        t.Errorf("элемент [%d] = %d не является положительным", i, value)
                                }
                        }
                })
        }
}

func TestMaximum(t *testing.T) {
        tests := []struct {
                name     string
                data     []int
                expected int
        }{
                {"nil-слайс", nil, 0},
                {"пустой слайс", []int{}, 0},
                {"один элемент", []int{42}, 42},
                {"несколько элементов", []int{3, 7, 2, 9, 5}, 9},
                {"максимум первый", []int{10, 1, 2, 3}, 10},
                {"максимум последний", []int{1, 2, 3, 10}, 10},
                {"все одинаковые", []int{5, 5, 5}, 5},
                {"отрицательные числа", []int{-3, -1, -7}, -1},
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        if got := maximum(tt.data); got != tt.expected {
                                t.Errorf("maximum(%v) = %d, хотим %d", tt.data, got, tt.expected)
                        }
                })
        }
}

func TestMaxChunks(t *testing.T) {
        tests := []struct {
                name     string
                data     []int
                expected int
        }{
                {"nil-слайс", nil, 0},
                {"пустой слайс", []int{}, 0},
                {"меньше восьми элементов", []int{1, 5, 3}, 5},
                {"один элемент", []int{7}, 7},
                {"ровно восемь элементов", []int{1, 2, 3, 4, 5, 6, 7, 8}, 8},
                {"не делится на восемь", []int{9, 1, 2, 3, 4, 5, 6, 7, 8, 10}, 10},
                {"максимум в последней части", makeMaxInLastChunk(), 1000},
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        if got := maxChunks(tt.data); got != tt.expected {
                                t.Errorf("maxChunks(%v) = %d, хотим %d", tt.data, got, tt.expected)
                        }
                })
        }
}

// makeMaxInLastChunk создаёт слайс из ста элементов, где максимум находится
// в последней восьмой части.
func makeMaxInLastChunk() []int {
        data := make([]int, 100)
        for i := range data {
                data[i] = i + 1
        }
        data[99] = 1000

        return data
}

// TestMaxChunksMatchesMaximum проверяет, что многопоточная и однопоточная версии
// дают одинаковый результат на больших данных.
func TestMaxChunksMatchesMaximum(t *testing.T) {
        const size = 1_000_000

        data := generateRandomElements(size)

        if len(data) != size {
                t.Fatalf("длина слайса = %d, хотим %d", len(data), size)
        }

        singleThreadMax := maximum(data)
        chunksMax := maxChunks(data)

        if singleThreadMax != chunksMax {
                t.Errorf("maximum = %d, maxChunks = %d — значения должны совпадать",
                        singleThreadMax, chunksMax)
        }
}

// Benchmark для функций поиска максимума.
var benchData = generateRandomElements(1_000_000)

func BenchmarkMaximum(b *testing.B) {
        for i := 0; i < b.N; i++ {
                maximum(benchData)
        }
}

func BenchmarkMaxChunks(b *testing.B) {
        for i := 0; i < b.N; i++ {
                maxChunks(benchData)
        }
}

// Проверка, что generateRandomElements(0) возвращает пустой, но не nil слайс.
func TestGenerateRandomElementsEmptyNotNil(t *testing.T) {
        data := generateRandomElements(0)
        if data == nil {
                t.Error("generateRandomElements(0) вернул nil, хотим пустой слайс")
        }
        if !reflect.DeepEqual(data, []int{}) {
                t.Errorf("generateRandomElements(0) = %v, хотим %v", data, []int{})
        }
}
