package alphabets

var Dictionary = map[string][]rune{
	"enUp":    []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	"enLow":   []rune("abcdefghijklmnopqrstuvwxyz"),
	"rusUp":   []rune("АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"),
	"rusLow":  []rune("абвгдеёжзийклмнопрстуфхцчшщъыьэюя"),
	"numbers": []rune("0123456789"),
	"special": []rune("\n !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"),
}

var KeyboardPrintableRunes = []rune(
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ" +
		"абвгдеёжзийклмнопрстуфхцчшщъыьэюя" +
		"0123456789" +
		" !\"#$%&№'()*+,-./:;<=>?@[\\]^_`{|}~\n",
)

var LenPrintableRunes = len(KeyboardPrintableRunes)

var EnglishReference = map[rune]float64{
	'E': 0.123, 'T': 0.096, 'A': 0.081, 'O': 0.079, 'N': 0.072,
	'I': 0.071, 'S': 0.066, 'R': 0.060, 'H': 0.051, 'L': 0.040,
	'D': 0.036, 'C': 0.032, 'U': 0.031, 'P': 0.023, 'F': 0.023,
	'M': 0.022, 'W': 0.020, 'Y': 0.019, 'B': 0.016, 'G': 0.016,
	'V': 0.009, 'K': 0.005, 'Q': 0.002, 'X': 0.002, 'J': 0.001, 'Z': 0.001,
}

var RussianReference = map[rune]float64{
	'О': 0.090, 'Е': 0.072, 'А': 0.062, 'И': 0.062, 'Н': 0.053,
	'Т': 0.053, 'С': 0.045, 'Р': 0.040, 'В': 0.038, 'Л': 0.035,
	'К': 0.028, 'М': 0.026, 'Д': 0.025, 'П': 0.023, 'У': 0.021,
	'Я': 0.018, 'Ы': 0.016, 'З': 0.016, 'Ъ': 0.014, 'Б': 0.014,
	'Г': 0.013, 'Ч': 0.012, 'Й': 0.010, 'Х': 0.009, 'Ж': 0.007,
	'Ю': 0.006, 'Ш': 0.006, 'Ц': 0.004, 'Щ': 0.003, 'Ф': 0.002,
	'Э': 0.003, 'Ь': 0.017, 'Ё': 0.001,
}

func FindSymbolInfo(char rune) (string, int) {
	for key, alphabet := range Dictionary {
		for i, r := range alphabet {
			if r == char {
				return key, i
			}
		}
	}
	return "", -1
}
