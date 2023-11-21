package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Guilherme-daCosta/go-mensageria-products/internal/infra/akafka"
	"github.com/Guilherme-daCosta/go-mensageria-products/internal/infra/repository"
	"github.com/Guilherme-daCosta/go-mensageria-products/internal/infra/web"
	"github.com/Guilherme-daCosta/go-mensageria-products/internal/usecase"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/products")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("aqui")
	repository := repository.NewProductRepositoryMysql(db)
	fmt.Println("aqui")
	createProductUseCase := usecase.NewCreateProductUseCase(repository)
	fmt.Println("aqui")
	listProductsUseCase := usecase.NewListProductsUseCase(repository)
	fmt.Println("aqui")

	productHandlers := web.NewProductHandlers(createProductUseCase, listProductsUseCase)
	fmt.Println("aqui")

	r := chi.NewRouter()
	fmt.Println("aqui")
	r.Post("/products", productHandlers.CreateProductHandler)
	fmt.Println("aqui")
	r.Get("/products", productHandlers.ListProductsHandler)
	fmt.Println("aqui")

	go http.ListenAndServe(":8000", r)
	fmt.Println("aqui")

	msgChan := make(chan *kafka.Message)
	fmt.Println("aqui")
	go akafka.Consume([]string{"products"}, "host.docker.internal:9094", msgChan)
	fmt.Println("aqui11")

	for msg := range msgChan {
		dto := usecase.CreateProductInputDto{}
		fmt.Println("aquiiii")
		err := json.Unmarshal(msg.Value, &dto)
		if err != nil {
			fmt.Println(err)
		}
		_, err = createProductUseCase.Execute(dto)
		if err != nil {
			fmt.Println(err)
		}
	}
}
