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
)

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
	re = regexp.MustCompile(`(?i)^rgba?\((\d+),(\d+),(\d+)`)
	result := re.FindStringSubmatch(str)
	if len(result) > 0 {
		r, _ := strconv.Atoi(result[1])
		g, _ := strconv.Atoi(result[2])
		b, _ := strconv.Atoi(result[3])
		return Color{r, g, b}, nil
	}

	// rgb(100%,100%,100%), rgba(100%,100%,100%,100%)
	re = regexp.MustCompile(`(?i)^rgba?\((\d+)%,(\d+)%,(\d+)%`)
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
	re = regexp.MustCompile(`(?i)^hsla?\((\d+),(\d+)%,(\d+)%`)
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
	re = regexp.MustCompile(`(?i)^hsla?\((\d+)%,(\d+)%,(\d+)%`)
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
	flag.Parse()
}

func run() error {
	var filename string
	if args := flag.Args(); len(args) > 0 {
		filename = args[0]
	}

	// パイプでの入力と、引数による処理を分ける。
	switch filename {
	case "", "-":
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

		// fmt.Println(str)
	default:
		oneColorFlagPtr := flag.Bool("one", false, "output only one color")

		// 1色表示フラグが立っているときは、空白を消して、1色のみ出力する。
		if *oneColorFlagPtr {
			str := strings.Join(flag.Args(), "")
			outputColor(str)
			return nil
		}

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
