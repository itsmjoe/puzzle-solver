package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Estado representa un estado del puzzle 8
type Estado struct {
	tablero    [9]int
	posVacio   int
	g          int // costo desde el inicio
	h          int // heurística
	f          int // g + h
	padre      *Estado
	movimiento string
}

// PriorityQueue implementa una cola de prioridad para A*
type PriorityQueue []*Estado

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].f < pq[j].f }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Estado))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// Aplicación principal
type PuzzleApp struct {
	window          fyne.Window
	tablero         [9]*widget.Button
	estadoActual    [9]int
	estadoFinal     [9]int
	solucion        []*Estado
	pasoActual      int
	infoLabel       *widget.Label
	algoritmoSelect *widget.Select
}

// Inicializar la aplicación
func NewPuzzleApp() *PuzzleApp {
	app := &PuzzleApp{
		estadoFinal: [9]int{1, 2, 3, 4, 5, 6, 7, 8, 0},
		pasoActual:  0,
	}

	// Estado inicial ordenado
	app.estadoActual = [9]int{1, 2, 3, 4, 5, 6, 7, 8, 0}

	return app
}

// Actualizar la visualización del tablero
func (app *PuzzleApp) actualizarTablero() {
	for i := 0; i < 9; i++ {
		if app.estadoActual[i] == 0 {
			app.tablero[i].SetText("")
			app.tablero[i].Importance = widget.LowImportance
		} else {
			app.tablero[i].SetText(strconv.Itoa(app.estadoActual[i]))
			app.tablero[i].Importance = widget.MediumImportance
		}
	}
	app.tablero[0].Refresh()
}

// Iniciar el puzzle con estado ordenado
func (app *PuzzleApp) iniciarPuzzle() {
	app.estadoActual = [9]int{1, 2, 3, 4, 5, 6, 7, 8, 0}
	app.solucion = nil
	app.pasoActual = 0
	app.actualizarTablero()
	app.infoLabel.SetText("Puzzle iniciado - Estado ordenado")
}

// Desordenar el tablero
func (app *PuzzleApp) desordenarTablero() {
	rand.Seed(time.Now().UnixNano())

	// Generar un estado aleatorio solucionable
	for {
		estado := [9]int{0, 1, 2, 3, 4, 5, 6, 7, 8}

		// Mezclar el array
		for i := len(estado) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			estado[i], estado[j] = estado[j], estado[i]
		}

		if app.esSolucionable(estado) {
			app.estadoActual = estado
			break
		}
	}

	app.solucion = nil
	app.pasoActual = 0
	app.actualizarTablero()
	app.infoLabel.SetText("Tablero desordenado - Listo para resolver")
}

// Verificar si un estado es solucionable
func (app *PuzzleApp) esSolucionable(estado [9]int) bool {
	inversiones := 0
	for i := 0; i < 8; i++ {
		if estado[i] == 0 {
			continue
		}
		for j := i + 1; j < 9; j++ {
			if estado[j] == 0 {
				continue
			}
			if estado[i] > estado[j] {
				inversiones++
			}
		}
	}
	return inversiones%2 == 0
}

// Resolver el puzzle
func (app *PuzzleApp) resolverPuzzle() {
	algoritmo := app.algoritmoSelect.Selected
	app.infoLabel.SetText("Resolviendo con " + algoritmo + "...")

	estadoInicial := &Estado{
		tablero:  app.estadoActual,
		posVacio: app.encontrarPosVacio(app.estadoActual),
		g:        0,
	}

	var solucion []*Estado
	var nodos int

	start := time.Now()

	switch algoritmo {
	case "A* (Manhattan)":
		solucion, nodos = app.aEstrellaManhattan(estadoInicial)
	case "A* (Euclidiana)":
		solucion, nodos = app.aEstrellaEuclidiana(estadoInicial)
	case "Búsqueda en Anchura":
		solucion, nodos = app.busquedaEnAnchura(estadoInicial)
	}

	duracion := time.Since(start)

	if solucion != nil {
		app.solucion = solucion
		app.pasoActual = 0
		app.infoLabel.SetText(fmt.Sprintf("¡Resuelto! Pasos: %d, Nodos: %d, Tiempo: %v ms",
			len(solucion)-1, nodos, duracion.Milliseconds()))
	} else {
		app.infoLabel.SetText("No se encontró solución")
	}
}

