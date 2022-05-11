package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

/*
ligne 62 congrate (print)
*/
type HangManData struct {
	Word             string   // Word composed of '_', ex: H_ll_
	ToFind           string   // Final word chosen by the program at the beginning. It is the word to find
	Attempts         int      // Number of attempts left
	HangmanPositions []string // It can be the array where the positions parsed in "hangman.txt" are stored
	index_tab_pendu  int      // ça va être là où on en est dans hangman positions.
}

func ToUpper(s string) string { // Ici on va créer une fonction toUpper qui va nous permettre de mettre les lettres minuscules dd'une string en majuscules.
	compteur := ""
	ch := []rune(s)
	for i := 0; i <= len(s)-1; i++ {
		if ch[i] >= 'a' && ch[i] <= 'z' {
			compteur = compteur + string(ch[i]-32)
		} else {
			compteur = compteur + string(ch[i])
		}
	}
	return compteur
}

func motatrouver() string { // la string qu'on va return sera un mot aléatoirement choisi dans words.txt
	var tabmots []string           // là où on va stocker les mots de words.txt
	f, err := os.Open("words.txt") // on va lire words.txt
	if err != nil {                // si erreur (words.txt vide ou autre)
		log.Fatal(err) // on va afficher l'erreur corespondante
	}

	scanner := bufio.NewScanner(f) // on va scanner le contenu de words.txt qui a été stocké dans f plus haut
	scanner.Split(bufio.ScanLines) // on va dire à notre scanner de scanner lignes par lignes (1 ligne = une place)
	for scanner.Scan() {           // tant que le scanner à du contenu à scanner
		tabmots = append(tabmots, (scanner.Text())) // on ajoute la ligne à  tabmots tant qu'il y'a des lignes dans words.txt
	}

	max := len(tabmots) // ici on définit les limites du tableau qui contient les mots pour en choisir un aléatoirement dedans
	min := 0
	rand.Seed(time.Now().UnixNano())
	rand := rand.Intn(max - min)  // On va stocker un int aléatoire choisi par la fonction rand.Intn situé entre les limites définies ici par max et min
	mot := ToUpper(tabmots[rand]) // Ici on va user de notre fonction ToUpper afin de mettre en majuscule le mot choisi dans tabmots par un int aléatoirement défini précédemment
	return mot                    // On retourne la variable mot dans laquelle est stocké notre mot choisi aléatoirement et mit ensuite en majuscules par ToUpper
}

func débutjeu(s string) string { // On va créer une fonction qui va servir a initialiser notre jeu.
	var motcaché []string // tableau contenant le mot dont certaines lettres vont être remplacées par des "_"
	for range s {
		motcaché = append(motcaché, "_") // ici on ajoute à motcaché un "_" pour chaque caractère de s
	}
	for i := 0; i < (len(s)/2)-1; i++ { // ici on va faire réveler de manière aléatoire un ou plusieurs des caractères du motcaché
		rand.Seed(time.Now().UnixNano())
		pos := rand.Intn(len(s))
		motcaché[pos] = string(s[pos])
	}
	res := "" // ici on créé juste une variable de type string à laquelle on va ajouter tous les caractères de notre motcaché fini
	for i := 0; i < len(motcaché); i++ {
		res = res + motcaché[i]
	}
	return res
}

func islettre(motcomplet string, lettrecherché string, motavecles_ string) bool {
	tab_motcomplet := []rune(motcomplet)       // on initialise un tableau de runes du mot complet quon cherche
	tabentrée := []rune(lettrecherché)         // On créé un tableau de runes contenant l'entrée utilisateur
	tabmotavecles_ := []rune(motavecles_)      // Tableau de runes de notre mot caché
	nb_de_fois_lalettre := 0                   // Compteur de fois où la lettre est dans notre mot
	for i := 0; i < len(tab_motcomplet); i++ { // On fait une boucle ici pour voir combien de fois la lettre donnée est présente dans notre mot complet
		if tabentrée[0] == tab_motcomplet[i] {
			nb_de_fois_lalettre++
		}
	}
	for i := 0; i < len(tab_motcomplet); i++ { // ici on va vérifier si on a déjà découvert la lettre qu'on cherche dans le mot
		if tab_motcomplet[i] == tabentrée[0] {
			compteur := 0 // ce compteur sert à voir si on a autant de fois la lettre cherchée dans le mot caché que dans le mot complet
			for i := 0; i < len(tabmotavecles_); i++ {
				if tabentrée[0] == tabmotavecles_[i] {
					compteur++
				}
				if compteur == nb_de_fois_lalettre { // si on a déjà découvert tous les "a" par exemple, on retourne donc false
					return false
				}
			}
			return true // sinon, on retourne true, la lettre entrée par l'utilisateur n'a pas été découverte et elle fait partie du mot qu'on cherche.
		}
	}
	return false
}

