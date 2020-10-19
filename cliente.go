package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MatiasMarchant/Prueba1/tree/master/chat"
)

//ingresarordenespymes Es la simulación de ingreso de órdenes en base al archivo "pymes.csv", que ingresa ordenes cada tiempoespera segundos, usa la variable c para hacer llamadas remotas
func ingresarordenespymes(nombreexcel string, tiempoespera string, c chat.ChatServiceClient) bool {
	csvfile, err := os.Open(nombreexcel)
	tiempoesperaint, _ := strconv.Atoi(strings.TrimSuffix(tiempoespera, "\n")) // DANGER
	if err != nil {
		log.Fatalln("No pude abrir el csv:", err)
	}
	defer csvfile.Close()
	var prioritarioBool bool
	r := csv.NewReader(csvfile)
	r.Read() // Saltarse la gracias primera linea
	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return true
		}
		valorint32, _ := strconv.ParseInt(row[2], 10, 32)

		if row[5] == "1" {
			prioritarioBool = true
		} else {
			prioritarioBool = false
		}
		orden := chat.Ordenclientepymes{
			Id:          row[0],
			Producto:    row[1],
			Valor:       int32(valorint32),
			Tienda:      row[3],
			Destino:     row[4],
			Prioritario: prioritarioBool,
		}

		response, err := c.RecibirOrdenPymes(context.Background(), &orden)
		if err != nil {
			log.Fatalf("Error usando RecibirOrdenPymes: %s", err)
		}

		log.Printf("Codigo de seguimiento: %s", response.Nordenseguimiento)
		time.Sleep(time.Second * time.Duration(int64(tiempoesperaint)))
	}

}

//ingresarordenespymes Es la simulación de ingreso de órdenes en base al archivo "retail.csv", que ingresa ordenes cada tiempoespera segundos
func ingresarordenesretail(nombreexcel string, tiempoespera string, c chat.ChatServiceClient) bool {
	csvfile, err := os.Open(nombreexcel)
	tiempoesperaint, _ := strconv.Atoi(strings.TrimSuffix(tiempoespera, "\n")) // CUIDADO \r POR WINDOWS
	if err != nil {
		log.Fatalln("No pude abrir el csv:", err)
	}
	defer csvfile.Close()

	r := csv.NewReader(csvfile)
	r.Read() // Saltarse la gracias primera linea
	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return true
		}
		valorint32, _ := strconv.ParseInt(row[2], 10, 32)

		orden := chat.Ordenclienteretail{
			Id:       row[0],
			Producto: row[1],
			Valor:    int32(valorint32),
			Tienda:   row[3],
			Destino:  row[4],
		}

		c.RecibirOrdenRetail(context.Background(), &orden)

		time.Sleep(time.Second * time.Duration(int64(tiempoesperaint)))
	}
}

//preguntasiniciales es la función que pregunta el tipo de tienda y tiempo de espera entre el envío de órdenes
func preguntasiniciales() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Preguntas iniciales:")
	fmt.Println("Pymes o Retail?")
	fmt.Printf("> ")
	tipotienda, _ := reader.ReadBytes('\n')
	fmt.Println("Tiempo en segundos de espera entre el envío de órdenes")
	fmt.Printf("> ")
	tiempoespera, _ := reader.ReadBytes('\n')
	return string(tipotienda), string(tiempoespera)
}

//main se conecta con logistica, llama a funciones definidas anteriormente y simula ingreso de códigos de seguimientos en el caso pymes
func main() {
	fmt.Println("Corriendo el sistema de cliente...\n")
	var tipotienda, tiempoespera string
	tipotienda, tiempoespera = preguntasiniciales()

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist37:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	if string(tipotienda) == "pymes\n" {
		go ingresarordenespymes("pymes.csv", tiempoespera, c)
	} else if string(tipotienda) == "retail\n" {
		ingresarordenesretail("retail.csv", tiempoespera, c)
	}

	time.Sleep(10 * time.Second)
	if string(tipotienda) == "pymes\n" {
		for true {
			numero := rand.Intn(10) // Arbitrario
			codigoseguimiento := strconv.Itoa(numero)
			orden := chat.Ordenseguimiento{
				Nordenseguimiento: codigoseguimiento,
			}
			respuesta, err := c.CodigoSeguimiento(context.Background(), &orden)
			if err != nil {
				log.Fatalf("Error usando c.CodigoSeguimiento")
			}
			fmt.Println("El estado de su pedido es:", respuesta.Estado)
			time.Sleep(8 * time.Second)
		}
	}
}
