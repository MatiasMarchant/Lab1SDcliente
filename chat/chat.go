package chat

import (
	//"fmt"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"time"
)

//Registro es
type Registro struct {
	timestamp   time.Time
	idpaquete   string
	tipo        string
	nombre      string
	valor       string
	origen      string
	destino     string
	seguimiento string
}

//Cola es
type Cola struct {
	idpaquete   string
	seguimiento string
	tipo        string
	valor       string
	intentos    string
	estado      string
	origen      string
	destino     string
}

//PaqueteEnMarcha es
type PaqueteEnMarcha struct {
	idpaquete     string
	estado        string
	idcamion      string
	idseguimiento string
	intentos      string
	origen        string
	destino       string
	timestamp     time.Time
}

//Server es
type Server struct {
	ListaRegistro    []Registro
	Seguimiento      string
	ColaRetail       []Cola
	ColaPrioritario  []Cola
	ColaNormal       []Cola
	PaquetesEnMarcha []PaqueteEnMarcha
}

//ActualizarRegistroPaqueteCamionRetail es
func (s *Server) ActualizarRegistroPaqueteCamionRetail(ctx context.Context, message *PaqueteEnviado) (*PaqueteEnviado, error) {
	// Actualizar PaquetesEnMarcha

	// Buscar indice
	indice := 0
	contador := 0
	for _, elem := range s.PaquetesEnMarcha {
		if message.Idpaquete == elem.idpaquete {
			indice = contador
		} else {
			contador++
		}
	}

	// Agregar paquete modificado al final
	paqueteModificado := PaqueteEnMarcha{
		idpaquete:     message.Idpaquete,
		estado:        message.Estado,
		idcamion:      message.Idcamion,
		idseguimiento: message.Seguimiento,
		intentos:      message.Intentos,
		origen:        message.Origen,
		destino:       message.Destino,
		timestamp:     time.Now(),
	}

	s.PaquetesEnMarcha = append(s.PaquetesEnMarcha, paqueteModificado)

	// Reemplazar indice por el ultimo, y acortar lista
	s.PaquetesEnMarcha[indice] = s.PaquetesEnMarcha[len(s.PaquetesEnMarcha)-1]
	s.PaquetesEnMarcha = s.PaquetesEnMarcha[:len(s.PaquetesEnMarcha)-1]

	// Notificar a financiero
	return message, nil
}

