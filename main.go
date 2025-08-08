/*
8-Puzzle Solver - Inteligencia Artificial I | USAC

DESCRIPCIÓN:
Esta aplicación gráfica resuelve el clásico problema del 8-puzzle usando algoritmos de búsqueda
informada y no informada. El 8-puzzle consiste en un tablero de 3x3 con 8 fichas numeradas
del 1 al 8 y un espacio vacío. El objetivo es ordenar las fichas desde cualquier configuración
inicial hasta el estado objetivo: 1,2,3,4,5,6,7,8,vacío.

ALGORITMOS IMPLEMENTADOS:
  - A* con Heurística Manhattan: Algoritmo de búsqueda informada que utiliza f(n) = g(n) + h(n)
    donde g(n) es el costo del camino y h(n) es la distancia Manhattan al objetivo.
  - Búsqueda en Anchura (BFS): Algoritmo de búsqueda no informada que explora nivel por nivel
    garantizando la solución óptima en número de movimientos.

CARACTERÍSTICAS PRINCIPALES:
- Interfaz gráfica moderna y profesional usando el framework Fyne
- Visualización en tiempo real del estado del puzzle y heurística Manhattan
- Animaciones suaves para mostrar el movimiento de las piezas
- Modo "Paso a Paso" para visualizar la solución completa
- Métricas de rendimiento: tiempo de ejecución, número de pasos, eficiencia
- Generación de configuraciones aleatorias garantizadas como solucionables
- Barra de progreso visual durante la ejecución de la solución

ARQUITECTURA DEL SISTEMA:
- Estado: Representación de una configuración del puzzle con información de búsqueda
- PuzzleButton: Widget personalizado con animación para cada celda del tablero
- PuzzleApp: Controlador principal que gestiona la lógica de negocio y la interfaz
- MyTheme: Tema visual personalizado para una experiencia profesional

TECNOLOGÍAS UTILIZADAS:
- Lenguaje: Go (Golang) 1.21+
- Framework GUI: Fyne v2.4.0
- Paradigma: Programación orientada a objetos y concurrente

CASOS DE USO:
1. Demostración académica de algoritmos de IA
2. Herramienta educativa para visualizar búsqueda heurística
3. Comparación de rendimiento entre algoritmos informados vs no informados
4. Análisis de complejidad computacional en problemas de búsqueda

COMPLEJIDAD COMPUTACIONAL:
- A*: O(b^d) donde b es el factor de ramificación y d la profundidad de la solución
- BFS: O(b^d) pero sin guía heurística, explora más estados
- Espacio: O(b^d) para almacenar los estados explorados

AUTOR: Joel Lombardo (itsmjoe)
INSTITUCIÓN: Universidad de San Carlos de Guatemala - Facultad de Ingeniería
CURSO: Inteligencia Artificial I
FECHA: Agosto 2025
LICENCIA: MIT License - Uso académico y educativo

REQUISITOS DEL SISTEMA:
- Go 1.21 o superior
- Sistema operativo: Linux, macOS, Windows
- Memoria RAM: Mínimo 512MB disponibles
- Resolución: 1024x768 o superior para óptima visualización

INSTRUCCIONES DE COMPILACIÓN:
1. go mod init 8-puzzle-solver
2. go get fyne.io/fyne/v2/app
3. go build -o 8-puzzle-solver main.go
4. ./8-puzzle-solver

REFERENCIAS ACADÉMICAS:
- Russell, S. & Norvig, P. "Artificial Intelligence: A Modern Approach"
- Hart, P. E., Nilsson, N. J., & Raphael, B. "A Formal Basis for the Heuristic Determination of Minimum Cost Paths"
*/
package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MyTheme implementa el interface fyne.Theme para definir un tema visual personalizado
// que mejora la experiencia del usuario con colores profesionales y consistentes.
type MyTheme struct{}

