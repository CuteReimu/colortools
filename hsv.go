package colortools

type HSV struct {
	H float64 // 0-360
	S float64 // 0-1
	V float64 // 0-1
	r, g, b uint32
	finished bool
}

func (H *HSV) RGBA() (uint32, uint32, uint32, uint32) {
	if H.finished {
		return H.r, H.g, H.b, 65535
	}
	H.H -= float64(int(H.H/360.0)) * 360.0
	if H.H < 0 {
		H.H += 360.0
	}
	r, g, b := 0.0, 0.0, 0.0
	i := int(H.H) / 60 % 6
	f := H.H/60 - float64(i)
	p := H.V * (1 - H.S)
	q := H.V * (1 - f*H.S)
	t := H.V * (1 - (1-f)*H.S)
	switch i {
	case 0:
		r, g, b = H.V, t, p
	case 1:
		r, g, b = q, H.V, p
	case 2:
		r, g, b = p, H.V, t
	case 3:
		r, g, b = p, q, H.V
	case 4:
		r, g, b = t, p, H.V
	case 5:
		r, g, b = H.V, p, q
	}
	H.r, H.g, H.b = uint32(r * 65535), uint32(g * 65535), uint32(b * 65535)
	H.finished = true
	return H.r, H.g, H.b, 65535
}