//EntregarPaqueteCamionRetail es
func (s *Server) EntregarPaqueteCamionRetail(ctx context.Context, message *IdCamion) (*ColaPaquete, error) {
	messageColaPaqueteError := ColaPaquete{
		Idpaquete:   "NoPaquetes",
		Seguimiento: "",
		Tipo:        "",
		Valor:       "",
		Intentos:    "",
		Estado:      "",
		Origen:      "",
		Destino:     "",
	}
	if len(s.ColaRetail) < 1 && len(s.ColaPrioritario) < 1 { // NoPaquetes ninguna cola
		return &messageColaPaqueteError, nil
	}

	indicemascaro := 0
	elemmascaro := 0
	iterador := 0
	/*
		fmt.Println("---")
		fmt.Println(len(s.ColaPrioritario))
		fmt.Println("---")
	*/

	if len(s.ColaRetail) >= 1 {
		//fmt.Println("Entre a ColaRetail")
		// Busqueda paquete mas caro
		for _, elem := range s.ColaRetail {
			value, _ := strconv.Atoi(elem.valor)
			// Buscar elemento más caro
			if value > elemmascaro {
				indicemascaro = iterador
				elemmascaro = value
			}
			iterador++
		}

		// Preparacion retorno paquete Retail mas caro

		messageColaPaquete := ColaPaquete{
			Idpaquete:   s.ColaRetail[indicemascaro].idpaquete,
			Seguimiento: s.ColaRetail[indicemascaro].seguimiento,
			Tipo:        s.ColaRetail[indicemascaro].tipo,
			Valor:       s.ColaRetail[indicemascaro].valor,
			Intentos:    s.ColaRetail[indicemascaro].intentos,
			Estado:      s.ColaRetail[indicemascaro].estado,
			Origen:      s.ColaRetail[indicemascaro].origen,
			Destino:     s.ColaRetail[indicemascaro].destino,
		}

		nuevoPaquete := PaqueteEnMarcha{
			idpaquete:     s.ColaRetail[indicemascaro].idpaquete,
			estado:        s.ColaRetail[indicemascaro].estado,
			idcamion:      message.Idcamion,
			idseguimiento: s.ColaRetail[indicemascaro].seguimiento,
			intentos:      s.ColaRetail[indicemascaro].intentos,
			origen:        s.ColaRetail[indicemascaro].origen,
			destino:       s.ColaRetail[indicemascaro].destino,
			timestamp:     time.Time{},
		}

		// Quitar elemento mas caro de la cola

		// https://yourbasic.org/golang/delete-element-slice/
		s.ColaRetail[indicemascaro] = s.ColaRetail[len(s.ColaRetail)-1]
		s.ColaRetail = s.ColaRetail[:len(s.ColaRetail)-1]
		//log.Println(&messageColaPaquete)

		//fmt.Println(s.ColaRetail)

		s.PaquetesEnMarcha = append(s.PaquetesEnMarcha, nuevoPaquete)

		return &messageColaPaquete, nil
	}

	indicemascaro = 0
	elemmascaro = 0
	iterador = 0
	if len(s.ColaPrioritario) >= 1 {
		// Busqueda paquete mas caro
		for _, elem := range s.ColaPrioritario {
			value, _ := strconv.Atoi(elem.valor)
			// Buscar elemento más caro
			if value > elemmascaro {
				indicemascaro = iterador
				elemmascaro = value
			}
			iterador++
		}

		// Preparacion retorno paquete Prioritario mas caro

		messageColaPaquete := ColaPaquete{
			Idpaquete:   s.ColaPrioritario[indicemascaro].idpaquete,
			Seguimiento: s.ColaPrioritario[indicemascaro].seguimiento,
			Tipo:        s.ColaPrioritario[indicemascaro].tipo,
			Valor:       s.ColaPrioritario[indicemascaro].valor,
			Intentos:    s.ColaPrioritario[indicemascaro].intentos,
			Estado:      s.ColaPrioritario[indicemascaro].estado,
			Origen:      s.ColaPrioritario[indicemascaro].origen,
			Destino:     s.ColaPrioritario[indicemascaro].destino,
		}

		nuevoPaquete := PaqueteEnMarcha{
			idpaquete:     s.ColaPrioritario[indicemascaro].idpaquete,
			estado:        s.ColaPrioritario[indicemascaro].estado,
			idcamion:      message.Idcamion,
			idseguimiento: s.ColaPrioritario[indicemascaro].seguimiento,
			intentos:      s.ColaPrioritario[indicemascaro].intentos,
			origen:        s.ColaPrioritario[indicemascaro].origen,
			destino:       s.ColaPrioritario[indicemascaro].destino,
			timestamp:     time.Time{},
		}

		// Quitar elemento mas caro de la cola

		// https://yourbasic.org/golang/delete-element-slice/
		s.ColaPrioritario[indicemascaro] = s.ColaPrioritario[len(s.ColaPrioritario)-1]
		s.ColaPrioritario = s.ColaPrioritario[:len(s.ColaPrioritario)-1]
		//log.Println(&messageColaPaquete)

		s.PaquetesEnMarcha = append(s.PaquetesEnMarcha, nuevoPaquete)

		return &messageColaPaquete, nil
	}
	return &messageColaPaqueteError, nil
}

//ActualizarRegistroPaqueteCamionNormal es
func (s *Server) ActualizarRegistroPaqueteCamionNormal(ctx context.Context, message *PaqueteEnviado) (*PaqueteEnviado, error) {
	// Actualizar PaquetesEnMarcha

	// Buscar indice
	indice := 0
	contador := 0
	for _, elem := range s.PaquetesEnMarcha {
		if message.Idpaquete == elem.idpaquete {
			indice = contador
		} else {
			contador++
		}
	}

	// Agregar paquete modificado al final
	paqueteModificado := PaqueteEnMarcha{
		idpaquete:     message.Idpaquete,
		estado:        message.Estado,
		idcamion:      message.Idcamion,
		idseguimiento: message.Seguimiento,
		intentos:      message.Intentos,
		origen:        message.Origen,
		destino:       message.Destino,
		timestamp:     time.Now(),
	}

	s.PaquetesEnMarcha = append(s.PaquetesEnMarcha, paqueteModificado)

	// Reemplazar indice por el ultimo, y acortar lista
	s.PaquetesEnMarcha[indice] = s.PaquetesEnMarcha[len(s.PaquetesEnMarcha)-1]
	s.PaquetesEnMarcha = s.PaquetesEnMarcha[:len(s.PaquetesEnMarcha)-1]

	// Notificar a financiero
	return message, nil
}