func (t MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Paleta cálida en tonos anaranjados, café y verdes
	// Inspirada en colores terracota y naturales para una experiencia acogedora
	switch name {
	case theme.ColorNamePrimary:
		return color.NRGBA{210, 120, 60, 255} // Naranja terracota principal
	case theme.ColorNameButton:
		return color.NRGBA{255, 248, 240, 255} // Crema suave para botones
	case theme.ColorNameForeground:
		return color.NRGBA{92, 64, 51, 255} // Café oscuro para texto
	case theme.ColorNameBackground:
		return color.NRGBA{253, 248, 243, 255} // Fondo crema muy claro
	case theme.ColorNameInputBackground:
		return color.NRGBA{255, 250, 245, 255} // Crema para inputs
	case theme.ColorNameSelection:
		return color.NRGBA{210, 120, 60, 40} // Naranja transparente para selección
	case theme.ColorNameHover:
		return color.NRGBA{210, 120, 60, 20} // Naranja muy transparente para hover
	case theme.ColorNameSuccess:
		return color.NRGBA{107, 142, 35, 255} // Verde musgo para RESOLVER
	case theme.ColorNameWarning:
		return color.NRGBA{139, 195, 74, 255} // Verde claro para animación (más visible)
	case theme.ColorNameError:
		return color.NRGBA{184, 92, 46, 255} // Café rojizo para errores
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	// Font retorna la fuente tipográfica por defecto del sistema.
	return theme.DefaultTheme().Font(style)
}

func (t MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	// Icon retorna los iconos del tema por defecto del sistema.
	return theme.DefaultTheme().Icon(name)
}

func (t MyTheme) Size(name fyne.ThemeSizeName) float32 {
	// Size retorna tamaños personalizados para diferentes elementos de texto.
	// Mejora la jerarquía visual y legibilidad de la interfaz.
	switch name {
	case theme.SizeNameText:
		return 14 // Tamaño base para texto normal
	case theme.SizeNameHeadingText:
		return 24 // Tamaño para títulos principales
	case theme.SizeNameSubHeadingText:
		return 18 // Tamaño para subtítulos
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// Estado representa un nodo en el árbol de búsqueda del problema del 8-puzzle.
// Contiene la configuración actual del tablero, referencias para reconstruir el camino,
// y metainformación para los algoritmos de búsqueda.
type Estado struct {
	tablero [9]int  // Configuración actual: posiciones 0-8, valor 0 representa espacio vacío
	padre   *Estado // Referencia al estado padre para reconstruir la solución
	costo   int     // g(n): Costo acumulado desde el estado inicial (profundidad)
	accion  string  // Acción realizada para llegar a este estado desde el padre
}

// PuzzleButton es un widget personalizado que extiende widget.Button de Fyne
// para representar cada celda del tablero del 8-puzzle con capacidades de animación.
// Implementa feedback visual para mostrar qué pieza se está moviendo durante la solución.
type PuzzleButton struct {
	widget.Button      // Hereda funcionalidad básica de botón de Fyne
	numero        int  // Valor numérico mostrado (1-8, 0 para vacío)
	esVacio       bool // Flag que indica si esta celda es el espacio vacío
	destacado     bool // Flag para animación temporal de movimiento
}

func NewPuzzleButton(numero int) *PuzzleButton {
	// NewPuzzleButton es el constructor que crea e inicializa un nuevo botón del puzzle.
	// Configura el widget base de Fyne y establece el estado inicial visual.
	// Parámetro numero: valor inicial del botón (0 para espacio vacío).
	btn := &PuzzleButton{numero: numero, destacado: false}
	btn.ExtendBaseWidget(btn) // Vincula el widget personalizado con el sistema de Fyne
	btn.esVacio = (numero == 0)
	btn.actualizarEstilo()
	return btn
}

func (pb *PuzzleButton) actualizarEstilo() {
	// actualizarEstilo actualiza el aspecto visual del botón según su estado actual.
	// Diferencia visualmente entre números, espacios vacíos y estados de animación.
	if pb.esVacio {
		pb.SetText("")                       // Completamente vacío, sin números
		pb.Importance = widget.LowImportance // Estilo tenue para el espacio vacío
	} else {
		pb.SetText(strconv.Itoa(pb.numero))
		if pb.destacado {
			pb.Importance = widget.WarningImportance // Naranja suave para indicar movimiento
		} else {
			pb.Importance = widget.HighImportance // Azul estándar para números
		}
	}
	pb.Refresh() // Fuerza la actualización visual del widget
}

func (pb *PuzzleButton) setNumero(numero int) {
	// setNumero actualiza el valor numérico mostrado en el botón y reestablece su estilo.
	// Utilizado para actualizar la representación visual del tablero.
	pb.numero = numero
	pb.esVacio = (numero == 0)
	pb.destacado = false
	pb.actualizarEstilo()
}

func (pb *PuzzleButton) destacar() {
	// destacar resalta temporalmente el botón para indicar que esta pieza se va a mover.
	if !pb.esVacio && pb.numero != 0 {
		// Solo destacar si no es espacio vacío y contiene un número válido
		pb.SetText(strconv.Itoa(pb.numero))
		pb.destacado = true
		pb.Importance = widget.WarningImportance // Verde claro brillante para animación
		pb.Refresh()

		// Quitar destaque después de 800ms usando una goroutine independiente
		go func() {
			time.Sleep(800 * time.Millisecond)
			fyne.DoAndWait(func() {
				pb.destacado = false
				if pb.numero != 0 {
					pb.SetText(strconv.Itoa(pb.numero))
				} else {
					pb.SetText("") // Mantener vacío si es 0
				}
				pb.Importance = widget.HighImportance
				pb.Refresh()
			})
		}()
	}
}

// PuzzleApp es la estructura principal que gestiona toda la aplicación del 8-puzzle.
// Implementa el patrón MVC (Modelo-Vista-Controlador) donde actúa como controlador,
// gestionando la lógica de negocio, el estado del puzzle y la interfaz gráfica.
type PuzzleApp struct {
	window       fyne.Window         // Ventana principal de la aplicación
	botones      [9]*PuzzleButton    // Array de botones representando el tablero 3x3
	estadoActual [9]int              // Estado actual del puzzle (modelo de datos)
	objetivo     [9]int              // Estado objetivo del puzzle [1,2,3,4,5,6,7,8,0]
	solucion     []Estado            // Secuencia de estados que resuelven el puzzle
	paso         int                 // Índice del paso actual en la visualización de la solución
	infoLabel    *widget.RichText    // Panel de información con formato enriquecido
	estadoLabel  *widget.Label       // Etiqueta de estado y heurística en tiempo real
	algoritmo    *widget.Select      // Selector de algoritmo de búsqueda
	progressBar  *widget.ProgressBar // Barra de progreso visual para la solución
}

func NuevaPuzzleApp() *PuzzleApp {
	// NuevaPuzzleApp es el constructor que inicializa la estructura principal de la aplicación.
	// Establece el estado objetivo estándar del 8-puzzle y valores iniciales.
	return &PuzzleApp{
		objetivo: [9]int{1, 2, 3, 4, 5, 6, 7, 8, 0}, // Configuración objetivo estándar
		paso:     0,
	}
}

func encontrarVacio(tablero [9]int) int {
	// encontrarVacio localiza y retorna la posición del espacio vacío (representado por 0) en el tablero.
	// Es una función auxiliar fundamental para generar movimientos válidos.
	// Retorna: índice de la posición vacía (0-8), o -1 si no se encuentra.
	for i := 0; i < 9; i++ {
		if tablero[i] == 0 {
			return i
		}
	}
	return -1 // No debería ocurrir en un puzzle válido
}

func esObjetivo(tablero [9]int, objetivo [9]int) bool {
	// esObjetivo verifica si la configuración actual del tablero coincide con el estado objetivo.
	// Es la condición de parada para los algoritmos de búsqueda.
	// Retorna: true si el puzzle está resuelto, false en caso contrario.
	for i := 0; i < 9; i++ {
		if tablero[i] != objetivo[i] {
			return false
		}
	}
	return true
}

func heuristicaManhattan(tablero [9]int) int {
	// heuristicaManhattan calcula la función heurística h(n) para el algoritmo A*.
	// La distancia Manhattan es la suma de distancias horizontales y verticales
	// de cada ficha desde su posición actual hasta su posición objetivo.
	// Esta heurística es admisible (nunca sobreestima) y consistente (monótona).
	//
	// Complejidad temporal: O(1) - siempre evalúa 9 posiciones
	// Complejidad espacial: O(1) - usa memoria constante
	//
	// Retorna: suma total de distancias Manhattan de todas las fichas mal ubicadas
	distancia := 0
	for i := 0; i < 9; i++ {
		if tablero[i] != 0 {
			// Calcular posición actual en coordenadas (fila, columna)
			fila_actual := i / 3
			col_actual := i % 3

			// Calcular posición objetivo en coordenadas (fila, columna)
			valor := tablero[i]
			fila_objetivo := (valor - 1) / 3 // valor-1 porque numeramos desde 1
			col_objetivo := (valor - 1) % 3

			// Sumar distancia Manhattan: |x1-x2| + |y1-y2|
			distancia += abs(fila_actual-fila_objetivo) + abs(col_actual-col_objetivo)
		}
	}
	return distancia
}

func abs(x int) int {
	// abs retorna el valor absoluto de un entero.
	// Función auxiliar para cálculos de distancia.
	if x < 0 {
		return -x
	}
	return x
}

func generarMovimientos(tablero [9]int) []Estado {
	// generarMovimientos genera todos los movimientos válidos desde el estado actual del tablero.
	// Un movimiento válido consiste en intercambiar el espacio vacío con una ficha adyacente
	// (arriba, abajo, izquierda, derecha).
	//
	// Retorna: slice de Estados representando todos los sucesores posibles.
	movimientos := []Estado{}
	posVacio := encontrarVacio(tablero)

	// Convertir posición lineal a coordenadas 2D
	fila := posVacio / 3
	col := posVacio % 3

	// Generar movimiento hacia arriba (intercambiar con ficha de arriba)
	if fila > 0 {
		nuevo := tablero
		nueva_pos := (fila-1)*3 + col
		nuevo[posVacio], nuevo[nueva_pos] = nuevo[nueva_pos], nuevo[posVacio]
		movimientos = append(movimientos, Estado{
			tablero: nuevo,
			accion:  "Arriba",
		})
	}

	// Generar movimiento hacia abajo (intercambiar con ficha de abajo)
	if fila < 2 {
		nuevo := tablero
		nueva_pos := (fila+1)*3 + col
		nuevo[posVacio], nuevo[nueva_pos] = nuevo[nueva_pos], nuevo[posVacio]
		movimientos = append(movimientos, Estado{
			tablero: nuevo,
			accion:  "Abajo",
		})
	}

	// Generar movimiento hacia izquierda (intercambiar con ficha de la izquierda)
	if col > 0 {
		nuevo := tablero
		nueva_pos := fila*3 + (col - 1)
		nuevo[posVacio], nuevo[nueva_pos] = nuevo[nueva_pos], nuevo[posVacio]
		movimientos = append(movimientos, Estado{
			tablero: nuevo,
			accion:  "Izquierda",
		})
	}

	// Generar movimiento hacia derecha (intercambiar con ficha de la derecha)
	if col < 2 {
		nuevo := tablero
		nueva_pos := fila*3 + (col + 1)
		nuevo[posVacio], nuevo[nueva_pos] = nuevo[nueva_pos], nuevo[posVacio]
		movimientos = append(movimientos, Estado{
			tablero: nuevo,
			accion:  "Derecha",
		})
	}

	return movimientos
}

func busquedaAEstrella(inicial [9]int, objetivo [9]int) []Estado {
	// busquedaAEstrella implementa el algoritmo A* (A-estrella) para encontrar la solución óptima.
	//
	// ALGORITMO A*:
	// 1. Mantiene dos listas: ABIERTA (nodos por explorar) y CERRADA (nodos explorados)
	// 2. Selecciona el nodo con menor f(n) = g(n) + h(n) de la lista ABIERTA
	// 3. Si es el objetivo, reconstruye y retorna el camino
	// 4. Si no, expande sus sucesores y los agrega a ABIERTA si no están en CERRADA
	// 5. Repite hasta encontrar solución o agotar posibilidades
	//
	// PROPIEDADES:
	// - Completitud: Siempre encuentra solución si existe
	// - Optimalidad: Garantiza la solución de menor costo con heurística admisible
	// - Complejidad temporal: O(b^d) donde b=factor ramificación, d=profundidad solución
	// - Complejidad espacial: O(b^d) para almacenar nodos en memoria
	//
	// PARÁMETROS:
	// - inicial: configuración inicial del tablero [9]int
	// - objetivo: configuración objetivo del tablero [9]int
	//
	// RETORNA: slice de Estados representando el camino solución (vacío si no hay solución)

	// Inicializar lista ABIERTA con el estado inicial
	abierta := []Estado{{tablero: inicial, costo: 0}}
	// Lista CERRADA para evitar reexplorar estados
	cerrada := []string{}

	for len(abierta) > 0 {
		// Encontrar el nodo con menor f(n) = g(n) + h(n) en la lista ABIERTA
		indice_mejor := 0
		mejor_f := abierta[0].costo + heuristicaManhattan(abierta[0].tablero)

		for i := 1; i < len(abierta); i++ {
			f := abierta[i].costo + heuristicaManhattan(abierta[i].tablero)
			if f < mejor_f {
				mejor_f = f
				indice_mejor = i
			}
		}

		// Extraer el mejor nodo de la lista ABIERTA
		actual := abierta[indice_mejor]
		abierta = append(abierta[:indice_mejor], abierta[indice_mejor+1:]...)

		// Verificar si alcanzamos el estado objetivo
		if esObjetivo(actual.tablero, objetivo) {
			// Reconstruir el camino desde el objetivo hasta el inicio
			camino := []Estado{}
			estado := &actual
			for estado != nil {
				camino = append([]Estado{*estado}, camino...)
				estado = estado.padre
			}
			return camino
		}

		// Agregar el estado actual a la lista CERRADA
		cerrada = append(cerrada, fmt.Sprintf("%v", actual.tablero))

		// Generar y evaluar todos los sucesores del estado actual
		for _, movimiento := range generarMovimientos(actual.tablero) {
			estado_str := fmt.Sprintf("%v", movimiento.tablero)

			// Verificar si el sucesor ya fue explorado (está en CERRADA)
			ya_explorado := false
			for _, cerrado := range cerrada {
				if cerrado == estado_str {
					ya_explorado = true
					break
				}
			}

			// Si no está explorado, agregarlo a la lista ABIERTA
			if !ya_explorado {
				movimiento.padre = &actual
				movimiento.costo = actual.costo + 1
				abierta = append(abierta, movimiento)
			}
		}
	}

	return []Estado{} // Retornar lista vacía si no hay solución
}

func busquedaAnchura(inicial [9]int, objetivo [9]int) []Estado {
	// busquedaAnchura implementa el algoritmo de Búsqueda en Anchura (BFS) para resolver el puzzle.
	//
	// ALGORITMO BFS:
	// 1. Utiliza una cola FIFO (First In, First Out) para explorar nodos nivel por nivel
	// 2. Explora todos los nodos a profundidad d antes de explorar nodos a profundidad d+1
	// 3. Mantiene lista de visitados para evitar ciclos infinitos
	// 4. Garantiza encontrar la solución con menor número de movimientos
	//
	// PROPIEDADES:
	// - Completitud: Siempre encuentra solución si existe y el espacio es finito
	// - Optimalidad: Garantiza solución óptima en número de movimientos (costo uniforme)
	// - Complejidad temporal: O(b^d) donde b=factor ramificación, d=profundidad solución
	// - Complejidad espacial: O(b^d) para almacenar todos los nodos de cada nivel
	//
	// DIFERENCIAS CON A*:
	// - No usa heurística (búsqueda ciega)
	// - Explora más nodos que A* en promedio
	// - Útil cuando no se dispone de heurística admisible
	// - Mejor para problemas donde todos los movimientos tienen el mismo costo
	//
	// PARÁMETROS:
	// - inicial: configuración inicial del tablero [9]int
	// - objetivo: configuración objetivo del tablero [9]int
	//
	// RETORNA: slice de Estados representando el camino solución (vacío si no hay solución)

	// Inicializar cola FIFO con el estado inicial
	cola := []Estado{{tablero: inicial}}
	// Lista de estados visitados para evitar ciclos
	visitados := []string{}

	for len(cola) > 0 {
		// Extraer el primer elemento de la cola (FIFO)
		actual := cola[0]
		cola = cola[1:]

		// Verificar si alcanzamos el estado objetivo
		if esObjetivo(actual.tablero, objetivo) {
			// Reconstruir el camino desde el objetivo hasta el inicio
			camino := []Estado{}
			estado := &actual
			for estado != nil {
				camino = append([]Estado{*estado}, camino...)
				estado = estado.padre
			}
			return camino
		}

		// Marcar el estado actual como visitado
		estado_str := fmt.Sprintf("%v", actual.tablero)
		visitados = append(visitados, estado_str)

		// Generar y evaluar todos los sucesores del estado actual
		for _, movimiento := range generarMovimientos(actual.tablero) {
			mov_str := fmt.Sprintf("%v", movimiento.tablero)

			// Verificar si el sucesor ya fue visitado
			ya_visitado := false
			for _, v := range visitados {
				if v == mov_str {
					ya_visitado = true
					break
				}
			}

			// Si no está visitado, agregarlo a la cola
			if !ya_visitado {
				movimiento.padre = &actual
				cola = append(cola, movimiento)
			}
		}
	}

	return []Estado{} // Retornar lista vacía si no hay solución
}

func (app *PuzzleApp) actualizarTablero() {
	// actualizarTablero sincroniza la interfaz gráfica con el estado actual del modelo de datos.
	// Actualiza cada botón del tablero según los valores en estadoActual.
	for i := 0; i < 9; i++ {
		app.botones[i].setNumero(app.estadoActual[i])
	}
	app.actualizarEstado()
}

func (app *PuzzleApp) actualizarEstado() {
	// actualizarEstado actualiza la información de estado mostrada al usuario.
	// Muestra si el puzzle está resuelto o en proceso, junto con la heurística Manhattan actual.
	if esObjetivo(app.estadoActual, app.objetivo) {
		app.estadoLabel.SetText("ESTADO: RESUELTO")
		app.estadoLabel.Importance = widget.SuccessImportance
	} else {
		manhattan := heuristicaManhattan(app.estadoActual)
		app.estadoLabel.SetText(fmt.Sprintf("ESTADO: EN PROCESO | Heurística Manhattan: %d", manhattan))
		app.estadoLabel.Importance = widget.MediumImportance
	}
}

func (app *PuzzleApp) iniciar() {
	// iniciar reinicia el puzzle al estado objetivo ordenado y limpia todas las variables de control.
	// Utilizado para comenzar una nueva sesión o reiniciar después de completar una solución.
	app.estadoActual = [9]int{1, 2, 3, 4, 5, 6, 7, 8, 0}
	app.solucion = []Estado{}
	app.paso = 0
	app.progressBar.SetValue(0)
	app.actualizarTablero()

	app.infoLabel.ParseMarkdown("## SISTEMA INICIALIZADO\n\n**Estado:** Puzzle ordenado correctamente\n\n**Acción:** Presiona 'Mezclar' para comenzar")
}

func (app *PuzzleApp) mezclar() {
	// mezclar genera una configuración aleatoria del puzzle realizando movimientos válidos.
	// Garantiza que la configuración resultante sea solucionable al partir del estado objetivo
	// y aplicar movimientos válidos únicamente.
	rand.Seed(time.Now().UnixNano())

	app.infoLabel.ParseMarkdown("## MEZCLANDO PUZZLE\n\n**Estado:** Generando configuración aleatoria...\n\n**Por favor espera**")

	// Partir del estado objetivo y aplicar movimientos aleatorios válidos
	app.estadoActual = app.objetivo
	for i := 0; i < 150; i++ {
		movimientos := generarMovimientos(app.estadoActual)
		if len(movimientos) > 0 {
			// Seleccionar un movimiento aleatorio de los disponibles
			mov := movimientos[rand.Intn(len(movimientos))]
			app.estadoActual = mov.tablero
		}
	}

	// Limpiar variables de control para nueva búsqueda
	app.solucion = []Estado{}
	app.paso = 0
	app.progressBar.SetValue(0)
	app.actualizarTablero()

	manhattan := heuristicaManhattan(app.estadoActual)
	app.infoLabel.ParseMarkdown(fmt.Sprintf("## PUZZLE MEZCLADO\n\n**Estado:** Configuración aleatoria generada\n\n**Heurística Manhattan:** %d\n\n**Acción:** Selecciona algoritmo y presiona 'Resolver'", manhattan))
}

func (app *PuzzleApp) resolver() {
	// resolver ejecuta el algoritmo de búsqueda seleccionado para encontrar la solución al puzzle.
	// Mide el tiempo de ejecución y proporciona métricas de rendimiento al usuario.
	algoritmo_seleccionado := app.algoritmo.Selected
	app.infoLabel.ParseMarkdown(fmt.Sprintf("## RESOLVIENDO PUZZLE\n\n**Algoritmo:** %s\n\n**Estado:** Buscando solución óptima...\n\n**Por favor espera**", algoritmo_seleccionado))

	// Medir tiempo de ejecución del algoritmo
	inicio := time.Now()

	// Ejecutar el algoritmo seleccionado
	if algoritmo_seleccionado == "A* con Heurística Manhattan" {
		app.solucion = busquedaAEstrella(app.estadoActual, app.objetivo)
	} else {
		app.solucion = busquedaAnchura(app.estadoActual, app.objetivo)
	}

	duracion := time.Since(inicio)

	if len(app.solucion) > 0 {
		// Solución encontrada - preparar para visualización paso a paso
		app.paso = 0
		app.progressBar.SetValue(0)

		// Clasificar eficiencia basada en número de pasos
		eficiencia := "Excelente"
		if len(app.solucion) > 20 {
			eficiencia = "Buena"
		}
		if len(app.solucion) > 30 {
			eficiencia = "Promedio"
		}

		app.infoLabel.ParseMarkdown(fmt.Sprintf("## PUZZLE RESUELTO\n\n**Algoritmo:** %s\n\n**Pasos de solución:** %d\n\n**Tiempo de ejecución:** %d ms\n\n**Eficiencia:** %s\n\n**Acción:** Usa 'Paso a Paso' para ver la solución",
			algoritmo_seleccionado, len(app.solucion)-1, duracion.Milliseconds(), eficiencia))
	} else {
		// No se encontró solución (caso improbable en 8-puzzle válido)
		app.infoLabel.ParseMarkdown("## ERROR\n\n**Estado:** No se encontró solución\n\n**Acción:** Intenta mezclar nuevamente")
	}
}

func (app *PuzzleApp) siguientePaso() {
	// siguientePaso avanza un paso en la visualización de la solución encontrada.
	// Implementa animación para mostrar qué pieza se mueve en cada transición.
	if len(app.solucion) == 0 {
		app.infoLabel.ParseMarkdown("## ADVERTENCIA\n\n**Estado:** No hay solución cargada\n\n**Acción:** Primero resuelve el puzzle")
		return
	}

	if app.paso < len(app.solucion) {
		estadoNuevo := app.solucion[app.paso].tablero

		// Actualizar barra de progreso
		progreso := float64(app.paso) / float64(len(app.solucion)-1)
		app.progressBar.SetValue(progreso)

		// Determinar la acción realizada
		accion := "Estado inicial"
		if app.paso > 0 {
			accion = app.solucion[app.paso].accion

			// Encontrar y destacar la pieza que se mueve para feedback visual
			for i := 0; i < 9; i++ {
				if app.estadoActual[i] != 0 && app.estadoActual[i] != estadoNuevo[i] {
					app.botones[i].destacar()
					break
				}
			}

			// Actualizar el tablero después de un pequeño delay para permitir la animación
			go func(nuevoEstado [9]int) {
				time.Sleep(300 * time.Millisecond)
				fyne.DoAndWait(func() {
					app.estadoActual = nuevoEstado
					app.actualizarTablero()
				})
			}(estadoNuevo)
		} else {
			// Primer paso: actualización directa sin animación
			app.estadoActual = estadoNuevo
			app.actualizarTablero()
		}

		// Actualizar información del paso actual
		app.infoLabel.ParseMarkdown(fmt.Sprintf("## EJECUTANDO SOLUCIÓN\n\n**Paso:** %d de %d\n\n**Movimiento:** %s\n\n**Progreso:** %.1f%%",
			app.paso+1, len(app.solucion), accion, progreso*100))

		app.paso++
	} else {
		// Solución completada
		app.progressBar.SetValue(1.0)
		app.infoLabel.ParseMarkdown("## SOLUCIÓN COMPLETADA\n\n**Estado:** Puzzle resuelto exitosamente\n\n**Felicitaciones:** El algoritmo funcionó correctamente")
	}
}

func main() {
	// main es la función principal que inicializa y ejecuta la aplicación gráfica.
	//
	// FLUJO DE EJECUCIÓN:
	// 1. Crea la aplicación Fyne y configura la ventana principal
	// 2. Construye el header con información institucional
	// 3. Inicializa la cuadrícula 3x3 del puzzle con botones personalizados
	// 4. Configura los controles de interacción (algoritmo, botones, progreso)
	// 5. Ensambla el layout completo usando containers de Fyne
	// 6. Inicializa el puzzle en estado ordenado
	// 7. Muestra la ventana y comienza el loop de eventos
	//
	// ARQUITECTURA DE LA UI:
	// - Header: Título y información institucional
	// - Puzzle Board: Cuadrícula 3x3 interactiva
	// - Control Panel: Selector de algoritmo y botones de acción
	// - Info Panel: Información de estado, progreso y resultados
	// - Progress Bar: Visualización del progreso durante solución paso a paso
	//
	// PATRÓN DE DISEÑO:
	// Implementa el patrón Observer donde la UI reacciona a cambios en el modelo de datos
	// Utiliza el patrón Command para encapsular acciones de usuario en métodos

	// Crear aplicación Fyne con tema personalizado mejorado
	myApp := app.New()
	myApp.Settings().SetTheme(&MyTheme{}) // Activar tema personalizado

	// Configurar ventana principal
	ventana := myApp.NewWindow("8-Puzzle Solver - Inteligencia Artificial I | USAC")
	ventana.Resize(fyne.NewSize(700, 800))
	ventana.CenterOnScreen()

	// Inicializar controlador principal
	puzzleApp := NuevaPuzzleApp()
	puzzleApp.window = ventana

	// Crear header institucional con paleta cálida
	headerCard := canvas.NewRectangle(color.NRGBA{250, 237, 225, 255}) // Fondo crema cálido
	headerCard.Resize(fyne.NewSize(700, 120))

	titulo := canvas.NewText("8-PUZZLE SOLVER", color.NRGBA{138, 75, 40, 255}) // Café oscuro elegante
	titulo.Alignment = fyne.TextAlignCenter
	titulo.TextSize = 28
	titulo.TextStyle.Bold = true

	subtitulo := canvas.NewText("Inteligencia Artificial I • Universidad de San Carlos de Guatemala", color.NRGBA{160, 100, 70, 255}) // Café medio
	subtitulo.Alignment = fyne.TextAlignCenter
	subtitulo.TextSize = 14

	headerContent := container.NewVBox(
		container.NewPadded(titulo),
		subtitulo,
	)

	header := container.NewStack(headerCard, headerContent)

	// Crear cuadrícula interactiva del puzzle con espaciado óptimo
	cuadricula := container.NewGridWithColumns(3)
	for i := 0; i < 9; i++ {
		btn := NewPuzzleButton(1) // Inicializar con número temporal
		btn.Resize(fyne.NewSize(100, 100))

		// Configuración inicial del tablero ordenado
		if i < 8 {
			// Números del 1 al 8
			btn.SetText(strconv.Itoa(i + 1))
			btn.Importance = widget.HighImportance
		} else {
			// Espacio vacío (sin mostrar número 0)
			btn.SetText("")
			btn.Importance = widget.LowImportance
		}
		btn.Refresh()

		puzzleApp.botones[i] = btn
		cuadricula.Add(btn)
	}

	// Contenedor con fondo cálido para la cuadrícula
	puzzleCard := canvas.NewRectangle(color.NRGBA{255, 251, 247, 255}) // Crema muy claro
	puzzleCardContainer := container.NewStack(puzzleCard, cuadricula)

	// Panel de estado en tiempo real
	puzzleApp.estadoLabel = widget.NewLabel("")
	puzzleApp.estadoLabel.Alignment = fyne.TextAlignCenter
	puzzleApp.estadoLabel.TextStyle.Bold = true

	// Barra de progreso para visualización paso a paso
	puzzleApp.progressBar = widget.NewProgressBar()
	puzzleApp.progressBar.TextFormatter = func() string {
		return fmt.Sprintf("Progreso: %.0f%%", puzzleApp.progressBar.Value*100)
	}

	// Selector de algoritmo de búsqueda
	etiquetaAlgoritmo := widget.NewLabel("ALGORITMO DE BÚSQUEDA")
	etiquetaAlgoritmo.TextStyle.Bold = true
	etiquetaAlgoritmo.Alignment = fyne.TextAlignCenter

	puzzleApp.algoritmo = widget.NewSelect(
		[]string{"A* con Heurística Manhattan", "Búsqueda en Anchura (BFS)"},
		nil,
	)
	puzzleApp.algoritmo.SetSelected("A* con Heurística Manhattan")

	// Botones de control principal con paleta cálida
	btnIniciar := widget.NewButton("INICIAR", puzzleApp.iniciar)
	btnIniciar.Importance = widget.LowImportance // Café claro para acción neutral

	btnMezclar := widget.NewButton("MEZCLAR", puzzleApp.mezclar)
	btnMezclar.Importance = widget.HighImportance // Naranja terracota para acción principal

	btnResolver := widget.NewButton("RESOLVER", puzzleApp.resolver)
	btnResolver.Importance = widget.SuccessImportance // Verde musgo para acción positiva

	btnPaso := widget.NewButton("PASO A PASO", puzzleApp.siguientePaso)
	btnPaso.Importance = widget.WarningImportance // Naranja cálido para visualización

	// Panel de información detallada con texto enriquecido
	etiquetaInfo := widget.NewLabel("INFORMACIÓN DEL SISTEMA")
	etiquetaInfo.TextStyle.Bold = true
	etiquetaInfo.Alignment = fyne.TextAlignCenter

	puzzleApp.infoLabel = widget.NewRichTextFromMarkdown("")

	// Área de scroll para información extensa
	infoScroll := container.NewScroll(puzzleApp.infoLabel)
	infoScroll.SetMinSize(fyne.NewSize(400, 150))

	// Organización de controles con flujo lógico mejorado
	// Etiquetas de pasos para guiar al usuario con colores suaves
	etiquetaPaso1 := widget.NewLabel("1. CONFIGURACIÓN INICIAL")
	etiquetaPaso1.TextStyle.Bold = true
	etiquetaPaso1.Alignment = fyne.TextAlignCenter

	etiquetaPaso2 := widget.NewLabel("2. RESOLUCIÓN Y VISUALIZACIÓN")
	etiquetaPaso2.TextStyle.Bold = true
	etiquetaPaso2.Alignment = fyne.TextAlignCenter

	// Fila 1: Configuración inicial
	filaConfiguracion := container.NewGridWithColumns(2,
		btnIniciar, btnMezclar,
	)

	// Fila 2: Resolución y visualización
	filaResolucion := container.NewGridWithColumns(2,
		btnResolver, btnPaso,
	)

	// Panel de controles reorganizado para mejor UX
	controles := container.NewVBox(
		etiquetaAlgoritmo,
		puzzleApp.algoritmo,
		widget.NewSeparator(),
		etiquetaPaso1,
		filaConfiguracion,
		widget.NewSeparator(),
		etiquetaPaso2,
		filaResolucion,
		widget.NewSeparator(),
		puzzleApp.progressBar,
	)

	// Panel de información completo
	infoPanel := container.NewVBox(
		etiquetaInfo,
		infoScroll,
	)

	// Layout principal vertical
	mainContent := container.NewVBox(
		puzzleApp.estadoLabel,
		puzzleCardContainer,
		widget.NewSeparator(),
		controles,
		widget.NewSeparator(),
		infoPanel,
	)

	// Ensamblaje final de la interfaz
	content := container.NewVBox(
		header,
		mainContent,
	)

	ventana.SetContent(content)

	// Inicializar en estado ordenado y comenzar loop de eventos
	puzzleApp.iniciar()
	ventana.ShowAndRun()
}