// Siguiente paso en la solución
func (app *PuzzleApp) siguientePaso() {
	if app.solucion == nil {
		app.infoLabel.SetText("Primero resuelve el puzzle")
		return
	}

	if app.pasoActual < len(app.solucion) {
		app.estadoActual = app.solucion[app.pasoActual].tablero
		app.actualizarTablero()

		movimiento := ""
		if app.pasoActual > 0 {
			movimiento = app.solucion[app.pasoActual].movimiento
		}

		app.infoLabel.SetText(fmt.Sprintf("Paso %d/%d - Movimiento: %s",
			app.pasoActual+1, len(app.solucion), movimiento))
		app.pasoActual++
	} else {
		app.infoLabel.SetText("¡Solución completada!")
	}
}

// Reiniciar al estado final
func (app *PuzzleApp) reiniciarPuzzle() {
	app.estadoActual = app.estadoFinal
	app.solucion = nil
	app.pasoActual = 0
	app.actualizarTablero()
	app.infoLabel.SetText("Puzzle reiniciado al estado final")
}

// Encontrar posición del espacio vacío
func (app *PuzzleApp) encontrarPosVacio(tablero [9]int) int {
	for i := 0; i < 9; i++ {
		if tablero[i] == 0 {
			return i
		}
	}
	return -1
}

// Heurística de Manhattan
func (app *PuzzleApp) heuristicaManhattan(tablero [9]int) int {
	distancia := 0
	for i := 0; i < 9; i++ {
		if tablero[i] != 0 {
			valor := tablero[i]
			posObjetivo := valor - 1

			filaActual := i / 3
			colActual := i % 3
			filaObjetivo := posObjetivo / 3
			colObjetivo := posObjetivo % 3

			distancia += int(math.Abs(float64(filaActual-filaObjetivo))) +
				int(math.Abs(float64(colActual-colObjetivo)))
		}
	}
	return distancia
}

// Heurística Euclidiana
func (app *PuzzleApp) heuristicaEuclidiana(tablero [9]int) int {
	distancia := 0.0
	for i := 0; i < 9; i++ {
		if tablero[i] != 0 {
			valor := tablero[i]
			posObjetivo := valor - 1

			filaActual := float64(i / 3)
			colActual := float64(i % 3)
			filaObjetivo := float64(posObjetivo / 3)
			colObjetivo := float64(posObjetivo % 3)

			distancia += math.Sqrt(math.Pow(filaActual-filaObjetivo, 2) +
				math.Pow(colActual-colObjetivo, 2))
		}
	}
	return int(distancia)
}

// Generar sucesores de un estado
func (app *PuzzleApp) generarSucesores(estado *Estado) []*Estado {
	sucesores := []*Estado{}
	posVacio := estado.posVacio

	movimientos := []struct {
		dx, dy int
		nombre string
	}{
		{-1, 0, "Arriba"},
		{1, 0, "Abajo"},
		{0, -1, "Izquierda"},
		{0, 1, "Derecha"},
	}

	fila := posVacio / 3
	col := posVacio % 3

	for _, mov := range movimientos {
		nuevaFila := fila + mov.dx
		nuevaCol := col + mov.dy

		if nuevaFila >= 0 && nuevaFila < 3 && nuevaCol >= 0 && nuevaCol < 3 {
			nuevaPos := nuevaFila*3 + nuevaCol

			nuevoTablero := estado.tablero
			nuevoTablero[posVacio], nuevoTablero[nuevaPos] = nuevoTablero[nuevaPos], nuevoTablero[posVacio]

			sucesor := &Estado{
				tablero:    nuevoTablero,
				posVacio:   nuevaPos,
				g:          estado.g + 1,
				padre:      estado,
				movimiento: mov.nombre,
			}

			sucesores = append(sucesores, sucesor)
		}
	}

	return sucesores
}

// Algoritmo A* con heurística Manhattan
func (app *PuzzleApp) aEstrellaManhattan(estadoInicial *Estado) ([]*Estado, int) {
	estadoInicial.h = app.heuristicaManhattan(estadoInicial.tablero)
	estadoInicial.f = estadoInicial.g + estadoInicial.h

	abiertos := &PriorityQueue{estadoInicial}
	heap.Init(abiertos)

	cerrados := make(map[string]bool)
	nodosExplorados := 0

	for abiertos.Len() > 0 {
		actual := heap.Pop(abiertos).(*Estado)
		nodosExplorados++

		if actual.tablero == app.estadoFinal {
			return app.reconstruirCamino(actual), nodosExplorados
		}

		clave := fmt.Sprintf("%v", actual.tablero)
		if cerrados[clave] {
			continue
		}
		cerrados[clave] = true

		for _, sucesor := range app.generarSucesores(actual) {
			claveSucesor := fmt.Sprintf("%v", sucesor.tablero)
			if !cerrados[claveSucesor] {
				sucesor.h = app.heuristicaManhattan(sucesor.tablero)
				sucesor.f = sucesor.g + sucesor.h
				heap.Push(abiertos, sucesor)
			}
		}
	}

	return nil, nodosExplorados
}