//EntregarPaqueteCamionNormal es
func (s *Server) EntregarPaqueteCamionNormal(ctx context.Context, message *IdCamion) (*ColaPaquete, error) {
	messageColaPaqueteError := ColaPaquete{
		Idpaquete:   "NoPaquetes",
		Seguimiento: "",
		Tipo:        "",
		Valor:       "",
		Intentos:    "",
		Estado:      "",
		Origen:      "",
		Destino:     "",
	}
	if len(s.ColaPrioritario) < 1 && len(s.ColaNormal) < 1 { // NoPaquetes ninguna cola
		return &messageColaPaqueteError, nil
	}

	indicemascaro := 0
	elemmascaro := 0
	iterador := 0
	/*
		fmt.Println("---")
		fmt.Println(len(s.ColaPrioritario))
		fmt.Println("---")
	*/

	if len(s.ColaPrioritario) >= 1 {
		//fmt.Println("Entre a ColaPrioritario")
		// Busqueda paquete mas caro
		for _, elem := range s.ColaPrioritario {
			value, _ := strconv.Atoi(elem.valor)
			// Buscar elemento más caro
			if value > elemmascaro {
				indicemascaro = iterador
				elemmascaro = value
			}
			iterador++
		}

		// Preparacion retorno paquete prioritario mas caro

		messageColaPaquete := ColaPaquete{
			Idpaquete:   s.ColaPrioritario[indicemascaro].idpaquete,
			Seguimiento: s.ColaPrioritario[indicemascaro].seguimiento,
			Tipo:        s.ColaPrioritario[indicemascaro].tipo,
			Valor:       s.ColaPrioritario[indicemascaro].valor,
			Intentos:    s.ColaPrioritario[indicemascaro].intentos,
			Estado:      s.ColaPrioritario[indicemascaro].estado,
			Origen:      s.ColaPrioritario[indicemascaro].origen,
			Destino:     s.ColaPrioritario[indicemascaro].destino,
		}

		nuevoPaquete := PaqueteEnMarcha{
			idpaquete:     s.ColaPrioritario[indicemascaro].idpaquete,
			estado:        s.ColaPrioritario[indicemascaro].estado,
			idcamion:      message.Idcamion,
			idseguimiento: s.ColaPrioritario[indicemascaro].seguimiento,
			intentos:      s.ColaPrioritario[indicemascaro].intentos,
			origen:        s.ColaPrioritario[indicemascaro].origen,
			destino:       s.ColaPrioritario[indicemascaro].destino,
			timestamp:     time.Time{},
		}

		// Quitar elemento mas caro de la cola

		// https://yourbasic.org/golang/delete-element-slice/
		s.ColaPrioritario[indicemascaro] = s.ColaPrioritario[len(s.ColaPrioritario)-1]
		s.ColaPrioritario = s.ColaPrioritario[:len(s.ColaPrioritario)-1]
		//log.Println(&messageColaPaquete)

		//fmt.Println(s.ColaPrioritario)

		s.PaquetesEnMarcha = append(s.PaquetesEnMarcha, nuevoPaquete)

		return &messageColaPaquete, nil
	}

	indicemascaro = 0
	elemmascaro = 0
	iterador = 0
	if len(s.ColaNormal) >= 1 {
		// Busqueda paquete mas caro
		for _, elem := range s.ColaNormal {
			value, _ := strconv.Atoi(elem.valor)
			// Buscar elemento más caro
			if value > elemmascaro {
				indicemascaro = iterador
				elemmascaro = value
			}
			iterador++
		}

		// Preparacion retorno paquete normal mas caro

		messageColaPaquete := ColaPaquete{
			Idpaquete:   s.ColaNormal[indicemascaro].idpaquete,
			Seguimiento: s.ColaNormal[indicemascaro].seguimiento,
			Tipo:        s.ColaNormal[indicemascaro].tipo,
			Valor:       s.ColaNormal[indicemascaro].valor,
			Intentos:    s.ColaNormal[indicemascaro].intentos,
			Estado:      s.ColaNormal[indicemascaro].estado,
			Origen:      s.ColaNormal[indicemascaro].origen,
			Destino:     s.ColaNormal[indicemascaro].destino,
		}

		nuevoPaquete := PaqueteEnMarcha{
			idpaquete:     s.ColaNormal[indicemascaro].idpaquete,
			estado:        s.ColaNormal[indicemascaro].estado,
			idcamion:      message.Idcamion,
			idseguimiento: s.ColaNormal[indicemascaro].seguimiento,
			intentos:      s.ColaNormal[indicemascaro].intentos,
			origen:        s.ColaNormal[indicemascaro].origen,
			destino:       s.ColaNormal[indicemascaro].destino,
			timestamp:     time.Time{},
		}

		// Quitar elemento mas caro de la cola

		// https://yourbasic.org/golang/delete-element-slice/
		s.ColaNormal[indicemascaro] = s.ColaNormal[len(s.ColaNormal)-1]
		s.ColaNormal = s.ColaNormal[:len(s.ColaNormal)-1]
		//log.Println(&messageColaPaquete)

		s.PaquetesEnMarcha = append(s.PaquetesEnMarcha, nuevoPaquete)

		return &messageColaPaquete, nil
	}
	return &messageColaPaqueteError, nil
}

