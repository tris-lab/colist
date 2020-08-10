package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var cssColors = map[string]string{
	"black":                "#000000",
	"silver":               "#c0c0c0",
	"gray":                 "#808080",
	"white":                "#ffffff",
	"maroon":               "#800000",
	"red":                  "#ff0000",
	"purple":               "#800080",
	"fuchsia":              "#ff00ff",
	"green":                "#008000",
	"lime":                 "#00ff00",
	"olive":                "#808000",
	"yellow":               "#ffff00",
	"navy":                 "#000080",
	"blue":                 "#0000ff",
	"teal":                 "#008080",
	"aqua":                 "#00ffff",
	"orange":               "#ffa500",
	"aliceblue":            "#f0f8ff",
	"antiquewhite":         "#faebd7",
	"aquamarine":           "#7fffd4",
	"azure":                "#f0ffff",
	"beige":                "#f5f5dc",
	"bisque":               "#ffe4c4",
	"blanchedalmond":       "#ffebcd",
	"blueviolet":           "#8a2be2",
	"brown":                "#a52a2a",
	"burlywood":            "#deb887",
	"cadetblue":            "#5f9ea0",
	"chartreuse":           "#7fff00",
	"chocolate":            "#d2691e",
	"coral":                "#ff7f50",
	"cornflowerblue":       "#6495ed",
	"cornsilk":             "#fff8dc",
	"crimson":              "#dc143c",
	"cyan":                 "#00ffff",
	"darkblue":             "#00008b",
	"darkcyan":             "#008b8b",
	"darkgoldenrod":        "#b8860b",
	"darkgray":             "#a9a9a9",
	"darkgreen":            "#006400",
	"darkgrey":             "#a9a9a9",
	"darkkhaki":            "#bdb76b",
	"darkmagenta":          "#8b008b",
	"darkolivegreen":       "#556b2f",
	"darkorange":           "#ff8c00",
	"darkorchid":           "#9932cc",
	"darkred":              "#8b0000",
	"darksalmon":           "#e9967a",
	"darkseagreen":         "#8fbc8f",
	"darkslateblue":        "#483d8b",
	"darkslategray":        "#2f4f4f",
	"darkslategrey":        "#2f4f4f",
	"darkturquoise":        "#00ced1",
	"darkviolet":           "#9400d3",
	"deeppink":             "#ff1493",
	"deepskyblue":          "#00bfff",
	"dimgray":              "#696969",
	"dimgrey":              "#696969",
	"dodgerblue":           "#1e90ff",
	"firebrick":            "#b22222",
	"floralwhite":          "#fffaf0",
	"forestgreen":          "#228b22",
	"gainsboro":            "#dcdcdc",
	"ghostwhite":           "#f8f8ff",
	"gold":                 "#ffd700",
	"goldenrod":            "#daa520",
	"greenyellow":          "#adff2f",
	"grey":                 "#808080",
	"honeydew":             "#f0fff0",
	"hotpink":              "#ff69b4",
	"indianred":            "#cd5c5c",
	"indigo":               "#4b0082",
	"ivory":                "#fffff0",
	"khaki":                "#f0e68c",
	"lavender":             "#e6e6fa",
	"lavenderblush":        "#fff0f5",
	"lawngreen":            "#7cfc00",
	"lemonchiffon":         "#fffacd",
	"lightblue":            "#add8e6",
	"lightcoral":           "#f08080",
	"lightcyan":            "#e0ffff",
	"lightgoldenrodyellow": "#fafad2",
	"lightgray":            "#d3d3d3",
	"lightgreen":           "#90ee90",
	"lightgrey":            "#d3d3d3",
	"lightpink":            "#ffb6c1",
	"lightsalmon":          "#ffa07a",
	"lightseagreen":        "#20b2aa",
	"lightskyblue":         "#87cefa",
	"lightslategray":       "#778899",
	"lightslategrey":       "#778899",
	"lightsteelblue":       "#b0c4de",
	"lightyellow":          "#ffffe0",
	"limegreen":            "#32cd32",
	"linen":                "#faf0e6",
	"magenta":              "#ff00ff",
	"mediumaquamarine":     "#66cdaa",
	"mediumblue":           "#0000cd",
	"mediumorchid":         "#ba55d3",
	"mediumpurple":         "#9370db",
	"mediumseagreen":       "#3cb371",
	"mediumslateblue":      "#7b68ee",
	"mediumspringgreen":    "#00fa9a",
	"mediumturquoise":      "#48d1cc",
	"mediumvioletred":      "#c71585",
	"midnightblue":         "#191970",
	"mintcream":            "#f5fffa",
	"mistyrose":            "#ffe4e1",
	"moccasin":             "#ffe4b5",
	"navajowhite":          "#ffdead",
	"oldlace":              "#fdf5e6",
	"olivedrab":            "#6b8e23",
	"orangered":            "#ff4500",
	"orchid":               "#da70d6",
	"palegoldenrod":        "#eee8aa",
	"palegreen":            "#98fb98",
	"paleturquoise":        "#afeeee",
	"palevioletred":        "#db7093",
	"papayawhip":           "#ffefd5",
	"peachpuff":            "#ffdab9",
	"peru":                 "#cd853f",
	"pink":                 "#ffc0cb",
	"plum":                 "#dda0dd",
	"powderblue":           "#b0e0e6",
	"rosybrown":            "#bc8f8f",
	"royalblue":            "#4169e1",
	"saddlebrown":          "#8b4513",
	"salmon":               "#fa8072",
	"sandybrown":           "#f4a460",
	"seagreen":             "#2e8b57",
	"seashell":             "#fff5ee",
	"sienna":               "#a0522d",
	"skyblue":              "#87ceeb",
	"slateblue":            "#6a5acd",
	"slategray":            "#708090",
	"slategrey":            "#708090",
	"snow":                 "#fffafa",
	"springgreen":          "#00ff7f",
	"steelblue":            "#4682b4",
	"tan":                  "#d2b48c",
	"thistle":              "#d8bfd8",
	"tomato":               "#ff6347",
	"turquoise":            "#40e0d0",
	"violet":               "#ee82ee",
	"wheat":                "#f5deb3",
	"whitesmoke":           "#f5f5f5",
	"yellowgreen":          "#9acd32",
	"rebeccapurple":        "#663399",
}

