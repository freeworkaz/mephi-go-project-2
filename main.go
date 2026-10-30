package main

import (
	"fmt"
    "math/rand"
    "sync"
    "time"
)

const (
    SIZE   = 10_000_000
    CHUNKS = 8
)

// generateRandomElements generates random elements.
func generateRandomElements(size int) []int {
    if size <= 0 {
            return []int{}
    }

    data := make([]int, size)
    for i := range data {
            // rand.Intn возвращает число из диапазона [0, n),
            // прибавляем 1, чтобы все элементы были строго положительными.
           data[i] = rand.Intn(1<<31-1) + 1
    }

    return data

}

// maximum returns the maximum number of elements.
func maximum(data []int) int {
	if len(data) == 0 {
            return 0
    }

    maxValue := data[0]
    for _, value := range data[1:] {
        if value > maxValue {
            maxValue = value
        }
    }
    return maxValue
}

// maxChunks returns the maximum number of elements in a chunks.
func maxChunks(data []int) int {
	if len(data) == 0 {
        return 0
    }

    chunkSize := len(data) / CHUNKS
    chunkMaxima := make([]int, CHUNKS)

    var wg sync.WaitGroup
    for i := 0; i < CHUNKS; i++ {
		startIndex := i * chunkSize
		endIndex := startIndex + chunkSize
        // Последняя часть включает все оставшиеся элементы.
        if i == CHUNKS-1 {
			endIndex = len(data)
		}

        wg.Add(1)
		go func(chunk []int, index int) {
			defer wg.Done()
            chunkMaxima[index] = maximum(chunk)
        }(data[startIndex:endIndex], i)
    }

    wg.Wait()

    return maximum(chunkMaxima)
}

func main() {
    fmt.Printf("Генерируем %d целых чисел\n", SIZE)
    data := generateRandomElements(SIZE)

	fmt.Println("Ищем максимальное значение в один поток")
    startTime := time.Now()
    maxValue := maximum(data)
    elapsed := time.Since(startTime)
    fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d мкс\n", maxValue, elapsed.Microseconds())

    fmt.Printf("Ищем максимальное значение в %d потоков\n", CHUNKS)
    startTime = time.Now()
    maxValue = maxChunks(data)
    elapsed = time.Since(startTime)
    fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d мкс\n", maxValue, elapsed.Microseconds())
}
