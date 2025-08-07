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

// Tema personalizado para la aplicación
type MyTheme struct{}

func (t MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return color.NRGBA{41, 98, 255, 255} // Azul profesional
	case theme.ColorNameButton:
		return color.NRGBA{248, 249, 250, 255} // Blanco suave
	case theme.ColorNameForeground:
		return color.NRGBA{33, 37, 41, 255} // Texto oscuro
	case theme.ColorNameBackground:
		return color.NRGBA{255, 255, 255, 255} // Fondo blanco
	case theme.ColorNameInputBackground:
		return color.NRGBA{255, 255, 255, 255} // Fondo blanco para inputs
	case theme.ColorNameSelection:
		return color.NRGBA{41, 98, 255, 40} // Azul transparente para selección
	case theme.ColorNameHover:
		return color.NRGBA{41, 98, 255, 20} // Azul muy transparente para hover
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t MyTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 24
	case theme.SizeNameSubHeadingText:
		return 18
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// Estado simple del puzzle
type Estado struct {
	tablero [9]int
	padre   *Estado
	costo   int
	accion  string
}

// Widget personalizado para los botones del puzzle con animación
type PuzzleButton struct {
	widget.Button
	numero    int
	esVacio   bool
	destacado bool
}

func NewPuzzleButton(numero int) *PuzzleButton {
	btn := &PuzzleButton{numero: numero, destacado: false}
	btn.ExtendBaseWidget(btn)
	btn.esVacio = (numero == 0)
	btn.actualizarEstilo()
	return btn
}

func (pb *PuzzleButton) actualizarEstilo() {
	if pb.esVacio {
		pb.SetText("")                       // Completamente vacío, sin números
		pb.Importance = widget.LowImportance // Estilo tenue para el espacio vacío
	} else {
		pb.SetText(strconv.Itoa(pb.numero))
		if pb.destacado {
			pb.Importance = widget.WarningImportance // Naranja suave para movimiento
		} else {
			pb.Importance = widget.HighImportance
		}
	}
	pb.Refresh()
}

func (pb *PuzzleButton) setNumero(numero int) {
	pb.numero = numero
	pb.esVacio = (numero == 0)
	pb.destacado = false
	pb.actualizarEstilo()
}

func (pb *PuzzleButton) destacar() {
	if !pb.esVacio && pb.numero != 0 {
		// Solo destacar si no es espacio vacío y no es 0
		pb.SetText(strconv.Itoa(pb.numero))
		pb.destacado = true
		pb.Importance = widget.WarningImportance
		pb.Refresh()

		// Quitar destaque después de 800ms usando una goroutine simple
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

// Aplicación del puzzle
type PuzzleApp struct {
	window       fyne.Window
	botones      [9]*PuzzleButton
	estadoActual [9]int
	objetivo     [9]int
	solucion     []Estado
	paso         int
	infoLabel    *widget.RichText
	estadoLabel  *widget.Label
	algoritmo    *widget.Select
	progressBar  *widget.ProgressBar
}

// Crear nueva aplicación
func NuevaPuzzleApp() *PuzzleApp {
	return &PuzzleApp{
		objetivo: [9]int{1, 2, 3, 4, 5, 6, 7, 8, 0},
		paso:     0,
	}
}

// Encontrar posición del 0 (vacío)
func encontrarVacio(tablero [9]int) int {
	for i := 0; i < 9; i++ {
		if tablero[i] == 0 {
			return i
		}
	}
	return -1
}

// Verificar si el estado es el objetivo
func esObjetivo(tablero [9]int, objetivo [9]int) bool {
	for i := 0; i < 9; i++ {
		if tablero[i] != objetivo[i] {
			return false
		}
	}
	return true
}

// Heurística Manhattan - básica para IA1
func heuristicaManhattan(tablero [9]int) int {
	distancia := 0
	for i := 0; i < 9; i++ {
		if tablero[i] != 0 {
			// Posición actual
			fila_actual := i / 3
			col_actual := i % 3

			// Posición objetivo (valor-1 porque empezamos en 1)
			valor := tablero[i]
			fila_objetivo := (valor - 1) / 3
			col_objetivo := (valor - 1) % 3

			// Distancia Manhattan
			distancia += abs(fila_actual-fila_objetivo) + abs(col_actual-col_objetivo)
		}
	}
	return distancia
}

// Función auxiliar para valor absoluto
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Generar movimientos posibles
func generarMovimientos(tablero [9]int) []Estado {
	movimientos := []Estado{}
	posVacio := encontrarVacio(tablero)

	fila := posVacio / 3
	col := posVacio % 3

	// Arriba
	if fila > 0 {
		nuevo := tablero
		nueva_pos := (fila-1)*3 + col
		nuevo[posVacio], nuevo[nueva_pos] = nuevo[nueva_pos], nuevo[posVacio]
		movimientos = append(movimientos, Estado{
			tablero: nuevo,
			accion:  "Arriba",
		})
	}

	// Abajo
	if fila < 2 {
		nuevo := tablero
		nueva_pos := (fila+1)*3 + col
		nuevo[posVacio], nuevo[nueva_pos] = nuevo[nueva_pos], nuevo[posVacio]
		movimientos = append(movimientos, Estado{
			tablero: nuevo,
			accion:  "Abajo",
		})
	}

	// Izquierda
	if col > 0 {
		nuevo := tablero
		nueva_pos := fila*3 + (col - 1)
		nuevo[posVacio], nuevo[nueva_pos] = nuevo[nueva_pos], nuevo[posVacio]
		movimientos = append(movimientos, Estado{
			tablero: nuevo,
			accion:  "Izquierda",
		})
	}

	// Derecha
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

// Búsqueda A* simple (versión para estudiantes de IA1)
func busquedaAEstrella(inicial [9]int, objetivo [9]int) []Estado {
	// Lista abierta (estados por explorar)
	abierta := []Estado{{tablero: inicial, costo: 0}}
	// Lista cerrada (estados ya explorados)
	cerrada := []string{}

	for len(abierta) > 0 {
		// Encontrar el estado con menor f = g + h
		indice_mejor := 0
		mejor_f := abierta[0].costo + heuristicaManhattan(abierta[0].tablero)

		for i := 1; i < len(abierta); i++ {
			f := abierta[i].costo + heuristicaManhattan(abierta[i].tablero)
			if f < mejor_f {
				mejor_f = f
				indice_mejor = i
			}
		}

		// Tomar el mejor estado
		actual := abierta[indice_mejor]
		abierta = append(abierta[:indice_mejor], abierta[indice_mejor+1:]...)

		// Si es el objetivo, reconstruir camino
		if esObjetivo(actual.tablero, objetivo) {
			camino := []Estado{}
			estado := &actual
			for estado != nil {
				camino = append([]Estado{*estado}, camino...)
				estado = estado.padre
			}
			return camino
		}

		// Agregar a cerrada
		cerrada = append(cerrada, fmt.Sprintf("%v", actual.tablero))

		// Generar sucesores
		for _, movimiento := range generarMovimientos(actual.tablero) {
			estado_str := fmt.Sprintf("%v", movimiento.tablero)

			// Verificar si ya está en cerrada
			ya_explorado := false
			for _, cerrado := range cerrada {
				if cerrado == estado_str {
					ya_explorado = true
					break
				}
			}

			if !ya_explorado {
				movimiento.padre = &actual
				movimiento.costo = actual.costo + 1
				abierta = append(abierta, movimiento)
			}
		}
	}

	return []Estado{} // No hay solución
}

// Búsqueda en anchura simple
func busquedaAnchura(inicial [9]int, objetivo [9]int) []Estado {
	cola := []Estado{{tablero: inicial}}
	visitados := []string{}

	for len(cola) > 0 {
		actual := cola[0]
		cola = cola[1:]

		if esObjetivo(actual.tablero, objetivo) {
			// Reconstruir camino
			camino := []Estado{}
			estado := &actual
			for estado != nil {
				camino = append([]Estado{*estado}, camino...)
				estado = estado.padre
			}
			return camino
		}

		estado_str := fmt.Sprintf("%v", actual.tablero)
		visitados = append(visitados, estado_str)

		for _, movimiento := range generarMovimientos(actual.tablero) {
			mov_str := fmt.Sprintf("%v", movimiento.tablero)

			ya_visitado := false
			for _, v := range visitados {
				if v == mov_str {
					ya_visitado = true
					break
				}
			}

			if !ya_visitado {
				movimiento.padre = &actual
				cola = append(cola, movimiento)
			}
		}
	}

	return []Estado{}
}

// Actualizar interfaz con animación de movimiento
func (app *PuzzleApp) actualizarTablero() {
	for i := 0; i < 9; i++ {
		app.botones[i].setNumero(app.estadoActual[i])
	}
	app.actualizarEstado()
}

// Actualizar información de estado
func (app *PuzzleApp) actualizarEstado() {
	if esObjetivo(app.estadoActual, app.objetivo) {
		app.estadoLabel.SetText("ESTADO: RESUELTO")
		app.estadoLabel.Importance = widget.SuccessImportance
	} else {
		manhattan := heuristicaManhattan(app.estadoActual)
		app.estadoLabel.SetText(fmt.Sprintf("ESTADO: EN PROCESO | Heurística Manhattan: %d", manhattan))
		app.estadoLabel.Importance = widget.MediumImportance
	}
}

// Inicializar puzzle ordenado
func (app *PuzzleApp) iniciar() {
	app.estadoActual = [9]int{1, 2, 3, 4, 5, 6, 7, 8, 0}
	app.solucion = []Estado{}
	app.paso = 0
	app.progressBar.SetValue(0)
	app.actualizarTablero()

	app.infoLabel.ParseMarkdown("## SISTEMA INICIALIZADO\n\n**Estado:** Puzzle ordenado correctamente\n\n**Acción:** Presiona 'Mezclar' para comenzar")
}

// Mezclar puzzle
func (app *PuzzleApp) mezclar() {
	rand.Seed(time.Now().UnixNano())

	app.infoLabel.ParseMarkdown("## MEZCLANDO PUZZLE\n\n**Estado:** Generando configuración aleatoria...\n\n**Por favor espera**")

	// Mezclar haciendo movimientos aleatorios
	app.estadoActual = app.objetivo
	for i := 0; i < 150; i++ {
		movimientos := generarMovimientos(app.estadoActual)
		if len(movimientos) > 0 {
			mov := movimientos[rand.Intn(len(movimientos))]
			app.estadoActual = mov.tablero
		}
	}

	app.solucion = []Estado{}
	app.paso = 0
	app.progressBar.SetValue(0)
	app.actualizarTablero()

	manhattan := heuristicaManhattan(app.estadoActual)
	app.infoLabel.ParseMarkdown(fmt.Sprintf("## PUZZLE MEZCLADO\n\n**Estado:** Configuración aleatoria generada\n\n**Heurística Manhattan:** %d\n\n**Acción:** Selecciona algoritmo y presiona 'Resolver'", manhattan))
}

// Resolver puzzle
func (app *PuzzleApp) resolver() {
	algoritmo_seleccionado := app.algoritmo.Selected
	app.infoLabel.ParseMarkdown(fmt.Sprintf("## RESOLVIENDO PUZZLE\n\n**Algoritmo:** %s\n\n**Estado:** Buscando solución óptima...\n\n**Por favor espera**", algoritmo_seleccionado))

	inicio := time.Now()

	if algoritmo_seleccionado == "A* con Heurística Manhattan" {
		app.solucion = busquedaAEstrella(app.estadoActual, app.objetivo)
	} else {
		app.solucion = busquedaAnchura(app.estadoActual, app.objetivo)
	}

	duracion := time.Since(inicio)

	if len(app.solucion) > 0 {
		app.paso = 0
		app.progressBar.SetValue(0)

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
		app.infoLabel.ParseMarkdown("## ERROR\n\n**Estado:** No se encontró solución\n\n**Acción:** Intenta mezclar nuevamente")
	}
}

// Mostrar siguiente paso con animación simple y sin errores de threading
func (app *PuzzleApp) siguientePaso() {
	if len(app.solucion) == 0 {
		app.infoLabel.ParseMarkdown("## ADVERTENCIA\n\n**Estado:** No hay solución cargada\n\n**Acción:** Primero resuelve el puzzle")
		return
	}

	if app.paso < len(app.solucion) {
		estadoNuevo := app.solucion[app.paso].tablero

		progreso := float64(app.paso) / float64(len(app.solucion)-1)
		app.progressBar.SetValue(progreso)

		accion := "Estado inicial"
		if app.paso > 0 {
			accion = app.solucion[app.paso].accion

			// Encontrar y destacar la pieza que se mueve (sin errores de threading)
			for i := 0; i < 9; i++ {
				if app.estadoActual[i] != 0 && app.estadoActual[i] != estadoNuevo[i] {
					app.botones[i].destacar()
					break
				}
			}

			// Actualizar el tablero después de un pequeño delay
			go func(nuevoEstado [9]int) {
				time.Sleep(300 * time.Millisecond)
				fyne.DoAndWait(func() {
					app.estadoActual = nuevoEstado
					app.actualizarTablero()
				})
			}(estadoNuevo)
		} else {
			// Primer paso: actualización directa
			app.estadoActual = estadoNuevo
			app.actualizarTablero()
		}

		app.infoLabel.ParseMarkdown(fmt.Sprintf("## EJECUTANDO SOLUCIÓN\n\n**Paso:** %d de %d\n\n**Movimiento:** %s\n\n**Progreso:** %.1f%%",
			app.paso+1, len(app.solucion), accion, progreso*100))

		app.paso++
	} else {
		app.progressBar.SetValue(1.0)
		app.infoLabel.ParseMarkdown("## SOLUCIÓN COMPLETADA\n\n**Estado:** Puzzle resuelto exitosamente\n\n**Felicitaciones:** El algoritmo funcionó correctamente")
	}
}

func main() {
	// Crear aplicación sin tema personalizado para evitar problemas con selector
	myApp := app.New()
	// myApp.Settings().SetTheme(&MyTheme{}) // Comentado temporalmente

	ventana := myApp.NewWindow("8-Puzzle Solver - Inteligencia Artificial I | USAC")
	ventana.Resize(fyne.NewSize(700, 800))
	ventana.CenterOnScreen()

	puzzleApp := NuevaPuzzleApp()
	puzzleApp.window = ventana

	// Header elegante sin tema personalizado
	headerCard := canvas.NewRectangle(color.NRGBA{240, 248, 255, 255})
	headerCard.Resize(fyne.NewSize(700, 120))

	titulo := canvas.NewText("8-PUZZLE SOLVER", color.NRGBA{25, 118, 210, 255})
	titulo.Alignment = fyne.TextAlignCenter
	titulo.TextSize = 28
	titulo.TextStyle.Bold = true

	subtitulo := canvas.NewText("Inteligencia Artificial I • Universidad de San Carlos de Guatemala", color.NRGBA{84, 110, 122, 255})
	subtitulo.Alignment = fyne.TextAlignCenter
	subtitulo.TextSize = 14

	headerContent := container.NewVBox(
		container.NewPadded(titulo),
		subtitulo,
	)

	header := container.NewStack(headerCard, headerContent)

	// Cuadrícula elegante del puzzle con espacio vacío real
	cuadricula := container.NewGridWithColumns(3)
	for i := 0; i < 9; i++ {
		btn := NewPuzzleButton(1) // Inicializar con número temporal
		btn.Resize(fyne.NewSize(100, 100))

		// Configuración manual inicial para asegurar visibilidad
		if i < 8 {
			// Números del 1 al 8
			btn.SetText(strconv.Itoa(i + 1))
			btn.Importance = widget.HighImportance
		} else {
			// Espacio vacío (sin número 0)
			btn.SetText("") // Completamente vacío
			btn.Importance = widget.LowImportance
		}
		btn.Refresh()

		puzzleApp.botones[i] = btn

		// Container con sombra simulada
		btnContainer := container.NewPadded(btn)
		cuadricula.Add(btnContainer)
	}

	// Card para la cuadrícula
	puzzleCard := canvas.NewRectangle(color.NRGBA{250, 250, 250, 255})
	puzzleCardContainer := container.NewStack(puzzleCard, container.NewPadded(cuadricula))

	// Panel de estado
	puzzleApp.estadoLabel = widget.NewLabel("")
	puzzleApp.estadoLabel.Alignment = fyne.TextAlignCenter
	puzzleApp.estadoLabel.TextStyle.Bold = true

	// Barra de progreso
	puzzleApp.progressBar = widget.NewProgressBar()
	puzzleApp.progressBar.TextFormatter = func() string {
		return fmt.Sprintf("Progreso: %.0f%%", puzzleApp.progressBar.Value*100)
	}

	// Selector de algoritmo elegante
	etiquetaAlgoritmo := widget.NewLabel("ALGORITMO DE BÚSQUEDA")
	etiquetaAlgoritmo.TextStyle.Bold = true
	etiquetaAlgoritmo.Alignment = fyne.TextAlignCenter

	puzzleApp.algoritmo = widget.NewSelect(
		[]string{"A* con Heurística Manhattan", "Búsqueda en Anchura (BFS)"},
		nil,
	)
	puzzleApp.algoritmo.SetSelected("A* con Heurística Manhattan")

	// Botones principales con texto limpio
	btnIniciar := widget.NewButton("INICIAR", puzzleApp.iniciar)
	btnIniciar.Importance = widget.MediumImportance

	btnMezclar := widget.NewButton("MEZCLAR", puzzleApp.mezclar)
	btnMezclar.Importance = widget.HighImportance

	btnResolver := widget.NewButton("RESOLVER", puzzleApp.resolver)
	btnResolver.Importance = widget.SuccessImportance

	btnPaso := widget.NewButton("PASO A PASO", puzzleApp.siguientePaso)
	btnPaso.Importance = widget.MediumImportance

	// Panel de información con RichText
	etiquetaInfo := widget.NewLabel("INFORMACIÓN DEL SISTEMA")
	etiquetaInfo.TextStyle.Bold = true
	etiquetaInfo.Alignment = fyne.TextAlignCenter

	puzzleApp.infoLabel = widget.NewRichTextFromMarkdown("")

	// Scroll para el texto de información
	infoScroll := container.NewScroll(puzzleApp.infoLabel)
	infoScroll.SetMinSize(fyne.NewSize(400, 150))

	// Layout de controles
	controlesGrid := container.NewGridWithColumns(2,
		btnIniciar, btnMezclar,
		btnResolver, btnPaso,
	)

	controles := container.NewVBox(
		etiquetaAlgoritmo,
		puzzleApp.algoritmo,
		widget.NewSeparator(),
		controlesGrid,
		widget.NewSeparator(),
		puzzleApp.progressBar,
	)

	// Panel de información
	infoPanel := container.NewVBox(
		etiquetaInfo,
		infoScroll,
	)

	// Layout central
	mainContent := container.NewVBox(
		puzzleApp.estadoLabel,
		puzzleCardContainer,
		widget.NewSeparator(),
		controles,
		widget.NewSeparator(),
		infoPanel,
	)

	// Container principal con padding
	content := container.NewBorder(
		header,
		nil,
		nil,
		nil,
		container.NewPadded(mainContent),
	)

	ventana.SetContent(content)

	// Inicializar y mostrar
	puzzleApp.iniciar()
	ventana.ShowAndRun()
}