// -----------------------------------------------------------------

// Color は、色情報を保持する構造体です。
type Color struct {
	R, G, B int
}

func (c Color) String() string {
	return fmt.Sprintf("R:%3d G:%3d B:%3d", c.R, c.G, c.B)
}

// NewColor は、文字列を元にしてColorを生成するコンストラクタです。
func NewColor(str string) (Color, error) {
	// #FFFFFF, #FFFFFFFFF
	re := regexp.MustCompile(`(?i)^#[0-9A-F]{6}`)
	if re.MatchString(str) {
		r := hex2Decimal(str[1:3])
		g := hex2Decimal(str[3:5])
		b := hex2Decimal(str[5:7])
		return Color{r, g, b}, nil
	}

	// #FFF, #FFFF
	re = regexp.MustCompile(`(?i)^#[0-9A-F]{3,4}$`)
	if re.MatchString(str) {
		r := hex2Decimal(str[1:2])
		g := hex2Decimal(str[2:3])
		b := hex2Decimal(str[3:4])
		return Color{r, g, b}, nil
	}

	// rgb(255,255,255), rgba(255,255,255,100%)
	re = regexp.MustCompile(`(?i)^rgba?\(\s*(\d+)[,\s]\s*(\d+)[,\s]\s*(\d+)`)
	result := re.FindStringSubmatch(str)
	if len(result) > 0 {
		r, _ := strconv.Atoi(result[1])
		g, _ := strconv.Atoi(result[2])
		b, _ := strconv.Atoi(result[3])
		return Color{r, g, b}, nil
	}

	// rgb(100%,100%,100%), rgba(100%,100%,100%,100%)
	re = regexp.MustCompile(`(?i)^rgba?\(\s*(\d+)%[,\s]\s*(\d+)%[,\s]\s*(\d+)%`)
	result = re.FindStringSubmatch(str)
	if len(result) > 0 {
		rp, _ := strconv.Atoi(result[1])
		gp, _ := strconv.Atoi(result[2])
		bp, _ := strconv.Atoi(result[3])
		r := rp * 255 / 100
		g := gp * 255 / 100
		b := bp * 255 / 100
		return Color{r, g, b}, nil
	}

	// hsl(360,100%,100%), hsla(360,100%,100%,100%), hsla(360,100%,100%/100%)
	re = regexp.MustCompile(`(?i)^hsla?\(\s*(\d+)[,\s]\s*(\d+)%[,\s]\s*(\d+)%`)
	result = re.FindStringSubmatch(str)
	if len(result) > 0 {
		h, _ := strconv.Atoi(result[1])
		s, _ := strconv.Atoi(result[2])
		l, _ := strconv.Atoi(result[3])
		r, g, b, e := hsl2Rgb(h, s, l)
		if e != nil {
			return Color{0, 0, 0}, errors.New("invalid argument for convert from hsl")
		}
		return Color{r, g, b}, nil
	}

	// hsl(100%,100%,100%), hsla(100%,100%,100%,100%), hsla(100%,100%,100%/100%)
	re = regexp.MustCompile(`(?i)^hsla?\(\s*(\d+)%[,\s]\s*(\d+)%[,\s]\s*(\d+)%`)
	result = re.FindStringSubmatch(str)
	if len(result) > 0 {
		h, _ := strconv.Atoi(result[1])
		s, _ := strconv.Atoi(result[2])
		l, _ := strconv.Atoi(result[3])
		h1 := h * 360 / 100
		r, g, b, e := hsl2Rgb(h1, s, l)
		if e != nil {
			return Color{0, 0, 0}, errors.New("invalid argument for convert from hsl")
		}
		return Color{r, g, b}, nil
	}

	// CSS色キーワードが渡されたかも。
	key := strings.ToLower(str)
	val, ok := cssColors[key]
	if ok {
		r := hex2Decimal(val[1:3])
		g := hex2Decimal(val[3:5])
		b := hex2Decimal(val[5:7])
		return Color{r, g, b}, nil
	}

	// どれにも当てはまらなかった場合は黒とエラーを返す。
	return Color{0, 0, 0}, errors.New("invalid argument")
}