func replace(entreeutilisateur string, motcomplet string, gruyère string) string { // fonction qui va nous servir à modifier le motcaché tout au long du jeu
	tabgruyère := []rune(gruyère)          // tableau de runes de notre mot caché
	tabentrée := []rune(entreeutilisateur) // tableau de runes de notre entrée utilisateur
	for i := 0; i < len(motcomplet); i++ { // on boucle sur la len de motcomplet
		if entreeutilisateur[0] == motcomplet[i] { // si l'entrée utilisateur est égale à une ou plusieurs lettres du mot complet
			tabgruyère[i] = tabentrée[0] // alors on modifie le ou les "_" de tabgruyères par la lettre en question selon sa ou ses positions dans le mot complet.
		}
	}
	return string(tabgruyère) // Enfin on retourne la string correspondante à notre résultat.
}

func ascii_letters(hangman HangManData) {
	err2 := os.Remove("resultat.txt") // Ici on vide le document texte qui stock le resultat qu'on va print pour éviter de print n'importe quoi
	if err2 != nil {                  // là on paramètre l'erreur qui sort si on a un soucis dans le remove
		fmt.Println(err2)
		fmt.Println("a nan là ça bug")
		return // On sort de la fonction si ça a bugé
	}
	file, err := os.OpenFile("resultat.txt", os.O_CREATE|os.O_WRONLY, 0600) // ici on ouvre notre résutat.txt qui est vide en premier lieu (CREATE pour le créer s'il n'existe pas déjà, WRONLY pour ouvrir le fichier ou on peut uniquement écrire dedans)
	defer file.Close()                                                      // Pour fermer le programme quand la fonction se termine
	var textascii []string
	// ----------------------------------------------------------------------//
	aff, err := os.Open("standard.txt")
	if err != nil {
		log.Fatal("failed to read standard.txt")
	}
	scanner := bufio.NewScanner(aff)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		textascii = append(textascii, scanner.Text())
	}
	// ----------------------------------------------------------------------// ouverture de standard.txt et stockage dans textascii.
	var slicealpha []string
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	compteur := 0
	var essai []string
	for _, a := range alphabet {
		slicealpha = append(slicealpha, string(a)) // là on ajoute une à une chaque lettre de notre alphabet dans la slicealpha
	}
	var hangmanword []string // tableau représentant chaque caractère de notre mot en cours de découverte
	for _, a := range hangman.Word {
		hangmanword = append(hangmanword, string(a))
	}
	var hangmantofind []string // tableau représentant chaque caractère du mot à trouver
	for _, a := range hangman.ToFind {
		hangmantofind = append(hangmantofind, string(a))
	}
	for i := range hangman.Word { // on passe le long de notre mot en cours de découverte
		for j := range slicealpha { // pour chacun des caractères de notre mot caché, on boucle sur tout l'alphabet
			compteur++                           // le compteur sert à savoir quelle lettre de l'alphabet nous allons devoir print en ascii art
			if hangmanword[i] == slicealpha[j] { // si le caractère du mot caché qu'on scan est identique à une des lettres de l'alphabet
				minimum := (compteur * 9) + 577 // première ligne dans standard.txt de la première lettre en minuscule (on ajt *9 psq chaque caractère dans standard.txt fait 9 lignes)
				maximum := (compteur * 9) + 586 // dernière ligne dans standard.txt de la première lettre en minuscule (same qu'au dessus, on fait ca pour obtenir le bon caractère en fonction de la position de la lettre dans l'alphabet)
				for _, a := range textascii[minimum:maximum] {
					_, err = file.WriteString("\n")
					_, err = file.WriteString(string(a)) // on écrit alors les caractères lignes par lignes dans file qui est le contenu de resultat.txt
					if err != nil {
						panic(err)
					}
				}
			}
		}

		if hangmanword[i] != hangmantofind[i] { // si le caractère quon souhaite print est l'underscore
			for _, a := range textascii[116:125] { // 116 à 125 car ce sont les lignes correspondantes à "_" dans standard.txt
				_, err = file.WriteString("\n") // on écrit alors les caractères lignes par lignes comme plus haut
				_, err = file.WriteString(string(a))
				if err != nil {
					panic(err)
				}
			}
		}
		compteur = 0 // ici on réinitialise le compteur  à la fin.
	}
	compteur1 := 9 // taille de la lettre ascii
	compteur2 := 1 // sert à savoir combien de lettres nous avons
	for j := 0; j < 9; j++ {
		for i := len(hangman.ToFind); i > 0; i-- {
			compteur2++
			n := 0
			n = 9*compteur2 - compteur1
			fileresult, _ := os.Open("resultat.txt") // on va rouvrir notre résultat.txt
			scanner := bufio.NewScanner(fileresult)  // scanner son contenu
			scanner.Split(bufio.ScanLines)           // split le scanner pour scanner lignes par lignes
			for scanner.Scan() {
				essai = append(essai, scanner.Text()) // on ajoute à essai les lignes de resultat une par une, dans des positions distinctes au sein tableau
			}
			fmt.Print(essai[n]) // on va print chaque ligne de chaque caractère côtes à côtes
			fileresult.Close()  // on close ensuite le fichier pour moins faire rammer
		}
		if j != 9 {
			fmt.Println("") // ici on passe à la ligne à chaque tour de boucle pour print la prochaine ligne de chaque caractère ensuite
		}
		compteur2 = 0 // On réinitialise le compteur2 pour le prochain tour de boucle
		compteur1--   // On décrémente le compteur 1 pour ensuite print la ligne suivante de chaque caractère choisi
	}
}

