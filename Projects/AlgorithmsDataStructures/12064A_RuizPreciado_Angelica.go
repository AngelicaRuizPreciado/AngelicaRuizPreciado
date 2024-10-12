package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//Define mandatory types

type piano struct {
	punti map[punto]*piastrella //ricerca O(1), ora è un grafo non orientato con lista di adiacenza
	R     *regole
}

//Define non mandatory types

type punto struct {
	x int
	y int
}

type piastrella struct {
	posizione    punto
	colore       string
	intensità    int
	circonvicine *[]*piastrella
}

type regolaType struct {
	condizioni []condizione
	beta       string
	consumo    int
}

type condizione struct {
	k      int
	colore string
}

type regole []*regolaType

//Define mandatory functions

func colora(p piano, x int, y int, alpha string, i int) {
	/*colora la piastrella di coordinate (x,y) con il colore alpha e l’intensità i.*/
	if p.punti[punto{x, y}] == nil {
		p.punti[punto{x, y}] = &piastrella{posizione: punto{x, y}, colore: alpha, intensità: i, circonvicine: &[]*piastrella{}}
		//aggiungere circonvicine
		p.aggiungiCirconvicine(x, y)
	} else {
		ricolora := p.punti[punto{x, y}]
		ricolora.colore = alpha
		ricolora.intensità = i
	}

}

func spegni(p piano, x int, y int) {
	/*spenge la piastrella di coordinate (x,y).*/
	if p.punti[punto{x, y}] == nil {
		return
	}
	for _, c := range *p.punti[punto{x, y}].circonvicine {
		//rimuove la piastrella da tutte le circonvicine
		for i, circonvicina := range *c.circonvicine {
			if circonvicina == p.punti[punto{x, y}] {
				*c.circonvicine = append((*c.circonvicine)[:i], (*c.circonvicine)[i+1:]...)
				break
			}
		}
	}
	delete(p.punti, punto{x, y})
}

func regola(p piano, r string) {
	/*aggiunge la regola r al sistema rappresentato da p.
	Se la regola ha sum(k) > 8 fa return senza aggiungere la regola*/

	var reg []string = strings.Split(r, " ")
	var regola regolaType = regolaType{
		condizioni: []condizione{},
		beta:       reg[0],
		consumo:    0,
	}
	var count int = 0
	for i := 1; i < len(reg); i++ {
		var k, err = strconv.Atoi(reg[i])
		if err == nil {
			count += k
			if count > 8 {
				fmt.Println("Invalid rule, sum of k's must be <= 8")
				return
			} else {
				regola.condizioni = append(regola.condizioni, condizione{k: k, colore: reg[i+1]})
			}
		}
	}
	*p.R = append(*p.R, &regola)
}

func stato(p piano, x int, y int) (string, int) {
	/*restituisce il colore e l’intensità della piastrella di coordinate (x,y).*/
	if val, ok := p.punti[punto{x, y}]; ok {
		fmt.Println(val.colore, val.intensità)
		return val.colore, val.intensità
	}
	return "", 0
}