//RecibirOrdenPymes es
func (s *Server) RecibirOrdenPymes(ctx context.Context, message *Ordenclientepymes) (*Ordenseguimiento, error) {
	var tipo string
	messageOrdenseguimiento := Ordenseguimiento{
		Nordenseguimiento: s.Seguimiento,
	}
	if message.Prioritario == true {
		tipo = "prioritario"
	} else {
		tipo = "normal"
	}
	nuevaEntrada := Registro{
		timestamp:   time.Now(),
		idpaquete:   message.Id,
		tipo:        tipo,
		nombre:      message.Producto,
		valor:       strconv.Itoa(int(message.Valor)),
		origen:      message.Tienda,
		destino:     message.Destino,
		seguimiento: s.Seguimiento,
	}

	nuevaCola := Cola{
		idpaquete:   message.Id,
		seguimiento: s.Seguimiento,
		tipo:        tipo,
		valor:       strconv.Itoa(int(message.Valor)),
		intentos:    "0",
		estado:      "En bodega",
		origen:      message.Tienda,
		destino:     message.Destino,
	}
	// Agregar paquete a una de las tres colas de Server
	if tipo == "prioritario" {
		// Add cola prioritario
		s.ColaPrioritario = append(s.ColaPrioritario, nuevaCola)
	} else {
		// Add cola normal
		s.ColaNormal = append(s.ColaNormal, nuevaCola)
	}

	// Actualizacion codigo de seguimiento (atributo Server)
	numerosiguiente, _ := strconv.Atoi(s.Seguimiento)
	numerosiguiente++
	s.Seguimiento = strconv.Itoa(numerosiguiente)
	s.ListaRegistro = append(s.ListaRegistro, nuevaEntrada)

	//log.Println(s.ColaNormal)
	//log.Println(s.ListaRegistro)
	return &messageOrdenseguimiento, nil
}

//RedecirOrdenPymes es
func (s *Server) RedecirOrdenPymes(ctx context.Context, message *Ordenclientepymes) (*Ordenclientepymes, error) {
	log.Println("Recibi la orden")
	log.Printf("Producto pymes: %s", message.Producto)
	return message, nil
}

//RecibirOrdenRetail es
func (s *Server) RecibirOrdenRetail(ctx context.Context, message *Ordenclienteretail) (*Ordenseguimiento, error) {
	messageOrdenseguimiento := Ordenseguimiento{
		Nordenseguimiento: "9999",
	}

	nuevaEntrada := Registro{
		timestamp:   time.Now(),
		idpaquete:   message.Id,
		tipo:        "retail",
		nombre:      message.Producto,
		valor:       strconv.Itoa(int(message.Valor)),
		origen:      message.Tienda,
		destino:     message.Destino,
		seguimiento: "9999",
	}

	nuevaCola := Cola{
		idpaquete:   message.Id,
		seguimiento: "9999",
		tipo:        "retail",
		valor:       strconv.Itoa(int(message.Valor)),
		intentos:    "0",
		estado:      "En bodega",
		origen:      message.Tienda,
		destino:     message.Destino,
	}

	// Agregar paquete a una de las tres colas de Server
	s.ColaRetail = append(s.ColaRetail, nuevaCola)
	//log.Println(s.ColaRetail)

	s.ListaRegistro = append(s.ListaRegistro, nuevaEntrada)
	//log.Println(s.ListaRegistro)
	return &messageOrdenseguimiento, nil
}

//RedecirOrdenRetail es
func (s *Server) RedecirOrdenRetail(ctx context.Context, message *Ordenclienteretail) (*Ordenclienteretail, error) {
	log.Println("Recibi la orden retail")
	log.Printf("Producto retail: %s", message.Producto)
	return message, nil
}

//CodigoSeguimiento es
func (s *Server) CodigoSeguimiento(ctx context.Context, message *Ordenseguimiento) (*Estado, error) {
	mensajeError := Estado{
		Estado: "En bodega",
	}
	for _, elem := range s.PaquetesEnMarcha {
		if message.Nordenseguimiento == elem.idseguimiento {
			mensajeFinal := Estado{
				Estado: elem.estado,
			}
			return &mensajeFinal, nil
		}
	}

	return &mensajeError, nil
}