// Luminance は、輝度を返します。0-255
func (c Color) Luminance() int {
	v := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
	return int(v)
}

// -----------------------------------------------------------------

// 16進数2桁の数字を10進数に変換します。
// FF -> 255, 0F -> 15, F -> FF -> 255, FFF -> 0 (Error!)
func hex2Decimal(str string) int {
	s := str
	l := len(s)

	// 文字列なし または 長過ぎるときは何もしない。
	if l == 0 || l > 2 {
		return 0
	}

	// 1文字の場合は、同じ文字が繰り返したものとして扱う。
	if l == 1 {
		s += s
	}

	v, _ := strconv.ParseInt(s, 16, 0)
	return int(v)
}

// HSLをRGBに変換する。変換できなかった場合は黒色を返す。
// hは度。0以上で、360を超えていても可。
// sとlは%。0-100
func hsl2Rgb(hi, si, li int) (int, int, int, error) {
	// 引数の範囲チェック。
	if hi < 0 || si < 0 || 100 < si || li < 0 || 100 < li {
		return 0, 0, 0, errors.New("hsl invalid argument")
	}

	// hを360度に収める。
	hi %= 360

	h := float64(hi)
	l := float64(li)
	s := float64(si)

	var max, min float64
	var r, g, b float64

	if li < 49 {
		max = 2.55 * (l + l*(s/100))
		min = 2.55 * (l - l*(s/100))
	} else {
		max = 2.55 * (l + (100-l)*(s/100))
		min = 2.55 * (l - (100-l)*(s/100))
	}

	if hi < 60 {
		r = max
		g = min + (max-min)*h/60
		b = min
	} else if hi < 120 {
		r = min + (max-min)*(120-h)/60
		g = max
		b = min
	} else if hi < 180 {
		r = min
		g = max
		b = min + (max-min)*(h-120)/60
	} else if hi < 240 {
		r = min
		g = min + (max-min)*(240-h)/60
		b = max
	} else if hi < 300 {
		r = min + (max-min)*(h-240)/60
		g = min
		b = max
	} else {
		r = max
		g = min
		b = min + (max-min)*(360-h)/60
	}

	return int(r), int(g), int(b), nil
}

// 引数で渡された色(文字列)をコンソールへ出力します。
func outputColor(str string) {
	space := "     "
	c, e := NewColor(str)

	// 色を生成できなかったときは、その旨を出力する。
	if e != nil {
		fmt.Printf("%s ", space)
		// fmt.Printf("\033[48;2;128;128;128m<invalid value> <%s>\033[0m", str)
		fmt.Printf("<invalid value> <%s>", str)
	} else {
		// 背景色を設定して色を描画。
		fmt.Printf("\033[48;2;%d;%d;%dm%s\033[0m ", c.R, c.G, c.B, space)

		// 色の輝度に応じて、文字列描画の背景色を生成。
		// l := c.Luminance()
		// bgc := 20
		// if l < 80 {
		// 	bgc = 230
		// }

		// 値を描画。
		// fmt.Printf("\033[48;2;%d;%d;%dm", bgc, bgc, bgc) // 背景色を設定
		// fmt.Printf("\033[38;2;%d;%d;%dm%s <%s>\033[0m", c.R, c.G, c.B, c, str)
		fmt.Printf("%s <%s>\033[0m", c, str)
	}

	fmt.Println("")
}

// -----------------------------------------------------------------

func init() {
	testing.Init()
	flag.Parse()
}

func run() error {
	var filename string
	if args := flag.Args(); len(args) > 0 {
		filename = args[0]
	}

	// パイプでの入力と、引数による処理を分ける。
	switch filename {
	case "", "-": // パイプでの実行時
		r := os.Stdin
		s, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		str := string(s)

		for _, color := range strings.Split(str, "\n") {
			// 末尾の改行のときなど、空行をスキップする。
			if color == "" {
				continue
			}

			outputColor(color)
		}

	default:
		// 通常時は渡された色でループする。
		for _, arg := range flag.Args() {
			// fmt.Println(arg)
			outputColor(arg)
		}
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
}
