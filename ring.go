/*

さて、では、もう一つ問題を出そう。
五つのビリヤードの玉を、真珠のネックレスのように、リングにつなげてみるとしよう。
玉には、それぞれナンバーが書いてあるな。
さて、この五つの玉のうち、幾つ取っても良いが、隣どうし連続したものしか取れないとしよう。
一つでも二つでも、五つ全部でも良い。
しかし、離れているものは取れない。
この条件で取った玉のナンバーを足し合わせて、1から21までのすべての数ができるようにしたい。
さあ、どのナンバーの玉をどのように並べて、ネックレスを作れば良いかな?
  -- 「笑わない数学者」より

*/

package main

import (
	"fmt"
)

const (
	SIZE  = 8                 //ビリヤードの玉の数
	TOTAL = SIZE*(SIZE-1) + 1 //全ての玉の合計
)

// ビリヤードの玉を配置するリングの状態
type Ring struct {
	value [SIZE]int   //ボールの数字
	flag  [TOTAL]bool //合計値の重複チェック
	size  int         //配置済の玉の個数
	sum   int         //配置済の玉の合計
}

func main() {
	r := Ring{}
	Search(r)
}

// 再帰的に探索
func Search(r Ring) {

	// ゴール
	if r.size == SIZE {
		err := r.Inspect()
		if err != nil {
			return
		}
		fmt.Printf("SIZE = %d\n", SIZE)
		fmt.Printf("%v\n", r.value)
		return
	}

	// 探索継続
	c := NewCandidate(&r)
	for i := 0; i < c.size; i++ {
		next, err := r.Set(c.value[i])
		if err == nil {
			Search(next)
		}
	}
}

// 重複チェックを行い、ボールを配置する
func (old Ring) Set(v int) (Ring, error) {

	// 戻し用
	r := old

	// 部分和を生成して、重複チェック
	r.value[r.size] = v
	for i := 0; i <= r.size; i++ {
		c := 0
		for j := i; j <= r.size; j++ {
			c += r.value[j]
		}
		if err := r.Check(c); err != nil {
			return old, err
		}
	}

	// ボールの配置
	r.size++
	r.sum += v
	return r, nil
}

// 重複チェック
func (r *Ring) Check(v int) error {
	if r.flag[v-1] == true {
		return fmt.Errorf("number %d aready exists.", v)
	}
	r.flag[v-1] = true
	return nil
}

// 全てのボールがセットされた状態で検査
func (r Ring) Inspect() error {

	if r.value[1] > r.value[SIZE-1] {
		return fmt.Errorf("mirror")
	}

	for i := 0; i < SIZE-2; i++ {
		// ローテート
		temp := r.value[0]
		for j := 0; j < SIZE-1; j++ {
			r.value[j] = r.value[j+1]
		}
		r.value[SIZE-1] = temp

		// 検査
		for j := 1; j < SIZE-i-1; j++ {
			c := 0
			for k := j; k < SIZE; k++ {
				c += r.value[k]
			}
			if err := r.Check(c); err != nil {
				return err
			}
		}
	}
	return nil
}

// 候補選定用ワーク
type Candidate struct {
	r     *Ring
	value [(TOTAL + 1) / 2]int //候補の数字
	size  int                  //選定済個数
	max   int                  //選定可能な数字の上限
}

// 候補選定
func NewCandidate(r *Ring) *Candidate {
	c := new(Candidate)
	c.r = r
	c.max = TOTAL - r.sum

	if r.size == 0 {
		c.Set(1)
	} else if r.size == SIZE-1 {
		c.Set(TOTAL - r.sum)
	} else {
		for i := 0; i < c.max; i++ {
			if r.flag[i] == false {
				c.Set(i + 1) //添え字＋１が対応する数になる
			}
		}
	}
	return c
}

// 候補にボールをセット
func (c *Candidate) Set(v int) {
	c.value[c.size] = v
	c.size++
	if c.size < SIZE-c.r.size {
		c.max -= v
	}
}
