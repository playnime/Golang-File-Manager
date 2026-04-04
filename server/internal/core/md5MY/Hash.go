package md5MY

func cut(text []byte) (int, []uint32) {
	lenght := uint64(len(text)) * 8
	text = append(text, 0x80)
	ost := ((56 - len(text)%64) + 64) % 64

	for i := 0; i < ost; i++ {
		text = append(text, 0x00)
	}

	for i := 0; i < 8; i++ {
		x := byte(lenght)
		text = append(text, x)
		lenght = lenght >> 8
	}

	blocks := len(text) / 64
	words := make([]uint32, blocks*16)

	for i := 0; i < len(text); i += 4 {
		var word uint32

		word |= uint32(text[i])
		word |= uint32(text[i+1]) << 8
		word |= uint32(text[i+2]) << 16
		word |= uint32(text[i+3]) << 24

		words[i/4] = word
	}

	return blocks, words
}

func select_word(i int) int {
	switch {
	case i >= 0 && i <= 15:
		return i
	case i >= 16 && i <= 31:
		return (5*i + 1) % 16
	case i >= 32 && i <= 47:
		return (3*i + 5) % 16
	case i >= 48 && i <= 63:
		return (7 * i) % 16
	}
	return 0
}

func f(b, c, d uint32, i int) uint32 {
	switch {
	case i >= 0 && i <= 15:
		return (b & c) | (^b & d) // F
	case i >= 16 && i <= 31:
		return (d & b) | (^d & c) // G
	case i >= 32 && i <= 47:
		return b ^ c ^ d // H
	case i >= 48 && i <= 63:
		return c ^ (b | ^d) // I
	}
	return 0
}

func combine(a uint32, b uint32, c uint32, d uint32, i int, word uint32) uint32 {

	t := a + f(b, c, d, i) + word + K[i]
	t = (t << S[i]) | (t >> (32 - S[i]))
	return t + b
}

func processing_one_block(words []uint32, block int, abcd []uint32) []uint32 {

	a, b, c, d := abcd[0], abcd[1], abcd[2], abcd[3]

	for i := 0; i < 64; i++ {
		word := words[block*16+select_word(i)]
		a_ := combine(a, b, c, d, i, word)
		a = d
		d = c
		c = b
		b = a_
	}

	abcd[0] += a
	abcd[1] += b
	abcd[2] += c
	abcd[3] += d

	return abcd
}

func HASH(text []byte) []byte {
	blocks, words := cut(text)

	A0 := uint32(0x67452301)
	B0 := uint32(0xefcdab89)
	C0 := uint32(0x98badcfe)
	D0 := uint32(0x10325476)

	state := []uint32{A0, B0, C0, D0}

	for i := 0; i < blocks; i++ {
		state = processing_one_block(words, i, state)
	}

	a, b, c, d := state[0], state[1], state[2], state[3]

	out := make([]byte, 16)

	// A
	out[0] = byte(a)
	out[1] = byte(a >> 8)
	out[2] = byte(a >> 16)
	out[3] = byte(a >> 24)

	// B
	out[4] = byte(b)
	out[5] = byte(b >> 8)
	out[6] = byte(b >> 16)
	out[7] = byte(b >> 24)

	// C
	out[8] = byte(c)
	out[9] = byte(c >> 8)
	out[10] = byte(c >> 16)
	out[11] = byte(c >> 24)

	// D
	out[12] = byte(d)
	out[13] = byte(d >> 8)
	out[14] = byte(d >> 16)
	out[15] = byte(d >> 24)

	return out
}
