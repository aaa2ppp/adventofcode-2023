// Code generated by "stringer -type=HandType"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[FiveOfKind-6]
	_ = x[FourOfKind-5]
	_ = x[FullHouse-4]
	_ = x[ThreeOfKind-3]
	_ = x[TwoPair-2]
	_ = x[OnePair-1]
	_ = x[HighCard-0]
}

const _HandType_name = "HighCardOnePairTwoPairThreeOfKindFullHouseFourOfKindFiveOfKind"

var _HandType_index = [...]uint8{0, 8, 15, 22, 33, 42, 52, 62}

func (i HandType) String() string {
	if i < 0 || i >= HandType(len(_HandType_index)-1) {
		return "HandType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _HandType_name[_HandType_index[i]:_HandType_index[i+1]]
}
