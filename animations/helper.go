package animations

// Wraps position back around to 0
func nextPos(pos uint, len uint) uint {

	if (pos + 1) >= len {
		return 0
	}
	return pos + 1
}
