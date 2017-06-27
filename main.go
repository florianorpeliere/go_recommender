package neighborhood_model

import (
	"errors"
	"math"
	"log"
)

type ModelUsersItems struct {
	data []User
}

type ItemRating struct {
	itemId float64
	rating float64
}

type User struct {
	Key int
	Similarity float64
	Rating []float64
}

func getPredictionItemBased(usersItems [][]float64, userId int, itemId int, topK int) (float64, error) {
	var usersItemsMeanCentered [][]float64

	// On recentre les notes des utilisateurs sur leur moyenne respective
	for i := 0; i < len(usersItems); i++ {
		currentUserItemsMeanCentered := MinusVec(usersItems[i], Average(filterZero(usersItems[i])))
		usersItemsMeanCentered = append(usersItemsMeanCentered, currentUserItemsMeanCentered)
	}

	// On créé le tableau inverse, items/users
	var itemsUsers [][]float64
	for i := 0; i < len(usersItems[0]); i++ {
		itemsUsers = append(itemsUsers, getColumn(usersItemsMeanCentered, i))
	}

	// On calcule la similarité de l'item itemId avec les autres
	itemsSimilarityArray := make(map[int]float64)
	for i := 0; i < len(itemsUsers); i++ {
		if itemId != i {
			itemsSimilarityArray[i] = AdjustedCosineSim(itemsUsers[itemId], itemsUsers[i])
		}
	}

	log.Print("Similarity with item ", itemId, ": ", itemsSimilarityArray)

	pairList := sortMapByValue(itemsSimilarityArray)
	sumSimRating := float64(0)
	sumSim := float64(0)

	items := getTopK(itemsUsers, pairList, topK, userId)
	for i := 0; i < len(items); i++ {
		sumSimRating += items[i].Similarity * usersItems[userId][items[i].Key]
		sumSim += items[i].Similarity
		log.Print(items[i].Similarity, usersItems[userId][items[i].Key])
	}

	return sumSimRating / sumSim, nil
}

func getPredictionUserBased(usersItems [][]float64, userId int, itemId int, topK int) (float64, error) {
	if (len(usersItems) - 1) < topK {
		return 0, errors.New("Cannot get more top users than all users")
	}

	targetUser := usersItems[userId]

	// On calcule la similarité de l'utilisateur userId avec les autres
	userSimilarityArray := make(map[int]float64)
	for i := 0; i < len(usersItems); i++ {
		if userId != i {
			userSimilarityArray[i] = PearsonSim(usersItems[i], targetUser)
		}
	}

	log.Print("Similarity: ", userSimilarityArray)

	// On récupere le top k des utilisateurs les plus similaires
	pairList := sortMapByValue(userSimilarityArray)
	users := getTopK(usersItems, pairList, topK, itemId)

	sumSimRating := float64(0)
	sumSim := float64(0)

	// On calcule le rating de l'élément itemId par l'utilisateur userId
	for i := 0; i < len(users); i++ {
		sumSimRating += users[i].Similarity * MinusVec(users[i].Rating[itemId:], Average(filterZero(users[i].Rating)))[0]
		sumSim += users[i].Similarity
	}

	return Average(filterZero(targetUser)) + sumSimRating / sumSim, nil
}

func getTopK(usersItems [][]float64, pairList PairList, topK int, onElementId int) ([]User) {
	var out []User
	for i := 0; i < len(pairList); i++ {
		if usersItems[pairList[i].Key][onElementId] != 0 {
			out = append(out, User{pairList[i].Key, pairList[i].Value, usersItems[pairList[i].Key]})
		}
		if len(out) == topK {
			break;
		}
	}

	return out
}

/*
 * Similarité de Pearson
 * a vecteurs des notes de l'utilisateur a
 * b vecteurs des notes de l'utilisateur b
 */
func PearsonSim(a, b []float64) (float64) {
	aCom, bCom, _ := CommuneValue(a, b)
	aComMeanCentered := MinusVec(aCom, Average(filterZero(a)))
	bComMeanCentered := MinusVec(bCom, Average(filterZero(b)))
	numerator, _ := DotProduct(aComMeanCentered, bComMeanCentered)
	denominator := math.Abs(Norm(aComMeanCentered) * Norm(bComMeanCentered))

	return numerator / denominator
}

/*
 * Similarité cosinus ajusté pour l'algorithme basé sur les items
 * a vecteurs des notes des utilisateurs pour l'item a
 * b vecteurs des notes des utilisateurs pour l'item b
 */
func AdjustedCosineSim(a []float64, b []float64) (float64) {
	aComMeanCentered, bComMeanCentered, _ := CommuneValue(a, b)

	numerator, _ := DotProduct(aComMeanCentered, bComMeanCentered)
	denominator := Norm(aComMeanCentered) * Norm(bComMeanCentered)

	return numerator / denominator
}

/*
 * Permet de récupérer la colonne d'une matrice.
 */
func getColumn(usersItems [][]float64, numColumn int) ([]float64){
	var out []float64
	for i := 0; i < len(usersItems); i++ {
		out = append(out, usersItems[i][numColumn])
	}

	return out
}
/*
 * Similarité cosinus
 */
func CosineSim(a, b []float64) (float64) {
	aCom, bCom, _ := CommuneValue(a, b)
	dp, _ := DotProduct(aCom, bCom)
	aNorm := Norm(aCom)
	bNorm := Norm(bCom)
	return dp / (aNorm * bNorm)
}

/*
 * Moyenne d'un vecteur
 */
func Average(a []float64) (float64) {
	prod := float64(0)
	for i := 0; i < len(a); i++ {
		prod += a[i];
	}

	return prod / float64(len(a));
}

/*
 * Produit scalaire de deux vecteurs
 */
func DotProduct(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, errors.New("Vectors have different length")
	}
	prod := float64(0)
	for i := 0; i < len(a); i++ {
		prod += a[i] * b[i]
	}

	return prod, nil
}

/*
 * Filtre les éléments 0 (absence de valeur) d'un vecteur
 */
func filterZero(a []float64) ([]float64) {
	var out []float64
	for i := 0; i < len(a); i++ {
		if a[i] != 0 {
			out = append(out, a[i])
		}
	}

	return out;
}

/*
 * Retourne les vecteurs a et b filtré avec les valeurs communes
 */
func CommuneValue(a, b []float64) ([]float64, []float64, error){
	if len(a) != len(b) {
		return nil, nil, errors.New("Vectors have different length")
	}
	var aCom []float64
	var bCom []float64
	for i := 0; i < len(a); i++ {
		if a[i] != 0 && b[i] != 0 {
			aCom = append(aCom, a[i])
			bCom = append(bCom, b[i])
		}
	}

	return aCom, bCom, nil
}

/*
 * Soustraction de la valeur sub à tous les éléments du vecteur vec
 */
func MinusVec(vec []float64, sub float64) ([]float64) {
	var out []float64
	for i := 0; i < len(vec); i++ {
		diff := float64(0)
		if vec[i] != 0 {
			diff = vec[i] - sub
		}
		out = append(out, diff)
	}

	return out
}

/*
 * Norme d'un vecteur
 */
func Norm(a []float64) (float64) {
	sum := float64(0)
	for i := 0; i < len(a); i++ {
		sum += a[i] * a[i]
	}

	return math.Sqrt(sum)
}

