package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fatih/color"
	greenapi "github.com/green-api/whatsapp-api-client-golang"
)

func main() {
	var url string = "https://api.green-api.com"
	var idInstance string
	var apiToken string
	var iterations int
	var scriptIterations int

	// color.Green("Введите ApiUrl:\n")
	// fmt.Scanln(&url)

	color.Green("Введите idInstance:\n")
	fmt.Scanln(&idInstance)

	color.Green("Введите ApiTokenInstance:\n")
	fmt.Scanln(&apiToken)

	color.Green("Сколько запросов отправить за одну итерацию:\n")
	fmt.Scanln(&iterations)

	color.Green("Сколько раз запустить скрипт:\n")
	fmt.Scanln(&scriptIterations)

	var wg sync.WaitGroup

	GreenAPI := greenapi.GreenAPI{
		Host:             url,
		IDInstance:       idInstance,
		APITokenInstance: apiToken,
	}

	//lagged := 0

	//fmt.Println((api.ApiUrl))
	color.Cyan(GreenAPI.Host)

	for x := 0; x < scriptIterations; x++ {
		allTime := 0
		lagged := 0
		color.Cyan("Итерация номер %v", x+1)
		for i := 0; i < iterations; i++ {
			startTime := time.Now()
			response, err := GreenAPI.Sending().SendMessage("79326980324@c.us", "Privet")
			if err != nil {
				log.Fatal(err)
			}
			duration := time.Since(startTime)
			allTime = allTime + int(duration.Milliseconds())

			fmt.Printf("Request %d took %d ms\n", i+1, duration.Milliseconds())
			fmt.Println(response)
			//fmt.Printf("%s - %v\n", response, time.Now())

			if int(duration.Milliseconds()) > 500 {
				lagged += 1
				color.Red("Запрос №%v , Долгое время ответа: %vms , %v", i+1, duration.Milliseconds(), response)
				allTime -= int(duration.Milliseconds())
			}

		}
		wg.Wait()
		color.Yellow("Отправлено %v запросов", iterations)
		color.Green("Среднее время ответа: %v ms", (allTime / iterations))
		color.Red("Кол-во подвисших запросов (>500ms): %v", lagged)
		time.Sleep(time.Second * 1)
	}

	// color.Cyan("Очищаю очередь сообщений...")
	// _, err := GreenAPI.Methods().Queues().ClearMessagesQueue()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// color.Green("Cleared message queue successfully")
}

//это не основная версия, не заливать никуда
