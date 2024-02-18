package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

// структура для хранеиения слов и их количества в итоговом слайсе
type words struct {
	word  string
	count int
}

// регулярка со знаками припенания.
// "–" в конце - это тире, а не дефис. дефис обрабатывается отдельно.
var rePunctSymbl = regexp.MustCompile(`[.,!?'"()\\/\[\]{}–]+`)

// функция возвращающает слайс с 10-ю наиболее часто встречаемыми во входном тексте словами
func Top10(inStr string) []string {
	var (
		wrdsArr []words
		outArr  []string
		mStr    = make(map[string]int)
	)

	// удаляем из текста знаки припенания (меняем их на пробел)
	inStr = rePunctSymbl.ReplaceAllString(strings.ToLower(inStr), " ")
	// переводим строку в слайс слов по пробелам
	sArr := strings.Fields(inStr)

	if len(sArr) == 0 {
		return outArr
	}

	// считаем слова через словарь
	for _, s := range sArr {
		if _, ok := mStr[s]; ok {
			mStr[s] = mStr[s] + 1
		} else if s != "-" { // если слово дефис - пропускаем
			mStr[s] = 1
		}
	}

	// переводим посчитанныц словарь в слайс объектов words
	for str, cnt := range mStr {
		wrdsArr = append(wrdsArr, words{word: str, count: cnt})
	}

	// сортируем лексикографически по словам
	sort.Slice(wrdsArr, func(i, j int) bool { return wrdsArr[i].word < wrdsArr[j].word })
	// сортируем по количеству слов с сохранением лексикографической соритровки слов в равновесных блоках
	sort.SliceStable(wrdsArr, func(i, j int) bool { return wrdsArr[i].count > wrdsArr[j].count })

	// если массив > 10 элементов, берем первые 10
	if len(wrdsArr) > 10 {
		wrdsArr = wrdsArr[:10]
	}

	// формируем массив строк из массива объектов
	for _, wrds := range wrdsArr {
		outArr = append(outArr, wrds.word)
	}

	return outArr
}
