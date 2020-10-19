########## Sistema Clientes ##########

Integrantes:
	Nicolás Durán 201673513-5
	Matías Marchant 201673556-9

Grupo: N&M's

Para ejecutar, ingrese a la carpeta “Lab1SDcliente” con el comando ‘cd Lab1SDcliente’ y luego ejecute 'make run' para ejecutar los archivos.
El sistema clientes le pedirá que ingrese si es cliente retail o pymes, las únicas opciones que funcionan es 'retail' ó 'pymes' en minúscula.
Se recomienda hacer la opción retail primero, que ingrese las órdenes y después cuando éste termine, hacer la opción pymes ya que la opción pymes no tiene
tiempo de finalización debido a que pide estado de paquetes indefinidamente.

La disposición de los roles de las máquinas virtuales es la siguiente:
Dist37 -> Logística
Dist38 -> Camiones
Dist39 -> Cliente
Dist40 -> Finanzas

Consideraciones extra:

En camion.go:
Si el paquete no es recibido, envés de intentar con el otro paquete y después con el no recibido, por simplicidad el camión espera el tiempo de envío de paquete definido al inicio del programa y vuelve a intentar con el primer paquete siguiendo las reglas del enunciado.

No es necesario que los camiones de retail tengan que llevar almenos un paquete de retail para poder tomar paquetes prioritarios, simplemente los toman (siempre prefiriendo los de retail)