func main() {
	var hangman HangManData // on créé une variable de type HangManData(voir en haut)
	hangman.Attempts = 10   // On initialise le nomrbe d'erreurs permises à 10
	args := os.Args[1:]     // On créé la variable args qui correspond à ce qu'on peut écrire après le go run hangman
	if len(os.Args) > 1 {   // regarde si la len de os.Args est supérieure à 1 (donc si args peut être utilisé)
		if args[0] == "--startWith" { // si le mot après go run hangman est "--startwith"
			a, _ := os.ReadFile("save.txt") // On va alors lire save.txt dans laquelle les valeurs de la structure HangMan ont été stockées plus bas si on écrit "stop" en input dans choose
			json.Unmarshal(a, &hangman)     // ici on modifie les valeurs de notre variable hangman de type HangMan par les valeurs stockées plus haut dans a
		}
	} else {
		hangman.ToFind = motatrouver()          // sinon on initialise le jeu à 0 avec la fonction motatrouver expliquée plus haut (elle trv un mot o pif dans words.txt qui va etre celui à trouver)
		hangman.Word = débutjeu(hangman.ToFind) // on va créer notre mot caché par rapport au résultat de motatrouver() qu'on a fait au dessus
		hangman.index_tab_pendu = 0             // là on met la variable qui va nous servir à aller chercher les positions du pendu dans notre tableau du pendu en fonction des attemps restantes
	}
	fmt.Println("Good Luck, you have ", hangman.Attempts, " attempts.") // là on print l'output de début de jeu avec le nombre d'attempts restantes
	ascii_letters(hangman)                                              // Ici on lance notre fonction ascii letter sur la struct hangman pour transformer les caractères en ascii art
	f, err := os.Open("hangman.txt")                                    // on ouvre hangman.txt et on stock son contenu dans f
	if err != nil {
		log.Fatal(err)
	}
	//var tabjoser []string          // ici on va créer un tab de string qui va stocker les positions de josé
	scanner := bufio.NewScanner(f) // On initialise un scanner sur f qui contient ttes les pos de notre pendu
	scanner.Split(bufio.ScanLines) // ici on dit au scanner de scanner ligne par ligne
	paragraphe := ""               // On créé une variable paragraphe vide qui va stocker ligne par ligne notre position de hangman
	for scanner.Scan() {           // on boucle tant que le scanner scan
		//compteur := 0

		paragraphe = paragraphe + scanner.Text() + "\n" // adds the value of scanner (that contains the characters from StylizedFile) to source

		/*
			if scanner.Text() == "=========" { // dès que le scanner tombe sur une ligne semblable à celle là
				paragraphe = paragraphe + "\n" + scanner.Text() // On ajoute cette ligne à notre paragraphe quand même
				tabjoser[compteur] = paragraphe                 // on ajoute le paragraphe dans le tableau à la position donnée par le compteur
				paragraphe = ""                                 // on vide le paragraphe pour le prochain tour de boucle
				compteur++                                      // On incrémente le compteur pour changer de position au prochain tour de boucle
			} else {
				paragraphe = paragraphe + "\n" + scanner.Text() // sinon on continue d'ajouter les lignes qu'on scan dans paragraphe
			}
		*/

		//hangman.HangmanPositions = tabjoser // On change la valeur de hangman.HangmanPositions par celle tabjoser
	}

	hangman.HangmanPositions = strings.Split(paragraphe, "=========")

	for index, v := range hangman.HangmanPositions {
		hangman.HangmanPositions[index] = v + "========="
	}

	hangman.HangmanPositions = hangman.HangmanPositions[:len(hangman.HangmanPositions)-1]

	for hangman.Attempts >= 0 { // Tant qu'on a pas perdu
		scanner := bufio.NewScanner(os.Stdin)                           // On créé un scanner qui va lire l'entréeutilisateur (l'ostdin)
		fmt.Print("\nChoose : ")                                        // On print le Choose à chaque tour pour indiquer qu'il faut choisir une lettre
		scanner.Scan()                                                  // On scan l'entréeutilisateur
		entreeUtilisateur := scanner.Text()                             // on stock le résultat du scanner dans cette variable
		if entreeUtilisateur[0] >= 'a' && entreeUtilisateur[0] <= 'z' { // si la lettre ecrite est en minuscule, on la met en maj avec toupper
			entreeUtilisateur = ToUpper(entreeUtilisateur)
		}
		if entreeUtilisateur == "STOP" || entreeUtilisateur == "stop" { // Si l'ostdin est stop ou STOP, on stock les valeurs en cour de hangman dans save.txt
			b, err := json.Marshal(HangManData{hangman.Word, hangman.ToFind, hangman.Attempts, hangman.HangmanPositions, hangman.index_tab_pendu}) // là on encode les valeurs dans b
			if err != nil {
				fmt.Println("Error:", err) // on paramètre l'erreur au cas où ca bug
			}
			ioutil.WriteFile("save.txt", b, 1)   // On écrit les valeurs stockées dans b dans save.txt
			fmt.Print("Game Saved in save.txt.") // On Print ça pour dire que l'opération a bien été effectuée
			break
		}
		if len(scanner.Text()) < 1 || len(scanner.Text()) > 1 { // si lentree utilisateur fait plus ou moins d'un caractère
			fmt.Print("PLEASE ENTER A VALID INPUT.") // Alors on print une erreur
		} else {
			if islettre(hangman.ToFind, entreeUtilisateur, hangman.Word) { // sinon, on passe par islettre qui vérifie si notre entrée est en rapport avec le mot complet cherché et si c'est le cas :
				hangman.Word = (replace(entreeUtilisateur, hangman.ToFind, hangman.Word)) // On ajoute la  ou les lettres manquantes au mot caché
				ascii_letters(hangman)                                                    // Ici on lance notre fonction ascii letter sur la struct hangman pour transformer les caractères en ascii art
				if hangman.Word == hangman.ToFind {                                       // si le mot qu'on a modifié avec replace est égal au mot qu'on cherche
					fmt.Println("Congrats !") // Alors on print congrats pour féliciter le joueur pour sa vicroire
					return                    // puis on sort de notre programme
				}
			} else {
				hangman.Attempts--         // Sinon, on enlève une chance au joueur
				if hangman.Attempts == 0 { // Si le joueur n'a plus de chance
					fmt.Print(hangman.HangmanPositions[hangman.index_tab_pendu], "\nGAME OVER") // on print la dernière position du pendu, puis GAME OVER
					return                                                                      // et on sort de notre fonction.
				} else {
					fmt.Print("Not present in the word, ", hangman.Attempts, " attempts remaining", "\n", hangman.HangmanPositions[hangman.index_tab_pendu], "\n") // sinon on print le message correspondant à une mauvaise réponse
					hangman.index_tab_pendu++                                                                                                                      // On augmente ensuite l'index du tableau du pendu pour print la prochaine position à la prochaine erreur.
				}
			}
		}
	}
}