// Algoritmo A* con heurística Euclidiana
func (app *PuzzleApp) aEstrellaEuclidiana(estadoInicial *Estado) ([]*Estado, int) {
	estadoInicial.h = app.heuristicaEuclidiana(estadoInicial.tablero)
	estadoInicial.f = estadoInicial.g + estadoInicial.h

	abiertos := &PriorityQueue{estadoInicial}
	heap.Init(abiertos)

	cerrados := make(map[string]bool)
	nodosExplorados := 0

	for abiertos.Len() > 0 {
		actual := heap.Pop(abiertos).(*Estado)
		nodosExplorados++

		if actual.tablero == app.estadoFinal {
			return app.reconstruirCamino(actual), nodosExplorados
		}

		clave := fmt.Sprintf("%v", actual.tablero)
		if cerrados[clave] {
			continue
		}
		cerrados[clave] = true

		for _, sucesor := range app.generarSucesores(actual) {
			claveSucesor := fmt.Sprintf("%v", sucesor.tablero)
			if !cerrados[claveSucesor] {
				sucesor.h = app.heuristicaEuclidiana(sucesor.tablero)
				sucesor.f = sucesor.g + sucesor.h
				heap.Push(abiertos, sucesor)
			}
		}
	}

	return nil, nodosExplorados
}

// Búsqueda en anchura
func (app *PuzzleApp) busquedaEnAnchura(estadoInicial *Estado) ([]*Estado, int) {
	cola := []*Estado{estadoInicial}
	visitados := make(map[string]bool)
	nodosExplorados := 0

	for len(cola) > 0 {
		actual := cola[0]
		cola = cola[1:]
		nodosExplorados++

		if actual.tablero == app.estadoFinal {
			return app.reconstruirCamino(actual), nodosExplorados
		}

		clave := fmt.Sprintf("%v", actual.tablero)
		if visitados[clave] {
			continue
		}
		visitados[clave] = true

		for _, sucesor := range app.generarSucesores(actual) {
			claveSucesor := fmt.Sprintf("%v", sucesor.tablero)
			if !visitados[claveSucesor] {
				cola = append(cola, sucesor)
			}
		}
	}

	return nil, nodosExplorados
}

// Reconstruir el camino de la solución
func (app *PuzzleApp) reconstruirCamino(estadoFinal *Estado) []*Estado {
	camino := []*Estado{}
	actual := estadoFinal

	for actual != nil {
		camino = append([]*Estado{actual}, camino...)
		actual = actual.padre
	}

	return camino
}

func main() {
	myApp := app.New()
	puzzleApp := NewPuzzleApp()
	puzzleApp.window = myApp.NewWindow("8-Puzzle Solver")
	puzzleApp.window.Resize(fyne.NewSize(500, 600))

	// Crear botones del tablero
	tableroContainer := container.NewGridWithColumns(3)
	for i := 0; i < 9; i++ {
		btn := widget.NewButton("", nil)
		btn.Resize(fyne.NewSize(80, 80))
		puzzleApp.tablero[i] = btn
		tableroContainer.Add(btn)
	}

	// Selector de algoritmo
	puzzleApp.algoritmoSelect = widget.NewSelect(
		[]string{"A* (Manhattan)", "A* (Euclidiana)", "Búsqueda en Anchura"},
		nil,
	)
	puzzleApp.algoritmoSelect.SetSelected("A* (Manhattan)")

	// Botones de control
	btnIniciar := widget.NewButton("Iniciar Puzzle", puzzleApp.iniciarPuzzle)
	btnDesordenar := widget.NewButton("Desordenar", puzzleApp.desordenarTablero)
	btnResolver := widget.NewButton("Resolver", puzzleApp.resolverPuzzle)
	btnPasoAPaso := widget.NewButton("Siguiente Paso", puzzleApp.siguientePaso)
	btnReiniciar := widget.NewButton("Reiniciar", puzzleApp.reiniciarPuzzle)

	// Label de información
	puzzleApp.infoLabel = widget.NewLabel("8-Puzzle Solver - Selecciona una opción")

	// Layout principal
	controles := container.NewVBox(
		widget.NewLabel("Algoritmo:"),
		puzzleApp.algoritmoSelect,
		container.NewGridWithColumns(2, btnIniciar, btnDesordenar),
		container.NewGridWithColumns(2, btnResolver, btnPasoAPaso),
		btnReiniciar,
		puzzleApp.infoLabel,
	)

	content := container.NewBorder(nil, controles, nil, nil, tableroContainer)
	puzzleApp.window.SetContent(content)

	puzzleApp.actualizarTablero()
	puzzleApp.window.ShowAndRun()
}
