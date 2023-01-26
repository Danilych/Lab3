package main

import (
	"fmt"
	"sync"
)

type token struct {
	data      string
	recipient int
	ttl       int
}

// Создание переменной для синхронизации, программа завершится когда все потоки завершат свою работу
var wg sync.WaitGroup

func send_message(current_thread_index int, recieved_channel chan token, to_send_channel chan token) {
	for { // бесконечный цикл для ожидания потоком сообщения по каналу
		message := <-recieved_channel
		fmt.Println("Поток ", current_thread_index, ": сообщение получено: ", message.data, ".")
		if message.recipient == current_thread_index {
			fmt.Println("Поток ", current_thread_index, ": сообщение достигло своего адреса!")
			wg.Done()
		} else {
			if message.ttl > 0 {
				message.ttl -= 1
				to_send_channel <- message //Отправка сообщения следующему каналу
			} else {
				fmt.Println("Поток ", current_thread_index, ": время жизни сообщения закончилось!")
				wg.Done()
			}
		}
	}
}

func main() {
	fmt.Println("Главный поток: введите число потоков")
	var threads_number int
	fmt.Scanln(&threads_number)

	// Создание каналов для общения потоков друг с другом
	var channels []chan token = make([]chan token, threads_number)
	for i := 0; i < threads_number; i++ {
		channels[i] = make(chan token)
	}
	wg.Add(1)
	// Создаем цепочку из потоков
	for i := 0; i < threads_number; i++ {
		if i == threads_number-1 {
			go send_message(i, channels[i], channels[i]) // Последний поток не отправляет сообщение следующему
		} else {
			go send_message(i, channels[i], channels[i+1])
		}
	}

	// Основная работа
	var user_token token

	fmt.Println("Главный поток: введите сообщение для отправки")
	fmt.Scanln(&user_token.data)
	fmt.Println("Главный поток: введите номер потока-получателя (от 0 до", threads_number-1, ")")
	fmt.Scanln(&user_token.recipient)
	if user_token.recipient > len(channels)-1 {
		user_token.recipient = len(channels) - 1
	}

	fmt.Println("Главный поток: введите длительность время жизни сообщения (количество впересылок)")
	fmt.Scanln(&user_token.ttl)

	//Отправка сообщения в начальный канал, чтоб начальный поток начал работу с сообщением
	channels[0] <- user_token

	// Ожидание конца работы всех потоков и их закрытие
	wg.Wait()
	for i := 0; i < threads_number; i++ {
		close(channels[i])
	}
}