func stampa(p piano) {
	/*stampa le regole contenute in p nell'ordine attuale. O(8n) con n regole*/
	fmt.Println("(")
	for _, regola := range *p.R {
		var rule regolaType = *regola
		fmt.Printf("%s: ", regola.beta)
		for i, condizione := range rule.condizioni {
			fmt.Printf("%d %s", condizione.k, condizione.colore)
			if i != len(rule.condizioni)-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Println(")")
}

// Define non mandatory functions
func (p *piano) aggiungiCirconvicine(x, y int) {
	//funzione strumentale per colora, aggiunge le circonvicine alla piastrella di coordinate (x,y)

	puntoCentrale := punto{x, y}
	piastrellaCentrale, _ := p.punti[puntoCentrale]

	for _, dx := range []int{-1, 0, 1} {
		for _, dy := range []int{-1, 0, 1} {
			if dx == 0 && dy == 0 {
				continue
			}
			puntoVicino := punto{x + dx, y + dy}
			piastrellaVicino, ok := p.punti[puntoVicino]
			if ok {
				*piastrellaCentrale.circonvicine = append(*piastrellaCentrale.circonvicine, piastrellaVicino)
				*piastrellaVicino.circonvicine = append(*piastrellaVicino.circonvicine, piastrellaCentrale)
			}
		}
	}
}

func blocco(p piano, x int, y int) {
	//calcola e stampa la somma delle intensità delle piastrelle del blocco della piastrella di coordinate (x,y).
	fmt.Println(visitaIntensita(p, x, y, false))
}

func bloccoOmog(p piano, x int, y int) {
	//calcola e stampa la somma delle intensità delle piastrelle del blocco omogeneo della piastrella di coordinate (x,y).
	//blocco omogeneo = blocco con tutte le piastrelle dello stesso colore
	fmt.Println(visitaIntensita(p, x, y, true))
}

func visitaIntensita(p piano, x int, y int, omogeneo bool) int {
	if p.punti[punto{x, y}] == nil {
		return 0
	}
	var sum int = 0
	var start piastrella = *p.punti[punto{x, y}]
	var Fr, Ext []piastrella
	Fr = append(Fr, start)
	for len(Fr) > 0 {
		var v = Fr[0]
		Fr = Fr[1:]
		Ext = append(Ext, v)
		sum += v.intensità

		for _, q := range *v.circonvicine {
			if !contains(Ext, *q) && !contains(Fr, *q) {
				if omogeneo && q.colore != start.colore {
					continue
				}
				Fr = append(Fr, *q)
			}
		}
	}
	return sum
}
func contains(p []piastrella, e piastrella) bool {
	for _, a := range p {
		if a == e {
			return true
		}
	}
	return false
}
func containsPointer(p []*piastrella, e *piastrella) bool {
	for _, a := range p {
		if a == e {
			return true
		}
	}
	return false
}

func propaga(p piano, x int, y int) {
	//applica alla piastrella (x,y) la prima regola di propagazione applicabile
	var accesa bool = true
	if p.punti[punto{x, y}] == nil { //se la piastrella è spenta bisogna accenderla prima, se non si applica niente bisogna spegnere
		accesa = false
		colora(p, x, y, "bianco", 1)
	}
	for _, regola := range *p.R { //O(n) con n regole
		var applicabile = VerificaRegola(p, *p.punti[punto{x, y}], regola)

		if applicabile {
			colora(p, x, y, regola.beta, p.punti[punto{x, y}].intensità)
			regola.consumo++
			accesa = true
			break
		}
	}
	if !accesa {
		spegni(p, x, y)
	}
}
func VerificaRegola(p piano, v piastrella, regola *regolaType) bool {
	// Verifica se la regola è applicabile
	var applicabile bool = true
	var vicini map[string]int = make(map[string]int)
	for _, circonvicina := range *v.circonvicine { //O(8) con massimo 8 circonvicine
		vicini[circonvicina.colore]++
	}

	for _, condizione := range regola.condizioni { //O(m) con m <= 8 condizioni
		// Questa parte Verifica se la condizione è soddisfatta
		if vicini[condizione.colore] < condizione.k {
			applicabile = false
			break
		}
		// Qui finisce la verifica
	}
	return applicabile
}

func propagaBlocco(p piano, x int, y int) {
	//propaga al blocco della piastrella di coordinate (x,y) la prima regola di propagazione applicabile
	if p.punti[punto{x, y}] == nil {
		return //se la piastrella è spenta non fa parte di nessun blocco
	}
	var start *piastrella = p.punti[punto{x, y}]
	var Fr, Ext []*piastrella
	Fr = append(Fr, start)
	var daApplicare map[*piastrella]*regolaType = make(map[*piastrella]*regolaType)
	for len(Fr) > 0 { //O(n) con n piastrelle (nel blocco)
		var v = Fr[0]
		Fr = Fr[1:]
		Ext = append(Ext, v)

		// sceglie la regola da applicare
		for _, regola := range *p.R { //O(n) con n regole
			// Se la regola è applicabile, la aggiunge alla lista delle regole da applicare
			if VerificaRegola(p, *v, regola) { //O(1)
				daApplicare[v] = regola
				break
			}
		}

		// Aggiunge le piastrelle vicine a Fr
		for _, q := range *v.circonvicine {
			// Se la piastrella q non è stata già visitata,
			// la aggiunge alla lista delle piastrelle da visitare
			if !containsPointer(Ext, q) && !containsPointer(Fr, q) {
				Fr = append(Fr, q)
			}
		}
	}
	// fine del ciclo, ora applica le regole
	for piast, reg := range daApplicare {
		piast.colore = reg.beta
		reg.consumo++
	}
}

func ordina(p piano) {
	//ordina le regole contenute in p in modo che quella di maggior consumo sia l'ultima (sort stabile) O(nlogn)
	var regs = *p.R
	var newRules = mergeSort(regs)
	*p.R = newRules
}
func mergeSort(arr []*regolaType) regole {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2
	left := arr[:mid]
	right := arr[mid:]

	left = mergeSort(left)
	right = mergeSort(right)

	return merge(left, right)
}

func merge(left, right []*regolaType) []*regolaType {
	result := make([]*regolaType, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i].consumo <= right[j].consumo {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}

func pista(p piano, x1, y1 int, s []string) {
	//Stampa la pista che parte da Piastrella(x, y) e segue la sequenza di direzioni s, se tale pista `e definita. Altrimenti non stampa nulla.
	if p.punti[punto{x1, y1}] == nil {
		return
	}
	var direzioni = s
	var pista []punto = []punto{punto{x1, y1}}
	var x, y int = x1, y1
	for _, direzione := range direzioni {
		switch direzione {
		case "NO":
			x--
			y++
		case "NN":
			y++
		case "NE":
			x++
			y++
		case "EE":
			x++
		case "SE":
			x++
			y--
		case "SS":
			y--
		case "SO":
			x--
			y--
		case "OO":
			x--
		//considero l'input sempre corretto, non serve un default
		}
		if p.punti[punto{x, y}] != nil {
			pista = append(pista, punto{x, y})
		} else {
			return
		}

	}
	fmt.Println("[")
	for _, punto := range pista {
		fmt.Print(punto.x, " ", punto.y, " ", p.punti[punto].colore, " ", p.punti[punto].intensità)
		fmt.Println()
	}
	fmt.Println("]")
}

func lung(p piano, x1, y1, x2, y2 int) int {
	//Determina la lunghezza della pista pi`u breve che parte da Piastrella(x1, y1) e arriva in Piastrella(x2, y2). Altrimenti non stampa nulla.
	if p.punti[punto{x1, y1}] == nil || p.punti[punto{x2, y2}] == nil {
		return -1
	}

	dist := visitaPistaBreve(p, x1, y1, x2, y2)
	if dist != -1 {
		fmt.Println(dist+1)
		return dist+1 //+1 per contare anche la piastrella di partenza
	}
	return dist
}
func visitaPistaBreve(p piano, x1, y1, x2, y2 int) int {
	if x1 == x2 && y1 == y2 {
		return 0
	}
	//crea le variabili di supporto
	visited := make(map[punto]bool)
	distance := make(map[punto]int)
	queue := []punto{{x1, y1}}

	//inizializza le distanze a 0 e le piastrelle come non visitate
	for pt := range p.punti {
		distance[pt] = -1
		visited[pt] = false
	}
	distance[punto{x1, y1}] = 0

	// BFS algorithm
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:] //dequeue

		if current.x == x2 && current.y == y2{ //trovata la seconda piastrella
			return distance[current]
		}

		for _, neighbor := range *p.punti[current].circonvicine {
			if !visited[neighbor.posizione] {
				visited[neighbor.posizione] = true
				distance[neighbor.posizione] = distance[current] + 1
				queue = append(queue, neighbor.posizione) //enqueue
			}
		}
	}

	return distance[punto{x2,y2}] //se non trova la piastrella di arrivo ritorna -1
}

// //

func checkInputLength(op []string, length int) bool {
	if len(op) != length {
		fmt.Println("Invalid number of arguments")
		return false
	}
	return true
}

func esegui(p piano, s string) {
	/*applica al sistema rappresentato da p l’operazione associata dalla stringa s, secondo quanto specificato nella Tabella 1.*/

	var op []string = strings.Split(s, " ")

	switch op[0] {
	case "C":
		// handle colora(x,y,alpha) operation
		if checkInputLength(op, 5) {
			var x, _ = strconv.Atoi(op[1])
			var y, _ = strconv.Atoi(op[2])
			var i, _ = strconv.Atoi(op[4])
			colora(p, x, y, op[3], i)
		}
	case "S":
		// handle spegni(x,y) operation
		if checkInputLength(op, 3) {
			var x, _ = strconv.Atoi(op[1])
			var y, _ = strconv.Atoi(op[2])
			spegni(p, x, y)
		}

	case "r":
		// handle regola(...) operation
		regola(p, s[2:])
	case "?":
		// handle stato(x,y) operation
		if checkInputLength(op, 3) {
			var x, _ = strconv.Atoi(op[1])
			var y, _ = strconv.Atoi(op[2])
			stato(p, x, y)
		}

	case "s":
		// handle stampa operation
		stampa(p)
	case "b":
		// handle blocco(x,y) operation
		if checkInputLength(op, 3) {
			var x, _ = strconv.Atoi(op[1])
			var y, _ = strconv.Atoi(op[2])
			blocco(p, x, y)
		}

	case "B":
		// handle bloccoOmog(x,y) operation
		if checkInputLength(op, 3) {
			var x, _ = strconv.Atoi(op[1])
			var y, _ = strconv.Atoi(op[2])
			bloccoOmog(p, x, y)
		}
	case "p":
		// handle propaga(x,y) operation
		if checkInputLength(op, 3) {
			var x, _ = strconv.Atoi(op[1])
			var y, _ = strconv.Atoi(op[2])
			propaga(p, x, y)
		}
	case "P":
		// handle propagaBlocco(x,y) operation
		if checkInputLength(op, 3) {
			var x, _ = strconv.Atoi(op[1])
			var y, _ = strconv.Atoi(op[2])
			propagaBlocco(p, x, y)
		}
	case "o":
		// handle ordina operation
		ordina(p)
	case "t":
		// handle pista(x,y,s) operation
		var x, _ = strconv.Atoi(op[1])
		var y, _ = strconv.Atoi(op[2])
		var dir = strings.Split(op[3], ",")
		pista(p, x, y, dir)
	case "L":
		// handle lung(x1,y1,x2,y2) operation
		if checkInputLength(op, 5) {
			var x1, _ = strconv.Atoi(op[1])
			var y1, _ = strconv.Atoi(op[2])
			var x2, _ = strconv.Atoi(op[3])
			var y2, _ = strconv.Atoi(op[4])
			lung(p, x1, y1, x2, y2)
		}
	default:
		fmt.Println(op[0], "Invalid operation")
	}

}

func main() {
	var regole regole = make(regole, 0)

	var p piano = piano{
		punti: make(map[punto]*piastrella),
		R:     &regole,
	}

	//consider adding an explanation of the input format and the menu options

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "q" {
			break
		}
		esegui(p, input)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}

}
